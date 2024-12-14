package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	pirsch "github.com/pirsch-analytics/pirsch-go-sdk/v2/pkg"
)

var client *pirsch.Client

func main() {
	basePath := ""

	if os.Getenv("BASE_PATH") != "" {
		basePath = os.Getenv("BASE_PATH")
	}

	static, err := template.ParseGlob(filepath.Join(basePath, "static/**.html"))

	if err != nil {
		slog.Error("Error loading templates", "err", err)
		return
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	client = pirsch.NewClient(clientID, clientSecret, nil)
	domain, err := client.Domain()

	if err != nil {
		slog.Error("Error loading domain", "err", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var from, to time.Time
		var pattern string
		var pages []pirsch.PageStats

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				log.Println("Error parsing form: ", err)
				w.WriteHeader(http.StatusBadRequest)
				errorResponse(w, err)
				return
			}

			from, err = time.Parse("2006-01-02", r.Form.Get("from"))

			if err != nil {
				log.Println("Error parsing date: ", err)
				w.WriteHeader(http.StatusBadRequest)
				errorResponse(w, err)
				return
			}

			to, err = time.Parse("2006-01-02", r.Form.Get("to"))

			if err != nil {
				log.Println("Error parsing date: ", err)
				w.WriteHeader(http.StatusBadRequest)
				errorResponse(w, err)
				return
			}

			pattern = r.Form.Get("pattern")
			filterPattern := strings.ReplaceAll(pattern, "/", "\\/")
			filterPattern = strings.ReplaceAll(filterPattern, "*", ".*")
			pages, err = client.Pages(&pirsch.Filter{
				DomainID: domain.ID,
				From:     from,
				To:       to,
				Pattern:  []string{filterPattern},
			})

			if err != nil {
				log.Println("Error loading pages: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				errorResponse(w, err)
				return
			}

			for i := range pages {
				pages[i].BounceRate = math.Round(pages[i].BounceRate*10000) / 100
				pages[i].RelativeVisitors = math.Round(pages[i].RelativeVisitors*10000) / 100
				pages[i].RelativeViews = math.Round(pages[i].RelativeViews*10000) / 100
			}
		}

		if err := static.ExecuteTemplate(w, "index.html", struct {
			From    string
			To      string
			Pattern string
			Pages   []pirsch.PageStats
		}{
			getDate(from),
			getDate(to),
			pattern,
			pages,
		}); err != nil {
			log.Println("Error rendering page: ", err)
		}
	})

	slog.Info("Starting server on localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Error starting server", "err", err)
	}
}

func errorResponse(w http.ResponseWriter, err error) {
	if _, e := w.Write([]byte(fmt.Sprintf("An error occurred: %s", err))); e != nil {
		slog.Error("Error writing response", "err", e)
	}
}

func getDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02")
}
