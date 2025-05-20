package services

import (
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

// detectFileType reads the first 512 bytes and detects the MIME type
func detectFileType(FileHeader *multipart.FileHeader) (string, *ServiceError) {
	file, err := FileHeader.Open()

	if err != nil {
		return "", systemErrString("error opening file header: " + err.Error())
	}
	defer file.Close()
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "", systemErrString("error reading uploaded file for detecting file type", err)
	}

	contentType := http.DetectContentType(buf)     // Detect MIME type
	if !strings.HasPrefix(contentType, "image/") { // Allow only images
		return "", clientErr("unsupported file type: " + contentType)
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
		return "", clientErr("unsupported file type: " + contentType)
	}
	return ext, nil
}

// func DeleteImageFile(uri string) error {
// 	// Construct the absolute path to the image file
// 	imagePath := filepath.Join(_UPLOAD_DIR, uri)

// 	// Ensure the file exists before attempting to delete it
// 	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
// 		return errors.New("file not found")
// 	}

// 	// Delete the file securely
// 	if err := os.Remove(imagePath); err != nil {
// 		return err
// 	}

// 	return nil
// }

type ImageUploader interface {
	UploadSingle(r *http.Request, key string) (string, *ServiceError)
	UploadMultiple(r *http.Request, key string) ([]string, *ServiceError)
}

type ImageUploadService struct {
	dir         string
	db          *gorm.DB
	maxMemoryMB int64 // default should be 10
}

func (service *ImageUploadService) UploadSingle(r *http.Request, formInputKey string) (string, *ServiceError) {
	// Limit upload size to 10MB
	r.ParseMultipartForm(service.maxMemoryMB << 20) // 10MB

	file, header, err := r.FormFile(formInputKey)
	if err != nil {
		return "", clientErr("Failed to read form file")
	}
	if header == nil {
		return "", clientErr("Fileheader is nil")
	}

	defer file.Close()

	ext, errs := detectFileType(header)
	if errs != nil {
		return "", errs
	}

	// Create destination file
	filename := uuid.NewString() + "." + ext
	uri := filepath.Join(service.dir, filename)
	dst, err := os.Create(uri)
	if err != nil {
		return "", systemErrString("Failed to create file", err)
	}

	defer dst.Close()

	if err := service.db.Create(&models.Image{
		Url:  filename,
		Size: header.Size,
	}).Error; err != nil {
		return "", systemErr(err)
	}

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", systemErrString("Failed to save file", err)
	}

	return filename, nil
}

func (service *ImageUploadService) UploadMultiple(r *http.Request, formInputKey string) ([]string, *ServiceError) {

	// Limit request size
	r.ParseMultipartForm(service.maxMemoryMB << 20) // 10MB max memory

	files := r.MultipartForm.File[formInputKey]
	if len(files) == 0 {
		return nil, clientErr("No files uploaded")
	}

	uris := make([]string, 0, len(files))
	for _, fileHeader := range files {
		uri, errs := service.uploadFileFromHeader(fileHeader)
		if errs != nil {
			return nil, errs
		}
		uris = append(uris, uri)
	}
	return uris, nil

}

func (service *ImageUploadService) uploadFileFromHeader(fileheader *multipart.FileHeader) (string, *ServiceError) {

	file, err := fileheader.Open()
	if err != nil {
		return "", systemErrString("Failed to open uploaded file", err)
	}
	defer file.Close()

	ext, errs := detectFileType(fileheader)
	if errs != nil {
		return "", errs
	}

	filename := uuid.NewString() + "." + ext
	dstPath := filepath.Join(service.dir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", systemErrString("Failed to create file on server", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", systemErrString("Failed to save uploaded file", err)
	}
	return dstPath, nil
}
