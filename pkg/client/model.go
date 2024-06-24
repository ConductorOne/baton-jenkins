package client

type NodesAPIData struct {
	Class    string     `json:"_class"`
	Computer []Computer `json:"computer"`
}

type Computer struct {
	Class               string           `json:"_class"`
	AssignedLabels      []AssignedLabels `json:"assignedLabels"`
	Description         string           `json:"description"`
	DisplayName         string           `json:"displayName"`
	Idle                bool             `json:"idle"`
	ManualLaunchAllowed bool             `json:"manualLaunchAllowed"`
}

type AssignedLabels struct {
	Name string `json:"name"`
}

type JobsAPIData struct {
	Class string `json:"_class"`
	Jobs  []Job  `json:"jobs"`
}

type Job struct {
	Class     string `json:"_class"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Buildable bool   `json:"buildable"`
	Color     string `json:"color"`
}

type ViewsAPIData struct {
	Class string `json:"_class"`
	Views []View `json:"views"`
}

type View struct {
	Class string `json:"_class"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type UsersAPIData struct {
	Class string  `json:"_class"`
	Users []Users `json:"users"`
}

type Users struct {
	LastChange interface{} `json:"lastChange"`
	Project    interface{} `json:"project"`
	User       User        `json:"user"`
}

type User struct {
	AbsoluteURL string      `json:"absoluteUrl"`
	Description interface{} `json:"description"`
	FullName    string      `json:"fullName"`
	ID          string      `json:"id"`
}

type RolesAPIData struct {
	RoleName   string `json:"RoleName"`
	RoleDetail []Role `json:"roles"`
}

type Role struct {
	Sid  string `json:"sid"`
	Type string `json:"type"`
}

type Group struct {
	AbsoluteURL string      `json:"absoluteUrl"`
	Description interface{} `json:"description"`
	FullName    string      `json:"fullName"`
	ID          string      `json:"id"`
}

type GroupsAPIData struct {
	Class  string  `json:"_class"`
	Groups []Group `json:"groups"`
}
