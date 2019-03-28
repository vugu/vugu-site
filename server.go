// +build !wasm

package main

//go:generate vugugen .

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vugu/vugu"
	"github.com/vugu/vugu/simplehttp"
)

func main() {
	dev := flag.Bool("dev", false, "Enable development features")
	dir := flag.String("dir", ".", "Project directory")
	httpl := flag.String("http", "127.0.0.1:8877", "Listen for HTTP on this host:port")
	flag.Parse()
	wd, _ := filepath.Abs(*dir)
	os.Chdir(wd)
	log.Printf("Starting HTTP Server at %q", *httpl)
	simplehttp.DefaultTemplateDataFunc = func(r *http.Request) interface{} {
		ret := make(map[string]interface{}, 3)
		ret["CSSFiles"] = []string{
			"/assets/css/vendor.css",
			"/assets/css/style.css",
		}
		ret["Request"] = r

		var buf bytes.Buffer
		rootInst, err := vugu.New(&Root{Router: NewServerRouter(*r.URL)}, nil)
		if err != nil {
			log.Printf("Error creating new component: %v", err)
			return ret
		}
		env := vugu.NewStaticHTMLEnv(&buf, rootInst, vugu.RegisteredComponentTypes())
		err = env.Render()
		if err != nil {
			log.Printf("Error creating new component: %v", err)
			return ret
		}
		ret["ServerRenderedOutput"] = template.HTML(buf.String())

		return ret
	}
	// simplehttp.DefaultStaticData["CSSFiles"] = []string{
	// 	"/assets/css/vendor.css",
	// 	"/assets/css/style.css",
	// }
	h := simplehttp.New(wd, *dev)
	log.Fatal(http.ListenAndServe(*httpl, h))
}
