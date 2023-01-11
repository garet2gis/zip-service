package dto

type FileEntry struct {
	// Путь до файла на хосте
	Path string `json:"path"`
	// Желаемый путь до файла в zip архиве
	ZipPath string `json:"zip_path"`
}

type ZipDescriptor struct {
	Files []FileEntry `json:"files"`
} // @name ZipDescriptor
