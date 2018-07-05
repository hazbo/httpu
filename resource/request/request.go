package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/hazbo/httpu/resource/request/headers"
	"github.com/hazbo/httpu/stash"
	utils "github.com/hazbo/httpu/utils/common"
)

// requestData represents the request body that will be sent.
type requestData struct {
	FromFile string `json:"fromFile"`
	contents []byte
}

// Strings returns the contents of the request body as a string.
func (rd requestData) String() string {
	return string(rd.contents)
}

// loadContents loads the contents of a text file to use as the request data
// when passed to a request.
func (rd *requestData) loadContents() error {
	if len(rd.contents) > 0 {
		return nil
	}
	b, err := ioutil.ReadFile(
		fmt.Sprintf("%s/%s", utils.ProjectPath, rd.FromFile))
	if err != nil {
		return err
	}
	rd.contents = b
	return nil
}

// RequestSpec represents the specification of the base request resource without
// any variants. A request can be made this way individually if there are no
// variants of it.
type RequestSpec struct {
	Uri         string            `json:"uri"`
	Method      string            `json:"method"`
	Data        requestData       `json:"data"`
	FormData    url.Values        `json:"formData"`
	Headers     http.Header       `json:"headers"`
	Variants    Variants          `json:"variants"`
	StashValues stash.StashValues `json:"stashValues"`
}

// UnmarshalJSON will ensure that the request headers that are by default passed
// through as string values will become type http.Header upon unmarshling the
// JSON.
func (rs *RequestSpec) UnmarshalJSON(j []byte) error {
	type Alias RequestSpec
	aux := &struct {
		Headers []struct {
			Header string `json:header"`
			Value  string `json:"value"`
		} `json:"headers"`
		FormData []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"formData"`
		*Alias
	}{
		Alias: (*Alias)(rs),
	}

	if err := json.Unmarshal(j, &aux); err != nil {
		return err
	}

	h := http.Header{}

	for _, hobj := range aux.Headers {
		h.Add(hobj.Header, hobj.Value)
	}

	rs.Headers = h

	f := url.Values{}
	for _, fdat := range aux.FormData {
		f.Add(fdat.Name, fdat.Value)
	}

	rs.FormData = f
	return nil
}

// Request represents a single resource which in itself can represent multiple
// types of requests for said resource, by creating RequestVariants.
type Request struct {
	Kind string      `json:"kind"`
	Name string      `json:"name"`
	Spec RequestSpec `json:"spec"`
}

// Variant checks for and returns a variant given by it's name.
func (r Request) Variant(n string) (Variant, error) {
	for _, v := range r.Spec.Variants {
		if v.Name == n {
			return v, nil
		}
	}
	return Variant{}, fmt.Errorf("Variant \"%s\" does not exist.", n)
}

// Variants returns a slice of Variant for each one that exists on the given
// Request resource.
func (r Request) Variants() Variants {
	return r.Spec.Variants
}

// httpRequest is an internal struct to store information about a spesefic
// request that will be made, regardless if there is a variant or not.
type httpRequest struct {
	url         string
	method      string
	headers     http.Header
	data        requestData
	formData    url.Values
	stashValues stash.StashValues
}

// Make makes a single HTTP request without a variant. Only the fields that come
// from the request will be used.
func (r *Request) Make(baseURL url.URL) (*http.Response, RequestStat, error) {
	// We updatet the request spec here before making a request to make sure it
	// has all needed data, both read from files and parsed from variables
	// contained within the config.
	r.Spec.Update()

	hr := httpRequest{
		url:         fmt.Sprintf("%s%s", baseURL.String(), r.Spec.Uri),
		method:      r.Spec.Method,
		headers:     r.Spec.Headers,
		data:        r.Spec.Data,
		formData:    r.Spec.FormData,
		stashValues: r.Spec.StashValues,
	}
	return hr.make()
}

// Make makes an HTTP request requing the baseURL and the variant of the request
// resource.
func (r *Request) MakeWithVariant(
	baseURL url.URL, v *Variant) (*http.Response, RequestStat, error) {
	r.Spec.Update()

	// The variant headers need the base request headers added to it before the
	// request is made.
	v.Headers = headers.Concat(r.Spec.Headers, v.Headers)
	hr := httpRequest{
		url: fmt.Sprintf(
			"%s%s%s", baseURL.String(), r.Spec.Uri, v.Path),
		method:      v.Method,
		headers:     v.Headers,
		data:        v.Data,
		formData:    v.FormData,
		stashValues: v.StashValues,
	}
	return hr.make()
}

type RequestStat struct {
	Total int
}

// make makes a request for either a standalone request, or a request with a
// variant. It doesn't care about which one, as long as the url, request method
// and headers are all passed through.
func (hr httpRequest) make() (*http.Response, RequestStat, error) {
	client := &http.Client{}

	// Prepare the request
	req, err := http.NewRequest(
		hr.method, hr.url, strings.NewReader(hr.requestBody()))

	if err != nil {
		return &http.Response{},
			RequestStat{},
			fmt.Errorf("Could not construct request: %s", err)
	}

	req.Header = hr.headers

	// Get the start time jsut before making the request
	start := time.Now()

	resp, err := client.Do(req)

	// Record the end time, even before catching any errors
	end := time.Now().Sub(start)

	if err != nil {
		return &http.Response{}, RequestStat{},
			fmt.Errorf("Error making request: %s", err)
	}

	rs := RequestStat{
		Total: int(end / time.Millisecond),
	}

	return hr.applyStash(resp), rs, nil
}

// requestBody checks to see if there is any form data present. If so, an
// encoded form of this is returned as a string. If not, the contents of
// Request.Spec.Data or Variant.Data is returned as a string instead.
func (hr httpRequest) requestBody() string {
	if len(hr.formData) > 0 {
		return hr.formData.Encode()
	}
	return hr.data.String()
}

// applyStash looks as the response body from the given request, checks to see
// if it can find the needed JSON path from the stash values, and store that in
// memory ready to be used within another request.
func (hr httpRequest) applyStash(r *http.Response) *http.Response {
	if len(hr.stashValues) == 0 {
		return r
	}

	rb, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(rb))
	for i, _ := range hr.stashValues {
		nsv, _, _, _ := jsonparser.Get(rb, hr.stashValues[i].JsonPath...)
		hr.stashValues[i].Value = string(nsv)
	}
	hr.stashValues.Push()
	return r
}

// Requests represents multiple Request resources.
type Requests []Request
