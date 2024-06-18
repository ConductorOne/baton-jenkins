package connector

import (
	"context"

	"github.com/conductorone/baton-jenkins/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type nodeBuilder struct {
	resourceType *v2.ResourceType
	client       *client.JenkinsClient
}

// Create a new connector resource for a 1Password group.
func nodeResource(ctx context.Context, role client.Computer, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"node_id":   role.AssignedLabels[0].Name,
		"node_name": role.AssignedLabels[0].Name,
	}

	groupTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	ret, err := rs.NewGroupResource(
		role.AssignedLabels[0].Name,
		resourceTypeNode,
		role.AssignedLabels[0].Name,
		groupTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (n *nodeBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return n.resourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (n *nodeBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var rv []*v2.Resource
	nodes, err := n.client.GetNodes(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, node := range nodes {
		nr, err := nodeResource(ctx, node, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, nr)
	}

	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (n *nodeBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (n *nodeBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newNodeBuilder(client *client.JenkinsClient) *nodeBuilder {
	return &nodeBuilder{
		resourceType: resourceTypeNode,
		client:       client,
	}
}
