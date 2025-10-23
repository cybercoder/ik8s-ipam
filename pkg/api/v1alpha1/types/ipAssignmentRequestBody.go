package types

type IpAssignmentRequestBody struct {
	Namespace          string `json:"namespace"`
	Name               string `json:"name"`
	ContainerInterface string `json:"containerInterface"`
	IpFamily           string `json:"ipFamily"`
}
