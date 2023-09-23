package e365_gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
)

func doRequest[T interface{}](cli *http.Client, req *http.Request) (T, error) {

	var nilT T

	r, err := cli.Do(req)
	if err != nil {
		return nilT, wrap(err.Error(), ErrSendRequest)
	}
	defer func() {
		_ = r.Body.Close()
	}()

	t := new(T)

	if r.StatusCode != http.StatusOK {
		bts, err := io.ReadAll(r.Body)
		if err != nil {
			return nilT, wrap(err.Error(), ErrReadResponseBody)
		}
		recoveredRespBody := io.NopCloser(bytes.NewBuffer(bts))
		def := new(respCommon)
		if err = json.NewDecoder(recoveredRespBody).Decode(def); err == nil {
			return nilT, wrap(fmt.Sprintf("%s: %s", r.Status, def.Error), ErrResponseStatusNotOK)
		}

		return nilT, wrap(string(bts), ErrResponseStatusNotOK)
	}

	if err = decodeStd(r.Body, t); err != nil {
		return nilT, wrap(err.Error(), ErrDecodeResponseBody)
	}

	return *t, nil

}

func decodeSonic(src io.Reader, dst interface{}) error {
	return sonic.ConfigFastest.NewDecoder(src).Decode(dst)
}

func decodeStd(src io.Reader, dst interface{}) error {
	if err := json.NewDecoder(src).Decode(dst); err != nil {
		return err
	}
	return nil
}
