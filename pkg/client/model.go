package client

type NodesAPIData struct {
	Class    string     `json:"_class,omitempty"`
	Computer []Computer `json:"computer,omitempty"`
}

type Computer struct {
	Class               string           `json:"_class,omitempty"`
	AssignedLabels      []AssignedLabels `json:"assignedLabels,omitempty"`
	Description         string           `json:"description,omitempty"`
	DisplayName         string           `json:"displayName,omitempty"`
	Idle                bool             `json:"idle,omitempty"`
	ManualLaunchAllowed bool             `json:"manualLaunchAllowed,omitempty"`
}

type AssignedLabels struct {
	Name string `json:"name,omitempty"`
}

type JobsAPIData struct {
	Class string `json:"_class,omitempty"`
	Jobs  []Job  `json:"jobs,omitempty"`
}

type Job struct {
	Class     string `json:"_class,omitempty"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
	Buildable bool   `json:"buildable,omitempty"`
	Color     string `json:"color,omitempty"`
}

type ViewsAPIData struct {
	Class string `json:"_class,omitempty"`
	Views []View `json:"views,omitempty"`
}

type View struct {
	Class string `json:"_class,omitempty"`
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
}

type UsersAPIData struct {
	Class string  `json:"_class,omitempty"`
	Users []Users `json:"users,omitempty"`
}

type Users struct {
	LastChange interface{} `json:"lastChange,omitempty"`
	Project    interface{} `json:"project,omitempty"`
	User       User        `json:"user,omitempty"`
}

type User struct {
	AbsoluteURL string      `json:"absoluteUrl,omitempty"`
	Description interface{} `json:"description,omitempty"`
	FullName    string      `json:"fullName,omitempty"`
	ID          string      `json:"id,omitempty"`
}

type RolesAPIData struct {
	RoleName   string `json:"RoleName,omitempty"`
	RoleDetail []Role `json:"roles,omitempty"`
}

type Role struct {
	Sid  string `json:"sid,omitempty"`
	Type string `json:"type,omitempty"`
}

type Group struct {
	AbsoluteURL string      `json:"absoluteUrl,omitempty"`
	Description interface{} `json:"description,omitempty"`
	FullName    string      `json:"fullName,omitempty"`
	ID          string      `json:"id,omitempty"`
}

type GroupsAPIData struct {
	Class  string  `json:"_class,omitempty"`
	Groups []Group `json:"groups,omitempty"`
}
