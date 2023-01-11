package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"zip_service/internal/apperror"
	"zip_service/internal/dto"
)

const (
	Download = "/download/"
)

type ZipService interface {
	StreamAllFiles(files []dto.FileEntry, destination io.Writer) error
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
	err := h.service.StreamAllFiles(zipDescriptor.Files, w)
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
