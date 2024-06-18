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
	baseUrl     = "http://localhost:8080"
)

func TestJenkinsClient_GetNodes(t *testing.T) {
	if userName == "" && password == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetNodes(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func TestJenkinsClient_GetJobs(t *testing.T) {
	if userName == "" && password == "" {
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
			bearerToken: "",
		},
		httpClient:   getClientForTesting(ctx),
		baseUrl:      baseUrl,
		jenkinsCache: GoCache{},
	}
}

func TestJenkinsClient_GetViews(t *testing.T) {
	if userName == "" && password == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetViews(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}

func TestJenkinsClient_GetUsers(t *testing.T) {
	if userName == "" && password == "" {
		t.Skip()
	}

	cli := getJenkinsClientForTesting()
	nodes, err := cli.GetUsers(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodes)
}
