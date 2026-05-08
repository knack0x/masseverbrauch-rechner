package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//go:embed assets
var assets embed.FS

//go:embed templates
var templateFiles embed.FS

var webVersion string

func init() {
	if data, err := os.ReadFile("VERSION"); err == nil {
		webVersion = strings.TrimSpace(string(data))
	} else {
		webVersion = "dev"
	}
}

func apiBaseURL() string {
	u := os.Getenv("API_URL")
	if u == "" {
		u = "http://localhost:8080/api/calculate"
	}
	return strings.TrimSuffix(u, "/api/calculate")
}

type IndexData struct {
	WebVersion string
	Slots      []int
}

type CalculateResultView struct {
	HauptmasseKG     string
	HauptmassePercent string
	Slots            []SlotResultView
	TotalKG          string
}

type SlotResultView struct {
	Name    string
	KG      string
	Percent string
}

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
	HauptmasseKG      float64      `json:"hauptmasse_kg"`
	HauptmassePercent float64      `json:"hauptmasse_percent"`
	Slots             []SlotResult `json:"slots"`
	TotalKG           float64      `json:"total_kg"`
}

type SlotResult struct {
	Name    string  `json:"name"`
	KG      float64 `json:"kg"`
	Percent float64 `json:"percent"`
}

func tmplSlots(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i + 1
	}
	return s
}

func tmplSub(a, b int) int {
	return a - b
}

func tmplNumberFormat(v float64, decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	s := fmt.Sprintf(format, v)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	var result []byte
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	if len(parts) > 1 {
		result = append(result, ',')
		result = append(result, parts[1]...)
	}
	return string(result)
}

func tmplTrSlotName(name string) string {
	tr := map[string]string{
		"Tower Slot 1": "Turmposition 1",
		"Tower Slot 2": "Turmposition 2",
		"Tower Slot 3": "Turmposition 3",
		"Tower Slot 4": "Turmposition 4",
		"Tower Slot 5": "Turmposition 5",
	}
	if v, ok := tr[name]; ok {
		return v
	}
	return name
}

var tmpls *template.Template

func initTemplates() error {
	funcMap := template.FuncMap{
		"slots":        tmplSlots,
		"sub":          tmplSub,
		"numberFormat": tmplNumberFormat,
		"trSlotName":   tmplTrSlotName,
	}
	var err error
	tmpls, err = template.New("").Funcs(funcMap).ParseFS(templateFiles, "templates/*.html")
	return err
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
		beforeStr := r.FormValue(fmt.Sprintf("slots[%d][before]", i))
		afterStr := r.FormValue(fmt.Sprintf("slots[%d][after]", i))
		if beforeStr == "" && afterStr == "" {
			continue
		}
		before, _ := strconv.ParseFloat(beforeStr, 64)
		after, _ := strconv.ParseFloat(afterStr, 64)
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

func Serve() {
	if err := initTemplates(); err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	assetsFS, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatalf("Failed to get assets sub FS: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))))
	mux.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := assets.ReadFile("assets/manifest.json")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/calculate", calculateHandler)
	mux.HandleFunc("/api/version", versionHandler)

	handler := cacheMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Web server starting on :%s (version: %s)", port, webVersion)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
