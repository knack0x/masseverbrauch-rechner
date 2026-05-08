package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") || r.URL.Path == "/manifest.json" {
			w.Header().Set("Cache-Control", "max-age=3600, must-revalidate")
		}
		next.ServeHTTP(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := IndexData{
		WebVersion: webVersion,
		Slots:      tmplSlots(5),
	}
	var buf bytes.Buffer
	if err := tmpls.ExecuteTemplate(&buf, "index.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	flow, _ := strconv.ParseFloat(r.FormValue("flow"), 64)
	runtimeMinutes, _ := strconv.ParseFloat(r.FormValue("runtime_minutes"), 64)

	var slots []Slot
	for i := 0; i < 5; i++ {
		before, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("slots[%d][before]", i)), 64)
		after, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("slots[%d][after]", i)), 64)
		slots = append(slots, Slot{Before: before, After: after})
	}

	req := CalculateRequest{
		Flow:           flow,
		RuntimeMinutes: runtimeMinutes,
		Slots:          slots,
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080/api/calculate"
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("API call failed: %v", err)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<div style="color: red; padding: 1rem;">Fehler bei der Berechnung: API nicht erreichbar</div>`)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		statusText := http.StatusText(resp.StatusCode)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<div style="color: red; padding: 1rem;">Fehler bei der Berechnung: %d %s</div>`, resp.StatusCode, statusText)
		return
	}

	var apiResult CalculateResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		log.Printf("API response decode failed: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	view := CalculateResultView{
		HauptmasseKG:      tmplNumberFormat(apiResult.HauptmasseKG, 2),
		HauptmassePercent: tmplNumberFormat(apiResult.HauptmassePercent, 2),
		TotalKG:           tmplNumberFormat(apiResult.TotalKG, 2),
	}
	for _, s := range apiResult.Slots {
		view.Slots = append(view.Slots, SlotResultView{
			Name:    tmplTrSlotName(s.Name),
			KG:      tmplNumberFormat(s.KG, 2),
			Percent: tmplNumberFormat(s.Percent, 2),
		})
	}

	var buf bytes.Buffer
	if err := tmpls.ExecuteTemplate(&buf, "calculate.html", view); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	baseURL := apiBaseURL()
	resp, err := http.Get(baseURL + "/api/version")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"version": "n/a"})
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
