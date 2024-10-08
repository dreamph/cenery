package echo

import (
	"errors"
	"io"
	"mime/multipart"

	"github.com/dreamph/cenery"
	"github.com/labstack/echo/v4"
)

func FormFile(c echo.Context, fileKey string) *cenery.FileData {
	file, _ := formFile(c, fileKey, false)
	return file
}

func formFile(c echo.Context, fileKey string, errorIfNotFound bool) (*cenery.FileData, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files, ok := form.File[fileKey]
	if !ok {
		if errorIfNotFound {
			return nil, errors.New("fileKey not found:" + fileKey)
		}
		return nil, nil
	}
	file := files[0]
	return &cenery.FileData{
		FileData:        fileHeaderToBytes(file),
		FileSize:        file.Size,
		FileContentType: file.Header["Content-Type"][0],
		FileName:        file.Filename,
	}, nil
}

func FormFiles(c echo.Context, fileKey string) *[]cenery.FileData {
	files, _ := formFiles(c, fileKey, false)
	return files
}

func formFiles(c echo.Context, fileKey string, errorIfNotFound bool) (*[]cenery.FileData, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files, ok := form.File[fileKey]
	if !ok {
		if errorIfNotFound {
			return nil, errors.New("fileKey not found:" + fileKey)
		}
		return nil, nil
	}

	var uploadFiles []cenery.FileData
	for _, file := range files {
		uploadFiles = append(uploadFiles, cenery.FileData{
			FileData:        fileHeaderToBytes(file),
			FileSize:        file.Size,
			FileContentType: file.Header["Content-Type"][0],
			FileName:        file.Filename,
		})
	}
	return &uploadFiles, nil
}

func fileHeaderToBytes(h *multipart.FileHeader) []byte {
	file, err := h.Open()
	if err != nil {
		return nil
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	data, err := io.ReadAll(file)
	if err != nil {
		return nil
	}
	return data
}
