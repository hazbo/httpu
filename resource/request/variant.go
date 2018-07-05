package request

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hazbo/httpu/stash"
)

// RequestVariant represents a given variant for an HTTP request. This could be
// related to the same resource, but using a different request method or path
// for example.
type Variant struct {
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Method      string            `json:"method"`
	Data        requestData       `json:"data"`
	FormData    url.Values        `json:"formData"`
	Headers     http.Header       `json:"headers"`
	StashValues stash.StashValues `json:"stashValues"`
}

func (v *Variant) UnmarshalJSON(j []byte) error {
	type Alias Variant
	aux := &struct {
		Headers []struct {
			Header string `json:"header"`
			Value  string `json:"value"`
		} `json:"headers"`
		FormData []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"formData"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(j, &aux); err != nil {
		return err
	}

	// Create the http headers
	h := http.Header{}

	for _, hobj := range aux.Headers {
		h.Add(hobj.Header, hobj.Value)
	}

	// Set Headers to type http.Header
	v.Headers = h

	f := url.Values{}
	for _, fdat := range aux.FormData {
		f.Add(fdat.Name, fdat.Value)
	}

	v.FormData = f

	// TODO: move this into the modifier
	if len(v.Data.contents) > 0 {
		return nil
	}

	v.Data.loadContents()
	return nil
}

// RequestVariants represent multiple RequestVariant
type Variants []Variant

func (vs Variants) Names() []string {
	var res []string
	for _, v := range vs {
		res = append(res, v.Name)
	}
	return res
}
