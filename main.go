package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const apiBase = "https://api.online.net/api/v1"

var httpClient = &http.Client{Timeout: 10 * time.Second}

// --- helpers ---------------------------------------------------------

func bearer() string {
	tok := os.Getenv("ONLINE_API_TOKEN")
	if tok == "" {
		log.Fatal("ONLINE_API_TOKEN not set")
	}
	return tok
}

// flatten any JSON object up to depth 3
func flatten(prefix string, v interface{}, depth int, out map[string]string) {
	if depth > 3 || v == nil {
		return
	}

	switch vv := v.(type) {

	case map[string]interface{}:
		for k, val := range vv {
			flatten(prefix+"_"+k, val, depth+1, out)
		}

	case []interface{}:
		for i, val := range vv {
			flatten(prefix+"_"+strconv.Itoa(i), val, depth+1, out)
		}

	case string, float64, bool:
		out[prefix] = fmt.Sprint(vv)
	}
}

func apiGET(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiBase+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+bearer())
	req.Header.Set("X-Pretty-JSON", "1")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API %s: %s", path, resp.Status)
	}
	return body, nil
}

// --- HTTP handler ----------------------------------------------------

type target struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func handleSD(w http.ResponseWriter, _ *http.Request) {
	raw, err := apiGET("/server")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	var paths []string
	if err := json.Unmarshal(raw, &paths); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	var results []target
	for _, p := range paths {
		sid := strings.TrimPrefix(p, "/api/v1/server/")
		detailRaw, err := apiGET("/server/" + sid)
		if err != nil {
			log.Printf("warn: %v", err)
			continue
		}

		var detail map[string]interface{}
		if err := json.Unmarshal(detailRaw, &detail); err != nil {
			continue
		}

		// first IP
		net := detail["network"].(map[string]interface{})
		ipArr := net["ip"].([]interface{})
		if len(ipArr) == 0 {
			continue
		}

		labels := map[string]string{"server_id": sid}
		flatten("meta", detail, 0, labels)

		results = append(results, target{
			Targets: []string{fmt.Sprint(ipArr[0]) + ":9100"},
			Labels:  labels,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}

// --- main ------------------------------------------------------------

func main() {
	http.HandleFunc("/scw-sd", handleSD)
	log.Println("Scaleway SD listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
