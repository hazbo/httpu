package request

import (
	"github.com/hazbo/httpu/env"
	"github.com/hazbo/httpu/resource/request/headers"
	"github.com/hazbo/httpu/stash"
)

type parser func(string) string

// Update modifies the request spec to include any data that has recently been
// added to the stash. If new values exist, the values within the request spec
// will be added at this point.
func (rs *RequestSpec) Update() {
	rs.addFormheader()
	rs.loadDataFiles()
	rs.parseEnvVars()
	rs.parseStashVars()
}

// addFormHeader adds a spesefic header to the request if form data has been
// attached to it. It checks to make sure that the form data slice is greater
// than 0 to do this.
func (rs *RequestSpec) addFormheader() {
	if len(rs.FormData) > 0 {
		rs.Headers[headers.ContentType] = []string{
			"application/x-www-form-urlencoded"}
	}
}

// loadDataFiles loads in any request data from seperate files for the given
// request spec.
func (rs *RequestSpec) loadDataFiles() {
	if rs.Data.FromFile != "" {
		rs.Data.loadContents()
	}
}

func parseVars(parse parser, rs *RequestSpec) {
	rs.Uri = parse(rs.Uri)
	rs.Method = parse(rs.Method)
	rs.Data.contents = []byte(parse(rs.Data.String()))

	for fdk, uv := range rs.FormData {
		for fdks, uvs := range uv {
			rs.FormData[fdk][fdks] = parse(uvs)
		}
	}

	for k, _ := range rs.Headers {
		rs.Headers[k][0] = parse(rs.Headers[k][0])
	}

	// Do the same as above, for all variants for the given request spec.
	for vi, _ := range rs.Variants {
		rs.Variants[vi].Path = parse(rs.Variants[vi].Path)
		rs.Variants[vi].Data.contents = []byte(parse(rs.Variants[vi].Data.String()))

		for i, _ := range rs.Variants[vi].Headers {
			for t, _ := range rs.Variants[vi].Headers[i] {
				rs.Variants[vi].Headers[i][t] = parse(rs.Variants[vi].Headers[i][t])
			}
		}

		for fdk, uv := range rs.Variants[vi].FormData {
			for fdks, uvs := range uv {
				rs.Variants[vi].FormData[fdk][fdks] = parse(uvs)
			}
		}
	}
}

// parseEnvVars goes through each possible instance of an environment variable
// existing and replaces the variable name with the new value if it does exist.
func (rs *RequestSpec) parseEnvVars() {
	parseVars(env.Parse, rs)
}

// parseStashVars goes through each possible instance of a stash variable
// existing and replaces the variable name with the new value if it does exist.
func (rs *RequestSpec) parseStashVars() {
	parseVars(stash.Parse, rs)
}
