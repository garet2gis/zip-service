package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"zip_service/internal/apperror"
	"zip_service/internal/dto"
)

const (
	Download = "/download/"
	Upload   = "/upload/"
)

type ZipService interface {
	GetFiles(files []dto.FileEntry, destination io.Writer) error
	UploadFile(fileName string, file *multipart.FileHeader) (err error)
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
	router.HandlerFunc(http.MethodPost, Upload, apperror.Middleware(h.Upload))
}

func (h *zipHandler) Upload(w http.ResponseWriter, req *http.Request) error {
	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("failed to parse a multipart %v", err)
		return err
	}
	file := req.MultipartForm
	for _, header := range file.File {
		err = h.service.UploadFile(header[0].Filename, header[0])
		if err != nil {
			return err
		}
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Download godoc
// @Summary Скачивание желаемых файлов в форме zip-архива
// @ID      download-zip
// @Param   user_id body dto.ZipDescriptor true "Zip Descriptor"
// @Tags    ZIP
// @Success 200 {string} string "ZIP file"
// @Failure 404 {object} apperror.AppError
// @Router  /download/ [post]
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
