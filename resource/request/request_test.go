package request

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hazbo/httpu/stash"
	utils "github.com/hazbo/httpu/utils/common"
	"github.com/stretchr/testify/assert"
)

func TestVariant(t *testing.T) {
	r := Request{
		Kind: "request",
		Name: "test-request",
		Spec: RequestSpec{
			Uri: "/",
			Variants: Variants{
				Variant{
					Name: "test-variant",
				},
			},
		},
	}

	v, err := r.Variant("test-variant")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "test-variant", v.Name, "it should be test-variant")
}

func TestVariants(t *testing.T) {
	r := Request{
		Kind: "request",
		Name: "test-request",
		Spec: RequestSpec{
			Uri: "/",
			Variants: Variants{
				Variant{
					Name: "test-variant-1",
				},
				Variant{
					Name: "test-variant-2",
				},
			},
		},
	}

	vs := r.Variants()
	expecting := 2
	if len(vs) != expecting {
		t.Error("Expecting:", expecting, "got:", len(vs))
	}
}

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func fixture() string {
	return `{"error": false}`
}

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	return func() {
		server.Close()
	}
}

func TestMake(t *testing.T) {
	teardown := setup()
	defer teardown()

	r := Request{
		Kind: "request",
		Name: "test-request",
		Spec: RequestSpec{
			Uri:    "/",
			Method: "GET",
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture())
	})

	u, err := url.Parse(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	resp, _, err := r.Make(*u)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, `{"error": false}`, string(b), "JSON encoded, error : false")
}

func TestMakeWithVariant(t *testing.T) {
	teardown := setup()
	defer teardown()

	r := Request{
		Kind: "request",
		Name: "test-request",
		Spec: RequestSpec{
			Uri: "/",
			Variants: Variants{
				Variant{
					Name:   "test-variant",
					Method: "GET",
					Path:   "/",
				},
			},
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture())
	})

	u, err := url.Parse(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	resp, _, err := r.MakeWithVariant(*u, &r.Spec.Variants[0])
	if err != nil {
		log.Fatal(err)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, `{"error": false}`, string(b), "JSON encoded, error : false")
}

func TestUMake(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture())
	})

	u, err := url.Parse(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	hr := httpRequest{
		url:         u.String(),
		method:      "GET",
		headers:     http.Header{},
		data:        requestData{},
		formData:    url.Values{},
		stashValues: stash.StashValues{},
	}

	resp, _, err := hr.make()
	if err != nil {
		log.Fatal(err)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, `{"error": false}`, string(b), "JSON encoded, error : false")
}

func TestRequestBody(t *testing.T) {
	u, err := url.Parse(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	hr1 := httpRequest{
		url:     u.String(),
		method:  "GET",
		headers: http.Header{},
		data: requestData{
			contents: []byte(fixture()),
		},
		formData:    url.Values{},
		stashValues: stash.StashValues{},
	}

	rb1 := hr1.requestBody()
	assert.Equal(t, `{"error": false}`, rb1, "JSON encoded, error : false")

	uv := url.Values{}
	uv.Set("foo", "bar")
	uv.Add("something", "else")
	hr2 := httpRequest{
		url:         u.String(),
		method:      "GET",
		headers:     http.Header{},
		data:        requestData{},
		formData:    uv,
		stashValues: stash.StashValues{},
	}

	rb2 := hr2.requestBody()
	assert.Equal(t, "foo=bar&something=else", rb2, "It should be an empty string")
}

func TestRequestDataString(t *testing.T) {
	rd := requestData{contents: []byte(`{"error": false}`)}
	assert.Equal(t, `{"error": false}`, rd.String())
}

func TestRequestDataLoadConents(t *testing.T) {
	// Messy!
	utils.ProjectPath = "../../"
	rd := requestData{FromFile: "projects/test_project/data/test.json"}
	rd.loadContents()
	assert.Equal(t, `{"error": "false"}`, string(rd.contents))
}
