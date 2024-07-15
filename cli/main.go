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
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	EventMetaKey string   `json:"event_meta_key"`
	EventPrinted string   `json:"event_printed"`
	EventClicked string   `json:"event_clicked"`
	Columns      []column `json:"columns"`
}

func main() {
	file, err := os.ReadFile("config.json")

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

	start, err := time.Parse(time.DateOnly, cfg.StartDate)

	if err != nil {
		log.Println("Error parsing start date: ", err)
		return
	}

	end, err := time.Parse(time.DateOnly, cfg.EndDate)

	if err != nil {
		log.Println("Error parsing end date: ", err)
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString("PeriodStart,PeriodEnd,CampaignName,Device,AdvertPrinted,AdvertClicked\n")

	for _, column := range cfg.Columns {
		statsPrinted, err := client.TotalVisitors(&pirsch.Filter{
			DomainID: domain.ID,
			From:     start,
			To:       end,
			Event:    []string{cfg.EventPrinted},
			EventMeta: map[string]string{
				cfg.EventMetaKey: column.Campaign,
			},
		})

		if err != nil {
			log.Println("Error getting visitors: ", err)
			return
		}

		statsClicked, err := client.TotalVisitors(&pirsch.Filter{
			DomainID: domain.ID,
			From:     start,
			To:       end,
			Event:    []string{cfg.EventClicked},
			EventMeta: map[string]string{
				cfg.EventMetaKey: column.Campaign,
			},
		})

		if err != nil {
			log.Println("Error getting visitors: ", err)
			return
		}

		buffer.WriteString(fmt.Sprintf("%s,%s,%s,%s,%d,%d\n", start.Format(time.DateOnly), end.Format(time.DateOnly), column.Campaign, column.Device, statsPrinted.Visitors, statsClicked.Visitors))
	}

	if err := os.WriteFile("report.csv", buffer.Bytes(), 0644); err != nil {
		log.Println("Error writing report.csv: ", err)
	}
}
