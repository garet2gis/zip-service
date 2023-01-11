package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"zip_service/internal/apperror"
	"zip_service/internal/dto"
)

const (
	Download = "/download/"
	Upload   = "/upload/"
)

type ZipService interface {
	GetFiles(files []dto.FileEntry, destination io.Writer) error
	UploadFile(fileName string, part *multipart.Part) (bytesWritten int64, partsWritten int64, err error)
	UnzipFile(fileName string) (err error)
}

type zipHandler struct {
	service ZipService
}

func NewZipHandler(s ZipService) Handler {
	return &zipHandler{
		service: s,
	}
}

func (h *zipHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, Download, apperror.Middleware(h.Download))
	router.HandlerFunc(http.MethodPost, Upload, apperror.Middleware(h.UploadZIP))
}

func (h *zipHandler) UploadZIP(w http.ResponseWriter, req *http.Request) error {
	multipartReader, err := req.MultipartReader()
	if err != nil {
		log.Printf("failed to get a multipart reader %v", err)
		return err
	}

	partBytes := int64(0)
	partCount := int64(0)
	for {
		part, partErr := multipartReader.NextPart()
		if partErr != nil {
			if partErr == io.EOF {
				break
			} else {
				return partErr
			}
		} else {
			if len(part.FileName()) > 0 && strings.HasSuffix(part.FileName(), ".zip") {
				fileName := "root" + "/" + part.FileName()
				partBytesIncr, partCountIncr, err := h.service.UploadFile(fileName, part)
				if err != nil {
					return err
				}
				err = h.service.UnzipFile(fileName)
				if err != nil {
					return err
				}
				partBytes += partBytesIncr
				partCount += partCountIncr
			}
		}
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *zipHandler) Download(w http.ResponseWriter, req *http.Request) error {
	var zipDescriptor dto.ZipDescriptor

	err := json.NewDecoder(req.Body).Decode(&zipDescriptor)
	if err != nil {
		return err
	}

	h.StreamEntries(&zipDescriptor, w)
	return nil
}

func (h *zipHandler) StreamEntries(zipDescriptor *dto.ZipDescriptor, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")
	w.WriteHeader(http.StatusOK)
	err := h.service.GetFiles(zipDescriptor.Files, w)
	if err != nil {
		// Close the connection so the client gets an error instead of 200 but invalid file
		closeForError(w)
	}
}

func closeForError(w http.ResponseWriter) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		return
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		return
	}

	err = conn.Close()
	if err != nil {
		return
	}
}
