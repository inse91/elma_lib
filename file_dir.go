package e365_gateway

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"time"
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

func (d Directory) SetClientTimeout(t time.Duration) {
	d.client.Timeout = t
}

// Upload создает файл в директории и загружет содержимое буфера в файл.
func (d Directory) Upload(ctx context.Context, buf *bytes.Buffer, name string) (File, error) {

	if buf == nil {
		return File{}, ErrNilItem
	}

	if buf.Len() == 0 {
		return File{}, ErrEmptyBuffer
	}

	if len(d.id) != uuid4Len {
		return File{}, wrap(d.id, ErrInvalidID)
	}

	hash := uuid.New().String()
	size := buf.Len()
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	formFile, err := w.CreateFormFile("file", name)
	if err != nil {
		return File{}, wrap(err.Error(), ErrCreateFormData)
	}
	if _, err = formFile.Write(buf.Bytes()); err != nil {
		return File{}, wrap(err.Error(), ErrWriteBytesBuffer)
	}
	if err = w.Close(); err != nil {
		return File{}, wrap(err.Error(), ErrCloseMultipartWriter)
	}

	url := d.stand.url() + pubV1ApiDiscDirInfo + d.id + methodUploadFile
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return File{}, wrap(err.Error(), ErrCreateRequest)
	}
	request.Header = d.stand.header()
	request.Header.Set("Content-Type", w.FormDataContentType())
	request.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", size, size))
	q := request.URL.Query()
	q.Set("hash", hash)
	request.URL.RawQuery = q.Encode()

	fr, err := doRequest[fileResponse](d.client, request)
	if err != nil {
		return File{}, err
	}
	if !fr.Success {
		return File{}, wrap(fr.Error, ErrResponseNotSuccess)
	}

	return fr.File, nil

}

// Info - получаетс информацию о директории.
func (d Directory) Info(ctx context.Context) (DirectoryInfo, error) {

	if len(d.id) != uuid4Len {
		return DirectoryInfo{}, wrap(d.id, ErrInvalidID)
	}

	url := d.stand.url() + pubV1ApiDiscDirInfo + d.id
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return DirectoryInfo{}, wrap(err.Error(), ErrCreateRequest)
	}
	request.Header = d.stand.header()

	di, err := doRequest[dirInfoResponse](d.client, request)
	if err != nil {
		return DirectoryInfo{}, err
	}
	if !di.Success {
		return DirectoryInfo{}, wrap(di.Error, ErrResponseNotSuccess)
	}

	return di.DirectoryInfo, nil

}
