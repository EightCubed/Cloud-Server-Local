package models

type Node struct {
	Data     string  `json:"data"`
	Adjacent []*Node `json:"adjacent"`
}

type DirectoryRequestBody struct {
	Directory string `json:"directory"` // Add JSON tags for better decoding
}
