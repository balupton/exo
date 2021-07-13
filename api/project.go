package api

import (
	"context"
	"net/http"

	"github.com/deref/exo/config"
	"github.com/deref/exo/josh"
)

type Project interface {
	// Deletes all of the components in the project, then deletes the project itself.
	Delete(context.Context, *DeleteInput) (*DeleteOutput, error)
	// Performs creates, updates, refreshes, disposes, as needed.
	Apply(context.Context, *ApplyInput) (*ApplyOutput, error)
	// Returns component descriptions.
	DescribeComponents(context.Context, *DescribeComponentsInput) (*DescribeComponentsOutput, error)
	// Creates a component and triggers an initialize lifecycle event.
	CreateComponent(context.Context, *CreateComponentInput) (*CreateComponentOutput, error)
	// Replaces the spec on a component and triggers an update lifecycle event.
	UpdateComponent(context.Context, *UpdateComponentInput) (*UpdateComponentOutput, error)
	// Triggers a refresh lifecycle event to update the component's state.
	RefreshComponent(context.Context, *RefreshComponentInput) (*RefreshComponentOutput, error)
	// Marks a component as disposed and triggers the dispose lifecycle event.
	// After being disposed, the component record will be deleted asynchronously.
	DisposeComponent(context.Context, *DisposeComponentInput) (*DisposeComponentOutput, error)
	// Disposes a component and then awaits the record to be deleted synchronously.
	DeleteComponent(context.Context, *DeleteComponentInput) (*DeleteComponentOutput, error)
}

type DeleteInput struct{}

type DeleteOutput struct{}

type ApplyInput struct {
	Config config.Config `json:"config"`
}

type ApplyOutput struct{}

type DescribeComponentsInput struct{}

type DescribeComponentsOutput struct {
	Components []ComponentDescription `json:"components"`
}

type ComponentDescription struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Spec        string  `json:"spec"`
	State       string  `json:"state"`
	Created     string  `json:"created"`
	Initialized *string `json:"initialized"`
	Disposed    *string `json:"disposed"`
}

type CreateComponentInput struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Spec string `json:"spec"`
}

type CreateComponentOutput struct {
	ID string `json:"id"`
}

type UpdateComponentInput struct {
	Name string `json:"name"`
	Spec string `json:"spec"`
}

type UpdateComponentOutput struct{}

type DisposeComponentInput struct {
	Name string `json:"name"`
}

type RefreshComponentInput struct {
	Name string `json:"name"`
}

type RefreshComponentOutput struct{}

type DisposeComponentOutput struct{}

type DeleteComponentInput struct {
	Name string `json:"name"`
}

type DeleteComponentOutput struct{}

func NewProjectMux(prefix string, project Project) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(prefix+"delete", josh.NewMethodHandler(project.Delete))
	mux.Handle(prefix+"apply", josh.NewMethodHandler(project.Apply))
	mux.Handle(prefix+"describe-components", josh.NewMethodHandler(project.DescribeComponents))
	mux.Handle(prefix+"create-component", josh.NewMethodHandler(project.CreateComponent))
	mux.Handle(prefix+"update-component", josh.NewMethodHandler(project.UpdateComponent))
	mux.Handle(prefix+"refresh-component", josh.NewMethodHandler(project.RefreshComponent))
	mux.Handle(prefix+"dispose-component", josh.NewMethodHandler(project.DisposeComponent))
	mux.Handle(prefix+"delete-component", josh.NewMethodHandler(project.DeleteComponent))
	return mux
}
