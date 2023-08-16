package render

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ServeHTML(prefix string, dir string) http.Handler {
	h := http.FileServer(http.Dir(dir))
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)

		if p == "" {
			http.ServeFile(w, r, "./pages/index.html")
		} else if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p + ".html"
			r2.URL.RawPath = rp
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})
}

func ServeStatic(prefix string, dir string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
}

// custom page server

func ServePages(p map[string]Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

type Page struct {
	Path string
	Type string
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
