package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

func SaveUploadedImage(r *http.Request, formFieldName, uploadDir string) (string, error) {
	file, header, err := r.FormFile(formFieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := validateImage(header, file); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	filename := uuid.New().String() + ext

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	dstPath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

// MaxImageSize is the maximum allowed image upload size (5 MB).
const MaxImageSize = 5 * 1024 * 1024

func validateImage(header *multipart.FileHeader, file multipart.File) error {
	// Enforce file size limit
	if header.Size > MaxImageSize {
		return errors.New("image exceeds the 5 MB size limit")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExtensions[ext] {
		return errors.New("unsupported image extension")
	}

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return err
	}

	contentType := http.DetectContentType(buffer[:n])
	if !allowedMimeTypes[contentType] {
		return errors.New("unsupported mime type")
	}

	_, err = file.Seek(0, io.SeekStart)
	return err
}
