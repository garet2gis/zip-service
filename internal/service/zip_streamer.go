package service

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"zip_service/internal/apperror"
	"zip_service/internal/dto"
)

const (
	ROOT_DIRECTORY    = "root/"
	UPLOADS_DIRECTORY = "root/uploads"
)

type ZipStreamer struct {
	bufferSize int
}

func NewZipStreamer() *ZipStreamer {
	return &ZipStreamer{
		// 1 KB
		bufferSize: 8 * 1024,
	}
}

func (s *ZipStreamer) GetFiles(files []dto.FileEntry, destination io.Writer) error {
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
			Method:   zip.Deflate,
			Modified: time.Now(),
		}

		entryWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		var partsWritten = int64(0)
		var bytesWritten = int64(0)
		var lastBytesRead = 0
		buffer := make([]byte, s.bufferSize)

		for lastBytesRead >= 0 {
			bytesRead, berr := file.Read(buffer)
			lastBytesRead = bytesRead
			if berr == io.EOF {
				break
			}
			if berr != nil {
				log.Printf("error reading data! %v", berr)
				return berr
			}
			if lastBytesRead > 0 {
				bytesWritten += int64(lastBytesRead)
				partsWritten++

				_, err = entryWriter.Write(buffer[:bytesRead])
				if err != nil {
					return err
				}
			}
		}

		success++
	}

	if success == 0 {
		return apperror.NewAppError(nil, "all files failed to archive", "")
	}

	return zipWriter.Close()
}

func (s *ZipStreamer) UploadFile(fileName string, fileHeader *multipart.FileHeader) (err error) {
	log.Printf("read file %s", fileName)

	if !strings.HasSuffix(fileName, ".zip") {
		return apperror.NewAppError(nil, "file must be in zip format", "")
	}

	aux, _ := fileHeader.Open()
	if err != nil {
		return err
	}

	fileName = path.Join(ROOT_DIRECTORY, fileName)

	drainTo, _ := os.Create(fileName)
	_, err = io.Copy(drainTo, aux)
	if err != nil {
		return err
	}

	// check zip file is correct
	f, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	defer func() { err = f.Close() }()

	// unzip file in goroutine
	go func() {
		_ = s.UnzipFile(fileName)
	}()

	return nil
}

func (s *ZipStreamer) UnzipFile(pathName string) (err error) {
	dst := UPLOADS_DIRECTORY
	archive, err := zip.OpenReader(pathName)
	if err != nil {
		return err
	}
	defer func() {
		err = archive.Close()
	}()
	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(filepath.Separator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err = io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		err = fileInArchive.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
