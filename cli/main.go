package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	pirsch "github.com/pirsch-analytics/pirsch-go-sdk/v2/pkg"
	"log"
	"os"
	"time"
)

type column struct {
	Campaign string `json:"campaign"`
	Device   string `json:"device"`
}

type config struct {
	ClientID     string        `json:"client_id"`
	ClientSecret string        `json:"client_secret"`
	Filter       pirsch.Filter `json:"filter"`
	EventMetaKey string        `json:"event_meta_key"`
	EventPrinted string        `json:"event_printed"`
	EventClicked string        `json:"event_clicked"`
	Columns      []column      `json:"columns"`
}

func main() {
	path := "./config.json"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	file, err := os.ReadFile(path)

	if err != nil {
		log.Println("config.json not found: ", err)
		return
	}

	var cfg config

	if err := json.Unmarshal(file, &cfg); err != nil {
		log.Println("Error loading config.json: ", err)
		return
	}

	client := pirsch.NewClient(cfg.ClientID, cfg.ClientSecret, nil)
	domain, err := client.Domain()

	if err != nil {
		log.Println("Error loading domain: ", err)
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString("PeriodStart,PeriodEnd,CampaignName,Device,AdvertPrinted,AdvertClicked\n")
	cfg.Filter.DomainID = domain.ID

	for _, column := range cfg.Columns {
		cfg.Filter.Event = []string{cfg.EventPrinted}
		cfg.Filter.EventMetaKey = []string{cfg.EventMetaKey}
		cfg.Filter.EventMeta = map[string]string{
			cfg.EventMetaKey: column.Campaign,
		}
		statsPrinted, err := client.EventMetadata(&cfg.Filter)

		if err != nil {
			log.Println("Error getting visitors: ", err)
			return
		}

		if len(statsPrinted) == 0 {
			continue
		}

		cfg.Filter.Event = []string{cfg.EventClicked}
		statsClicked, err := client.EventMetadata(&cfg.Filter)

		if err != nil {
			log.Println("Error getting visitors: ", err)
			return
		}

		if len(statsClicked) == 0 {
			continue
		}

		buffer.WriteString(fmt.Sprintf("%s,%s,%s,%s,%d,%d\n", cfg.Filter.From.Format(time.DateOnly), cfg.Filter.To.Format(time.DateOnly), column.Campaign, column.Device, statsPrinted[0].Visitors, statsClicked[0].Visitors))
	}

	if err := os.WriteFile("report.csv", buffer.Bytes(), 0644); err != nil {
		log.Println("Error writing report.csv: ", err)
	}
}
