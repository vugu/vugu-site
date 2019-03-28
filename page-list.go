package main

import (
	"path"
	"strings"
)

type sitePage struct {
	Path        string
	ShortTitle  string
	LongTitle   string
	Description string
}

type sitePageList []sitePage

// find by path prefix?

func (l sitePageList) FindByPath(p string) *sitePage {
	p = path.Clean("/" + p)
	for i, pg := range l {
		if pg.Path == p {
			return &l[i]
		}
	}
	return nil
}

func (l sitePageList) FindByPathOrEmpty(p string) sitePage {
	ptr := l.FindByPath(p)
	if ptr != nil {
		return *ptr
	}
	return sitePage{}
}

func (l sitePageList) SiteDocPages() []sitePage {
	ret := make([]sitePage, 0, 16)
	for _, pg := range l {
		if strings.HasPrefix(pg.Path, "/doc/") || pg.Path == "/doc" {
			ret = append(ret, pg)
		}
	}
	return ret
}

func (l sitePageList) FindNextDocPage(p string) *sitePage {
	p = path.Clean("/" + p)
	var foundIdx = -1
	for i, pg := range l {
		if pg.Path == p {
			foundIdx = i
			goto found
		}
	}
	return nil
found:
	if foundIdx >= len(l)-1 {
		return nil // no next page possible
	}
	ret := &l[foundIdx+1]
	if !strings.HasPrefix(ret.Path, "/doc") {
		return nil
	}
	return ret
}

// list of site pages, controls static page output and some navigation
var allPages = sitePageList{
	sitePage{
		Path:        "/",
		ShortTitle:  "Vugu Home",
		LongTitle:   "Vugu: A modern UI library for Go+WebAssembly",
		Description: "Pure Go. Targets WebAssembly (and/or server). Most modern browsers supported. Experimental, for now. Really cool.",
	},
	sitePage{
		Path:        "/faq",
		ShortTitle:  "Vugu FAQ",
		LongTitle:   "Vugu Frequently Asked Questions",
		Description: "What browsers are supported? Do it use any JS? Stability and compatibility? Roadmap? How does it build?",
	},
	sitePage{
		Path:        "/doc",
		ShortTitle:  "What is Vugu?",
		LongTitle:   "What is Vugu?",
		Description: "Vugu is Go library which makes it easy to write web user interfaces in Go. You write UI components in .vugu files, converted to .go, compiled to wasm.  Browse DOM sychronization.",
	},
	sitePage{
		Path:        "/doc/start",
		ShortTitle:  "Getting Started",
		LongTitle:   "Getting Started",
		Description: "Let's make a basic working Vugu application that runs in your browser. It only takes three small files to start.",
	},
	sitePage{
		Path:        "/doc/files",
		ShortTitle:  "Vugu Files",
		LongTitle:   "Vugu Files - Overview",
		Description: "Vugu files have three sections: Markup, Style and Code. Markup is the HTML element which is the display portion of your file. Style is a regular style tag. Code is Go language code.",
	},
	sitePage{
		Path:        "/doc/files/markup",
		ShortTitle:  "Markup (HTML/Go)",
		LongTitle:   "Vugu Files - Markup (HTML/Go)",
		Description: "The Markup section is an element which has the HTML that is displayed for this file. In addition to regular HTML, some specific attributes have special meaning in Vugu and allow you to introduce logic into your component's display.",
	},
	sitePage{
		Path:        "/doc/files/style",
		ShortTitle:  "Styles (CSS)",
		LongTitle:   "Vugu Files - Styles (CSS)",
		Description: "Style blocks are simply a way to express CSS that corresponds to your component and is output along with your component markup. Example...",
	},
	sitePage{
		Path:        "/doc/files/code",
		ShortTitle:  "Code (Go)",
		LongTitle:   "Vugu Files - Code (Go)",
		Description: "Go code can be included in your component with in a script tag. This code is copied from your .vugu file into the resulting code generated .go file. It is the appropriate place to include structs, methods and imports needed by your component.",
	},
	sitePage{
		Path:        "/doc/dom-events",
		ShortTitle:  "DOM Events",
		LongTitle:   "DOM Events",
		Description: "Events can be attached to HTML elements by providing an attribute of the event name prefixed with '@'. The name of the event correspoonds to a regular DOM event as would be provided to addEventListener.",
	},
	sitePage{
		Path:        "/doc/components",
		ShortTitle:  "Using Components",
		LongTitle:   "Using Components",
		Description: "Components are individual files which are used to organize your user interface code. Each component lives in a .vugu file. Each .vugu file is processed to produce a .go file.",
	},
	//  sitePage {
	//    Path: "/doc/components-in-depth",
	//    ShortTitle: "Components in Depth",
	//    LongTitle: "Components in Depth",
	//  },
	sitePage{
		Path:        "/doc/program",
		ShortTitle:  "Program Structure",
		LongTitle:   "Vugu Program Structure",
		Description: "WebAssembly main: Like any Go program, it starts with main function. To render a page to HTML, you need a root component - lives in root.vugu by default. Once we have a root component instance, we need an environment.",
	},
	sitePage{
		Path:        "/doc/build-and-dist",
		ShortTitle:  "Building and Distribution",
		LongTitle:   "Building and Distribution",
		Description: "To make a proper server suitable for staging or production, or to start adding more server-side functionality to, you'll want to create a server.go file, and place separate web server code in here.",
	},
}
