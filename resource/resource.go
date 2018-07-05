package resource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/buger/jsonparser"
	"github.com/hazbo/httpu/resource/request"
	utils "github.com/hazbo/httpu/utils/common"
)

// RequestMap represents a list of all requests associated by name with a
// string.
type RequestMap map[string]request.Request

// String returns a list of all requests using the request map, along with any
// variants that exist within the request spec.
func (rm RequestMap) String() string {
	var b bytes.Buffer
	b.WriteString("Requests:\n\n")

	for _, r := range allReqVars() {
		b.WriteString(fmt.Sprintf("%s\n", r))
	}

	return b.String()
}

var (
	// Requests is a map of each Request resource accesseble by it's name.
	Requests = RequestMap{}
)

// Search searches through the loaded resources to see if there is a match for
// the given query string. This will in turn call searchRequestVariants to find
// the variants for that resource.
//
// TODO: Ideally we should not be running this so often as we are (currently on
// each keypress to CmdBar in DefaultMode for default UI). These return values
// should be stored somewhere so they can be quickly accessed without having to
// keep searching through each request and each of it's variants.
//
// Also we might not want to conflate retuning each variant for each request in
// what's named SearchRequests.
func SearchRequests(query string) []string {
	var res []string
	if query == "" {
		return res
	}
	for _, rv := range allReqVars() {
		if len(query) <= len(rv) && rv[:len(query)] == query {
			res = append(res, rv)
		}
	}
	return res
}

// allReqVars gets all requests and variants then returns them in the format
// {request}.{variant}, akin to how the user interacts with command bar
// included with the default user interface.
func allReqVars() []string {
	var reqVars []string
	for name, r := range Requests {
		if r.Spec.Method != "" && r.Spec.Uri != "" {
			reqVars = append(reqVars, fmt.Sprintf("%s", name))
		}
		for _, v := range r.Variants() {
			reqVars = append(reqVars, fmt.Sprintf("%s.%s", name, v.Name))
		}
	}
	sort.Strings(reqVars)
	return reqVars
}

// loadRequest loads the config for a request resource and unmarshals it to it's
// native type.
func loadRequest(cfg []byte) (request.Request, error) {
	var r request.Request
	err := json.Unmarshal(cfg, &r)
	if err != nil {
		return request.Request{}, err
	}
	return r, nil
}

// FilePath represents a given filepath for a resource
type FilePath string

// Load attempts to load a given resource file and create a new resource based
// on the 'kind' specefied within the config.
func (fp FilePath) Load() error {
	var (
		res []byte
		err error
	)

	res, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", utils.ProjectPath, string(fp)))
	if err != nil {
		return err
	}

	kind, _, _, err := jsonparser.Get(res, "kind")
	if err != nil {
		return err
	}

	name, _, _, err := jsonparser.Get(res, "name")
	if err != nil {
		return err
	}

	// Check the kind and assign the name for the given resource to be stored
	switch string(kind) {
	case "request":
		req, err := loadRequest(res)
		if err != nil {
			return err
		}
		req.Spec.Update()
		Requests[string(name)] = req
	}
	return nil
}

// FilePaths represents multiple filepaths.
type FilePaths []FilePath
