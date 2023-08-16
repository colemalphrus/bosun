package server

import (
	"github.com/colemalphrus/bosun/mux"
	"net/http"
)

type Config struct {
	SkipPages    bool
	SkipStatic   bool
	SkipLiveness bool
	SkipAPI      bool
	PageDir      string
	PageRoot     string
	StaticDir    string
	StaticRoot   string
	Multiplexer  *mux.RouteMux
	SubPackages  []SubPackage
}

type SubPackage interface {
	Register(routeMux *mux.RouteMux)
}

func (c *Config) Initialize() {
	c.Multiplexer = mux.New()
}

func (c *Config) Serve(port string) {
	if !c.SkipStatic {
		c.staticSetup()
	}

	if !c.SkipLiveness {
		c.livenessSetup()
	}

	if !c.SkipPages {
		c.pagesSetup()
	}
	http.ListenAndServe(port, c.Multiplexer)
}

func (c *Config) RegisterSubPackage(sub SubPackage) {
	sub.Register(c.Multiplexer)
}
