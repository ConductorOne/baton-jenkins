package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-jenkins/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	memberEntitlement = "member"
	adminEntitlement  = "admin"
)

type groupBuilder struct {
	resourceType *v2.ResourceType
	client       *client.JenkinsClient
}

var groupEntitlementAccessLevels = []string{
	memberEntitlement,
	adminEntitlement,
}

// groupResource gets a new connector resource for a Jenkins group.
func groupResource(ctx context.Context, group client.Group, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"group_id":   group.ID,
		"group_name": group.ID,
	}

	groupTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	ret, err := rs.NewGroupResource(
		group.ID,
		resourceTypeGroup,
		group.ID,
		groupTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (g *groupBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return g.resourceType
}

// List returns all the groups from the database as resource objects.
// Groups include a GroupTrait because they are the 'shape' of a standard group.
func (g *groupBuilder) List(ctx context.Context, parentId *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var ret []*v2.Resource
	defaultGroup := client.Group{
		Description: "Default group",
		FullName:    "Authenticated Users",
		ID:          "authenticated",
	}
	groups, err := g.client.GetGroups(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	groups = append(groups, defaultGroup)
	for _, group := range groups {
		res, err := groupResource(ctx, group, parentId)
		if err != nil {
			return nil, "", nil, err
		}

		ret = append(ret, res)
	}

	return ret, "", nil, nil
}

func (g *groupBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	for _, level := range groupEntitlementAccessLevels {
		rv = append(rv, ent.NewPermissionEntitlement(resource, level,
			ent.WithDisplayName(fmt.Sprintf("%s Group %s", resource.DisplayName, titleCase(level))),
			ent.WithDescription(fmt.Sprintf("Access to %s group in Jenkins", resource.DisplayName)),
			ent.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("group:%s:role:%s", resource.Id.Resource, level),
			}),
			ent.WithGrantableTo(resourceTypeGroup),
		))
	}

	return rv, "", nil, nil
}

func (g *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, token *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var rv []*v2.Grant
	return rv, "", nil, nil
}

func (g *groupBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	return nil, nil
}

func (g *groupBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	return nil, nil
}

func newGroupBuilder(client *client.JenkinsClient) *groupBuilder {
	return &groupBuilder{
		resourceType: resourceTypeGroup,
		client:       client,
	}
}
