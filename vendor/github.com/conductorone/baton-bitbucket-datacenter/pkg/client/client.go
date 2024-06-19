package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type DataCenterClient struct {
	auth       *auth
	httpClient *uhttp.BaseHttpClient
	baseUrl    string
}

type BitbucketError struct {
	ErrorMessage     string                   `json:"error"`
	ErrorDescription string                   `json:"error_description"`
	ErrorCode        int                      `json:"errorCode,omitempty"`
	ErrorSummary     string                   `json:"errorSummary,omitempty" toml:"error_description"`
	ErrorLink        string                   `json:"errorLink,omitempty"`
	ErrorId          string                   `json:"errorId,omitempty"`
	ErrorCauses      []map[string]interface{} `json:"errorCauses,omitempty"`
}

func (b *BitbucketError) Error() string {
	return b.ErrorMessage
}

// GET - http://{baseurl}/rest/api/latest/users
// GET - http://{baseurl}/rest/api/latest/projects
// GET - http://{baseurl}/rest/api/latest/admin/groups
// GET - http://{baseurl}/rest/api/latest/admin/groups/more-members?context={context}
// GET - http://{baseurl}/rest/api/latest/admin/permissions/users
// GET - http://{baseurl}/rest/api/latest/projects/{projectKey}/repos/{repositorySlug}/permissions/users
// GET - http://{baseurl}/rest/api/latest/projects/{projectKey}/repos/{repositorySlug}/permissions/users?name={name}&permission={permission}
// POST - http://{baseurl}/rest/api/latest/admin/users/add-groups
// POST - http://{baseurl}/rest/api/latest/admin/users/remove-group
// PUT - http://{baseurl}/rest/api/latest/projects/{projectKey}/repos/{repositorySlug}/permissions/users?name={name}&permission={permission}
// PUT - http://{baseurl}/rest/api/latest/projects/{projectKey}/permissions/groups
// PUT - http://{baseurl}/rest/api/latest/projects/{projectKey}/permissions/users'
// DEL - http://{baseurl}/rest/api/latest/projects/{projectKey}/repos/{repositorySlug}/permissions/users?name={name}
// DEL - http://{baseurl}/rest/api/latest/projects/{projectKey}/permissions/users
// DEL - http://{baseurl}/rest/api/latest/projects/{projectKey}/permissions/groups

const (
	allUsersEndpoint                      = "rest/api/latest/users"
	allProjectsEndpoint                   = "rest/api/latest/projects"
	allRepositoriesEndpoint               = "rest/api/latest/repos"
	allGroupsEndpoint                     = "rest/api/latest/groups"
	allGroupMembersEndpoint               = "rest/api/latest/admin/groups/more-members"
	allUsersWithGlobalPermissionEndpoint  = "rest/api/latest/admin/permissions/users"
	allGroupsWithGlobalPermissionEndpoint = "rest/api/latest/admin/permissions/groups"
	usersWithPermission                   = "permissions/users"
	groupsWithPermission                  = "permissions/groups"
	addUserToGroupsEndpoint               = "rest/api/latest/admin/users/add-groups"
	removeUserFromGroupEndpoint           = "rest/api/latest/admin/users/remove-group"
)

type auth struct {
	user, password string
	bearerToken    string
}

func NewClient() *DataCenterClient {
	return &DataCenterClient{
		httpClient: &uhttp.BaseHttpClient{},
		baseUrl:    "http://localhost:7990",
		auth: &auth{
			user:        "",
			password:    "",
			bearerToken: "",
		},
	}
}

func (d *DataCenterClient) WithUser(bitbucketUsername string) *DataCenterClient {
	d.auth.user = bitbucketUsername
	return d
}

func (d *DataCenterClient) WithPassword(bitbucketPassword string) *DataCenterClient {
	d.auth.password = bitbucketPassword
	return d
}

func (d *DataCenterClient) WithBearerToken(bitbucketToken string) *DataCenterClient {
	d.auth.bearerToken = bitbucketToken
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

func (d *DataCenterClient) getToken() string {
	return d.auth.bearerToken
}

func (d *DataCenterClient) getUser() string {
	return d.auth.user
}

func (d *DataCenterClient) getPWD() string {
	return d.auth.password
}

func (d *DataCenterClient) CheckCredentials() bool {
	if d.IsBasicAuthentication() || d.getToken() != "" {
		return true
	}

	return false
}

func (d *DataCenterClient) IsBasicAuthentication() bool {
	if d.getUser() != "" && d.getPWD() != "" {
		return true
	}

	return false
}

func (d *DataCenterClient) IsTokenAuthentication() bool {
	return d.getPWD() != ""
}

func isValidUrl(baseUrl string) bool {
	u, err := url.Parse(baseUrl)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func New(ctx context.Context, baseUrl string, bitbucketClient *DataCenterClient) (*DataCenterClient, error) {
	var (
		clientId     = bitbucketClient.getUser()
		clientSecret = bitbucketClient.getPWD()
		clientToken  = bitbucketClient.getToken()
	)
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	cli := uhttp.NewBaseHttpClient(httpClient)
	if !isValidUrl(baseUrl) {
		return nil, fmt.Errorf("the url : %s is not valid", baseUrl)
	}

	// basic authentication or bearerToken
	dc := DataCenterClient{
		httpClient: cli,
		baseUrl:    baseUrl,
		auth: &auth{
			user:        clientId,
			password:    clientSecret,
			bearerToken: clientToken,
		},
	}

	return &dc, nil
}

// GetUsers
// Get all users. Only authenticated users may call this resource.
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-system-maintenance/#api-api-latest-users-get
func (d *DataCenterClient) GetUsers(ctx context.Context, startPage, limit string) ([]Users, Page, error) {
	var (
		userData     UsersAPIData
		page         Page
		sPage, nPage = "0", "0"
	)
	endpointUrl, err := url.JoinPath(d.baseUrl, allUsersEndpoint)
	if err != nil {
		return nil, Page{}, err
	}

	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&userData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(userData.Start)
	nPage = strconv.Itoa(userData.NextPageStart)
	if !userData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(userData.Size),
		}
	}

	return userData.Users, page, nil
}

func (d *DataCenterClient) ListUsers(ctx context.Context, opts PageOptions) ([]Users, string, error) {
	var nextPageToken string = ""
	users, page, err := d.GetUsers(ctx, strconv.Itoa(opts.Page), strconv.Itoa(opts.PerPage))
	if err != nil {
		return users, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return users, nextPageToken, err
}

// GetProjects
// Get projects. Only projects for which the authenticated user has the PROJECT_VIEW permission will be returned.
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-project/#api-api-latest-projects-get
func (d *DataCenterClient) GetProjects(ctx context.Context, startPage, limit string) ([]Projects, Page, error) {
	var (
		projectData  ProjectsAPIData
		page         Page
		sPage, nPage = "0", "0"
	)
	endpointUrl, err := url.JoinPath(d.baseUrl, allProjectsEndpoint)
	if err != nil {
		return nil, Page{}, err
	}

	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&projectData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(projectData.Start)
	nPage = strconv.Itoa(projectData.NextPageStart)
	if !projectData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(projectData.Size),
		}
	}

	return projectData.Projects, page, nil
}

func (d *DataCenterClient) ListProjects(ctx context.Context, opts PageOptions) ([]Projects, string, error) {
	var nextPageToken string = ""
	projects, page, err := d.GetProjects(ctx, strconv.Itoa(opts.Page), strconv.Itoa(opts.PerPage))
	if err != nil {
		return projects, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return projects, nextPageToken, err
}

func (d *DataCenterClient) GetRepos(ctx context.Context, startPage, limit string) ([]Repos, Page, error) {
	var (
		repoData     ReposAPIData
		page         Page
		sPage, nPage = "0", "0"
	)
	endpointUrl, err := url.JoinPath(d.baseUrl, allRepositoriesEndpoint)
	if err != nil {
		return nil, Page{}, err
	}

	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&repoData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(repoData.Start)
	nPage = strconv.Itoa(repoData.NextPageStart)
	if !repoData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(repoData.Size),
		}
	}

	return repoData.Repos, page, nil
}

func (d *DataCenterClient) ListRepos(ctx context.Context, opts PageOptions) ([]Repos, string, error) {
	var nextPageToken string = ""
	repos, page, err := d.GetRepos(ctx, strconv.Itoa(opts.Page), strconv.Itoa(opts.PerPage))
	if err != nil {
		return repos, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return repos, nextPageToken, err
}

// GetGroups
// Get groups. The authenticated user must have LICENSED_USER permission or higher to call this resource.
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-admin-groups-get
func (d *DataCenterClient) GetGroups(ctx context.Context, startPage, limit string) ([]string, Page, error) {
	var (
		groupData    GroupsAPIData
		page         Page
		sPage, nPage = "0", "0"
	)
	endpointUrl, err := url.JoinPath(d.baseUrl, allGroupsEndpoint)
	if err != nil {
		return nil, Page{}, err
	}

	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&groupData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(groupData.Start)
	nPage = strconv.Itoa(groupData.NextPageStart)
	if !groupData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(groupData.Size),
		}
	}

	return groupData.Groups, page, nil
}

func (d *DataCenterClient) ListGroups(ctx context.Context, opts PageOptions) ([]string, string, error) {
	var nextPageToken string = ""
	groups, page, err := d.GetGroups(ctx, strconv.Itoa(opts.Page), strconv.Itoa(opts.PerPage))
	if err != nil {
		return groups, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return groups, nextPageToken, err
}

// GetGroupMembers
// Get group members. The authenticated user must have the LICENSED_USER permission to call this resource.
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-admin-groups-more-members-get
func (d *DataCenterClient) GetGroupMembers(ctx context.Context, startPage, limit, groupName string) ([]Members, Page, error) {
	var (
		memberData   MembersAPIData
		page         Page
		sPage, nPage = "0", "0"
	)
	endpointUrl := fmt.Sprintf("%s/%s?context=%s", d.baseUrl, allGroupMembersEndpoint, groupName)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&memberData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(memberData.Start)
	nPage = strconv.Itoa(memberData.NextPageStart)
	if !memberData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(memberData.Size),
		}
	}

	return memberData.Members, page, nil
}

// setRawQuery
// Query parameters.
// Start : number for the page (inclusive). If not passed, first page is assumed.
// limit : Number of items to return. If not passed, a page size of 25 is used.
func setRawQuery(uri *url.URL, sPage string, limit string) {
	q := uri.Query()
	q.Set("start", sPage)
	q.Set("limit", limit)
	uri.RawQuery = q.Encode()
}

func (d *DataCenterClient) ListGroupMembers(ctx context.Context, opts PageOptions, groupName string) ([]Members, string, error) {
	var nextPageToken string = ""
	members, page, err := d.GetGroupMembers(ctx,
		strconv.Itoa(opts.Page),
		strconv.Itoa(opts.PerPage),
		groupName,
	)
	if err != nil {
		return members, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return members, nextPageToken, err
}

// GetGlobalUserPermissions
// Get users with a global permission
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-admin-permissions-users-get
func (d *DataCenterClient) GetGlobalUserPermissions(ctx context.Context, startPage, limit string) ([]UsersPermissions, Page, error) {
	var (
		permissionsData GlobalPermissionsAPIData
		page            Page
		sPage, nPage    = "0", "0"
	)
	endpointUrl, err := url.JoinPath(d.baseUrl, allUsersWithGlobalPermissionEndpoint)
	if err != nil {
		return nil, Page{}, err
	}

	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&permissionsData))
	if err != nil {
		return nil, Page{}, &BitbucketError{
			ErrorMessage:     err.Error(),
			ErrorDescription: err.Error(),
			ErrorCode:        resp.StatusCode,
			ErrorSummary:     fmt.Sprint(resp.Body),
			ErrorLink:        resp.Status,
		}
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(permissionsData.Start)
	nPage = strconv.Itoa(permissionsData.NextPageStart)
	if !permissionsData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(permissionsData.Size),
		}
	}

	return permissionsData.UsersPermissions, page, nil
}

func (d *DataCenterClient) GetGlobalGroupPermissions(ctx context.Context, startPage, limit string) ([]GroupsPermissions, Page, error) {
	var (
		permissionsData GlobalGroupPermissionsAPIData
		page            Page
		sPage, nPage    = "0", "0"
	)
	endpointUrl, err := url.JoinPath(d.baseUrl, allGroupsWithGlobalPermissionEndpoint)
	if err != nil {
		return nil, Page{}, err
	}

	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&permissionsData))
	if err != nil {
		return nil, Page{}, &BitbucketError{
			ErrorMessage:     err.Error(),
			ErrorDescription: err.Error(),
			ErrorCode:        resp.StatusCode,
			ErrorSummary:     fmt.Sprint(resp.Body),
		}
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(permissionsData.Start)
	nPage = strconv.Itoa(permissionsData.NextPageStart)
	if !permissionsData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(permissionsData.Size),
		}
	}

	return permissionsData.GroupsPermissions, page, nil
}

func (d *DataCenterClient) ListGlobalUserPermissions(ctx context.Context, opts PageOptions) ([]UsersPermissions, string, error) {
	var nextPageToken string = ""
	usersPermissions, page, err := d.GetGlobalUserPermissions(ctx,
		strconv.Itoa(opts.Page),
		strconv.Itoa(opts.PerPage),
	)
	if err != nil {
		return []UsersPermissions{}, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return usersPermissions, nextPageToken, err
}

func (d *DataCenterClient) ListGlobalGroupPermissions(ctx context.Context, opts PageOptions) ([]GroupsPermissions, string, error) {
	var nextPageToken string = ""
	groupsPermissions, page, err := d.GetGlobalGroupPermissions(ctx,
		strconv.Itoa(opts.Page),
		strconv.Itoa(opts.PerPage),
	)
	if err != nil {
		return groupsPermissions, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return groupsPermissions, nextPageToken, err
}

// GetUserRepositoryPermissions
// Get users with permission to repository
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-projects-projectkey-repos-repositoryslug-permissions-users-get
func (d *DataCenterClient) GetUserRepositoryPermissions(ctx context.Context, startPage, limit, projectKey, repositorySlug string) ([]UsersPermissions, Page, error) {
	var (
		permissionData GlobalPermissionsAPIData
		page           Page
		sPage, nPage   = "0", "0"
	)
	endpointUrl := fmt.Sprintf("%s/%s/%s/repos/%s/%s", d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		repositorySlug,
		usersWithPermission,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&permissionData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(permissionData.Start)
	nPage = strconv.Itoa(permissionData.NextPageStart)
	if !permissionData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(permissionData.Size),
		}
	}

	return permissionData.UsersPermissions, page, nil
}

func (d *DataCenterClient) ListUserRepositoryPermissions(ctx context.Context, opts PageOptions, projectKey, repositorySlug string) ([]UsersPermissions, string, error) {
	var nextPageToken string = ""
	permissions, page, err := d.GetUserRepositoryPermissions(ctx,
		strconv.Itoa(opts.Page),
		strconv.Itoa(opts.PerPage),
		projectKey,
		repositorySlug,
	)
	if err != nil {
		return permissions, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return permissions, nextPageToken, err
}

func (d *DataCenterClient) GetUserProjectsPermissions(ctx context.Context, startPage, limit, projectKey string) ([]UsersPermissions, Page, error) {
	var (
		permissionData GlobalPermissionsAPIData
		page           Page
		sPage, nPage   = "0", "0"
	)
	endpointUrl := fmt.Sprintf("%s/%s/%s/%s", d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		usersWithPermission,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&permissionData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(permissionData.Start)
	nPage = strconv.Itoa(permissionData.NextPageStart)
	if !permissionData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(permissionData.Size),
		}
	}

	return permissionData.UsersPermissions, page, nil
}

func (d *DataCenterClient) GetGroupProjectsPermissions(ctx context.Context, startPage, limit, projectKey string) ([]GroupsPermissions, Page, error) {
	var (
		permissionData GlobalGroupPermissionsAPIData
		page           Page
		sPage, nPage   = "0", "0"
	)
	strUrl := fmt.Sprintf("%s/%s/%s/%s", d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		groupsWithPermission,
	)
	uri, err := url.Parse(strUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&permissionData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(permissionData.Start)
	nPage = strconv.Itoa(permissionData.NextPageStart)
	if !permissionData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(permissionData.Size),
		}
	}

	return permissionData.GroupsPermissions, page, nil
}

func (d *DataCenterClient) ListUserProjectsPermissions(ctx context.Context, opts PageOptions, projectKey string) ([]UsersPermissions, string, error) {
	var nextPageToken string = ""
	permissions, page, err := d.GetUserProjectsPermissions(ctx,
		strconv.Itoa(opts.Page),
		strconv.Itoa(opts.PerPage),
		projectKey,
	)
	if err != nil {
		return permissions, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return permissions, nextPageToken, err
}

func (d *DataCenterClient) ListGroupProjectsPermissions(ctx context.Context, opts PageOptions, projectKey string) ([]GroupsPermissions, string, error) {
	var nextPageToken string = ""
	permissions, page, err := d.GetGroupProjectsPermissions(ctx,
		strconv.Itoa(opts.Page),
		strconv.Itoa(opts.PerPage),
		projectKey,
	)
	if err != nil {
		return permissions, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return permissions, nextPageToken, err
}

func (d *DataCenterClient) GetGroupRepositoryPermissions(ctx context.Context, startPage, limit, projectKey, repositorySlug string) ([]GroupsPermissions, Page, error) {
	var (
		permissionData GroupPermissionsAPIData
		page           Page
		sPage, nPage   = "0", "0"
	)
	strUrl := fmt.Sprintf("%s/%s/%s/repos/%s/%s", d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		repositorySlug,
		groupsWithPermission,
	)
	uri, err := url.Parse(strUrl)
	if err != nil {
		return nil, Page{}, err
	}

	if startPage != "" {
		sPage = startPage
	}

	setRawQuery(uri, sPage, limit)
	req, err := d.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return nil, Page{}, err
	}

	resp, err := d.httpClient.Do(req, uhttp.WithJSONResponse(&permissionData))
	if err != nil {
		return nil, Page{}, err
	}

	defer resp.Body.Close()
	sPage = strconv.Itoa(permissionData.Start)
	nPage = strconv.Itoa(permissionData.NextPageStart)
	if !permissionData.IsLastPage {
		page = Page{
			PreviousPage: &sPage,
			NextPage:     &nPage,
			Count:        int64(permissionData.Size),
		}
	}

	return permissionData.GroupsPermissions, page, nil
}

func (d *DataCenterClient) ListGroupRepositoryPermissions(ctx context.Context, opts PageOptions, projectKey, repositorySlug string) ([]GroupsPermissions, string, error) {
	var nextPageToken string = ""
	permissions, page, err := d.GetGroupRepositoryPermissions(ctx,
		strconv.Itoa(opts.Page),
		strconv.Itoa(opts.PerPage),
		projectKey,
		repositorySlug,
	)
	if err != nil {
		return permissions, "", err
	}

	if page.HasNext() {
		nextPageToken = *page.NextPage
	}

	return permissions, nextPageToken, err
}

// AddUserToGroups
// Add a user to one or more groups.
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-admin-users-add-groups-post
func (d *DataCenterClient) AddUserToGroups(ctx context.Context, groupName, userName string) error {
	var (
		body struct {
			Groups []string `json:"groups"`
			User   string   `json:"user"`
		}
		payload = []byte(fmt.Sprintf(`{"groups": ["%s"], "user": "%s"}`, groupName, userName))
	)
	endpointUrl, err := url.JoinPath(d.baseUrl, addUserToGroupsEndpoint)
	if err != nil {
		return err
	}

	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	err = json.Unmarshal(payload, &body)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodPost,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
		uhttp.WithJSONBody(body),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("user not added")
	}

	return nil
}

// RemoveUserFromGroup
// Remove user from group
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-admin-users-remove-group-post
func (d *DataCenterClient) RemoveUserFromGroup(ctx context.Context, userName, groupName string) error {
	var (
		body struct {
			Context  string `json:"context"`
			ItemName string `json:"itemName"`
		}
		payload = []byte(fmt.Sprintf(`{"context": "%s", "itemName": "%s"}`, userName, groupName))
	)
	endpointUrl, err := url.JoinPath(d.baseUrl, removeUserFromGroupEndpoint)
	if err != nil {
		return err
	}

	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	err = json.Unmarshal(payload, &body)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodPost,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
		uhttp.WithJSONBody(body),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("user not removed")
	}

	return nil
}

// UpdateUserRepositoryPermission
// Update user repository permission
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-projects-projectkey-repos-repositoryslug-permissions-users-put
func (d *DataCenterClient) UpdateUserRepositoryPermission(ctx context.Context, projectKey, repositorySlug, userName, permission string) error {
	endpointUrl := fmt.Sprintf("%s/%s/%s/repos/%s/%s?name=%s&permission=%s",
		d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		repositorySlug,
		usersWithPermission,
		userName,
		permission,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodPut,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("user not added")
	}

	return nil
}

// UpdateGroupRepositoryPermission
// Update group repository permission
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-projects-projectkey-repos-repositoryslug-permissions-groups-put
func (d *DataCenterClient) UpdateGroupRepositoryPermission(ctx context.Context, projectKey, repositorySlug, groupName, permission string) error {
	endpointUrl := fmt.Sprintf("%s/%s/%s/repos/%s/%s?name=%s&permission=%s",
		d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		repositorySlug,
		groupsWithPermission,
		groupName,
		permission,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodPut,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("group not added")
	}

	return nil
}

// RevokeGroupRepositoryPermission
// Revoke group repository permission
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-projects-projectkey-repos-repositoryslug-permissions-groups-delete
func (d *DataCenterClient) RevokeGroupRepositoryPermission(ctx context.Context, projectKey, repositorySlug, groupName string) error {
	endpointUrl := fmt.Sprintf("%s/%s/%s/repos/%s/%s?name=%s",
		d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		repositorySlug,
		groupsWithPermission,
		groupName,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodDelete,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("group not added")
	}

	return nil
}

// RevokeUserRepositoryPermission
// Revoke user repository permission
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-permission-management/#api-api-latest-projects-projectkey-repos-repositoryslug-permissions-users-delete
func (d *DataCenterClient) RevokeUserRepositoryPermission(ctx context.Context, projectKey, repositorySlug, userName string) error {
	endpointUrl := fmt.Sprintf("%s/%s/%s/repos/%s/%s?name=%s",
		d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		repositorySlug,
		usersWithPermission,
		userName,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodDelete,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("user not added")
	}

	return nil
}

// RevokeUserProjectPermission
// Revoke user project permission
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-project/#api-api-latest-projects-projectkey-permissions-users-delete
func (d *DataCenterClient) RevokeUserProjectPermission(ctx context.Context, projectKey, userName string) error {
	endpointUrl := fmt.Sprintf("%s/%s/%s/%s?name=%s",
		d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		usersWithPermission,
		userName,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodDelete,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("user not added")
	}

	return nil
}

// RevokeGroupProjectPermission
// Revoke group project permission.
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-project/#api-api-latest-projects-projectkey-permissions-groups-delete
func (d *DataCenterClient) RevokeGroupProjectPermission(ctx context.Context, projectKey, groupName string) error {
	endpointUrl := fmt.Sprintf("%s/%s/%s/%s?name=%s",
		d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		groupsWithPermission,
		groupName,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodDelete,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("user not added")
	}

	return nil
}

// UpdateUserProjectPermission
// Update user project permission. Available project permissions are: PROJECT_READ, PROJECT_WRITE, PROJECT_ADMIN
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-project/#api-api-latest-projects-projectkey-permissions-users-put
func (d *DataCenterClient) UpdateUserProjectPermission(ctx context.Context, projectKey, userName, permission string) error {
	endpointUrl := fmt.Sprintf("%s/%s/%s/%s?name=%s&permission=%s",
		d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		usersWithPermission,
		userName,
		permission,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodPut,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("project not updated")
	}

	return nil
}

// UpdateGroupProjectPermission
// Update group project permission. Available project permissions are: PROJECT_READ, PROJECT_WRITE, PROJECT_ADMIN
// https://developer.atlassian.com/server/bitbucket/rest/v819/api-group-project/#api-api-latest-projects-projectkey-permissions-groups-put
func (d *DataCenterClient) UpdateGroupProjectPermission(ctx context.Context, projectKey, groupName, permission string) error {
	endpointUrl := fmt.Sprintf("%s/%s/%s/%s?name=%s&permission=%s",
		d.baseUrl,
		allProjectsEndpoint,
		projectKey,
		groupsWithPermission,
		groupName,
		permission,
	)
	uri, err := url.Parse(endpointUrl)
	if err != nil {
		return err
	}

	req, err := d.httpClient.NewRequest(ctx,
		http.MethodPut,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithAuthorization(d.getUser(), d.getPWD(), d.getToken()),
	)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("project not updated")
	}

	return nil
}
