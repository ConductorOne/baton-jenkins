package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type JenkinsClient struct {
	auth         *auth
	httpClient   *uhttp.BaseHttpClient
	baseUrl      string
	jenkinsCache GoCache
}

type JenkinsError struct {
	ErrorMessage     string                   `json:"error"`
	ErrorDescription string                   `json:"error_description"`
	ErrorCode        int                      `json:"errorCode,omitempty"`
	ErrorSummary     string                   `json:"errorSummary,omitempty" toml:"error_description"`
	ErrorLink        string                   `json:"errorLink,omitempty"`
	ErrorId          string                   `json:"errorId,omitempty"`
	ErrorCauses      []map[string]interface{} `json:"errorCauses,omitempty"`
}

func (b *JenkinsError) Error() string {
	return b.ErrorMessage
}

// GET - http://{baseurl}/computer/api/json
const (
	allNodes = "computer/api/json?pretty&tree=computer[displayName,description,idle,manualLaunchAllowed,assignedLabels[name]]"
	allJobs  = "api/json?pretty&tree=jobs[name,url,color,buildable]"
)

type auth struct {
	user, password string
	bearerToken    string
}

func NewClient() *JenkinsClient {
	return &JenkinsClient{
		httpClient: &uhttp.BaseHttpClient{},
		baseUrl:    "http://localhost:8080",
		auth: &auth{
			user:        "",
			password:    "",
			bearerToken: "",
		},
		jenkinsCache: NewGoCache(10, 10),
	}
}

func (d *JenkinsClient) WithUser(jenkinsUsername string) *JenkinsClient {
	d.auth.user = jenkinsUsername
	return d
}

func (d *JenkinsClient) WithPassword(jenkinsPassword string) *JenkinsClient {
	d.auth.password = jenkinsPassword
	return d
}

func (d *JenkinsClient) WithBearerToken(jenkinsToken string) *JenkinsClient {
	d.auth.bearerToken = jenkinsToken
	return d
}

func WithAuthorizationBearerHeader(token string) uhttp.RequestOption {
	return uhttp.WithHeader("Authorization", "Bearer "+token)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func WithSetBasicAuthHeader(username, password string) uhttp.RequestOption {
	return uhttp.WithHeader("Authorization", "Basic "+basicAuth(username, password))
}

func WithSetBearerAuthHeader(token string) uhttp.RequestOption {
	return uhttp.WithHeader("Authorization", "Bearer "+token)
}

func WithAuthorization(username, password, token string) uhttp.RequestOption {
	if token != "" {
		return WithSetBearerAuthHeader(token)
	}

	return WithSetBasicAuthHeader(username, password)
}

func (d *JenkinsClient) getToken() string {
	return d.auth.bearerToken
}

func (d *JenkinsClient) getUser() string {
	return d.auth.user
}

func (d *JenkinsClient) getPWD() string {
	return d.auth.password
}

func (d *JenkinsClient) CheckCredentials() bool {
	if d.IsBasicAuthentication() || d.getToken() != "" {
		return true
	}

	return false
}

func (d *JenkinsClient) IsBasicAuthentication() bool {
	if d.getUser() != "" && d.getPWD() != "" {
		return true
	}

	return false
}

func (d *JenkinsClient) IsTokenAuthentication() bool {
	return d.getPWD() != ""
}

func isValidUrl(baseUrl string) bool {
	u, err := url.Parse(baseUrl)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func New(ctx context.Context, baseUrl string, jenkinsClient *JenkinsClient) (*JenkinsClient, error) {
	var (
		clientId     = jenkinsClient.getUser()
		clientSecret = jenkinsClient.getPWD()
		clientToken  = jenkinsClient.getToken()
	)
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	cli := uhttp.NewBaseHttpClient(httpClient)
	if !isValidUrl(baseUrl) {
		return nil, fmt.Errorf("the url : %s is not valid", baseUrl)
	}

	// basic authentication or token
	jc := JenkinsClient{
		httpClient: cli,
		baseUrl:    baseUrl,
		auth: &auth{
			user:        clientId,
			password:    clientSecret,
			bearerToken: clientToken,
		},
		jenkinsCache: NewGoCache(10, 10),
	}

	return &jc, nil
}

// GetNodes
// Get all nodes. Only authenticated users may call this resource.
func (d *JenkinsClient) GetNodes(ctx context.Context) ([]Computer, error) {
	var nodeData NodesAPIData
	endpointUrl := fmt.Sprintf("%s/%s", d.baseUrl, allNodes)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		uhttp.WithHeader("Accept", "application/xml"),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&nodeData))
	if err != nil {
		return nil, &JenkinsError{
			ErrorMessage:     err.Error(),
			ErrorDescription: err.Error(),
			ErrorCode:        resp.StatusCode,
			ErrorSummary:     fmt.Sprint(resp.Body),
			ErrorLink:        endpointUrl,
		}
	}

	defer resp.Body.Close()

	return nodeData.Computer, nil
}

// GetJobs
// Get all jobs. Only authenticated users may call this resource.
func (d *JenkinsClient) GetJobs(ctx context.Context) ([]Job, error) {
	var jobData JobsAPIData
	endpointUrl := fmt.Sprintf("%s/%s", d.baseUrl, allJobs)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		uhttp.WithHeader("Accept", "application/xml"),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&jobData))
	if err != nil {
		return nil, &JenkinsError{
			ErrorMessage:     err.Error(),
			ErrorDescription: err.Error(),
			ErrorCode:        resp.StatusCode,
			ErrorSummary:     fmt.Sprint(resp.Body),
			ErrorLink:        endpointUrl,
		}
	}

	defer resp.Body.Close()

	return jobData.Jobs, nil
}
