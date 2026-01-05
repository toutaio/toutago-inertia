package main

import (
	"encoding/json"
	"net/http"

	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

func main() {
	config := inertia.Config{
		RootView: "templates/app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	if err != nil {
		panic(err)
	}

	mgr.Share("appName", "Basic Example")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, _ := mgr.Render("Home", map[string]interface{}{
			"greeting": "Welcome!",
		}, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	})

	handler := mgr.Middleware()(mux)
	http.ListenAndServe(":3000", handler)
}
