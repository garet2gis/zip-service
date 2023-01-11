package dto

type FileEntry struct {
	Path    string `json:"path"`
	ZipPath string `json:"zip_path"`
}

type ZipDescriptor struct {
	Files []FileEntry `json:"files"`
}
