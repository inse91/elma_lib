package e365_gateway

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
)

const (
	pubV1ApiDiscDirList = "/pub/v1/disk/directory/list"
	pubV1ApiDiscDirInfo = "/pub/v1/disk/directory/"
	methodUploadFile    = "/upload"
)

type Directory struct {
	stand  Stand
	client *http.Client
	id     string
}

func (d Directory) Upload(buf *bytes.Buffer, name string) (File, error) {

	if buf == nil {
		return File{}, ErrNilItem
	}

	hash := uuid.New().String()
	size := buf.Len()
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	formFile, err := w.CreateFormFile("file", name)
	if err != nil {
		return File{}, fmt.Errorf("failed creating form data: %w", err)
	}
	if _, err = formFile.Write(buf.Bytes()); err != nil {
		return File{}, fmt.Errorf("failed writing to form data from file: %w", err)
	}
	if err = w.Close(); err != nil {
		return File{}, fmt.Errorf("failed closing writer: %w", err)
	}

	url := d.stand.url() + pubV1ApiDiscDirInfo + d.id + methodUploadFile
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return File{}, fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = d.stand.header()
	request.Header.Set("Content-Type", w.FormDataContentType())
	request.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", size, size))
	q := request.URL.Query()
	q.Set("hash", hash)
	request.URL.RawQuery = q.Encode()

	response, err := d.client.Do(request)
	if err != nil {
		return File{}, fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	fr := new(fileResponse)
	if err = decodeStd(response.Body, fr); err != nil {
		return File{}, err
	}

	if response.StatusCode != http.StatusOK {
		return File{}, fmt.Errorf("%w: %d (%s)", ErrResponseNotOK, response.StatusCode, fr.Error)
	}

	if !fr.Success {
		return File{}, fmt.Errorf("%w: %s", ErrResponseNotSuccess, fr.Error)
	}

	return fr.File, nil

}

func (d Directory) Info() (DirectoryInfo, error) {

	url := d.stand.url() + pubV1ApiDiscDirInfo + d.id
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return DirectoryInfo{}, fmt.Errorf("failed creating new request: %s", err)
	}
	request.Header = d.stand.header()

	response, err := d.client.Do(request)
	if err != nil {
		return DirectoryInfo{}, fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	di := new(dirInfoResponse)
	if err = decodeStd(response.Body, di); err != nil {
		return DirectoryInfo{}, fmt.Errorf("failed decoding response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return DirectoryInfo{}, fmt.Errorf("%w: %d (%s)", ErrResponseNotOK, response.StatusCode, di.Error)
	}

	if !di.Success {
		return DirectoryInfo{}, fmt.Errorf("%w: %s", ErrResponseNotSuccess, di.Error)
	}

	return di.Directory, nil

}
