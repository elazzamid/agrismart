package farm

import (
	"encoding/json"
	"errors"
	"io"
)

// decodeSingleJSON accepts exactly one JSON value and rejects trailing data.
func decodeSingleJSON(dec *json.Decoder, dst any) error {
	if err := dec.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}
