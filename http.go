package e365_gateway

import (
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
		def := new(respCommon)
		if err = json.NewDecoder(r.Body).Decode(def); err == nil {
			return nilT, wrap(fmt.Sprintf("%s: %s", r.Status, def.Error), ErrResponseStatusNotOK)
		}
		bts, _ := io.ReadAll(r.Body)
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
