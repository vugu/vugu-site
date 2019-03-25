// +build ignore

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vugu/vugu"
)

func main() {

	wd, _ := os.Getwd()
	l := ":8877"
	log.Printf("Starting HTTP Server at %q", l)
	h := vugu.NewDevHTTPHandler(wd, http.Dir(wd))

	h.IndexTemplate = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
<link rel="stylesheet" href="/assets/css/vendor.css" />
<link rel="stylesheet" href="/assets/css/style.css" />
<title>Vugu: The Go+WebAssembly UI Framework</title>
<meta charset="utf-8">
<script src="/wasm_exec.js"></script>
</head>
<body>
<div id="root_mount_parent">
<center>...</center>
</div>
<script>
const go = new Go();
WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject).then((result) => {
	go.run(result.instance);
});
</script>
</body>
</html>`

	log.Fatal(http.ListenAndServe(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if r.URL.Path == "/api/sample" {
		// 	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
		// 	items := make([]string, 0, n)
		// 	for i := 0; i < n; i++ {
		// 		items = append(items, fmt.Sprintf("api sample item %d (ts=%v)", i, time.Now()))
		// 	}
		// 	w.Header().Set("Content-Type", "application/json")
		// 	json.NewEncoder(w).Encode(items)
		// 	return
		// }
		h.ServeHTTP(w, r)
	})))

}
