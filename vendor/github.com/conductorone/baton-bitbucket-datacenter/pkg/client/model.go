package client

type UsersAPIData struct {
	Size          int     `json:"size,omitempty"`
	Limit         int     `json:"limit,omitempty"`
	IsLastPage    bool    `json:"isLastPage,omitempty"`
	Users         []Users `json:"values,omitempty"`
	Start         int     `json:"start,omitempty"`
	NextPageStart int     `json:"nextPageStart,omitempty"`
}

type Users struct {
	Name         string   `json:"name,omitempty"`
	EmailAddress string   `json:"emailAddress,omitempty"`
	Active       bool     `json:"active,omitempty"`
	DisplayName  string   `json:"displayName,omitempty"`
	ID           int      `json:"id,omitempty"`
	Slug         string   `json:"slug,omitempty"`
	Type         string   `json:"type,omitempty"`
	Links        UserSelf `json:"links,omitempty"`
}

type UserSelf struct {
	Self any `json:"self,omitempty"`
}

type ProjectsAPIData struct {
	Size          int        `json:"size,omitempty"`
	Limit         int        `json:"limit,omitempty"`
	IsLastPage    bool       `json:"isLastPage,omitempty"`
	Projects      []Projects `json:"values,omitempty"`
	Start         int        `json:"start,omitempty"`
	NextPageStart int        `json:"nextPageStart,omitempty"`
}

type Projects struct {
	Key    string      `json:"key,omitempty"`
	ID     int         `json:"id,omitempty"`
	Name   string      `json:"name,omitempty"`
	Public bool        `json:"public,omitempty"`
	Type   string      `json:"type,omitempty"`
	Links  ProjectSelf `json:"links,omitempty"`
}

type ProjectSelf struct {
	Self any `json:"self,omitempty"`
}

type ReposAPIData struct {
	Size          int     `json:"size,omitempty"`
	Limit         int     `json:"limit,omitempty"`
	IsLastPage    bool    `json:"isLastPage,omitempty"`
	Repos         []Repos `json:"values,omitempty"`
	Start         int     `json:"start,omitempty"`
	NextPageStart int     `json:"nextPageStart,omitempty"`
}

type Repos struct {
	Slug          string   `json:"slug,omitempty"`
	ID            int      `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	HierarchyId   string   `json:"hierarchyId,omitempty"`
	ScmId         string   `json:"scmId,omitempty"`
	State         string   `json:"state,omitempty"`
	StatusMessage string   `json:"statusMessage,omitempty"`
	Forkable      bool     `json:"forkable,omitempty"`
	Project       Projects `json:"project,omitempty"`
}

type GroupsAPIData struct {
	Size          int      `json:"size,omitempty"`
	Limit         int      `json:"limit,omitempty"`
	IsLastPage    bool     `json:"isLastPage,omitempty"`
	Groups        []string `json:"values,omitempty"`
	Start         int      `json:"start,omitempty"`
	NextPageStart int      `json:"nextPageStart,omitempty"`
}

type MembersAPIData struct {
	Size          int       `json:"size,omitempty"`
	Limit         int       `json:"limit,omitempty"`
	IsLastPage    bool      `json:"isLastPage,omitempty"`
	Members       []Members `json:"values,omitempty"`
	Start         int       `json:"start,omitempty"`
	NextPageStart int       `json:"nextPageStart,omitempty"`
}

type Members struct {
	Name           string   `json:"name,omitempty"`
	EmailAddress   string   `json:"emailAddress,omitempty"`
	Active         bool     `json:"active,omitempty"`
	DisplayName    string   `json:"displayName,omitempty"`
	ID             int      `json:"id,omitempty"`
	Slug           string   `json:"slug,omitempty"`
	Type           string   `json:"type,omitempty"`
	DirectoryName  string   `json:"directoryName,omitempty"`
	Deletable      bool     `json:"deletable,omitempty"`
	MutableDetails bool     `json:"mutableDetails,omitempty"`
	MutableGroups  bool     `json:"mutableGroups,omitempty"`
	Links          UserSelf `json:"links,omitempty"`
}

type GlobalPermissionsAPIData struct {
	Size             int                `json:"size,omitempty"`
	Limit            int                `json:"limit,omitempty"`
	IsLastPage       bool               `json:"isLastPage,omitempty"`
	UsersPermissions []UsersPermissions `json:"values,omitempty"`
	Start            int                `json:"start,omitempty"`
	NextPageStart    int                `json:"nextPageStart,omitempty"`
}

type GlobalGroupPermissionsAPIData struct {
	Size              int                 `json:"size,omitempty"`
	Limit             int                 `json:"limit,omitempty"`
	IsLastPage        bool                `json:"isLastPage,omitempty"`
	GroupsPermissions []GroupsPermissions `json:"values,omitempty"`
	Start             int                 `json:"start,omitempty"`
	NextPageStart     int                 `json:"nextPageStart,omitempty"`
}

type UsersPermissions struct {
	User       Users  `json:"user,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type GroupPermissionsAPIData struct {
	Size              int                 `json:"size,omitempty"`
	Limit             int                 `json:"limit,omitempty"`
	IsLastPage        bool                `json:"isLastPage,omitempty"`
	GroupsPermissions []GroupsPermissions `json:"values,omitempty"`
	Start             int                 `json:"start,omitempty"`
	NextPageStart     int                 `json:"nextPageStart,omitempty"`
}

type GroupsPermissions struct {
	Group      Groups `json:"group,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type Groups struct {
	Name string `json:"name,omitempty"`
}
