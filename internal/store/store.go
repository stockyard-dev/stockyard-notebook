package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

type Notebook struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Color     string `json:"color,omitempty"`
	CreatedAt string `json:"created_at"`
	NoteCount int    `json:"note_count"`
}

type Note struct {
	ID         string   `json:"id"`
	NotebookID string   `json:"notebook_id"`
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	Tags       []string `json:"tags"`
	Pinned     bool     `json:"pinned"`
	Archived   bool     `json:"archived"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	WordCount  int      `json:"word_count"`
}

type NoteFilter struct {
	NotebookID string
	Tag        string
	Search     string
	Pinned     string // "true", "false", ""
	Archived   string // "true", "all", or default (only non-archived)
	SortBy     string // created, updated, title
	SortDir    string
	Limit      int
	Offset     int
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	dsn := filepath.Join(dataDir, "notebook.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	for _, q := range []string{
		`CREATE TABLE IF NOT EXISTS notebooks (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT UNIQUE NOT NULL,
			color TEXT DEFAULT '#c45d2c',
			created_at TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			notebook_id TEXT DEFAULT '',
			title TEXT NOT NULL,
			body TEXT DEFAULT '',
			tags_json TEXT DEFAULT '[]',
			pinned INTEGER DEFAULT 0,
			archived INTEGER DEFAULT 0,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_notebook ON notes(notebook_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_pinned ON notes(pinned)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_updated ON notes(updated_at)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_archived ON notes(archived)`,
		`CREATE TABLE IF NOT EXISTS extras (
			resource TEXT NOT NULL,
			record_id TEXT NOT NULL,
			data TEXT NOT NULL DEFAULT '{}',
			PRIMARY KEY(resource, record_id)
		)`,
	} {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func wordCount(s string) int {
	return len(strings.Fields(s))
}

// ─── Notebooks ────────────────────────────────────────────────────

func (d *DB) CreateNotebook(nb *Notebook) error {
	nb.ID = genID()
	nb.CreatedAt = now()
	if nb.Slug == "" {
		nb.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(nb.Name), " ", "-"))
	}
	if nb.Color == "" {
		nb.Color = "#c45d2c"
	}
	_, err := d.db.Exec(
		`INSERT INTO notebooks (id, name, slug, color, created_at) VALUES (?, ?, ?, ?, ?)`,
		nb.ID, nb.Name, nb.Slug, nb.Color, nb.CreatedAt,
	)
	return err
}

func (d *DB) GetNotebook(id string) *Notebook {
	var nb Notebook
	err := d.db.QueryRow(
		`SELECT id, name, slug, color, created_at FROM notebooks WHERE id=?`, id,
	).Scan(&nb.ID, &nb.Name, &nb.Slug, &nb.Color, &nb.CreatedAt)
	if err != nil {
		return nil
	}
	d.db.QueryRow(`SELECT COUNT(*) FROM notes WHERE notebook_id=? AND archived=0`, id).Scan(&nb.NoteCount)
	return &nb
}

func (d *DB) ListNotebooks() []Notebook {
	rows, err := d.db.Query(`SELECT id, name, slug, color, created_at FROM notebooks ORDER BY name ASC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Notebook
	for rows.Next() {
		var nb Notebook
		if err := rows.Scan(&nb.ID, &nb.Name, &nb.Slug, &nb.Color, &nb.CreatedAt); err != nil {
			continue
		}
		d.db.QueryRow(`SELECT COUNT(*) FROM notes WHERE notebook_id=? AND archived=0`, nb.ID).Scan(&nb.NoteCount)
		out = append(out, nb)
	}
	return out
}

func (d *DB) UpdateNotebook(id string, nb *Notebook) error {
	_, err := d.db.Exec(
		`UPDATE notebooks SET name=?, slug=?, color=? WHERE id=?`,
		nb.Name, nb.Slug, nb.Color, id,
	)
	return err
}

// DeleteNotebook removes the notebook and orphans its notes (sets their
// notebook_id to ”). Note extras for the orphaned notes are kept.
// Notebook's own extras must be cleaned up by the caller.
func (d *DB) DeleteNotebook(id string) error {
	d.db.Exec(`UPDATE notes SET notebook_id='' WHERE notebook_id=?`, id)
	_, err := d.db.Exec(`DELETE FROM notebooks WHERE id=?`, id)
	return err
}

// ─── Notes ────────────────────────────────────────────────────────

func (d *DB) CreateNote(n *Note) error {
	n.ID = genID()
	n.CreatedAt = now()
	n.UpdatedAt = n.CreatedAt
	if n.Tags == nil {
		n.Tags = []string{}
	}
	n.WordCount = wordCount(n.Body)
	tj, _ := json.Marshal(n.Tags)
	pinned, archived := 0, 0
	if n.Pinned {
		pinned = 1
	}
	if n.Archived {
		archived = 1
	}
	_, err := d.db.Exec(
		`INSERT INTO notes (id, notebook_id, title, body, tags_json, pinned, archived, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.ID, n.NotebookID, n.Title, n.Body, string(tj), pinned, archived, n.CreatedAt, n.UpdatedAt,
	)
	return err
}

func (d *DB) scanNote(s interface{ Scan(...any) error }) *Note {
	var n Note
	var tj string
	var pinned, archived int
	if err := s.Scan(&n.ID, &n.NotebookID, &n.Title, &n.Body, &tj, &pinned, &archived, &n.CreatedAt, &n.UpdatedAt); err != nil {
		return nil
	}
	json.Unmarshal([]byte(tj), &n.Tags)
	if n.Tags == nil {
		n.Tags = []string{}
	}
	n.Pinned = pinned == 1
	n.Archived = archived == 1
	n.WordCount = wordCount(n.Body)
	return &n
}

const noteCols = `id, notebook_id, title, body, tags_json, pinned, archived, created_at, updated_at`

func (d *DB) GetNote(id string) *Note {
	return d.scanNote(d.db.QueryRow(`SELECT `+noteCols+` FROM notes WHERE id=?`, id))
}

func (d *DB) ListNotes(f NoteFilter) ([]Note, int) {
	where := []string{"1=1"}
	args := []any{}
	if f.NotebookID != "" {
		where = append(where, "notebook_id=?")
		args = append(args, f.NotebookID)
	}
	if f.Tag != "" {
		where = append(where, `tags_json LIKE ?`)
		args = append(args, `%"`+f.Tag+`"%`)
	}
	if f.Search != "" {
		where = append(where, "(title LIKE ? OR body LIKE ?)")
		s := "%" + f.Search + "%"
		args = append(args, s, s)
	}
	if f.Pinned == "true" {
		where = append(where, "pinned=1")
	}
	if f.Archived == "true" {
		where = append(where, "archived=1")
	} else if f.Archived != "all" {
		where = append(where, "archived=0")
	}
	w := strings.Join(where, " AND ")

	var total int
	d.db.QueryRow("SELECT COUNT(*) FROM notes WHERE "+w, args...).Scan(&total)

	order := "updated_at"
	switch f.SortBy {
	case "created":
		order = "created_at"
	case "title":
		order = "title"
	}
	dir := "DESC"
	if f.SortDir == "asc" {
		dir = "ASC"
	}
	if f.Limit <= 0 {
		f.Limit = 50
	}
	q := fmt.Sprintf(
		"SELECT %s FROM notes WHERE %s ORDER BY pinned DESC, %s %s LIMIT ? OFFSET ?",
		noteCols, w, order, dir,
	)
	args = append(args, f.Limit, f.Offset)
	rows, err := d.db.Query(q, args...)
	if err != nil {
		return nil, 0
	}
	defer rows.Close()
	var out []Note
	for rows.Next() {
		if n := d.scanNote(rows); n != nil {
			out = append(out, *n)
		}
	}
	return out, total
}

func (d *DB) UpdateNote(id string, n *Note) error {
	n.UpdatedAt = now()
	n.WordCount = wordCount(n.Body)
	tj, _ := json.Marshal(n.Tags)
	pinned, archived := 0, 0
	if n.Pinned {
		pinned = 1
	}
	if n.Archived {
		archived = 1
	}
	_, err := d.db.Exec(
		`UPDATE notes SET notebook_id=?, title=?, body=?, tags_json=?, pinned=?, archived=?, updated_at=?
		 WHERE id=?`,
		n.NotebookID, n.Title, n.Body, string(tj), pinned, archived, n.UpdatedAt, id,
	)
	return err
}

func (d *DB) DeleteNote(id string) error {
	_, err := d.db.Exec(`DELETE FROM notes WHERE id=?`, id)
	return err
}

func (d *DB) PinNote(id string, pinned bool) error {
	v := 0
	if pinned {
		v = 1
	}
	_, err := d.db.Exec(`UPDATE notes SET pinned=?, updated_at=? WHERE id=?`, v, now(), id)
	return err
}

func (d *DB) ArchiveNote(id string, archived bool) error {
	v := 0
	if archived {
		v = 1
	}
	_, err := d.db.Exec(`UPDATE notes SET archived=?, updated_at=? WHERE id=?`, v, now(), id)
	return err
}

// ─── Tags ─────────────────────────────────────────────────────────

func (d *DB) AllTags() []string {
	rows, err := d.db.Query(`SELECT DISTINCT tags_json FROM notes WHERE tags_json != '[]'`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	seen := map[string]bool{}
	for rows.Next() {
		var j string
		rows.Scan(&j)
		var tags []string
		json.Unmarshal([]byte(j), &tags)
		for _, t := range tags {
			seen[t] = true
		}
	}
	out := make([]string, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	return out
}

// ─── Export ───────────────────────────────────────────────────────

func (d *DB) ExportMarkdown(id string) string {
	n := d.GetNote(id)
	if n == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("# " + n.Title + "\n\n")
	if len(n.Tags) > 0 {
		sb.WriteString("Tags: " + strings.Join(n.Tags, ", ") + "\n")
	}
	sb.WriteString("Created: " + n.CreatedAt + "\n")
	sb.WriteString("Updated: " + n.UpdatedAt + "\n\n---\n\n")
	sb.WriteString(n.Body)
	return sb.String()
}

func (d *DB) ExportAll() string {
	notes, _ := d.ListNotes(NoteFilter{Limit: 10000, Archived: "all"})
	var sb strings.Builder
	for i, n := range notes {
		if i > 0 {
			sb.WriteString("\n\n---\n\n")
		}
		sb.WriteString("# " + n.Title + "\n\n")
		if len(n.Tags) > 0 {
			sb.WriteString("Tags: " + strings.Join(n.Tags, ", ") + "\n")
		}
		sb.WriteString(n.Body)
	}
	return sb.String()
}

// ─── Stats ────────────────────────────────────────────────────────

type Stats struct {
	Notes     int `json:"notes"`
	Notebooks int `json:"notebooks"`
	Pinned    int `json:"pinned"`
	Archived  int `json:"archived"`
	Tags      int `json:"tags"`
	Words     int `json:"words"`
}

func (d *DB) Stats() Stats {
	var s Stats
	d.db.QueryRow(`SELECT COUNT(*) FROM notes WHERE archived=0`).Scan(&s.Notes)
	d.db.QueryRow(`SELECT COUNT(*) FROM notebooks`).Scan(&s.Notebooks)
	d.db.QueryRow(`SELECT COUNT(*) FROM notes WHERE pinned=1 AND archived=0`).Scan(&s.Pinned)
	d.db.QueryRow(`SELECT COUNT(*) FROM notes WHERE archived=1`).Scan(&s.Archived)
	s.Tags = len(d.AllTags())
	rows, _ := d.db.Query(`SELECT body FROM notes WHERE archived=0`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var b string
			rows.Scan(&b)
			s.Words += wordCount(b)
		}
	}
	return s
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
