package connector

import (
	"context"

	"github.com/conductorone/baton-jenkins/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type viewBuilder struct {
	resourceType *v2.ResourceType
	client       *client.JenkinsClient
}

// Create a new connector resource for a 1Password group.
func viewResource(ctx context.Context, view client.View, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"node_id":   view.Name,
		"node_name": view.Name,
	}

	groupTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	ret, err := rs.NewGroupResource(
		view.Name,
		resourceTypeView,
		view.Name,
		groupTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (v *viewBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return v.resourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (v *viewBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var rv []*v2.Resource
	views, err := v.client.GetViews(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, view := range views {
		nr, err := viewResource(ctx, view, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, nr)
	}

	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (v *viewBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (v *viewBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newViewBuilder(client *client.JenkinsClient) *viewBuilder {
	return &viewBuilder{
		resourceType: resourceTypeView,
		client:       client,
	}
}
