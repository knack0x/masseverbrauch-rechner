package main

// IndexData is passed to the index.html template to render the main page.
type IndexData struct {
	WebVersion string
	Slots      []int // 1..5, rendered as "Turmposition 1" .. "Turmposition 5"
}

// CalculateResultView is passed to the calculate.html template to render the
// result dialog after a successful API response.
type CalculateResultView struct {
	HauptmasseKG      string
	HauptmassePercent string
	Slots             []SlotResultView
	TotalKG           string
}

// SlotResultView represents one slot row in the result dialog.
type SlotResultView struct {
	Name    string // e.g. "Turmposition 1"
	KG      string // German-formatted number
	Percent string // German-formatted number
}

// CalculateRequest is the JSON body sent to the POST /api/calculate endpoint.
type CalculateRequest struct {
	Flow           float64 `json:"flow"`
	RuntimeMinutes float64 `json:"runtime_minutes"`
	Slots          []Slot  `json:"slots"`
}

// Slot represents one tower-position measurement sent to the API.
type Slot struct {
	Before float64 `json:"before"`
	After  float64 `json:"after"`
}

// CalculateResponse is the JSON body returned by the POST /api/calculate endpoint.
type CalculateResponse struct {
	HauptmasseKG      float64      `json:"hauptmasse_kg"`
	HauptmassePercent float64      `json:"hauptmasse_percent"`
	Slots             []SlotResult `json:"slots"`
	TotalKG           float64      `json:"total_kg"`
}

// SlotResult is one slot entry in the API response.
type SlotResult struct {
	Name    string  `json:"name"`
	KG      float64 `json:"kg"`
	Percent float64 `json:"percent"`
}
