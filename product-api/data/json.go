package data

import (
	"encoding/json"
	"io"
)

func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)

}

func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}
