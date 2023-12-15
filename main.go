package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zhinea/sps/utils"
	"net/http"
	"os"
	"strings"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(3))

	db = database.GetConnection()
	ctx := context.Background()
	MasterDomain := os.Getenv("MASTER_DOMAIN")

	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")

			host := r.Host

			//check master host in exists in redis
			if strings.HasPrefix(host, MasterDomain) {
				handler.ServeHTTP(w, r)
				return
			}

		})
	})
}
