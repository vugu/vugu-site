// +build !wasm

package main

//go:generate vugugen .

import (
	"bytes"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/vugu/vugu"
	"github.com/vugu/vugu/simplehttp"
)

func init() {
	simplehttp.DefaultPageTemplateSource = `<!doctype html>
<html>
<head>

<link rel="apple-touch-icon" sizes="180x180" href="/assets/images/favicon/apple-touch-icon.png">
<link rel="icon" type="image/png" sizes="32x32" href="/assets/images/favicon/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/assets/images/favicon/favicon-16x16.png">
<!-- link rel="manifest" href="/assets/images/favicon/site.webmanifest" -->
<link rel="mask-icon" href="/assets/images/favicon/safari-pinned-tab.svg" color="#5bbad5">
<meta name="msapplication-TileColor" content="#da532c">
<meta name="theme-color" content="#ffffff">

{{if .Title}}
<title>{{.Title}}</title>
{{else}}
<title>Vugu Dev - {{.Request.URL.Path}}</title>
{{end}}
<meta charset="utf-8"/>
{{if .MetaTags}}{{range $k, $v := .MetaTags}}
<meta name="{{$k}}" content="{{$v}}"/>
{{end}}{{end}}
<meta name="viewport" content="width=device-width, initial-scale=1">
{{if .CSSFiles}}{{range $f := .CSSFiles}}
<link rel="stylesheet" href="{{$f}}" />
{{end}}{{end}}
<script src="https://cdn.jsdelivr.net/npm/text-encoding@0.7.0/lib/encoding.min.js"></script> <!-- MS Edge polyfill -->
<script src="/wasm_exec.js"></script>

<!-- Global site tag (gtag.js) - Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=UA-137242493-1"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());

  gtag('config', 'UA-137242493-1');
</script>

</head>
<body>
<div id="root_mount_parent">
{{if .ServerRenderedOutput}}{{.ServerRenderedOutput}}{{else}}
<img style="position: absolute; top: 50%; left: 50%;" src="https://cdnjs.cloudflare.com/ajax/libs/galleriffic/2.0.1/css/loader.gif">
{{end}}
</div>
<script>
var wasmSupported = (typeof WebAssembly === "object");
if (wasmSupported) {
	if (!WebAssembly.instantiateStreaming) { 
		WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		};
	}
	const go = new Go();
	WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject).then((result) => {
		go.run(result.instance);
	});
} else {
	//document.getElementById("root_mount_parent").innerHTML = 'This application requires WebAssembly support.  Please upgrade your browser.';
	console.log("Full functionality requires WebAssembly support.  Please upgrade your browser.")
}
</script>
</body>
</html>
`
}

func main() {
	dev := flag.Bool("dev", false, "Enable development features")
	dir := flag.String("dir", ".", "Project directory")
	httpl := flag.String("http", "127.0.0.1:8877", "Listen for HTTP on this host:port")
	staticout := flag.String("staticout", "", "Write static HTML files to this folder and then exit (instead of running the webserver)")
	flag.Parse()
	wd, _ := filepath.Abs(*dir)
	os.Chdir(wd)
	simplehttp.DefaultTemplateDataFunc = func(r *http.Request) interface{} {
		ret := make(map[string]interface{}, 3)
		ret["CSSFiles"] = []string{
			"/assets/css/vendor.css",
			"/assets/css/style.css",
		}
		ret["Request"] = r

		sitePage := allPages.FindByPath(r.URL.Path)
		if sitePage == nil {
			return nil // this makes it 404
		}
		ret["Title"] = sitePage.LongTitle
		ret["MetaTags"] = map[string]string{
			"description": sitePage.Description,
		}

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

	h := simplehttp.New(wd, *dev)

	// if staticout then we render all the pages to a directory and exit
	if *staticout != "" {
		outdir, err := filepath.Abs(*staticout)
		if err != nil {
			panic(err)
		}
		for _, pg := range allPages {
			// pgoutf, err := os.OpenFile(pgoutpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			// if err != nil {
			// 	panic(err)
			// }
			// defer pgoutf.Close()
			req, err := http.NewRequest("GET", pg.Path, nil)
			if err != nil {
				panic(err)
			}
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			res := rr.Result()
			if res.StatusCode != 200 {
				log.Printf("Got back bad status code for page %v", pg.Path)
				log.Fatal(res)
			}
			defer res.Body.Close()

			ppath := pg.Path
			if ppath == "/" {
				ppath = "/index.html"
			} else {
				ppath = ppath + ".html"
			}
			pgoutpath := filepath.Join(outdir, ppath)
			os.MkdirAll(filepath.Dir(pgoutpath), 0755)
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(pgoutpath, b, 0644)
			if err != nil {
				panic(err)
			}

		}
		return
	}

	log.Printf("Starting HTTP Server at %q", *httpl)
	log.Fatal(http.ListenAndServe(*httpl, h))
}
