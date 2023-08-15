package bosun

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
	"os"
)

// Component Generation

type ComponentGenerator interface {
	Build(w http.ResponseWriter, r *http.Request) (Component, error)
}

type Component struct {
	Path string
	Data interface{}
}

func (c *Component) Render() string {
	templateContent, err := os.ReadFile(c.Path)
	if err != nil {
		fmt.Println(err.Error())
	}
	return mustache.Render(string(templateContent), c.Data)
}

//Component Configuration

type ComponentConfig struct {
	COMPONENTS map[string]ComponentGenerator
}

func (c *ComponentConfig) RegisterComponent(tag string, generator ComponentGenerator) {
	c.COMPONENTS[tag] = generator
}

func NewComponentConfig() ComponentConfig {
	return ComponentConfig{
		COMPONENTS: make(map[string]ComponentGenerator),
	}
}

// Opinionated component server

func (c *ComponentConfig) ServeComponents(w http.ResponseWriter, r *http.Request) {
	componentID := r.URL.Query().Get("id")
	component, err := c.COMPONENTS[componentID].Build(w, r)
	if err != nil {
		return
	}
	text := component.Render()
	w.Write([]byte(text))
}
