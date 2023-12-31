package mux

import (
	"net/http"
	"regexp"
	"strings"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, Context)

//type Context map[string]string

type Context struct {
	PathParams map[string]string
	MWData     map[string]string
	MWErrors   []error
}

type route struct {
	methods []string
	regex   *regexp.Regexp
	params  map[int]string
	handler HandlerFunc
}

type RouteMux struct {
	routes     []*route
	middleware []HandlerFunc
}

func New() *RouteMux {
	return &RouteMux{}
}

func (m *RouteMux) HandleFunc(pattern string, handler HandlerFunc) *route {
	return m.AddRoute(pattern, handler)
}

func (m *RouteMux) Handle(pattern string, handler http.Handler) *route {
	return m.AddRoute(pattern, func(w http.ResponseWriter, r *http.Request, m Context) {
		handler.ServeHTTP(w, r)
	})
}

func (r *route) Methods(methods ...string) {
	r.methods = methods
}

func (m *RouteMux) Middleware(filter HandlerFunc) {
	m.middleware = append(m.middleware, filter)
}

func (m *RouteMux) AddRoute(pattern string, handler HandlerFunc) *route {

	//split the url into sections
	parts := strings.Split(pattern, "/")

	//find params that start with ":"
	//replace with regular expressions
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			params[j] = part[1:]
			parts[i] = "([^/]+)"
			j++
		}
	}

	//recreate the url pattern, with parameters replaced
	//by regular expressions. then compile the regex
	pattern = strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
	}

	//now create the Route
	route := &route{}
	route.regex = regex
	route.handler = handler
	route.params = params

	//and finally append to the list of Routes
	m.routes = append(m.routes, route)

	return route
}

func (m *RouteMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestPath := r.URL.Path
	context := Context{
		PathParams: make(map[string]string),
		MWData:     make(map[string]string),
		MWErrors:   nil,
	}

	//find a matching Route
	for _, route := range m.routes {

		if !validateMethod(route.methods, r.Method) {
			continue
		}

		if !route.regex.MatchString(requestPath) {
			continue
		}

		//get path params
		matches := route.regex.FindStringSubmatch(requestPath)

		if len(route.params) > 0 {
			for i, match := range matches[1:] {
				context.PathParams[route.params[i]] = match
			}
		}

		//execute middleware
		for _, filter := range m.middleware {
			filter(w, r, context)
		}

		//Invoke the request handler
		route.handler(w, r, context)
		break
	}
}

func validateMethod(s []string, str string) bool {
	if len(s) == 0 {
		return true
	}
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
