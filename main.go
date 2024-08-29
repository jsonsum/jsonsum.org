package main

import (
	"jsonsum.org/jsonsum.org/public"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	githubAPIToken := os.Getenv("GITHUB_API_TOKEN")
	if githubAPIToken == "" {
		panic("GITHUB_API_TOKEN must be set")
	}

	// Main site
	siteFS := http.FileServer(http.FS(public.FS))
	http.Handle("/", siteFS)

	// Maven repo
	http.HandleFunc("/org/jsonsum/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		proxy := httputil.ReverseProxy{Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(&url.URL{
				Scheme: "https",
				Host:   "maven.pkg.github.com",
				Path:   "/jsonsum/jsonsum-kotlin/", // this is the base path only
			})
			pr.Out.SetBasicAuth("", githubAPIToken)
		}}
		proxy.ServeHTTP(w, r)
	})

	// Go magic
	http.HandleFunc("/jsonsum-go", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/jsonsum-go.html"
		siteFS.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
