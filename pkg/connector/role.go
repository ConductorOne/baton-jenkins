package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-jenkins/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	gr "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type roleBuilder struct {
	resourceType *v2.ResourceType
	client       *client.JenkinsClient
}

const NF = -1

// Create a new connector resource for a jenkins role.
func roleResource(ctx context.Context, role string, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"node_id":   role,
		"node_name": role,
	}

	groupTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	ret, err := rs.NewGroupResource(
		role,
		resourceTypeRole,
		role,
		groupTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return r.resourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (r *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var (
		rv []*v2.Resource
	)
	roles, err := r.client.GetAllRoles(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, role := range roles {
		nr, err := roleResource(ctx, role.RoleName, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, nr)
	}

	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (r *roleBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	// role does not have an id
	permission := resource.DisplayName
	// create entitlement for each role
	permissionOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeUser, resourceTypeGroup),
		ent.WithDisplayName(fmt.Sprintf("%s Role %s", resource.DisplayName, permission)),
		ent.WithDescription(fmt.Sprintf("%s access to %s - %s role in Jenkins", titleCase(permission), resource.Id.Resource, resource.DisplayName)),
	}
	rv = append(rv, ent.NewPermissionEntitlement(
		resource,
		permission,
		permissionOptions...,
	))

	return rv, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var (
		err       error
		rv        []*v2.Grant
		userType  = "USER"
		groupType = "GROUP"
	)
	roles, err := r.client.GetAllRoles(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, role := range roles {
		if role.RoleName != resource.Id.Resource {
			continue
		}

		for _, rd := range role.RoleDetail {
			switch rd.Type {
			case userType:
				user := client.Users{
					User: client.User{
						ID: rd.Sid,
					},
				}
				ur, err := userResource(ctx, user, resource.Id)
				if err != nil {
					return nil, "", nil, fmt.Errorf("error creating user resource for role %s: %w", resource.Id.Resource, err)
				}

				tr := gr.NewGrant(resource, role.RoleName, ur.Id)
				rv = append(rv, tr)
			case groupType:
				group := client.Group{
					ID: rd.Sid,
				}
				ur, err := groupResource(ctx, group, resource.Id)
				if err != nil {
					return nil, "", nil, fmt.Errorf("error creating user resource for role %s: %w", resource.Id.Resource, err)
				}

				tr := gr.NewGrant(resource, role.RoleName, ur.Id)
				rv = append(rv, tr)
			}
		}
	}

	return rv, "", nil, nil
}

func (r *roleBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	if principal.Id.ResourceType != resourceTypeUser.Id && principal.Id.ResourceType != resourceTypeGroup.Id {
		l.Warn(
			"jenkins-connector: only users or groups can be granted repo membership",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("jenkins-connector: only users or groups can be granted repo membership")
	}

	return nil, nil
}

func (r *roleBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	principal := grant.Principal
	entitlement := grant.Entitlement
	principalIsUser := principal.Id.ResourceType == resourceTypeUser.Id
	principalIsGroup := principal.Id.ResourceType == resourceTypeGroup.Id
	if !principalIsUser && !principalIsGroup {
		l.Warn(
			"jenkins-connector: only users and groups can have repository permissions revoked",
			zap.String("principal_id", principal.Id.Resource),
			zap.String("principal_type", principal.Id.ResourceType),
		)

		return nil, fmt.Errorf("jenkins-connector: only users and groups can have repository permissions revoked")
	}

	_, _, err := ParseEntitlementID(entitlement.Id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func newRoleBuilder(client *client.JenkinsClient) *roleBuilder {
	return &roleBuilder{
		resourceType: resourceTypeRole,
		client:       client,
	}
}
