package search

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/blevesearch/bleve"
)

func TestSearch1(t *testing.T) {

	tmpDir, err := ioutil.TempDir("", "TestSearch1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(filepath.Join(tmpDir, "example.bleve"), mapping)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := struct {
		Name string
	}{
		Name: "text",
	}

	// index some data
	index.Index("id", data)

	// search for some text
	query := bleve.NewMatchQuery("text")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", searchResults)

}
