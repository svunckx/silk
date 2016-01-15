package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// EchoHandler gets an http.Handler that echos request data
// back in the response.
func EchoHandler() http.Handler {
	return http.HandlerFunc(handleEcho)
}

// EchoDataHandler gets an http.Handler that echos request data
// back in the response in JSON format.
func EchoDataHandler() http.Handler {
	return http.HandlerFunc(handleEchoData)
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	// set Server header
	w.Header().Set("Server", "EchoHandler")
	// write summary of request
	fmt.Fprintln(w, strings.ToUpper(r.Method), r.URL.Path)
	// put in the Content-Length
	var bodybuf bytes.Buffer
	if _, err := io.Copy(&bodybuf, r.Body); err != nil {
		log.Println("copying request into buffer failed:", err)
	}
	r.Header.Set("Content-Length", strconv.Itoa(bodybuf.Len()))
	// write headers
	writeSortedHeaders(w, r.Header)
	// write body
	if _, err := io.Copy(w, &bodybuf); err != nil {
		log.Println("copying request into response failed:", err)
	}
}

func handleEchoData(w http.ResponseWriter, r *http.Request) {
	// set Server header
	w.Header().Set("Server", "EchoDataHandler")

	out := make(map[string]interface{})
	out["method"] = r.Method
	out["path"] = r.URL.Path

	var bodybuf bytes.Buffer
	if _, err := io.Copy(&bodybuf, r.Body); err != nil {
		log.Println("copying request into buffer failed:", err)
	}
	r.Header.Set("Content-Length", strconv.Itoa(bodybuf.Len()))
	for k := range r.Header {
		for _, v := range r.Header[k] {
			out[k] = v
		}
	}
	out["bodystr"] = bodybuf.String()
	var bodyData interface{}
	if err := json.NewDecoder(&bodybuf).Decode(&bodyData); err != nil {
		out["bodyerr"] = err.Error()
	}
	out["body"] = bodyData
	if err := json.NewEncoder(w).Encode(out); err != nil {
		panic(err)
	}
}

func writeSortedHeaders(w io.Writer, headers http.Header) {
	// get header keys
	var keys []string
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range headers[k] {
			vb, err := json.Marshal(v)
			if err != nil {
				log.Println("silk/testutil: cannot marshal header value:", err)
				continue
			}
			fmt.Fprintln(w, "* "+k+":", string(vb))
		}
	}
}
