package models

type DirectoryRequestBody struct {
	Directory string `json:"directory"` // Add JSON tags for better decoding
}

type ResponseDataFileDirectory struct {
	Data    Node     `json:"data"`
	Path    []string `json:"path"`
	Message string   `json:"message"` // Corrected the typo from 'messaeg' to 'Message'
}
