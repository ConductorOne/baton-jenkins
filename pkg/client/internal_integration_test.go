package client

import (
	"context"
	"os"
	"testing"

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

func TestJenkinsClient_GetNodes(t *testing.T) {
	if userName == "" && password == "" && token == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetNodes(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func TestJenkinsClient_GetJobs(t *testing.T) {
	if userName == "" && password == "" && token == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetJobs(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func getClientForTesting(ctx context.Context) *uhttp.BaseHttpClient {
	httpClient, _ := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	return uhttp.NewBaseHttpClient(httpClient)
}

func getJenkinsClientForTesting() *JenkinsClient {
	return &JenkinsClient{
		auth: &auth{
			user:        userName,
			password:    password,
			bearerToken: token,
		},
		httpClient: getClientForTesting(ctx),
		baseUrl:    baseUrl,
	}
}

func TestJenkinsClient_GetViews(t *testing.T) {
	if userName == "" && password == "" && token == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetViews(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func TestJenkinsClient_GetUsers(t *testing.T) {
	if userName == "" && password == "" && token == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetUsers(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func TestJenkinsClient_GetRoles(t *testing.T) {
	if userName == "" && password == "" && token == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetAllRoles(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func TestJenkinsClient_GetGroups(t *testing.T) {
	if userName == "" && password == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetGroups(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func TestJenkinsClient_GetAllRoles(t *testing.T) {
	if userName == "" && password == "" && token == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetAllRoles(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func TestJenkinsClient_SetRoles(t *testing.T) {
	if userName == "" && password == "" && token == "" {
		t.Skip()
	}

	roleName := "reviewer"
	userName := "localuser"
	cli := getJenkinsClientForTesting()
	roles, err := cli.AssignUserRole(ctx, roleName, userName)
	assert.Nil(t, err)
	assert.NotNil(t, roles)
}
