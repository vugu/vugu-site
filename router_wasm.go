package main

import (
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/vugu/vugu"

	js "syscall/js"
)

// Router - let's try making a router that doesn't make my head explode and my soul weep for longing of simpler times.
// Nope, failed.  Nice idea but there are too many edge cases to make this work.  Can't serve a real 404 page from
// the server without knowing the paths; can't easily do "highlight this link on a sub page without also highlighting
// the one for the parent" (see top links Getting Started and Documentation).  Hacking this in for now but needs more
// thought before we try to make a generic router.
type Router struct {
	u       url.URL // the most recent URL - starts with the one the browser gives us an gets updated from BrowseTo/ReplaceTo
	checkOK bool    // set to true when a path check matches
}

var router = &Router{} // singleton instance - there's only one browser

func init() {
	// TODO: should be able to use this thing in JS or not

	browserURL := js.Global().Get("document").Get("location").String()
	u, err := url.Parse(browserURL)
	if err != nil {
		log.Printf("Error parsing URL from browser(%q): %v", browserURL, err)
		return
	}
	router.u = *u
}

func (r *Router) Path() string {
	return r.u.Path
}

func (r *Router) PathCheckText(p string, txt string) string {
	if r.PathCheck(p) {
		return txt
	}
	return ""
}

func (r *Router) PathCheckText2(p string, txtok string, txtnot string) string {
	if r.PathCheck(p) {
		return txtok
	}
	return txtnot
}

// PathCheck returns true if a the path matches this string exactly after cleaning.
func (r *Router) PathCheck(p string) bool {
	if r.u.Path == path.Clean("/"+p) {
		r.checkOK = true
		return true
	}
	return false
}

func (r *Router) PathPrefixCheckText(p string, txt string) string {
	if r.PathPrefixCheck(p) {
		return txt
	}
	return ""
}

func (r *Router) PathPrefixCheckText2(p string, txtok string, txtnot string) string {
	if r.PathPrefixCheck(p) {
		return txtok
	}
	return txtnot
}

// PathPrefixCheck returns true if a path prefix matches.  /a will match /a, /a/ and /a/b etc.
func (r *Router) PathPrefixCheck(p string) bool {
	p = path.Clean("/" + p)
	if r.u.Path == p {
		r.checkOK = true
		return true
	}
	if strings.HasPrefix(r.u.Path, p+"/") {
		r.checkOK = true
		return true
	}
	return false
}

// PathChecksFailed returns true if no PathCheck or PathPrefixCheck calls have succeeded since the last URL change.
// Used to implement the page not found case.
func (r *Router) PathChecksFailed() bool {
	return !r.checkOK
}

// BrowseTo goes to a page.  It updates the browser's URL using history.pushState and updates the internal path data.
func (r *Router) BrowseTo(p string, event ...*vugu.DOMEvent) error {
	newU, err := url.Parse(p)
	if err != nil {
		return err
	}
	r.u = *newU
	js.Global().Get("window").Get("history").Call("pushState", nil, "", p)
	r.checkOK = false

	for _, e := range event {
		e.PreventDefault()
	}

	return nil
}

// ReplaceTo is like BrowseTo but uses history.replaceState and so does not add an entry in the browser's history.
func (r *Router) ReplaceTo(p string, event ...*vugu.DOMEvent) error {
	newU, err := url.Parse(p)
	if err != nil {
		return err
	}
	r.u = *newU
	js.Global().Get("window").Get("history").Call("replaceState", nil, "", p)
	r.checkOK = false

	for _, e := range event {
		e.PreventDefault()
	}

	return nil
}

// Reload tells the browser to reload the current page - will do a full page refresh.
func (r *Router) Reload() error {
	js.Global().Get("document").Get("location").Call("reload")
	return nil
}

// ReloadTo causes the browser to go to a URL with a full page reload.
func (r *Router) ReloadTo(p string) error {
	newU, err := url.Parse(p)
	if err != nil {
		return err
	}
	r.u = *newU
	js.Global().Get("document").Set("location", p)
	return nil
}
