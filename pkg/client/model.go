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
