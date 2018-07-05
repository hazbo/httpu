package api

import (
	"log"
	"net/http"

	"github.com/hazbo/httpu/api/requests"
)

func StartServer() {
	http.HandleFunc("/requests", requests.GetRequests)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
