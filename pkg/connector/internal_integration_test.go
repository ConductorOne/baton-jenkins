package connector

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/conductorone/baton-jenkins/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/stretchr/testify/assert"
)

var (
	ctx         = context.Background()
	userName, _ = os.LookupEnv("BATON_JENKINS_USERNAME")
	password, _ = os.LookupEnv("BATON_JENKINS_PASSWORD")
	token, _    = os.LookupEnv("BATON_JENKINS_TOKEN")
	baseUrl     = "http://localhost:8080"
)

func TestResourceTypeGrantFails(t *testing.T) {
	var roleEntitlement, roleId, userId string
	if userName == "" && (password == "" || token == "") {
		t.Skip()
	}

	grantEntitlement := "role:reviewer:reviewer"
	grantPrincipal := "localuser"
	grantPrincipalType := "user"
	_, data, err := ParseEntitlementID(grantEntitlement)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	roleId = data[1]
	roleEntitlement = data[2]
	userId = grantPrincipal
	users := getUserForTesting(userId, userId)
	principal, err := userResource(ctx, *users, nil)
	assert.Nil(t, err)
	resource, err := roleResource(ctx, roleId, nil)
	assert.Nil(t, err)
	entitlement := getEntitlementForTesting(resource, grantPrincipalType, roleEntitlement)
	cli := getJenkinsClientForTesting()
	roleBuilder := getRoleBuilderForTesting(cli)
	_, err = roleBuilder.Grant(ctx, principal, entitlement)
	assert.NotNil(t, err)
	errMsg := fmt.Sprintf("jenkins-connector: user %s already has this role permission", principal.DisplayName)
	assert.Equal(t, err.Error(), errMsg, errMsg)
}

func TestResourceTypeGrant(t *testing.T) {
	var roleEntitlement, roleId, userId string
	if userName == "" && (password == "" || token == "") {
		t.Skip()
	}

	grantEntitlement := "role:reviewer:reviewer"
	grantPrincipal := "localuser"
	grantPrincipalType := "user"
	_, data, err := ParseEntitlementID(grantEntitlement)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	roleId = data[1]
	roleEntitlement = data[2]
	userId = grantPrincipal
	users := getUserForTesting(userId, userId)
	principal, err := userResource(ctx, *users, nil)
	assert.Nil(t, err)
	resource, err := roleResource(ctx, roleId, nil)
	assert.Nil(t, err)
	entitlement := getEntitlementForTesting(resource, grantPrincipalType, roleEntitlement)
	cli := getJenkinsClientForTesting()
	roleBuilder := getRoleBuilderForTesting(cli)
	_, err = roleBuilder.Grant(ctx, principal, entitlement)
	assert.Nil(t, err)
}

func TestResourceTypeRevokeFails(t *testing.T) {
	// --revoke-grant "role:reviewer:reviewer:user:localuser"
	var roleId, userId string
	if userName == "" && (password == "" || token == "") {
		t.Skip()
	}

	revokeGrant := "role:reviewer:reviewer:user:localuser"
	_, roleData, err := ParseGrantID(revokeGrant)
	assert.Nil(t, err)
	assert.NotNil(t, roleData)
	grantEntitlement := fmt.Sprintf("%s:%s:%s", roleData[0], roleData[1], roleData[2])
	grantPrincipal := roleData[4]
	_, data, err := ParseEntitlementID(grantEntitlement)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	roleId = data[1]
	userId = grantPrincipal
	users := getUserForTesting(userId, userId)
	principal, err := userResource(ctx, *users, nil)
	assert.Nil(t, err)
	resource, err := roleResource(ctx, roleId, nil)
	assert.Nil(t, err)
	cli := getJenkinsClientForTesting()
	roleBuilder := getRoleBuilderForTesting(cli)
	gr := grant.NewGrant(resource, roleId, principal.Id)
	annos := annotations.Annotations(gr.Annotations)
	v1Identifier := &v2.V1Identifier{
		Id: V1GrantID(V1MembershipEntitlementID(roleId), userId),
	}
	annos.Update(v1Identifier)
	gr.Annotations = annos
	_, err = roleBuilder.Revoke(ctx, gr)
	assert.NotNil(t, err)
	errMsg := fmt.Sprintf("jenkins-connector: user %s does not have this role", userId)
	assert.Equal(t, err.Error(), errMsg, errMsg)
}

func TestResourceTypeRevoke(t *testing.T) {
	// --revoke-grant "role:reviewer:reviewer:user:localuser"
	var roleId, userId string
	if userName == "" && (password == "" || token == "") {
		t.Skip()
	}

	revokeGrant := "role:reviewer:reviewer:user:localuser"
	_, roleData, err := ParseGrantID(revokeGrant)
	assert.Nil(t, err)
	assert.NotNil(t, roleData)
	grantEntitlement := fmt.Sprintf("%s:%s:%s", roleData[0], roleData[1], roleData[2])
	grantPrincipal := roleData[4]
	_, data, err := ParseEntitlementID(grantEntitlement)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	roleId = data[1]
	userId = grantPrincipal
	users := getUserForTesting(userId, userId)
	principal, err := userResource(ctx, *users, nil)
	assert.Nil(t, err)
	resource, err := roleResource(ctx, roleId, nil)
	assert.Nil(t, err)
	cli := getJenkinsClientForTesting()
	roleBuilder := getRoleBuilderForTesting(cli)
	gr := grant.NewGrant(resource, roleId, principal.Id)
	annos := annotations.Annotations(gr.Annotations)
	v1Identifier := &v2.V1Identifier{
		Id: V1GrantID(V1MembershipEntitlementID(roleId), userId),
	}
	annos.Update(v1Identifier)
	gr.Annotations = annos
	_, err = roleBuilder.Revoke(ctx, gr)
	assert.Nil(t, err)
}

func getRoleBuilderForTesting(client *client.JenkinsClient) *roleBuilder {
	return &roleBuilder{
		resourceType: resourceTypeRole,
		client:       client,
	}
}

func getUserForTesting(userId, fullName string) *client.Users {
	return &client.Users{
		User: client.User{
			ID:       userId,
			FullName: fullName,
		},
	}
}

func getJenkinsClientForTesting() *client.JenkinsClient {
	cli := client.NewClient()
	cli.WithUser(userName).WithPassword(password).WithBearerToken(token).WithBaseUrl(baseUrl)
	cli.SetClient(getClientForTesting(ctx))
	return cli
}

func getClientForTesting(ctx context.Context) *uhttp.BaseHttpClient {
	httpClient, _ := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	return uhttp.NewBaseHttpClient(httpClient)
}

func getEntitlementForTesting(resource *v2.Resource, resourceDisplayName, roleEntitlement string) *v2.Entitlement {
	options := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeRole),
		ent.WithDisplayName(fmt.Sprintf("%s Role %s", resourceDisplayName, roleEntitlement)),
		ent.WithDescription(fmt.Sprintf("%s of %s Jenkins role", roleEntitlement, resourceDisplayName)),
	}

	return ent.NewAssignmentEntitlement(resource, roleEntitlement, options...)
}
