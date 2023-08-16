package e365_gateway

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	"io"
)

func decodeSonic(src io.Reader, dst interface{}) error {

	if err := sonic.ConfigFastest.NewDecoder(src).Decode(dst); err != nil {
		return err
	}

	return nil
}

func decodeStd(src io.Reader, dst interface{}) error {
	if err := json.NewDecoder(src).Decode(dst); err != nil {
		return err
	}
	return nil
}
