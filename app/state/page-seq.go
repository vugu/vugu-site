package state

import (
	"strings"
)

// PageSeqRef is a reference to PageSeq that can be embedded in components that need it.
type PageSeqRef struct{ *PageSeq }

// PageSeqSet implements PageSeqSetter.
func (r *PageSeqRef) PageSeqSet(v *PageSeq) { r.PageSeq = v }

// PageSeqSetter is implemented by PageSeqRef to follow the reference setter pattern.
type PageSeqSetter interface{ PageSeqSet(v *PageSeq) }

// PageSeq keeps track of the sequence of pages so we can implement previous and next links.
type PageSeq struct {
	PathList []string
	PageMap  map[string]interface{} // so we can get the actual title data (using PageInfoFrom)
}

// Next returns the next link given a path.  Empty string return means either no match found or p is first.
func (ps *PageSeq) Next(p string) string {
	for i, pth := range ps.PathList {
		if pth == p {
			if i == len(ps.PathList)-1 {
				return ""
			}
			return ps.PathList[i+1]
		}
	}
	return ""
}

// Prev returns the previous link given a path.  Empty string return means either no match found or p is last.
func (ps *PageSeq) Prev(p string) string {
	for i, pth := range ps.PathList {
		if pth == p {
			if i == 0 {
				return ""
			}
			return ps.PathList[i-1]
		}
	}
	return ""
}

// WithPrefix returns a new PageSeq filtered to the paths matching the specified prefix.
func (ps *PageSeq) WithPrefix(pfx string) *PageSeq {
	ret := &PageSeq{PageMap: ps.PageMap}
	for _, p := range ps.PathList {
		if p == pfx || strings.HasPrefix(p, pfx+"/") {
			ret.PathList = append(ret.PathList, p)
		}
	}
	return ret
}
