package models

type DirectoryRequestBody struct {
	Directory string `json:"directory"` // Add JSON tags for better decoding
}

type BreadCrumbType struct {
	Title        string `json:"title"`
	AbsolutePath string `json:"absolutePath"`
}

type ResponseDataFileDirectory struct {
	Data    Node             `json:"data"`
	Path    []BreadCrumbType `json:"path"`
	Message string           `json:"message"` // Corrected the typo from 'messaeg' to 'Message'
}

type DeleteRequestBody struct {
	FilesToBeDeleted []Node `json:"filesToBeDeleted"`
}

type DeleteResponseBody struct {
	SuccessCount int    `json:"success_count"`
	FailureCount int    `json:"failure_count"`
	Message      string `json:"message"`
}
