package resource

import (
	"testing"

	"github.com/hazbo/httpu/resource/request"
	"github.com/stretchr/testify/assert"
)

func TestSearchRequests(t *testing.T) {
	Requests = map[string]request.Request{
		"testrequest1": request.Request{
			Name: "testrequest1",
			Kind: "request",
			Spec: request.RequestSpec{
				Variants: request.Variants{
					request.Variant{
						Name: "testreqvar1",
					},
				},
			}},
		"testrequest2": request.Request{
			Name: "testrequest2",
			Kind: "request",
			Spec: request.RequestSpec{
				Variants: request.Variants{
					request.Variant{
						Name: "testreqvar2",
					},
				},
			}},
		"anotherRequest": request.Request{
			Name: "anotherRequest",
			Kind: "request",
			Spec: request.RequestSpec{
				Variants: request.Variants{
					request.Variant{
						Name: "testreqvar3",
					},
				},
			}},
	}

	assert.Equal(t, 2, len(SearchRequests("test")), "there should be two request / variants")
	assert.Equal(t, 1, len(SearchRequests("anotherRequest.testr")), "there should be one request / variants")
	assert.Equal(t, 1, len(SearchRequests("testrequest2.t")), "there should be one request / variants")
	assert.Equal(t, 1, len(SearchRequests("ano")), "there should be one request / variants")
	assert.Equal(t, 0, len(SearchRequests("nothing")), "there should be zero request / variants")
}
