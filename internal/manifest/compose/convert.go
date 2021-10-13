package compose

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Converter struct {
	// ProjectName is used as a prefix for the resources created by this importer.
	ProjectName string
}

func (c *Converter) Convert(bs []byte) (*hcl.File, hcl.Diagnostics) {
	project, err := compose.Parse(bytes.NewBuffer(bs))
	if err != nil {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  err.Error(),
			},
		}
	}

	b := exohcl.NewBuilder(bs)
	var diags hcl.Diagnostics

	// Since containers reference networks and volumes by their docker-compose name, but the
	// Docker components will have a namespaced name, so we need to keep track of which
	// volumes/components a service references.
	networksByComposeName := map[string]string{}
	volumesByComposeName := map[string]string{}

	for originalName, volume := range project.Volumes {
		name := exohcl.MangleName(originalName)
		if originalName != name {
			var subject *hcl.Range
			diags = append(diags, exohcl.NewRenameWarning(originalName, name, subject))
		}

		if volume.Name == "" {
			volume.Name = c.prefixedName(originalName, "")
		}
		volumesByComposeName[originalName] = volume.Name

		b.AddComponentBlock(makeComponentBlock("volume", name, volume))
	}

	// Set up networks.
	hasDefaultNetwork := false
	for originalName, network := range project.Networks {
		if originalName == "default" {
			hasDefaultNetwork = true
		}
		name := exohcl.MangleName(originalName)
		if originalName != name {
			var subject *hcl.Range
			diags = append(diags, exohcl.NewRenameWarning(originalName, name, subject))
		}

		// If a `name` key is specified in the network configuration (usually used in conjunction with `external: true`),
		// then we should honor that as the docker network name. Otherwise, we should set the name as
		// `<project_name>_<network_key>`.
		if network.Name == "" {
			network.Name = c.prefixedName(originalName, "")
		}
		networksByComposeName[originalName] = network.Name

		if network.Driver == "" {
			network.Driver = "bridge"
		}

		b.AddComponentBlock(makeComponentBlock("network", name, network))
	}
	// TODO: Docker Compose only creates the default network if there is at least 1 service that does not
	// specify a network. We should do the same.
	if !hasDefaultNetwork {
		componentName := "default"
		name := c.prefixedName(componentName, "")
		networksByComposeName[componentName] = name

		b.AddComponentBlock(makeComponentBlock("network", componentName, map[string]string{
			"name":   name,
			"driver": "bridge",
		}))
	}

	for originalName, service := range project.Services {
		name := exohcl.MangleName(originalName)
		if originalName != name {
			var subject *hcl.Range
			diags = append(diags, exohcl.NewRenameWarning(originalName, name, subject))
		}
		component := exohcl.Component{
			Name: name,
			Type: "container",
		}

		if service.ContainerName == "" {
			// The generated container name intentionally matches the container name generated by Docker Compose
			// when the scale is set at 1. When we address scaling containers, this will need to be updated to
			// use a different suffix for each container.
			service.ContainerName = c.prefixedName(originalName, "1")
		}

		if service.Labels == nil {
			service.Labels = make(compose.Dictionary)
		}
		for k := range service.Labels {
			if strings.HasPrefix(k, "com.docker.compose") {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("service may not specify labels with prefix \"com.docker.compose\", but %q specified %q", originalName, k),
				})
				return nil, diags
			}
		}
		service.Labels["com.docker.compose.project"] = &c.ProjectName
		service.Labels["com.docker.compose.service"] = &originalName

		// Map the docker-compose network name to the name of the docker network that is created.
		defaultNetworkName := networksByComposeName["default"]
		if len(service.Networks) > 0 {
			mappedNetworks := make([]compose.ServiceNetwork, len(service.Networks))
			for i, network := range service.Networks {
				networkName, ok := networksByComposeName[network.Network]
				if !ok {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  fmt.Sprintf("unknown network: %q", network),
					})
					continue
				}
				mappedNetworks[i] = network
				component.DependsOn = append(component.DependsOn, networkName)
			}
			service.Networks = mappedNetworks
		} else {
			service.Networks = []compose.ServiceNetwork{
				{
					Network: defaultNetworkName,
				},
			}
			component.DependsOn = append(component.DependsOn, "default")
		}

		if len(service.Volumes) > 0 {
			for i, volumeMount := range service.Volumes {
				if volumeMount.Type != "volume" {
					continue
				}
				if volumeName, ok := volumesByComposeName[volumeMount.Source]; ok {
					originalName := volumeMount.Source
					service.Volumes[i].Source = volumeName
					component.DependsOn = append(component.DependsOn, originalName)
				}
				// If the volume was not listed under the top-level "volumes" key, then the docker engine
				// will create a new volume that will not be namespaced by the Compose project name.
			}
		}

		for _, dependency := range service.DependsOn.Services {
			if dependency.Condition != "service_started" {
				var subject *hcl.Range
				diags = append(diags, exohcl.NewUnsupportedFeatureWarning(
					fmt.Sprintf("service condition %q", dependency.Service),
					"only service_started is currently supported",
					subject,
				))
			}
			component.DependsOn = append(component.DependsOn, exohcl.MangleName(dependency.Service))
		}

		for idx, link := range service.Links {
			var linkService, linkAlias string
			parts := strings.Split(link, ":")
			switch len(parts) {
			case 1:
				linkService = parts[0]
				linkAlias = parts[0]
			case 2:
				linkService = parts[0]
				linkAlias = parts[1]
			default:
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("expected SERVICE or SERVICE:ALIAS for link, but got: %q", link),
				})
				return nil, diags
			}
			// NOTE [RESOLVING SERVICE CONTAINERS]:
			// There are several locations in a compose definition where a service may reference another service
			// by the compose name. We currently handle these situations by rewriting these locations to reference
			// a container named `<project>_<mangled_service_name>_1` with the assumption that a container will
			// be created by that name. However, this will break when the referenced service specifies a non-default
			// container name. Additionally, we may want to handle cases where a service is scaled past a single
			// container.
			// Some of these values could/should be resolved at runtime, and we should do it when we have the entire
			// project graph available.

			// See https://github.com/docker/compose/blob/v2.0.0-rc.3/compose/service.py#L836 for how compose configures
			// links.
			mangledServiceName := exohcl.MangleName(linkService)
			containerName := c.prefixedName(mangledServiceName, "1")
			service.Links[idx] = fmt.Sprintf("%s:%s", containerName, linkAlias)
			component.DependsOn = append(component.DependsOn, mangledServiceName)
		}

		component.Spec = yamlutil.MustMarshalString(service)
		b.AddComponentBlock(makeComponentBlock("volume", name, component))
	}

	return b.Build(), diags
}

func (c *Converter) prefixedName(name string, suffix string) string {
	var out strings.Builder
	out.WriteString(c.ProjectName)
	out.WriteByte('_')
	out.WriteString(name)
	if suffix != "" {
		out.WriteByte('_')
		out.WriteString(suffix)
	}

	return out.String()
}

func makeComponentBlock(typ string, name string, spec interface{}) *hclsyntax.Block {
	obj := yamlToHCL(spec).(*hclsyntax.ObjectConsExpr)
	attrs := make([]*hclsyntax.Attribute, len(obj.Items))
	for i, item := range obj.Items {
		key := item.KeyExpr.(*hclsyntax.TemplateExpr)
		if len(key.Parts) != 1 {
			panic("unexpected multi-part template")
		}
		val := item.ValueExpr
		attrs[i] = &hclsyntax.Attribute{
			Name:      key.Parts[0].(*hclsyntax.LiteralValueExpr).Val.AsString(),
			NameRange: key.SrcRange,
			Expr:      val,
			SrcRange:  hcl.RangeBetween(key.SrcRange, val.Range()),
		}
	}
	return &hclsyntax.Block{
		Type:   typ,
		Labels: []string{name},
		Body: &hclsyntax.Body{
			Attributes: exohcl.NewAttributes(attrs...),
		},
	}
}

func yamlToHCL(v interface{}) hclsyntax.Expression {
	marshaller, ok := v.(yaml.InterfaceMarshaler)
	if ok {
		v, err := marshaller.MarshalYAML()
		if err != nil {
			panic(err)
		}
		return yamlToHCL(v)
	}

	switch v := v.(type) {
	case nil:
		return exohcl.NewNullLiteral(hcl.Range{})
	case string:
		return exohcl.NewStringLiteral(v, hcl.Range{})
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return exohcl.NewNullLiteral(hcl.Range{})
			}
			rv = rv.Elem()

		}
		switch rv.Kind() {
		case reflect.Struct:

			typ := rv.Type()

			numField := rv.NumField()
			items := make([]hclsyntax.ObjectConsItem, 0, numField)
			for i := 0; i < numField; i++ {
				fld := typ.Field(i)
				tag := fld.Tag.Get("yaml")
				if tag == "" {
					continue
				}
				options := strings.Split(tag, ",")
				name := options[0]
				omitempty := false
				for _, option := range options[1:] {
					switch option {
					case "omitempty":
						omitempty = true
					default:
						panic(fmt.Errorf("unsupported yaml field tag option: %q", option))
					}
				}

				fldV := rv.Field(i)

				if omitempty && fldV.IsZero() {
					continue
				}

				item := hclsyntax.ObjectConsItem{
					KeyExpr:   exohcl.NewStringLiteral(name, hcl.Range{}),
					ValueExpr: yamlToHCL(fldV.Interface()),
				}
				items = append(items, item)
			}
			return &hclsyntax.ObjectConsExpr{
				Items: items,
			}

		case reflect.Map:
			items := make([]hclsyntax.ObjectConsItem, 0, rv.Len())
			iter := rv.MapRange()
			for iter.Next() {
				item := hclsyntax.ObjectConsItem{
					KeyExpr:   yamlToHCL(iter.Key().Interface()),
					ValueExpr: yamlToHCL(iter.Value().Interface()),
				}
				items = append(items, item)
			}
			return &hclsyntax.ObjectConsExpr{
				Items: items,
			}

		default:
			panic(fmt.Errorf("unexpected yaml type: %T", v))
		}
	}
}
