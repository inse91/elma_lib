package e365_gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	pubV1ApiDisc  = "/pub/v1/disk/file/"
	methodGetLink = "/get-link"
)

type FileAdapter struct {
	stand  Stand
	client *http.Client
	//header http.Header
}

func NewFileAdapter(s Stand) FileAdapter {
	return FileAdapter{
		stand: s,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (fa FileAdapter) UploadFile(data io.Reader) (File, error) {

	// TODO - implement
	panic("not implemented")

	return File{}, nil

}

func (fa FileAdapter) DownloadFile(id string) (io.ReadCloser, error) {
	link, err := fa.GetDownloadLink(id)
	if err != nil {
		return nil, err
	}

	//response, err := fa.client.Get(link)
	response, err := http.DefaultClient.Get(link)
	if err != nil {
		return nil, fmt.Errorf("failed downloading file: %w", err)
	}

	if response.StatusCode == http.StatusOK {
		return response.Body, nil
	}

	defer func() {
		_ = response.Body.Close()
	}()
	errBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading error response body: %w", err)
	}

	return nil, fmt.Errorf("%w: %s", ErrResponseNotSuccess, string(errBody))

}

func (fa FileAdapter) GetDownloadLink(id string) (string, error) {

	if len(id) != uuid4Len {
		return "", fmt.Errorf("%s: %w", id, ErrInvalidID)
	}

	url := fa.stand.url() + pubV1ApiDisc + id + methodGetLink
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = fa.stand.header()

	response, err := fa.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	fr := new(getFileLinkResp)
	if err = decodeStd(response.Body, fr); err != nil {
		return "", fmt.Errorf("failed decoding response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %d (%s)", ErrResponseStatusNotOK, response.StatusCode, fr.Error)
	}

	if !fr.Success {
		return "", fmt.Errorf("%w: %s", ErrResponseNotSuccess, fr.Error)
	}

	return fr.Link, nil

}

func (fa FileAdapter) getDirectoriesList() (string, error) {
	bts, err := json.Marshal(filter{
		From:   0,
		Size:   100,
		Active: true,
		SearchFilter: SearchFilter{
			Fields: Fields{
				"system": false,
				//"__deletedAt": nil,
			},
			IDs: nil,
			SortExpressions: []SortExpression{
				{
					Ascending: false,
					Field:     "__name",
				},
			},
			AtStatus:      nil,
			StatusGroupId: "",
		},
	})

	url := fa.stand.url() + pubV1ApiDiscDirList
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bts))
	if err != nil {
		return "", fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = fa.stand.header()

	response, err := fa.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	bt, err := io.ReadAll(response.Body)
	if err != nil {

	}

	_ = bt

	fr := new(getFileLinkResp)
	if err = decodeStd(response.Body, fr); err != nil {
		return "", fmt.Errorf("failed decoding response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %d (%s)", ErrResponseStatusNotOK, response.StatusCode, fr.Error)
	}

	if !fr.Success {
		return "", fmt.Errorf("%w: %s", ErrResponseNotSuccess, fr.Error)
	}

	return fr.Link, nil

}

func (fa FileAdapter) NewDirectory(id string) Directory {
	return Directory{
		id:     id,
		stand:  fa.stand,
		client: fa.client,
	}
}
