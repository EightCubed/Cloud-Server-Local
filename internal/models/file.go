package models

type FileType string

const (
	FileTypeFile   FileType = "file"
	FileTypeFolder FileType = "folder"
)

type Node struct {
	File     File    `json:"file"`     // Change this to use the File struct
	Adjacent []*Node `json:"adjacent"` // Recursive structure for the tree
}

type File struct {
	FileName string   `json:"filename"`
	FileType FileType `json:"filetype"` // This can be FileTypeFile or FileTypeFolder
}
