// +build ignore

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vugu/vugu/simplehttp"
)

func main() {
	wd, _ := os.Getwd()
	l := "127.0.0.1:8877"
	log.Printf("Starting HTTP Server at %q", l)
	h := simplehttp.New(wd, true)
	simplehttp.DefaultStaticData["CSSFiles"] = []string{
		"/assets/css/vendor.css",
		"/assets/css/style.css",
	}
	log.Fatal(http.ListenAndServe(l, h))
}
