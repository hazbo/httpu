package request

import (
	"github.com/hazbo/httpu/env"
	"github.com/hazbo/httpu/resource/request/headers"
	"github.com/hazbo/httpu/stash"
)

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

// TODO: This is essentially going to be the same for stash vars, stop this code
// repetition.
func (rs *RequestSpec) parseEnvVars() {
	rs.Uri = env.Parse(rs.Uri)
	rs.Method = env.Parse(rs.Method)
	rs.Data.contents = []byte(env.Parse(rs.Data.String()))

	for fdk, uv := range rs.FormData {
		for fdks, uvs := range uv {
			rs.FormData[fdk][fdks] = env.Parse(uvs)
		}
	}

	for k, _ := range rs.Headers {
		rs.Headers[k][0] = env.Parse(rs.Headers[k][0])
	}

	// Do the same as above, for all variants for the given request spec.
	for vi, _ := range rs.Variants {
		rs.Variants[vi].Path = env.Parse(rs.Variants[vi].Path)
		rs.Variants[vi].Data.contents = []byte(env.Parse(rs.Variants[vi].Data.String()))

		for i, _ := range rs.Variants[vi].Headers {
			// TODO: don't just do this for [0]
			rs.Variants[vi].Headers[i][0] = env.Parse(rs.Variants[vi].Headers[i][0])
		}

		for fdk, uv := range rs.Variants[vi].FormData {
			for fdks, uvs := range uv {
				rs.Variants[vi].FormData[fdk][fdks] = env.Parse(uvs)
			}
		}
	}
}

// parseStashVars goes through each possible instance of a stash variable
// existing and replaces the variable name with the newvalue if it does exist.
func (rs *RequestSpec) parseStashVars() {
	// Parse any variables within the URI, Method and request data.
	rs.Uri = stash.Parse(rs.Uri)
	rs.Method = stash.Parse(rs.Method)
	rs.Data.contents = []byte(stash.Parse(rs.Data.String()))

	for fdk, uv := range rs.FormData {
		for fdks, uvs := range uv {
			rs.FormData[fdk][fdks] = stash.Parse(uvs)
		}
	}

	for k, _ := range rs.Headers {
		rs.Headers[k][0] = stash.Parse(rs.Headers[k][0])
	}

	// Do the same as above, for all variants for the given request spec.
	for vi, _ := range rs.Variants {
		rs.Variants[vi].Path = stash.Parse(rs.Variants[vi].Path)
		rs.Variants[vi].Data.contents = []byte(stash.Parse(rs.Variants[vi].Data.String()))

		for i, _ := range rs.Variants[vi].Headers {
			// TODO: don't just do this for [0]
			rs.Variants[vi].Headers[i][0] = stash.Parse(rs.Variants[vi].Headers[i][0])
		}

		for fdk, uv := range rs.Variants[vi].FormData {
			for fdks, uvs := range uv {
				rs.Variants[vi].FormData[fdk][fdks] = stash.Parse(uvs)
			}
		}
	}
}
