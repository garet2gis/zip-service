package service

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path"
	"time"
	"zip_service/internal/apperror"
	"zip_service/internal/dto"
)

const (
	ROOT_DIRECTORY = "root/"
)

type ZipStreamer struct {
}

func NewZipStreamer() *ZipStreamer {
	return &ZipStreamer{}
}

func (s *ZipStreamer) StreamAllFiles(files []dto.FileEntry, destination io.Writer) error {
	if len(files) == 0 {
		return apperror.NewAppError(nil, "must have at least 1 file", "")
	}

	zipWriter := zip.NewWriter(destination)
	success := 0

	for _, entry := range files {
		file, err := os.Open(path.Join(ROOT_DIRECTORY, entry.Path))
		if err != nil {
			return err
		}

		header := &zip.FileHeader{
			Name:     entry.ZipPath,
			Method:   zip.Store,
			Modified: time.Now(),
		}

		entryWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(entryWriter, file)
		if err != nil {
			return err
		}

		err = zipWriter.Flush()
		if err != nil {
			return err
		}

		flushingWriter, ok := destination.(http.Flusher)
		if ok {
			flushingWriter.Flush()
		}

		success++
	}

	if success == 0 {
		return apperror.NewAppError(nil, "all files failed to archive", "")
	}

	return zipWriter.Close()
}
