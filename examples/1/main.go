package main

import (
	"net/http"

	"github.com/gopherzz/simplerouter"
)

func main() {
	r := simplerouter.New()
	r.GET("/hello/:name/second/:second", func(rw http.ResponseWriter, r *http.Request) {
		first := simplerouter.GetParam(r, "name")
		second := simplerouter.GetParam(r, "second")
		rw.Write([]byte("Hello " + first + " " + second))
	})
	http.ListenAndServe(":80", r)
}
