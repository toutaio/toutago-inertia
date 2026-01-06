package handlers

import (
	"net/http"

	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

type AdminDashboardProps struct {
	Stats struct {
		Users int `json:"users"`
		Todos int `json:"todos"`
	} `json:"stats"`
}

func AdminDashboard(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := inertia.NewContext(w, r, i)

		props := AdminDashboardProps{}
		props.Stats.Users = 42  // Mock data
		props.Stats.Todos = 128 // Mock data

		ctx.Render("admin/Dashboard", props)
	}
}
