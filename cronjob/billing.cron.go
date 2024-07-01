package cronjob

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/zhinea/sps/controllers/gtagcontroller"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: 60 * time.Second,
	}
)

type ContainerLogs struct {
	ContainerID int16
	Total       int
}

type RequestPayload struct {
	Checksum string
	Logs     []ContainerLogs
}

func BillingSchedule() {
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	var results []ContainerLogs
	database.DB.Model(&gtagcontroller.RequestLog{}).
		Where("created_at >= ?", fiveMinutesAgo).
		Select("container_id, COUNT(*) as total").
		Group("container_id").
		Scan(&results)

	var sumTotal int

	for _, result := range results {
		sumTotal += result.Total

		fmt.Printf("ContainerID: %s, Total Logs: %d\n", result.ContainerID, result.Total)
	}

	checksumHash := md5.Sum([]byte(strconv.Itoa(sumTotal) + ".PayloadRune:key"))

	checksum := hex.EncodeToString(checksumHash[:])

	sendRequest(RequestPayload{
		Checksum: checksum,
		Logs:     results,
	})
}

func sendRequest(payload RequestPayload) {
	log.Println(payload.Checksum)

	url := utils.Cfg.Container.ServerUrl + "/api/v1/poolback-estung-tung/" + utils.Cfg.Container.ServerID + "/billings"

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	r, errR := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println(errR)
		return
	}

	r.Header.Add("Content-Type", "application/json")

	res, err2 := httpClient.Do(r)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	defer res.Body.Close()

	postResult := &struct {
		status int
	}{}

	derr := json.NewDecoder(res.Body).Decode(postResult)
	if derr != nil {
		fmt.Println(derr)
		return
	}

	if res.StatusCode != http.StatusCreated {
		fmt.Println(strconv.Itoa(res.StatusCode) + " Res Billing poolback")
		return
	}

	log.Println("Billing poolback success")
}
