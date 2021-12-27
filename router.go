package simplerouter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Param struct {
	Key string
	Val string
	Idx int
}

type Params []Param

var AllParams map[string]Params

type Router struct {
	routes map[string]*Route
}

type Route struct {
	method string
	path   string
	f      http.HandlerFunc
}

func New() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}

func parseParams(path string) Params {
	var params Params
	splited := strings.Split(path[1:], "/")
	// fmt.Println(splited)
	for idx, key := range splited {
		if strings.HasPrefix(key, ":") {
			params = append(params, Param{Key: key[1:], Idx: idx})
		}
	}
	return params
}

func getParamsWithValues(p Params, reqPath string) Params {
	var params Params
	splited := strings.Split(reqPath[1:], "/")
	for i := 0; i < len(p); i++ {
		p[i].Val = splited[p[i].Idx]
		params = append(params, p[i])
	}
	return params
}

func (r *Router) GET(path string, f http.HandlerFunc) {
	r.routes[path] = &Route{
		method: "GET",
		path:   path,
		f:      f,
	}
}

func (r *Router) POST(path string, f http.HandlerFunc) {
	r.routes[path] = &Route{
		method: "POST",
		path:   path,
		f:      f,
	}
}

func (r *Router) PUT(path string, f http.HandlerFunc) {
	r.routes[path] = &Route{
		method: "PUT",
		path:   path,
		f:      f,
	}
}

func (r *Router) PATCH(path string, f http.HandlerFunc) {
	r.routes[path] = &Route{
		method: "PATCH",
		path:   path,
		f:      f,
	}
}

func (r *Router) DELETE(path string, f http.HandlerFunc) {
	r.routes[path] = &Route{
		method: "DELETE",
		path:   path,
		f:      f,
	}
}

func GetParam(r *http.Request, name string) string {
	params := r.Context().Value(ContextKey("params")).(Params)
	for _, param := range params {
		if param.Key == name {
			return param.Val
		}
	}
	return ""
}

type ContextKey string

// Parse 2 paths one is the path with params and the other is the path without params
func comparePathsWithParams(path1, path2 string) bool {
	if path1 == path2 {
		return true
	}
	// Remove the first / and the last / if they exist
	path1 = strings.Trim(path1, "/")
	path2 = strings.Trim(path2, "/")
	// Split the paths
	splited1 := strings.Split(path1, "/")
	splited2 := strings.Split(path2, "/")
	fmt.Println(splited1, splited2)
	if len(splited1) != len(splited2) {
		return false
	}
	for i := 0; i < len(splited1); i++ {
		fmt.Println(splited1[i], splited2[i])
		if !strings.HasPrefix(splited1[i], ":") && splited1[i] != splited2[i] {
			return false
		}
	}
	return true
}

func (r *Router) getRoute(req *http.Request) *Route {
	reqPath := req.URL.Path
	for path, route := range r.routes {
		if comparePathsWithParams(path, reqPath) {
			return route
		}
	}
	return nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route := r.getRoute(req)
	fmt.Println(route)
	if route == nil {
		http.NotFound(w, req)
		return
	}
	if route.method != req.Method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parsedParams := parseParams(route.path)
	params := getParamsWithValues(parsedParams, req.URL.Path)
	req = req.WithContext(context.WithValue(req.Context(), ContextKey("params"), params))
	// fmt.Println(req.Context().Value(ContextKey("params")).(Params))
	route.f(w, req)
}
