package gtagcontroller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	utils2 "github.com/gofiber/fiber/v2/utils"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/handler"
	"github.com/zhinea/sps/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
	analyticURL     = "https://google-analytics.com/g/collect"
	pathProject, _  = os.Getwd()
	gtagFilePath    = pathProject + "/storage/gtag.min.js"
	regexDomain     = regexp.MustCompile(`PROXY_DOMAIN`)
	regexGTAGConfig = regexp.MustCompile(`G-7NJG7X7KP1`)
	fileMutex       sync.Mutex // Mutex to handle file read/write race conditions
)

type RequestLog struct {
	ContainerID int       `json:"container_id"`
	IPAddr      string    `json:"ip_addr"`
	Domain      string    `json:"domain"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetScripts(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/javascript")

	ctx := c.Context()
	ipAddr := c.IP()

	config := ctx.Value("domain").(handler.Domain)

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

	scriptName := fmt.Sprintf("compiled-%s.js", utils.MD5Hash([]byte(config.Domain)))
	path := filepath.Join(pathProject, "storage/compiled", scriptName)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		log.Println("Creating compiled file", scriptName)

		file, err := os.ReadFile(gtagFilePath)
		if err != nil {
			log.Println(err)

			return err
		}

		subDomain := "http://" + config.Domain + "/_gg"

		editedBody := regexDomain.ReplaceAllString(string(file), subDomain)
		editedBody = regexGTAGConfig.ReplaceAllString(editedBody, subDomain)

		go func() {
			defer utils.Recover()

			errWrite := os.WriteFile(path, []byte(editedBody), 0644)
			if errWrite != nil {
				log.Println("Error writing", errWrite.Error())
			}
		}()

		return c.Send([]byte(editedBody))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)

		return err
	}
	log.Println("Found compiled file", scriptName)

	return c.Send(file)
}

func HandleTrackData(c *fiber.Ctx) error {
	payload := utils2.CopyString(c.Query("cache"))
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
			return
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
