package httpG

import (
	"fmt"
	"github.com/aseemsethi/iotus/db"
	"net/http"
)

func ApiCustomers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	str := fmt.Sprintf("%#v", db.C.Customers)
	fmt.Fprintf(w, str)
}
