package state

// PageInfoRef is a reference to PageInfo that can be embedded in components that need it.
type PageInfoRef struct{ *PageInfo }

// PageInfoSet implements PageInfoSetter.
func (r *PageInfoRef) PageInfoSet(v *PageInfo) { r.PageInfo = v }

// PageInfoSetter is implemented by PageInfoRef to follow the reference setter pattern.
type PageInfoSetter interface{ PageInfoSet(v *PageInfo) }

// PageInfo holds general page information.
type PageInfo struct {
	Path string // URL path

	Title           string // title tag
	MetaDescription string // meta description tag

	LongTitle  string // longer title for h1 on page
	ShortTitle string // abbreviated title for side bar
}

// PageInfoFrom creates a PageInfo from an object.
func PageInfoFrom(p string, i interface{}) PageInfo {
	var ret PageInfo

	ret.Path = p

	//log.Printf("PageInfoFrom called on: %#v", i)

	if titler, ok := i.(titler); ok {
		ret.Title = titler.Title()
	}
	if metaDescriptioner, ok := i.(metaDescriptioner); ok {
		ret.MetaDescription = metaDescriptioner.MetaDescription()
	}
	if longTitler, ok := i.(longTitler); ok {
		ret.LongTitle = longTitler.LongTitle()
	}
	if shortTitler, ok := i.(shortTitler); ok {
		ret.ShortTitle = shortTitler.ShortTitle()
	}

	if ret.Title == "" {
		ret.Title = "Vugu"
	}

	if ret.LongTitle == "" {
		ret.LongTitle = ret.Title
	}
	if ret.ShortTitle == "" {
		ret.ShortTitle = ret.Title
	}
	if ret.MetaDescription == "" {
		ret.MetaDescription = "Vugu: Pure Go. Targets WebAssembly (and/or server). Most modern browsers supported. Experimental, for now. Really cool."
	}

	return ret
}

type titler interface{ Title() string }
type metaDescriptioner interface{ MetaDescription() string }
type longTitler interface{ LongTitle() string }
type shortTitler interface{ ShortTitle() string }
