package gtagcontroller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	utils2 "github.com/gofiber/fiber/v2/utils"
	"github.com/patrickmn/go-cache"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/handler"
	"github.com/zhinea/sps/utils"
	"io"
	"log"
	"net/http"
	urlUtils "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
	pathProject, _ = os.Getwd()
	injectFilePath = pathProject + "/storage/injects/gtag.inject.js"
	injectRegex    = regexp.MustCompile(`var\s+data\s+=\s+{`)
	regexDomain    = regexp.MustCompile(`PROXY_DOMAIN`)
	scriptCache    *cache.Cache
)

type RequestLog struct {
	ContainerID int       `json:"container_id"`
	IPAddr      string    `json:"ip_addr"`
	Domain      string    `json:"domain"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

func init() {
	scriptCache = cache.New(12*time.Hour, 3*time.Hour)
}

func GetScripts(c *fiber.Ctx) error {
	start := time.Now()
	c.Set("Content-Type", "application/javascript")

	config := c.Context().Value("domain").(handler.Domain)
	ipAddr := c.IP()

	// Menggunakan Go routines untuk menyimpan data request log ke database
	go logRequest(config, ipAddr, "scripts")

	scriptName := fmt.Sprintf("gtag-%s.js", utils.MD5Hash([]byte(config.Domain)))

	if scriptData, found := scriptCache.Get(scriptName); found {
		log.Printf("Serving cached script %s, time took: %s\n", scriptName, time.Since(start))
		return c.Send(scriptData.([]byte))
	}

	// If not in cache, compile the script
	scriptData, err := compileScript(c, config, scriptName)
	if err != nil {
		return err
	}

	// Store the compiled script in cache
	scriptCache.Set(scriptName, scriptData, cache.DefaultExpiration)

	log.Printf("Compiled and cached script %s, time took: %v\n", scriptName, time.Since(start))
	return c.Send(scriptData)

}

func compileScript(c *fiber.Ctx, config handler.Domain, scriptName string) ([]byte, error) {
	path := filepath.Join(pathProject, "storage/compiled", scriptName)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		injectFile, err := os.ReadFile(injectFilePath)

		if err != nil {
			log.Println("Error reading gtag inject file", err)
			return nil, err
		}

		requestURL := "https://www.googletagmanager.com/gtag/js?id=" + config.GtagID

		headers := c.GetReqHeaders()
		req, err := createRequest(&createRequestConfig{
			method:    c.Method(),
			url:       requestURL,
			userAgent: strings.Join(headers["User-Agent"], " "),
			referrer:  strings.Join(headers["Referer"], " "),
			ipAddress: c.IP(),
		})
		if err != nil {
			log.Println("Error creating request:", err)
			return nil, errors.New("error while pre-getting scripts")
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error sending request to Google Analytics:", err)
			return nil, errors.New("error while getting scripts")
		}
		defer resp.Body.Close()

		file, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return nil, err
		}

		// parsing subdomain like https://domain.com/[random string]/ga
		subDomain := fmt.Sprintf("%s://%s/%s/ga", c.Protocol(), config.Domain, utils.RandString(10))

		editedBody := injectRegex.ReplaceAllString(string(file), string(injectFile)+`var data={`)
		editedBody = regexDomain.ReplaceAllString(editedBody, subDomain)

		script := []byte(editedBody)

		go func() {
			if err := os.WriteFile(path, script, 0644); err != nil {
				log.Println("Error writing compiled script:", err)
			}
		}()

		return script, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error while reading compiled scripts", err)
		return nil, err
	}

	return file, nil
}

func HandleTrackData(c *fiber.Ctx) error {
	payload := utils2.CopyString(c.Query("cached"))
	config := c.Context().Value("domain").(handler.Domain)

	if payload == "" {
		return c.Status(http.StatusUnprocessableEntity).SendString("err: missing paramater")
	}

	url, err := base64.StdEncoding.DecodeString(payload)

	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("err: err when decoding payload")
	}

	_, errParsing := urlUtils.ParseRequestURI(string(url))
	if errParsing != nil {
		log.Println(errParsing)
		return errParsing
	}

	ipAddress := c.IP()
	requestURL := fmt.Sprintf("%s&uip=%s&_uip=%s", url, ipAddress, ipAddress)

	headers := c.GetReqHeaders()
	req, err := createRequest(&createRequestConfig{
		method:    c.Method(),
		url:       requestURL,
		body:      nil,
		userAgent: strings.Join(headers["User-Agent"], " "),
		referrer:  strings.Join(headers["Referer"], " "),
		ipAddress: ipAddress,
	})

	if err != nil {
		log.Println("Error creating request:", err)
		return c.SendString("err: parsing request")
	}

	// Menggunakan Go routines untuk melakukan pengiriman request ke Google Analytics
	go sendTrackRequest(req, config, ipAddress)

	return c.SendStatus(200)
}

func sendTrackRequest(req *http.Request, config handler.Domain, ipAddr string) {
	defer utils.Recover()

	resp, err := client.Do(req)
	if err != nil {
		log.Println("err: can't send to google analytics", err)
		return
	}
	defer resp.Body.Close()

	logRequest(config, ipAddr, "track")
}

func logRequest(config handler.Domain, ipAddr, reqType string) {
	defer utils.Recover()

	if err := database.DB.Create(&RequestLog{
		ContainerID: config.ContainerID,
		IPAddr:      ipAddr,
		Type:        reqType,
		Domain:      config.Domain,
		CreatedAt:   time.Now(),
	}).Error; err != nil {
		log.Println("err: can't insert request log to database", err)
	}
}

type createRequestConfig struct {
	method    string
	url       string
	body      io.Reader
	userAgent string
	referrer  string
	ipAddress string
}

func createRequest(config *createRequestConfig) (*http.Request, error) {
	req, err := http.NewRequest(config.method, config.url, config.body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.userAgent)
	req.Header.Set("Referer", config.referrer)
	req.Header.Set("X-Forwarded-For", config.referrer)
	req.Header.Set("X-Real-IP", config.ipAddress)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	return req, err
}
