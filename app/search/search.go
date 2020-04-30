package search

import "fmt"

type Searcher struct {
	// TBD
}

type ResultSet struct {
	ResultList []Result
	// TODO: when we get into paging we'll need to see what is available and add it here
}

type Result struct {
	URL     string // the link URL (for now this is just a path like "/doc/start" but we may need make it a full URL later)
	Title   string // title of search result (html)
	Snippet string // the short snippet displayed directly under the title (html)
	Detail  string // the longer detailed data for when the user selects "more" (html)
}

func New(filePath string) (*Searcher, error) {
	panic(fmt.Errorf("not yet implemented"))
}

func (s *Searcher) Search(query string) (interface{}, error) {
	panic(fmt.Errorf("not yet implemented"))
}
