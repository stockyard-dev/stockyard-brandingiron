package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stockyard-dev/stockyard-brandingiron/internal/store"
)

type Server struct {
	db     *store.DB
	mux    *http.ServeMux
	port   int
	limits Limits
}

func New(db *store.DB, port int, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), port: port, limits: limits}
	s.mux.HandleFunc("GET /api/og", s.hGenerateOG)
	s.mux.HandleFunc("POST /api/og", s.hGenerateOGPost)

	s.mux.HandleFunc("POST /api/templates", s.hCreateTemplate)
	s.mux.HandleFunc("GET /api/templates", s.hListTemplates)
	s.mux.HandleFunc("GET /api/templates/{name}", s.hGetTemplate)
	s.mux.HandleFunc("DELETE /api/templates/{id}", s.hDelTemplate)

	s.mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) })
	s.mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]string{"status": "ok"}) })
	s.mux.HandleFunc("GET /ui", s.handleUI)
	s.mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{"product": "stockyard-brandingiron", "version": "0.1.0"})
	})
	return s
}

func (s *Server) Start() error {
	log.Printf("[brandingiron] :%d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

// GET /api/og?title=Hello&subtitle=World&template=default&bg=1a1410&fg=f0e6d3&accent=e8753a
func (s *Server) hGenerateOG(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	title := q.Get("title")
	subtitle := q.Get("subtitle")
	templateName := q.Get("template")
	if templateName == "" {
		templateName = "default"
	}

	tpl, _ := s.db.GetTemplate(templateName)
	if tpl == nil {
		tpl = &store.Template{Width: 1200, Height: 630, BgColor: "#1a1410", TextColor: "#f0e6d3", AccentColor: "#e8753a", FontSize: 48}
	}

	// Allow query param overrides
	if bg := q.Get("bg"); bg != "" {
		tpl.BgColor = "#" + bg
	}
	if fg := q.Get("fg"); fg != "" {
		tpl.TextColor = "#" + fg
	}
	if accent := q.Get("accent"); accent != "" {
		tpl.AccentColor = "#" + accent
	}

	svg := generateSVG(title, subtitle, tpl, !s.limits.RemoveWatermark)
	s.db.RecordGeneration(templateName, title, subtitle)

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write([]byte(svg))
}

// POST /api/og with JSON body
func (s *Server) hGenerateOGPost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Subtitle    string `json:"subtitle"`
		Template    string `json:"template"`
		BgColor     string `json:"bg_color"`
		TextColor   string `json:"text_color"`
		AccentColor string `json:"accent_color"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	templateName := req.Template
	if templateName == "" {
		templateName = "default"
	}

	tpl, _ := s.db.GetTemplate(templateName)
	if tpl == nil {
		tpl = &store.Template{Width: 1200, Height: 630, BgColor: "#1a1410", TextColor: "#f0e6d3", AccentColor: "#e8753a", FontSize: 48}
	}

	if req.BgColor != "" {
		tpl.BgColor = req.BgColor
	}
	if req.TextColor != "" {
		tpl.TextColor = req.TextColor
	}
	if req.AccentColor != "" {
		tpl.AccentColor = req.AccentColor
	}

	svg := generateSVG(req.Title, req.Subtitle, tpl, !s.limits.RemoveWatermark)
	s.db.RecordGeneration(templateName, req.Title, req.Subtitle)

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write([]byte(svg))
}

func generateSVG(title, subtitle string, tpl *store.Template, watermark bool) string {
	w := tpl.Width
	h := tpl.Height
	bg := tpl.BgColor
	fg := tpl.TextColor
	accent := tpl.AccentColor
	fontSize := tpl.FontSize

	// Wrap title if too long
	titleLines := wrapText(title, 30)
	if len(titleLines) == 0 {
		titleLines = []string{""}
	}

	titleY := h/2 - (len(titleLines)*fontSize/2 + 10)
	if subtitle != "" {
		titleY -= 20
	}

	var titleSVG strings.Builder
	for i, line := range titleLines {
		y := titleY + i*(fontSize+8)
		titleSVG.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-family="system-ui, -apple-system, sans-serif" font-size="%d" font-weight="700" fill="%s" text-anchor="middle">%s</text>`,
			w/2, y, fontSize, svgEsc(fg), svgEsc(line)))
	}

	subtitleSVG := ""
	if subtitle != "" {
		subY := titleY + len(titleLines)*(fontSize+8) + 16
		subtitleSVG = fmt.Sprintf(`<text x="%d" y="%d" font-family="system-ui, -apple-system, sans-serif" font-size="%d" fill="%s" text-anchor="middle" opacity="0.7">%s</text>`,
			w/2, subY, fontSize/2, svgEsc(fg), svgEsc(subtitle))
	}

	watermarkSVG := ""
	if watermark {
		watermarkSVG = fmt.Sprintf(`<text x="%d" y="%d" font-family="monospace" font-size="11" fill="%s" text-anchor="end" opacity="0.3">stockyard.dev</text>`,
			w-20, h-16, svgEsc(fg))
	}

	return fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">
<rect width="%d" height="%d" fill="%s"/>
<rect x="0" y="0" width="%d" height="4" fill="%s"/>
<rect x="0" y="%d" width="%d" height="4" fill="%s"/>
<line x1="60" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1" opacity="0.15"/>
<line x1="60" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1" opacity="0.15"/>
%s
%s
%s
</svg>`,
		w, h, w, h,
		w, h, svgEsc(bg),
		w, svgEsc(accent),
		h-4, w, svgEsc(accent),
		h/2-80, w-60, h/2-80, svgEsc(accent),
		h/2+60, w-60, h/2+60, svgEsc(accent),
		titleSVG.String(),
		subtitleSVG,
		watermarkSVG)
}

func wrapText(text string, maxChars int) []string {
	if len(text) <= maxChars {
		return []string{text}
	}
	words := strings.Fields(text)
	var lines []string
	current := ""
	for _, word := range words {
		if current == "" {
			current = word
		} else if len(current)+1+len(word) <= maxChars {
			current += " " + word
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func svgEsc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// Template CRUD
func (s *Server) hCreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		BgColor     string `json:"bg_color"`
		TextColor   string `json:"text_color"`
		AccentColor string `json:"accent_color"`
		FontSize    int    `json:"font_size"`
		Layout      string `json:"layout"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Name == "" {
		wj(w, 400, map[string]string{"error": "name required"})
		return
	}
	t, err := s.db.CreateTemplate(req.Name, req.Width, req.Height, req.BgColor, req.TextColor, req.AccentColor, req.FontSize, req.Layout)
	if err != nil {
		wj(w, 500, map[string]string{"error": err.Error()})
		return
	}
	wj(w, 201, map[string]any{"template": t})
}

func (s *Server) hListTemplates(w http.ResponseWriter, r *http.Request) {
	ts, _ := s.db.ListTemplates()
	if ts == nil {
		ts = []store.Template{}
	}
	wj(w, 200, map[string]any{"templates": ts, "count": len(ts)})
}

func (s *Server) hGetTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := s.db.GetTemplate(r.PathValue("name"))
	if err != nil {
		wj(w, 404, map[string]string{"error": "template not found"})
		return
	}
	wj(w, 200, map[string]any{"template": t})
}

func (s *Server) hDelTemplate(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteTemplate(r.PathValue("id"))
	wj(w, 200, map[string]string{"status": "deleted"})
}

func wj(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
