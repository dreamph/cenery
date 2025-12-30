package chi

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/dreamph/cenery"
)

func FormFile(r *http.Request, fileKey string) *cenery.FileData {
	file, _ := formFile(r, fileKey, false)
	return file
}

func formFile(r *http.Request, fileKey string, errorIfNotFound bool) (*cenery.FileData, error) {
	form, err := multipartForm(r)
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
	contentType := "application/octet-stream" // default
	if ct := file.Header.Get("Content-Type"); ct != "" {
		contentType = ct
	}
	return &cenery.FileData{
		FileData:        fileHeaderToBytes(file),
		FileSize:        file.Size,
		FileContentType: contentType,
		FileName:        file.Filename,
	}, nil
}

func FormFiles(r *http.Request, fileKey string) *[]cenery.FileData {
	files, _ := formFiles(r, fileKey, false)
	return files
}

func formFiles(r *http.Request, fileKey string, errorIfNotFound bool) (*[]cenery.FileData, error) {
	form, err := multipartForm(r)
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

	uploadFiles := make([]cenery.FileData, len(files))
	for i, file := range files {
		contentType := "application/octet-stream" // default
		if ct := file.Header.Get("Content-Type"); ct != "" {
			contentType = ct
		}
		uploadFiles[i] = cenery.FileData{
			FileData:        fileHeaderToBytes(file),
			FileSize:        file.Size,
			FileContentType: contentType,
			FileName:        file.Filename,
		}
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

// FormFileStream returns a streaming file upload (no memory allocation)
func FormFileStream(r *http.Request, fileKey string) (*cenery.FileStream, error) {
	form, err := multipartForm(r)
	if err != nil {
		return nil, err
	}

	files, ok := form.File[fileKey]
	if !ok {
		return nil, errors.New("fileKey not found:" + fileKey)
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	contentType := "application/octet-stream" // default
	if ct := fileHeader.Header.Get("Content-Type"); ct != "" {
		contentType = ct
	}

	return &cenery.FileStream{
		File:            file,
		FileName:        fileHeader.Filename,
		FileSize:        fileHeader.Size,
		FileContentType: contentType,
	}, nil
}

// FormFilesStream returns multiple streaming file uploads (no memory allocation)
func FormFilesStream(r *http.Request, fileKey string) ([]*cenery.FileStream, error) {
	form, err := multipartForm(r)
	if err != nil {
		return nil, err
	}

	fileHeaders, ok := form.File[fileKey]
	if !ok {
		return nil, errors.New("fileKey not found:" + fileKey)
	}

	streams := make([]*cenery.FileStream, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			// Close previously opened files on error
			for _, s := range streams {
				_ = s.File.Close()
			}
			return nil, err
		}

		contentType := "application/octet-stream" // default
		if ct := fileHeader.Header.Get("Content-Type"); ct != "" {
			contentType = ct
		}

		streams = append(streams, &cenery.FileStream{
			File:            file,
			FileName:        fileHeader.Filename,
			FileSize:        fileHeader.Size,
			FileContentType: contentType,
		})
	}

	return streams, nil
}

func multipartForm(r *http.Request) (*multipart.Form, error) {
	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return nil, err
		}
	}
	return r.MultipartForm, nil
}
