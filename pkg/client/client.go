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
	auth       *auth
	httpClient *uhttp.BaseHttpClient
	baseUrl    string
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

// GET - http://{baseurl}/computer/api/json?pretty&tree=computer[displayName,description,idle,manualLaunchAllowed,assignedLabels[name]]
// GET - http://{baseurl}/api/json?pretty&tree=jobs[name,url,color,buildable]
// GET - http://{baseurl}/api/json?pretty&tree=views[name,url]
// GET - http://{baseurl}/asynchPeople/api/json?pretty&depth=3
// GET - http://{baseurl}/role-strategy/strategy/getAllRoles?type=globalRoles
const (
	allNodes        = "computer/api/json?pretty&tree=computer[displayName,description,idle,manualLaunchAllowed,assignedLabels[name]]"
	allJobs         = "api/json?pretty&tree=jobs[name,url,color,buildable]"
	allViews        = "api/json?pretty&tree=views[name,url]"
	allUsers        = "asynchPeople/api/json?pretty&depth=3"
	allGlobalRoles  = "role-strategy/strategy/getAllRoles?type=globalRoles"
	allProjectRoles = "role-strategy/strategy/getAllRoles?type=projectRoles"
	allSlaveRoles   = "role-strategy/strategy/getAllRoles?type=slaveRoles"
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
	}

	return &jc, nil
}

func getRequest(ctx context.Context, cli *JenkinsClient, baseUrl, apiUrl string) (*http.Request, string, error) {
	endpointUrl := fmt.Sprintf("%s/%s", baseUrl, apiUrl)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, "", err
	}

	req, err := cli.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		uhttp.WithHeader("Accept", "application/xml"),
		WithAuthorization(cli.getUser(), cli.getPWD(), cli.getToken()),
	)
	if err != nil {
		return nil, "", err
	}

	return req, endpointUrl, nil
}

func getCustomError(err error, resp *http.Response, endpointUrl string) *JenkinsError {
	return &JenkinsError{
		ErrorMessage:     err.Error(),
		ErrorDescription: err.Error(),
		ErrorCode:        resp.StatusCode,
		ErrorSummary:     fmt.Sprint(resp.Body),
		ErrorLink:        endpointUrl,
	}
}

// GetNodes
// Get all nodes. Only authenticated users may call this resource.
func (d *JenkinsClient) GetNodes(ctx context.Context) ([]Computer, error) {
	var nodeData NodesAPIData
	req, endpointUrl, err := getRequest(ctx, d, d.baseUrl, allNodes)
	if err != nil {
		return nil, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&nodeData))
	if err != nil {
		return nil, getCustomError(err, resp, endpointUrl)
	}

	defer resp.Body.Close()

	return nodeData.Computer, nil
}

// GetJobs
// Get all jobs. Only authenticated users may call this resource.
func (d *JenkinsClient) GetJobs(ctx context.Context) ([]Job, error) {
	var jobData JobsAPIData
	req, endpointUrl, err := getRequest(ctx, d, d.baseUrl, allJobs)
	if err != nil {
		return nil, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&jobData))
	if err != nil {
		return nil, getCustomError(err, resp, endpointUrl)
	}

	defer resp.Body.Close()

	return jobData.Jobs, nil
}

// GetViews
// Get all views. Only authenticated users may call this resource.
func (d *JenkinsClient) GetViews(ctx context.Context) ([]View, error) {
	var viewData ViewsAPIData
	req, endpointUrl, err := getRequest(ctx, d, d.baseUrl, allViews)
	if err != nil {
		return nil, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&viewData))
	if err != nil {
		return nil, getCustomError(err, resp, endpointUrl)
	}

	defer resp.Body.Close()

	return viewData.Views, nil
}

// GetUsers
// Get all users. Only authenticated users may call this resource.
func (d *JenkinsClient) GetUsers(ctx context.Context) ([]Users, error) {
	var userData UsersAPIData
	req, endpointUrl, err := getRequest(ctx, d, d.baseUrl, allUsers)
	if err != nil {
		return nil, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&userData))
	if err != nil {
		return nil, getCustomError(err, resp, endpointUrl)
	}

	defer resp.Body.Close()

	return userData.Users, nil
}

// GetRoles
// Get all global roles. Only authenticated users may call this resource.
func (d *JenkinsClient) GetRoles(ctx context.Context, apiUrl string) ([]RolesAPIData, error) {
	var (
		rolesAPIData []RolesAPIData
		roleData     map[string]any
		roleDetail   []any
		ok           bool
	)
	req, endpointUrl, err := getRequest(ctx, d, d.baseUrl, apiUrl)
	if err != nil {
		return nil, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&roleData))
	if err != nil {
		return nil, getCustomError(err, resp, endpointUrl)
	}

	defer resp.Body.Close()
	for roleName, roleDetails := range roleData {
		var roles []Role
		if roleDetail, ok = roleDetails.([]any); !ok {
			return nil, err
		}

		for _, itemDetails := range roleDetail {
			item := itemDetails.(map[string]any)
			roles = append(roles, Role{
				Sid:  fmt.Sprint(item["sid"]),
				Type: fmt.Sprint(item["type"]),
			})
		}
		rolesAPIData = append(rolesAPIData, RolesAPIData{
			RoleName:   roleName,
			RoleDetail: roles,
		})
	}

	return rolesAPIData, nil
}

// GetAllRoles
// Get all roles. Only authenticated users may call this resource.
func (d *JenkinsClient) GetAllRoles(ctx context.Context) ([]RolesAPIData, error) {
	var allRoles []RolesAPIData
	roles, err := d.GetRoles(ctx, allGlobalRoles)
	if err != nil {
		return nil, err
	}

	allRoles = append(allRoles, roles...)
	roles, err = d.GetRoles(ctx, allProjectRoles)
	if err != nil {
		return nil, err
	}

	allRoles = append(allRoles, roles...)
	roles, err = d.GetRoles(ctx, allSlaveRoles)
	if err != nil {
		return nil, err
	}

	allRoles = append(allRoles, roles...)
	return allRoles, nil
}

// GetGroups
// Get all groups.
func (d *JenkinsClient) GetGroups(ctx context.Context) ([]Group, error) {
	var (
		groupData GroupsAPIData
		arrIDs    []string
	)
	groups, err := d.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		for _, item := range group.RoleDetail {
			if item.Type != "GROUP" {
				continue
			}
			arrIDs = append(arrIDs, item.Sid)
		}
	}

	groupIDs := removeDuplicates(arrIDs)
	groupData = GroupsAPIData{
		Groups: groupIDs,
	}

	return groupData.Groups, nil
}

func removeDuplicates(groupIDs []string) []Group {
	var groups []Group
	keys := make(map[string]bool)
	// If the key(value of the slice) is not equal we append it. else we jump on another element.
	for _, id := range groupIDs {
		if _, value := keys[id]; !value {
			keys[id] = true
			groups = append(groups, Group{
				ID: id,
			})
		}
	}

	return groups
}
