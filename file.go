package e365_gateway

import (
	"context"
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
}

func NewFileAdapter(s Stand) FileAdapter {
	return FileAdapter{
		stand: s,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (fa FileAdapter) SetClientTimeout(t time.Duration) {
	fa.client.Timeout = t
}

func (fa FileAdapter) DownloadFile(ctx context.Context, id string) (io.ReadCloser, error) {
	link, err := fa.GetDownloadLink(ctx, id)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return nil, wrap(err.Error(), ErrCreateRequest)
	}
	// TODO какой клиент лучше использовать?
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, wrap(err.Error(), ErrSendRequest)
	}

	if response.StatusCode == http.StatusOK {
		return response.Body, nil
	}

	defer func() {
		_ = response.Body.Close()
	}()
	errBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, wrap(err.Error(), ErrResponseStatusNotOK)
	}

	return nil, wrap(string(errBody), ErrResponseStatusNotOK)

}

func (fa FileAdapter) GetDownloadLink(ctx context.Context, id string) (string, error) {

	if len(id) != uuid4Len {
		return "", wrap(id, ErrInvalidID)
	}

	url := fa.stand.url() + pubV1ApiDisc + id + methodGetLink
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", wrap(err.Error(), ErrCreateRequest)
	}
	request.Header = fa.stand.header()

	fr, err := doRequest[getFileLinkResp](fa.client, request)
	if err != nil {
		return "", err
	}

	if !fr.Success {
		return "", wrap(fr.Error, ErrResponseNotSuccess)
	}

	return fr.Link, nil

}

func (fa FileAdapter) NewDirectory(id string) Directory {
	return Directory{
		id:     id,
		stand:  fa.stand,
		client: &http.Client{Timeout: time.Second * 5},
	}
}

//func (fa FileAdapter) uploadFile(data io.Reader) (File, error) {
//
//	// TODO - implement
//	panic("not implemented")
//
//	return File{}, nil
//
//}
//
//func (fa FileAdapter) getDirectoriesList() (string, error) {
//	bts, err := json.Marshal(filter{
//		From:   0,
//		Size:   100,
//		Active: true,
//		SearchFilter: SearchFilter{
//			Fields: Fields{
//				"system": false,
//				//"__deletedAt": nil,
//			},
//			IDs: nil,
//			SortExpressions: []SortExpression{
//				{
//					Ascending: false,
//					Field:     "__name",
//				},
//			},
//			AtStatus:      nil,
//			StatusGroupId: "",
//		},
//	})
//
//	url := fa.stand.url() + pubV1ApiDiscDirList
//	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bts))
//	if err != nil {
//		return "", fmt.Errorf("failed creating request: %w", err)
//	}
//	request.Header = fa.stand.header()
//
//	response, err := fa.client.Do(request)
//	if err != nil {
//		return "", fmt.Errorf("failed sending request: %w", err)
//	}
//	defer func() {
//		_ = response.Body.Close()
//	}()
//
//	bt, err := io.ReadAll(response.Body)
//	if err != nil {
//
//	}
//
//	_ = bt
//
//	fr := new(getFileLinkResp)
//	if err = decodeStd(response.Body, fr); err != nil {
//		return "", fmt.Errorf("failed decoding response body: %w", err)
//	}
//
//	if response.StatusCode != http.StatusOK {
//		return "", fmt.Errorf("%w: %d (%s)", ErrResponseStatusNotOK, response.StatusCode, fr.Error)
//	}
//
//	if !fr.Success {
//		return "", fmt.Errorf("%w: %s", ErrResponseNotSuccess, fr.Error)
//	}
//
//	return fr.Link, nil
//
//}
