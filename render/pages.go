package render

import (
	"fmt"
	"github.com/colemalphrus/bosun/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Path string
	Type string
}

func ServePages(p map[string]Page) mux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ctx mux.Context) {
		urlPath := r.URL.Path
		if len(r.URL.Path) > 2 {
			urlPath = strings.TrimSuffix(urlPath, "/")
		}
		path, exists := p[urlPath]
		if !exists {
			w.WriteHeader(404)
			w.Write([]byte("404 page not found"))
			return
		}
		http.ServeFile(w, r, path.Path)
	}
}

func ExtractPages(dirname string) map[string]Page {
	return cleanPages(DigPages(dirname), dirname)
}

func DigPages(dirname string) map[string]Page {

	pages := make(map[string]Page)

	// Open the directory
	dir, err := os.Open(dirname)
	if err != nil {
		log.Fatalf("Failed opening directory: %s", err)
	}
	//defer dir.Close()

	// Read all files from the directory
	list, _ := dir.Readdir(0)

	// Loop through the files and print them
	for _, file := range list {
		path := dirname + "/" + file.Name()

		if !file.IsDir() {
			fileType := strings.Split(file.Name(), ".")
			pages[path] = Page{
				Path: path,
				Type: fileType[len(fileType)-1],
			}
		} else {
			m := DigPages(path)
			pages = mergeMaps(pages, m)
		}
	}

	return pages
}

func cleanPages(p map[string]Page, prefix string) map[string]Page {
	cp := make(map[string]Page)
	for k, v := range p {
		k = strings.TrimPrefix(k, prefix)
		k = strings.TrimSuffix(k, ".html")
		k = strings.Replace(k, "index", "", -1)
		cp[k] = v
		fmt.Println(k)
		//fmt.Println(v.Path)
	}

	return cp
}

func mergeMaps(m1, m2 map[string]Page) map[string]Page {
	merged := make(map[string]Page)

	// Copy m1 to merged
	for k, v := range m1 {
		merged[k] = v
	}

	// Copy m2 to merged, potentially overwriting values from m1
	for k, v := range m2 {
		merged[k] = v
	}

	return merged
}
