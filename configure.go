package httpu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/hazbo/httpu/resource"
	"github.com/hazbo/httpu/resource/request"
	utils "github.com/hazbo/httpu/utils/common"
)

// Config is a wrapper for the Base config.
type Config struct {
	Project Project `json:"project"`
}

// the base URL which is then extended in "request" Resources using a URI and a
// path within a Variant.
//
// Headers that are set here will apply to all requests, but can be overridden
// within an individual request / variant.
//
// The resources are read as []string from the JSON, however as they are being
// unmarsheled, the JSON for the filename is fetched and is unmarsheled into the
// Resources that exist within Base.
type Project struct {
	URL           url.URL            `json:"url"`
	ResourceFiles resource.FilePaths `json:"resourceFiles"`
	ProjectPath   string

	Requests map[string]request.Request
}

const (
	packagesDir     = ".httpu/packages"
	projectFileName = "project.json"
)

// ConfigureFromFile reads in a base JSON config file and decodes it into Config
func ConfigureFromFile(filePath string) error {
	var (
		cfg []byte
		err error
		c   Config
	)

	// Set a default base project path
	utils.ProjectPath = "./"

	if cfg, err = ioutil.ReadFile(
		fmt.Sprintf("./%s/%s", filePath, projectFileName)); err != nil {

		// Look for project in the home packages directory
		hd, _ := utils.HomeDir()
		cfg, err = ioutil.ReadFile(fmt.Sprintf(
			"%s/%s/%s/%s", hd, packagesDir, filePath, projectFileName))

		if err != nil {
			return fmt.Errorf("Error loading config: %s", err)
		}

		// Override the project path to the home directory path
		utils.ProjectPath = fmt.Sprintf("%s/%s", hd, packagesDir)
	}

	err = json.Unmarshal(cfg, &c)
	if err != nil {
		return fmt.Errorf("Unable to parse JSON: %s", err)
	}

	c.Project.ProjectPath = fmt.Sprintf("%s/%s", utils.ProjectPath, filePath)

	for _, rf := range c.Project.ResourceFiles {
		err := rf.Load()
		if err != nil {
			return err
		}
	}

	c.Project.Requests = resource.Requests

	session = c.Project

	return nil
}

var session Project

func Session() Project {
	return session
}

// UnmarshalJSON is an implemenation of json.Unmarshaler and is used to parse
// the URL into a native url.URL type and the Headers into http.Header.
//
// If an invalid URL is passed, an error will occur at this point.
func (p *Project) UnmarshalJSON(j []byte) error {
	type Alias Project
	aux := &struct {
		URL string `json:"url"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(j, &aux); err != nil {
		return err
	}

	// Parse the URL to check it is valid
	urlp, err := url.Parse(aux.URL)

	if err != nil {
		return fmt.Errorf("url must be a valid URL: %s", err)
	}

	// Set the URL to a parsed url of type url.URL
	p.URL = *urlp

	return nil
}
