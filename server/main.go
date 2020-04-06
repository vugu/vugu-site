package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/vugu/vgrouter"

	"github.com/d0sbit/werr"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu-site/app"
	"github.com/vugu/vugu/staticrender"
)

func main() {

	dev := flag.Bool("dev", false, "Enable development server mode")
	build := flag.Bool("build", false, "Build static output")

	devhttp := flag.String("devhttp", "127.0.0.1:8921", "In dev mode the host:port to listen on")

	flag.Parse()

	switch {
	case *dev:
		chdirToGomod()
		runDevServer(*devhttp)

	case *build:
		panic(fmt.Errorf("not yet implemented"))
	default:
		log.Fatal("A mode must be specified, try with -h for help")
	}

}

func chdirToGomod() {

	dir := "."

	for i := 0; i < 10; i++ {
		_, err := os.Stat(filepath.Join(dir, "go.mod"))
		if os.IsNotExist(err) {
			dir = filepath.Join(dir, "..")
			continue
		} else if err != nil {
			log.Fatalf("Error while looking for go.mod: %v", err)
		}
		err = os.Chdir(dir)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Fatalf("Could not find go.mod in current directory or parents")

}

func runDevServer(addr string) {

	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting dev server at %q in %q", addr, dir)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		werr.WriteError(w, func() error {

			p := path.Clean("/" + r.URL.Path)

			// if no file extension then we assume it's a page and use the vugu rendering workflow
			if path.Ext(p) == "" {

				buildEnv, err := vugu.NewBuildEnv()
				if err != nil {
					return err
				}
				var rbuf bytes.Buffer
				renderer := staticrender.New(&rbuf)
				app, rootBuilder := app.VuguSetup(buildEnv, renderer.EventEnv())

				notFound := false
				nfh := app.Router.GetNotFound()
				app.Router.SetNotFound(vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
					notFound = true
					if nfh != nil {
						nfh.RouteHandle(rm)
					}
				}))

				app.Router.ProcessRequest(r)

				buildResults := buildEnv.RunBuild(rootBuilder)
				// log.Printf("CSS len = %d", len(buildResults.Out.CSS))
				err = renderer.Render(buildResults)
				if err != nil {
					return err
				}

				if notFound {
					w.WriteHeader(404)
				}
				w.Write(rbuf.Bytes())

				return nil
			}

			// otherwise it's a static file
			f, err := http.Dir(filepath.Join(dir, "static")).Open(p)
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return nil
			} else if err != nil {
				return err
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				return err
			}

			http.ServeContent(w, r, p, fi.ModTime(), f)

			return nil
		}())

	})

	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	log.Fatal(srv.ListenAndServe())

}

// buildEnv, err := vugu.NewBuildEnv()
// if err != nil {
// 	panic(err)
// }

// renderer, err := domrender.NewJSRenderer(*mountPoint)
// if err != nil {
// 	panic(err)
// }
// defer renderer.Release()

// rootBuilder := app.VuguSetup(buildEnv, renderer.EventEnv())

// for ok := true; ok; ok = renderer.EventWait() {

// 	buildResults := buildEnv.RunBuild(rootBuilder)

// 	err = renderer.Render(buildResults)
// 	if err != nil {
// 		panic(err)
// 	}
// }
