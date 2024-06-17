package client

type NodesAPIData struct {
	Class    string     `json:"_class"`
	Computer []Computer `json:"computer"`
}

type Computer struct {
	Class               string `json:"_class"`
	Description         string `json:"description"`
	DisplayName         string `json:"displayName"`
	Idle                bool   `json:"idle"`
	ManualLaunchAllowed bool   `json:"manualLaunchAllowed"`
}
