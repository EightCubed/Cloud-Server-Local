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
