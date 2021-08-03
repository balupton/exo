package container

import (
	"github.com/deref/exo/providers/docker/compose"
	"github.com/deref/exo/util/logging"
	docker "github.com/docker/docker/client"
)

type Container struct {
	ComponentID string
	Spec
	State

	Logger        logging.Logger
	WorkspaceRoot string
	Docker        *docker.Client
	SyslogPort    int
}

type Spec compose.Service

type State struct {
	ImageID     string `json:"imageId"`
	ContainerID string `json:"containerId"`
	Running     bool   `json:"running"`
}
