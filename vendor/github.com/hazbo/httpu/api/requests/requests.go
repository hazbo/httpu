package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hazbo/httpu/resource"
	"github.com/hazbo/httpu/resource/request"
)

func GetRequests(rw http.ResponseWriter, r *http.Request) {
	var rs request.Requests
	for _, r := range resource.Requests {
		rs = append(rs, r)
	}
	j, err := json.Marshal(rs)
	if err != nil {
		log.Fatal(err)
		return
	}

	var out bytes.Buffer
	json.Indent(&out, j, "", "  ")

	fmt.Fprint(rw, string(out.String()))
}
