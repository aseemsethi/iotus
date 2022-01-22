package main

import (
	"fmt"
	"github.com/aseemsethi/iotus/utils"
	"net/http"
)

func api(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Fprintf(w, "hello\n")
}

func main() {
	fmt.Printf("\nIOTUS Tool Starting..")
	utils.Mqtt_init()
	http.HandleFunc("/api", api)
	http.ListenAndServe(":8090", nil)
}
