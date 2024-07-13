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
		setHeaders(w)

		log.Println("Access from ", r.URL.Path)

		host := r.Host

		if isSystemPath(r.URL.Path) {
			handler.ServeHTTP(w, r)
			return
		}

		domain, err := getDomain(r.Context(), host)
		if err != nil {
			handleError(w, fmt.Sprintf("Error retrieving domain: %v", err))
			return
		}

		if domain == nil {
			handleError(w, "Domain or host not registered.")
			return
		}

		if domain.IsActive == 0 {
			handleError(w, "Container Paused.")
			return
		}

		ctx := context.WithValue(r.Context(), "domain", *domain)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=604800")
}

func isSystemPath(path string) bool {
	if strings.Contains(path, utils.Cfg.Server.SystemPath) {
		return true
	}

	if strings.Contains(path, "health") {
		return true
	}

	return false
}

func getDomain(ctx context.Context, host string) (*Domain, error) {
	domain, err := fetchDomainFromRedis(ctx, host)
	if err != nil {
		// If there's an error fetching from Redis, log it and continue to DB
		fmt.Printf("Redis error: %v\n", err)
	} else if domain != nil {
		return domain, nil
	}

	// If domain is not in Redis or there was an error, fetch from DB
	domain, err = fetchDomainFromDB(ctx, host)
	if err != nil {
		return nil, err
	}

	if domain != nil {
		go cacheDomainInRedis(ctx, host, *domain)
	}

	return domain, nil
}

func fetchDomainFromRedis(ctx context.Context, host string) (*Domain, error) {
	res, err := database.Redis.Get(ctx, "host:"+host).Result()
	if err != nil {
		return nil, err
	}

	if res == "" {
		return nil, nil
	}

	var domain Domain
	err = json.Unmarshal([]byte(res), &domain)
	if err != nil {
		return nil, err
	}

	return &domain, nil
}

func fetchDomainFromDB(ctx context.Context, host string) (*Domain, error) {
	var domain Domain
	err := database.DB.
		Model(&Domain{}).
		Select("domains.domain, domains.container_id, containers.id, containers.is_active, JSON_UNQUOTE(JSON_EXTRACT(containers.config, '$.options.gtag_id')) as gtag_id").
		Joins("JOIN containers ON domains.container_id = containers.id").
		Where("domains.domain = ?", host).
		First(&domain).
		Error

	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil // Domain not found, but it's not an error
		}
		return nil, err
	}

	return &domain, nil
}

func cacheDomainInRedis(ctx context.Context, host string, domain Domain) {
	domainJSON, err := json.Marshal(domain)
	if err == nil {
		err = database.Redis.Set(ctx, "host:"+host, domainJSON, time.Minute*30).Err()
		if err != nil {
			fmt.Printf("Error caching domain in Redis: %v\n", err)
		}
	} else {
		fmt.Printf("Error marshalling domain: %v\n", err)
	}
}

func handleError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "private")
	w.Header().Set("X-Rellic-Message", message)
	w.Write([]byte(fmt.Sprintf("setTimeout(()=>{console.log('[Rellic] %s');}, 100)", message)))
}
