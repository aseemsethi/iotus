package httpG

import (
	"fmt"
	"github.com/aseemsethi/iotus/db"
	"net/http"
	"strconv"
)

// Main HTTP GET
func ApiCustomers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	param1 := req.URL.Query().Get("cid")
	if param1 != "" {
		fmt.Printf("\nQuery param cid is %s", param1)
		for _, v := range db.C.Customers {
			if strconv.Itoa(v.Cid) == param1 {
				str := fmt.Sprintf("%#v", v)
				fmt.Fprintf(w, str)
			}
		}
	} else {
		fmt.Printf("Query param cid is nil")
		str := fmt.Sprintf("%#v", db.C.Customers)
		fmt.Fprintf(w, str)
	}
}
