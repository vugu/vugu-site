package app

import (
	"github.com/vugu/vgrouter"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu-site/app/components"
	"github.com/vugu/vugu-site/app/pages"
	"github.com/vugu/vugu-site/app/state"
)

//go:generate vugugen -r -skip-go-mod -skip-main ./components
//go:generate vugugen -s -r ./pages
//go:generate vgrgen -r ./pages

// VuguSetup performs UI setup and wiring.
func VuguSetup(buildEnv *vugu.BuildEnv, eventEnv vugu.EventEnv) (*App, vugu.Builder) {

	app := &App{
		Router:   vgrouter.New(eventEnv),
		PageInfo: &state.PageInfo{},
	}

	pageMap := pages.MakeRoutes().WithRecursive(true).WithClean(true).Map()
	pageSeq := &state.PageSeq{
		PageMap:  pageMap,
		PathList: siteNavPathList,
	}
	app.PageSeq = pageSeq

	buildEnv.SetWireFunc(func(b vugu.Builder) {

		if c, ok := b.(vgrouter.NavigatorSetter); ok {
			c.NavigatorSet(app.Router)
		}

		if c, ok := b.(state.PageInfoSetter); ok {
			c.PageInfoSet(app.PageInfo)
		}

		if c, ok := b.(state.PageSeqSetter); ok {
			c.PageSeqSet(app.PageSeq)
		}

	})

	root := &components.Root{}
	buildEnv.WireComponent(root)

	// changes by section
	app.Router.MustAddRoute("/doc", vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		root.Hero = &components.NavHero{}
		root.Sidebar = &components.DocSidebar{}
	}))

	// pages - add automatically from generated routes
	for path, inst := range pageMap {
		instBuilder := inst.(vugu.Builder)
		app.Router.MustAddRouteExact(path, vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
			root.Body = instBuilder
		}))
	}

	app.Router.MustAddRouteExact("/", vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		root.Body = nil
		root.FullBody = &pages.Index{}
	}))

	app.Router.SetNotFound(vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		root.Body = &components.NotFound{}
	}))

	// add another route at the end that always runs and handles the page info
	app.Router.MustAddRoute("/", vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		*app.PageInfo = state.PageInfoFrom(rm.Path, root.Body) // overwrite PageInfo
	}))

	if app.Router.BrowserAvail() {
		err := app.Router.ListenForPopState()
		if err != nil {
			panic(err)
		}

		err = app.Router.Pull()
		if err != nil {
			panic(err)
		}
	}

	return app, root
}

// App holds overall application state.
type App struct {
	*vgrouter.Router
	*state.PageInfo
	*state.PageSeq
}

// sequence of the previous and next links and the doc sidebar uses this sequence for the things under /doc
var siteNavPathList = []string{
	"/doc",
	"/doc/start",
	"/doc/files",
	"/doc/files/markup",
	"/doc/files/style",
	"/doc/files/code",
	"/doc/dom-events",
	"/doc/components",
	"/doc/program",
	"/doc/build-and-dist",
}
