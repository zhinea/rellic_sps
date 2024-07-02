package handler

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

type Domain struct {
	ContainerID int    `json:"container_id"`
	Domain      string `json:"domain"`
	GtagID      string `json:"gtag_id"`
	ID          int    `json:"id"`
	IsActive    int    `json:"is_active"`
}

func AppMiddleware(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "public, max-age=604800")

		ctx := context.Background()
		host := r.Host

		if strings.Contains(r.URL.Path, utils.Cfg.Server.SystemPath) {
			log.Println("Access route system detected")
			handler.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.URL.Path, "health") {
			log.Println("Access route system detected")
			handler.ServeHTTP(w, r)
			return
		}

		//check host in exists in redis
		//if strings.Contains(utils.Cfg.Server.Domain, host) {
		//	log.Println("Master domain detected", utils.Cfg.Server.Domain, host)
		//	handler.ServeHTTP(w, r)
		//	return
		//}

		// Check if host exists in cache
		res, _ := database.Redis.Get(ctx, "host:"+host).Result()

		var domain Domain

		if res != "" {
			err := json.Unmarshal([]byte(res), &domain)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
			}

			log.Println("Domain found in cache:", domain.Domain)

		} else {
			// Retrieve data from the database if not found in cache
			err := database.DB.
				Model(&Domain{}).
				//Table("domains").
				//Select("domain, container_id, server_id").
				Select("domains.domain, domains.container_id, containers.id, containers.is_active, JSON_UNQUOTE(JSON_EXTRACT(containers.config, '$.options.gtag_id')) as gtag_id").
				Joins("JOIN containers ON domains.container_id = containers.id").
				Where("domains.domain = ?", host).
				First(&domain).
				Error

			if err != nil {
				w.Header().Set("Content-Type", "application/javascript")
				w.Write([]byte("setTimeout(()=>{console.log('Domain or host not registered on rellic.app. Please insert your custom domain first, and make it to primary domain!');alert('rellic actication error');}, 500)"))
				return
			}

			if domain.IsActive == 0 {
				w.Header().Set("Content-Type", "application/javascript")
				w.Write([]byte("setTimeout(()=>{console.log('The owner this domain has paused the container.');}, 100)"))
				return
			}

			// Store data in cache
			go func() {
				defer utils.Recover()

				domainJSON, err := json.Marshal(domain)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
				} else {
					err := database.Redis.Set(ctx, "host:"+host, domainJSON, time.Minute*30).Err()
					if err != nil {
						fmt.Println("Error storing data in cache:", err)
					}
				}
			}()
		}

		ctxs := context.WithValue(r.Context(), "domain", domain)

		log.Println("ACCEPTED Request from:", host, "~> [CID", domain.ContainerID, "] [", domain.GtagID, "]")

		handler.ServeHTTP(w, r.WithContext(ctxs))
	})
}
