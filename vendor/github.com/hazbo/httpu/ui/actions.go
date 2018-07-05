package ui

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hazbo/httpu"
	"github.com/hazbo/httpu/resource"
	"github.com/hazbo/httpu/resource/request"
	"github.com/hazbo/httpu/ui/printer"
	"github.com/jroimartin/gocui"
)

// makeRequest wraps httpu's req.Make, updates the relevant views and updates
// the request spec.
func makeRequest(
	req *request.Request) (*http.Response, request.RequestStat, error) {

	s := httpu.Session()

	resp, stat, err := req.Make(s.URL)
	if err != nil {
		return &http.Response{}, request.RequestStat{}, err
	}

	RequestView.Clear()
	req.Spec.Update()

	return resp, stat, nil
}

// makeRequestWithVariant makes a request using a variant, in a similar way to
// makeRequest.
func makeRequestWithVariant(
	req *request.Request,
	v *request.Variant) (*http.Response, request.RequestStat, error) {

	s := httpu.Session()
	resp, stat, err := req.MakeWithVariant(s.URL, v)
	if err != nil {
		return &http.Response{}, request.RequestStat{}, err
	}

	RequestView.Clear()
	req.Spec.Update()

	return resp, stat, nil
}

func writeRequestDataVariant(r *request.Request, v *request.Variant) {
	b := bytes.NewBufferString(printer.Color("Request:\n", printer.ColorGreen))
	b.WriteString(fmt.Sprintf("%s: %s%s\n\n", v.Method, r.Spec.Uri, v.Path))

	if len(v.Headers) > 0 {
		b.WriteString(printer.Color("Headers:\n", printer.ColorGreen))
	}

	for h, val := range v.Headers {
		b.WriteString(fmt.Sprintf("%s: %s\n", h, val[0]))
	}

	if len(v.FormData) == 0 && len(v.Data.String()) == 0 {
		fmt.Fprint(RequestView, b.String())
		return
	}

	b.WriteString(printer.Color("\nData:\n", printer.ColorGreen))

	jp := printer.NewJSONPrinter()

	if len(v.FormData) > 0 {
		b.WriteString(fmt.Sprintf("%s", v.FormData.Encode()))
	} else {
		jp.PrintString(b, v.Data.String())
	}

	fmt.Fprint(RequestView, b.String())
}

// writeRequestData write the data being sent into the request view.
func writeRequestData(r *request.Request) {
	b := bytes.NewBufferString(printer.Color("Request:\n", printer.ColorGreen))
	b.WriteString(fmt.Sprintf("%s: %s\n\n", r.Spec.Method, r.Spec.Uri))

	if len(r.Spec.Headers) > 0 {
		b.WriteString(printer.Color("Headers:\n", printer.ColorGreen))
	}

	for h, val := range r.Spec.Headers {
		b.WriteString(fmt.Sprintf("%s: %s\n", h, val[0]))
	}

	if len(r.Spec.FormData) == 0 && len(r.Spec.Data.String()) == 0 {
		fmt.Fprint(RequestView, b.String())
		return
	}

	b.WriteString(printer.Color("\nData:\n", printer.ColorGreen))

	jp := printer.NewJSONPrinter()

	if len(r.Spec.FormData) > 0 {
		b.WriteString(fmt.Sprintf("%s", r.Spec.FormData.Encode()))
	} else {
		jp.PrintString(b, r.Spec.Data.String())
	}

	fmt.Fprint(RequestView, b.String())
}

func writeResponseData(r *http.Response, stat request.RequestStat) {
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	ResponseView.Clear()

	jp := printer.NewJSONPrinter()
	jp.PrintString(ResponseView, string(b))

	StatusCodeView.Clear()
	RequestTimeView.Clear()

	switch true {
	case r.StatusCode <= 299 && r.StatusCode >= 200:
		StatusCodeView.BgColor = gocui.ColorGreen
	case r.StatusCode <= 499 && r.StatusCode >= 400:
		StatusCodeView.BgColor = gocui.ColorYellow
	case r.StatusCode <= 599 && r.StatusCode >= 500:
		StatusCodeView.BgColor = gocui.ColorRed
	}

	fmt.Fprintf(StatusCodeView, "  %d", r.StatusCode)
	fmt.Fprintf(RequestTimeView, "%dms", stat.Total)
}

// defaultKeyPress is called at the end of each keypress in default mode
// to search for resources that are autocompleated and then shown in the
// request view.
func defaultKeyPress(v *gocui.View) error {
	if HttpuMode == CommandMode {
		return nil
	}
	RequestView.Clear()

	if cmdBarBuffer() == "" {
		fmt.Fprint(RequestView, resource.Requests)
		return nil
	}

	for _, rc := range resource.SearchRequests(cmdBarBuffer()) {
		fmt.Fprintf(RequestView, "%s\n", rc)
	}
	return nil
}
