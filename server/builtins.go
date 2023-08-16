package server

import (
	"fmt"
	"github.com/colemalphrus/bosun/render"
	"net/http"
)

func (c *Config) pagesSetup() {
	fmt.Println("Setting Up Pages")
	if c.PageDir == "" {
		c.PageDir = "./pages"
	}
	if c.PageRoot == "" {
		c.PageRoot = "/"
	}

	pages := render.ExtractPages(c.PageDir)
	c.Multiplexer.HandleFunc(c.PageRoot, render.ServePages(pages))
}

func (c *Config) staticSetup() {
	fmt.Println("Setting Up Static Pages")
	if c.StaticDir == "" {
		c.StaticDir = "./static"
	}
	if c.StaticRoot == "" {
		c.StaticRoot = "/static/"
	}

	//c.Multiplexer.Handle(c.StaticRoot, render.ServeStatic(c.StaticRoot, c.StaticDir))
	c.Multiplexer.Handle("/static/", render.ServeStatic("/static/", "./static"))
}

func (c *Config) livenessSetup() {
	//c.Multiplexer.Handle(c.StaticRoot, render.ServeStatic(c.StaticRoot, c.StaticDir))
	c.Multiplexer.HandleFunc("/liveness", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("200::OK"))
	})
}
