package models

type FileType string

const (
	FileTypeFile   FileType = "file"
	FileTypeFolder FileType = "folder"
)

type Node struct {
	File            File    `json:"file"`
	Children        []*Node `json:"children"`
	FilePath        string  `json:"filepath"`
	ParentDirectory string  `json:"parentdirectory"`
}

type File struct {
	FileName         string   `json:"filename"`
	FileType         FileType `json:"filetype"`
	AbsoluteFilePath string   `json:"absolutefilepath"`
}

type FilePathRequest struct {
	FilePath string `json:"filepath"`
}
