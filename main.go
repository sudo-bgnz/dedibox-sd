// scw_sd.go
//
// Tiny HTTP-SD endpoint that converts Online/Scaleway-Dedibox servers
// into Prometheus scrape targets.
//
//	export ONLINE_API_TOKEN=xxxxxxxx
//	go run scw_sd.go
//
// Prometheus http_sd_configs URL →  http://<host>:8000/scw-sd
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const apiBase = "https://api.online.net/api/v1"

var httpClient = &http.Client{Timeout: 10 * time.Second}

type serverDetail struct {
	Network struct {
		IP []string `json:"ip"`
	} `json:"network"`
	Tags []string `json:"tags"`
}

type target struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func bearer() string {
	tok := os.Getenv("ONLINE_API_TOKEN")
	if tok == "" {
		log.Fatal("ONLINE_API_TOKEN not set")
	}
	return tok
}

func apiGET(path string, v any) error {
	req, err := http.NewRequest("GET", apiBase+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+bearer())
	req.Header.Set("X-Pretty-JSON", "1")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API %s: %s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func handleSD(w http.ResponseWriter, _ *http.Request) {
	var paths []string
	if err := apiGET("/server", &paths); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	var result []target

	for _, p := range paths {
		sid := strings.TrimPrefix(p, "/api/v1/server/")

		// build correct detail path without duplicating /api/v1
		detailPath := "/server/" + sid

		var detail serverDetail
		if err := apiGET(detailPath, &detail); err != nil {
			log.Printf("warn: %v", err)
			continue
		}
		if len(detail.Network.IP) == 0 {
			continue
		}
		labels := map[string]string{"server_id": sid}
		for _, t := range detail.Tags {
			labels["tag_"+t] = "1"
		}
		result = append(result, target{
			Targets: []string{detail.Network.IP[0] + ":9100"},
			Labels:  labels,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/scw-sd", handleSD)
	addr := ":8000"
	log.Printf("Scaleway SD listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
