package handlers

import (
	"errors"
	"fmt"
	"io"
	"linn221/shop/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func respondError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, models.ErrNotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, models.ErrBadRequest) {
		status = http.StatusBadRequest
	}
	http.Error(w, err.Error(), status)
}

func HandleImageUploadSingle(db *gorm.DB, dir string) http.HandlerFunc {
	var maxMemoryMB int64 = 10
	formInputKey := "image"

	return func(w http.ResponseWriter, r *http.Request) {

		// Limit upload size to 10MB
		r.ParseMultipartForm(maxMemoryMB << 20) // 10MB

		file, header, err := r.FormFile(formInputKey)
		if err != nil {
			http.Error(w, "Failed to read form file", http.StatusBadRequest)
			return
		}
		if header == nil {
			http.Error(w, "File header is nil", http.StatusBadRequest)
		}

		defer file.Close()

		ext, err := detectFileType(header)
		if err != nil {
			respondError(w, err)
			return
		}

		// Create destination file
		filename := uuid.NewString() + "." + ext
		uri := filepath.Join(dir, filename)
		dst, err := os.Create(uri)
		if err != nil {
			http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		if err := db.Create(&models.Image{
			Url:  filename,
			Size: header.Size,
		}).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy uploaded file to destination
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "failed to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// detectFileType reads the first 512 bytes and detects the MIME type
func detectFileType(FileHeader *multipart.FileHeader) (string, error) {
	file, err := FileHeader.Open()

	if err != nil {
		return "", fmt.Errorf("error opening file header: %v", err)
	}
	defer file.Close()
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("error reading uploaded file for detecting file type: %v", err)
	}

	contentType := http.DetectContentType(buf)     // Detect MIME type
	if !strings.HasPrefix(contentType, "image/") { // Allow only images
		return "", badRequest("unsupported file type: " + contentType)
	}
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	case "image/bmp":
		ext = "bmp"
	case "image/webp":
		ext = "webp"
	case "image/svg+xml":
		ext = "svg"
	case "image/x-icon":
		ext = "ico"
	default:
		return "", badRequest("unsupported file type: " + contentType)
	}
	return ext, nil
}

// func (service *ImageUploadService) UploadMultiple(r *http.Request, formInputKey string) ([]string, error) {

// 	// Limit request size
// 	r.ParseMultipartForm(service.maxMemoryMB << 20) // 10MB max memory

// 	files := r.MultipartForm.File[formInputKey]
// 	if len(files) == 0 {
// 		return nil, services.ClientErr("No files uploaded")
// 	}

// 	uris := make([]string, 0, len(files))
// 	for _, fileHeader := range files {
// 		uri, errs := service.uploadFileFromHeader(fileHeader)
// 		if errs != nil {
// 			return nil, errs
// 		}
// 		uris = append(uris, uri)
// 	}
// 	return uris, nil

// }

// func (service *ImageUploadService) uploadFileFromHeader(fileheader *multipart.FileHeader) (string, error) {

// 	file, err := fileheader.Open()
// 	if err != nil {
// 		return "", services.SystemErrString("Failed to open uploaded file", err)
// 	}
// 	defer file.Close()

// 	ext, errs := detectFileType(fileheader)
// 	if errs != nil {
// 		return "", errs
// 	}

// 	filename := uuid.NewString() + "." + ext
// 	dstPath := filepath.Join(service.dir, filename)
// 	dst, err := os.Create(dstPath)
// 	if err != nil {
// 		return "", services.SystemErrString("Failed to create file on server", err)
// 	}
// 	defer dst.Close()

// 	if _, err := io.Copy(dst, file); err != nil {
// 		return "", services.SystemErrString("Failed to save uploaded file", err)
// 	}
// 	return dstPath, nil
// }
