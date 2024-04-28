package utils

import "net/http"

func ServeHtml(r *http.Request, path string) string {
	p := "." + r.URL.Path

	if p == "./" || p == "./"+path {
		p = "./web/html/" + path + ".html"
	}

	return p
}
