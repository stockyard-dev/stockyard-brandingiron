package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"
	_ "modernc.org/sqlite"
)

type DB struct{ conn *sql.DB }

func Open(dataDir string) (*DB, error) {
	os.MkdirAll(dataDir, 0755)
	conn, err := sql.Open("sqlite", filepath.Join(dataDir, "brandingiron.db"))
	if err != nil { return nil, err }
	conn.Exec("PRAGMA journal_mode=WAL"); conn.Exec("PRAGMA busy_timeout=5000"); conn.SetMaxOpenConns(4)
	db := &DB{conn: conn}; return db, db.migrate()
}
func (db *DB) Close() error { return db.conn.Close() }

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
CREATE TABLE IF NOT EXISTS templates (
    id TEXT PRIMARY KEY, name TEXT NOT NULL UNIQUE, width INTEGER DEFAULT 1200,
    height INTEGER DEFAULT 630, bg_color TEXT DEFAULT '#1a1410',
    text_color TEXT DEFAULT '#f0e6d3', accent_color TEXT DEFAULT '#e8753a',
    font_size INTEGER DEFAULT 48, layout TEXT DEFAULT 'centered',
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS generations (
    id TEXT PRIMARY KEY, template_id TEXT DEFAULT '', title TEXT DEFAULT '',
    subtitle TEXT DEFAULT '', bg_color TEXT DEFAULT '', text_color TEXT DEFAULT '',
    accent_color TEXT DEFAULT '', width INTEGER DEFAULT 1200, height INTEGER DEFAULT 630,
    created_at TEXT DEFAULT (datetime('now'))
);
`)
	// Default template
	db.conn.Exec(`INSERT OR IGNORE INTO templates (id,name) VALUES ('default','default')`)
	return err
}

type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	BgColor     string `json:"bg_color"`
	TextColor   string `json:"text_color"`
	AccentColor string `json:"accent_color"`
	FontSize    int    `json:"font_size"`
	Layout      string `json:"layout"`
	CreatedAt   string `json:"created_at"`
}

func (db *DB) CreateTemplate(name string, width, height int, bgColor, textColor, accentColor string, fontSize int, layout string) (*Template, error) {
	id := "tpl_" + gid(6); now := time.Now().UTC().Format(time.RFC3339)
	if width <= 0 { width = 1200 }; if height <= 0 { height = 630 }
	if bgColor == "" { bgColor = "#1a1410" }; if textColor == "" { textColor = "#f0e6d3" }
	if accentColor == "" { accentColor = "#e8753a" }; if fontSize <= 0 { fontSize = 48 }
	if layout == "" { layout = "centered" }
	_, err := db.conn.Exec("INSERT INTO templates (id,name,width,height,bg_color,text_color,accent_color,font_size,layout,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)",
		id, name, width, height, bgColor, textColor, accentColor, fontSize, layout, now)
	if err != nil { return nil, err }
	return &Template{ID: id, Name: name, Width: width, Height: height, BgColor: bgColor, TextColor: textColor, AccentColor: accentColor, FontSize: fontSize, Layout: layout, CreatedAt: now}, nil
}

func (db *DB) ListTemplates() ([]Template, error) {
	rows, err := db.conn.Query("SELECT id,name,width,height,bg_color,text_color,accent_color,font_size,layout,created_at FROM templates ORDER BY name")
	if err != nil { return nil, err }; defer rows.Close()
	var out []Template
	for rows.Next() { var t Template; rows.Scan(&t.ID, &t.Name, &t.Width, &t.Height, &t.BgColor, &t.TextColor, &t.AccentColor, &t.FontSize, &t.Layout, &t.CreatedAt); out = append(out, t) }
	return out, rows.Err()
}

func (db *DB) GetTemplate(name string) (*Template, error) {
	var t Template
	err := db.conn.QueryRow("SELECT id,name,width,height,bg_color,text_color,accent_color,font_size,layout,created_at FROM templates WHERE name=? OR id=?", name, name).
		Scan(&t.ID, &t.Name, &t.Width, &t.Height, &t.BgColor, &t.TextColor, &t.AccentColor, &t.FontSize, &t.Layout, &t.CreatedAt)
	return &t, err
}

func (db *DB) DeleteTemplate(id string) { db.conn.Exec("DELETE FROM templates WHERE id=? AND name != 'default'", id) }

func (db *DB) RecordGeneration(templateID, title, subtitle string) {
	id := "gen_" + gid(8)
	db.conn.Exec("INSERT INTO generations (id,template_id,title,subtitle) VALUES (?,?,?,?)", id, templateID, title, subtitle)
}

func (db *DB) Stats() map[string]any {
	var templates, generations int
	db.conn.QueryRow("SELECT COUNT(*) FROM templates").Scan(&templates)
	db.conn.QueryRow("SELECT COUNT(*) FROM generations").Scan(&generations)
	return map[string]any{"templates": templates, "generations": generations}
}

func gid(n int) string { b := make([]byte, n); rand.Read(b); return hex.EncodeToString(b) }
