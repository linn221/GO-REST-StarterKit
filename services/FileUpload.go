package services

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

var _UPLOAD_DIR string

// detectFileType reads the first 512 bytes and detects the MIME type
func detectFileType(FileHeader multipart.FileHeader) (string, *ServiceError) {
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
		return "", clientErr(fmt.Sprintf("unsupported file type: %s", contentType))
	}
	return ext, nil
}

func UploadImageFile(_, fileHeader *multipart.FileHeader, id int) (string, string, *ServiceError) {

	// filename := filepath.Base(file.Filename)
	// buffer := make([]byte, 512)
	// _, err := file.Read(buffer)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	if fileHeader == nil {
		panic("fileHeader is nil")
	}

	ext, errs := detectFileType(*fileHeader)
	if errs != nil {
		return "", "", errs
	}
	fileName := fmt.Sprintf("%d.%s", id, ext)
	// filePath := filepath.Join(_UPLOAD_DIR, fileName)
	//upload part
	// if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
	// 	return "", "", myerror.NewWithMessage("error saving uploaded file", err)
	// }
	return fileName, ext, nil
}

// func DeleteImageFile(imgId int, ext string) *myerror.ServiceError {
// 	// Construct the absolute path to the image file
// 	imagePath := filepath.Join(_UPLOAD_DIR, fmt.Sprintf("%d.%s", imgId, ext))

// 	// Ensure the file exists before attempting to delete it
// 	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
// 		return &myerror.ServiceError{
// 			Err:        err,
// 			Message:    "file not found",
// 			StatusCode: http.StatusNotFound,
// 		}
// 	}

// 	// Delete the file securely
// 	if err := os.Remove(imagePath); err != nil {
// 		return systemErrString("failed to delete file", err)
// 	}

// 	return nil
// }
