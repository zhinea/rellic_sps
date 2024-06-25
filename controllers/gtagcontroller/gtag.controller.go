package gtagcontroller

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/handler"
	"github.com/zhinea/sps/utils"
	"log"
	"net/http"
	"os"
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
	ContainerID int       `json:"container_id"`
	IPAddr      string    `json:"ip_addr"`
	Domain      string    `json:"domain"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetScripts(c *fiber.Ctx) error {
	ctx := c.Context()
	config := ctx.Value("domain").(handler.Domain)
	c.Set("Content-Type", "application/javascript")

	filePath, _ := os.Getwd()
	filePath = filePath + "/storage/gtag.min.js"
	ipAddr := c.IP()

	// Menggunakan Go routines untuk menyimpan data request log ke database
	go func() {
		defer utils.Recover()

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

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)

		return err
	}

	subDomain := "http://" + config.Domain + "/_gg"

	re := regexp.MustCompile(`PROXY_DOMAIN`)

	editedBody := re.ReplaceAllString(string(file), subDomain)

	return c.Send([]byte(editedBody))
}

func HandleTrackData(c *fiber.Ctx) error {
	payload := c.Query("cache")
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
		defer utils.Recover()

		resp, err := client.Do(req)
		if err != nil {
			log.Println("err: can't send to google analytics", err)
		}
		defer resp.Body.Close()

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
