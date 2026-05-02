package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CalculateRequest struct {
	Flow           float64 `json:"flow"`
	RuntimeMinutes float64 `json:"runtime_minutes"`
	Slots          []Slot  `json:"slots"`
}

type Slot struct {
	Before float64 `json:"before"`
	After  float64 `json:"after"`
}

type CalculateResponse struct {
	HauptmasseKG     float64      `json:"hauptmasse_kg"`
	HauptmassePercent float64     `json:"hauptmasse_percent"`
	Slots            []SlotResult `json:"slots"`
	TotalKG          float64      `json:"total_kg"`
}

type SlotResult struct {
	Name    string  `json:"name"`
	KG      float64 `json:"kg"`
	Percent float64 `json:"percent"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("Request data: flow=%.2f, runtime=%.2fmin, slots=%d", req.Flow, req.RuntimeMinutes, len(req.Slots))

	hauptmasseKG := req.Flow * (req.RuntimeMinutes * 60)
	totalKG := hauptmasseKG

	var results []SlotResult
	for i, slot := range req.Slots {
		kg := slot.After - slot.Before
		if kg < 0 {
			kg = 0
		}
		totalKG += kg
		results = append(results, SlotResult{
			Name: "Tower Slot " + strconv.Itoa(i+1),
			KG:   kg,
		})
	}

	hauptmassePercent := (hauptmasseKG / totalKG) * 100
	for i := range results {
		results[i].Percent = (results[i].KG / totalKG) * 100
	}

	resp := CalculateResponse{
		HauptmasseKG:     hauptmasseKG,
		HauptmassePercent: hauptmassePercent,
		Slots:            results,
		TotalKG:          totalKG,
	}

	log.Printf("Calculation complete: total=%.2fkg, hauptmasse=%.2fkg (%.1f%%)", totalKG, hauptmasseKG, hauptmassePercent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/calculate", calculateHandler)
	log.Println("API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
