package gtagcontroller

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/handler"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
	analyticURL = "https://google-analytics.com/g/collect"
)

type RequestLog struct {
	ContainerID string    `json:"container_id"`
	IPAddr      string    `json:"ip_addr"`
	Domain      string    `json:"domain"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetScripts(c *fiber.Ctx) error {
	config := c.Context().Value("domain").(handler.Domain)
	ctx := c.Context()
	c.Set("Content-Type", "application/javascript")

	gtagURL := "https://www.googletagmanager.com/gtag/js?id=" + config.GtagID
	cacheKey := "gtag:" + config.ID
	ipAddr := c.IP()

	// Menggunakan Go routines untuk menyimpan data request log ke database
	go func() {
		request := RequestLog{
			ContainerID: config.ContainerID,
			IPAddr:      ipAddr,
			Type:        "scripts",
			Domain:      config.Domain,
			CreatedAt:   time.Now(),
		}

		if err := database.DB.Create(&request).Error; err != nil {
			log.Println("err: can't insert request log to database", err)
		}
	}()

	cachedScripts, err := database.Redis.Get(ctx, cacheKey).Result()

	if err == nil && cachedScripts != "" {
		return c.SendString(cachedScripts)
	}

	resp, errClient := client.Get(gtagURL)

	if errClient != nil {
		log.Fatal(errClient)
		return c.SendString("console.log('Error loading gtag.js');")
	}

	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)

	if errRead != nil {
		log.Fatal(errRead)
		return c.SendString("console.log('Error loading gtag.js');")
	}

	subDomain := "http://" + config.Domain

	re := regexp.MustCompile(`b=(\w+)\.sendBeacon&&(\w+)\.sendBeacon\(a\)`)

	editedBody := strings.Replace(string(body), `https://"+a+".google-analytics.com/g/collect`, subDomain+"/proxy", -1)
	editedBody = re.ReplaceAllString(editedBody, `let Ev=new URL(a),et=btoa(Ev.search);a=new URL("?p="+et+"&d="+Ev.searchParams.get("tid"),Ev.origin+Ev.pathname).href; b=false;`)

	errSet := database.Redis.Set(ctx, cacheKey, editedBody, time.Hour*24*7).Err()

	if errSet != nil {
		log.Fatal(errSet)
	}

	return c.Send([]byte(editedBody))
}

func HandleTrackData(c *fiber.Ctx) error {
	payload := c.Query("p")
	config := c.Context().Value("domain").(handler.Domain)

	if payload == "" {
		c.Status(http.StatusUnprocessableEntity)
		return c.SendString("err: missing paramater")
	}

	params, err := base64.StdEncoding.DecodeString(payload)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.SendString("err: err when decoding payload")
	}

	ipAddress := c.IP()
	requestURL := analyticURL + string(params) + "&uip=" + ipAddress + "&_uip=" + ipAddress

	req, errPreReq := http.NewRequest(c.Method(), requestURL, nil)

	if errPreReq != nil {
		log.Fatal(errPreReq)
		return c.SendString("err: parsing request")
	}

	headers := c.GetReqHeaders()
	userAgent := headers["User-Agent"]

	req.Header.Set("User-Agent", strings.Join(userAgent, " "))
	req.Header.Set("Referer", strings.Join(headers["Referer"], " "))
	req.Header.Set("X-Forwarded-For", ipAddress)
	req.Header.Set("X-Real-IP", ipAddress)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	// Menggunakan Go routines untuk melakukan pengiriman request ke Google Analytics
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()

		resp, err := client.Do(req)
		if err != nil {
			log.Println("err: can't send to google analytics", err)
		}

		defer resp.Body.Close()
	}()

	// Menggunakan Go routines untuk menyimpan data request log ke database
	go func() {
		request := RequestLog{
			ContainerID: config.ContainerID,
			IPAddr:      ipAddress,
			Type:        "track",
			Domain:      config.Domain,
			CreatedAt:   time.Now(),
		}

		if err := database.DB.Create(&request).Error; err != nil {
			log.Println("err: can't insert request log to database", err)
		}
	}()

	return c.SendStatus(200)
}
