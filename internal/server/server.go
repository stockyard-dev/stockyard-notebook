package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/stockyard-dev/stockyard-notebook/internal/store"
)

type Server struct {
	db      *store.DB
	mux     *http.ServeMux
	limMu   sync.RWMutex // guards limits, hot-reloadable by /api/license/activate
	limits  Limits
	dataDir string
	pCfg    map[string]json.RawMessage
}

func New(db *store.DB, limits Limits, dataDir string) *Server {
	s := &Server{
		db:      db,
		mux:     http.NewServeMux(),
		limits:  limits,
		dataDir: dataDir,
	}
	s.loadPersonalConfig()

	// Notebooks
	s.mux.HandleFunc("GET /api/notebooks", s.listNotebooks)
	s.mux.HandleFunc("POST /api/notebooks", s.createNotebook)
	s.mux.HandleFunc("GET /api/notebooks/{id}", s.getNotebook)
	s.mux.HandleFunc("PUT /api/notebooks/{id}", s.updateNotebook)
	s.mux.HandleFunc("DELETE /api/notebooks/{id}", s.deleteNotebook)

	// Notes
	s.mux.HandleFunc("GET /api/notes", s.listNotes)
	s.mux.HandleFunc("POST /api/notes", s.createNote)
	s.mux.HandleFunc("GET /api/notes/{id}", s.getNote)
	s.mux.HandleFunc("PUT /api/notes/{id}", s.updateNote)
	s.mux.HandleFunc("DELETE /api/notes/{id}", s.deleteNote)
	s.mux.HandleFunc("POST /api/notes/{id}/pin", s.pinNote)
	s.mux.HandleFunc("POST /api/notes/{id}/unpin", s.unpinNote)
	s.mux.HandleFunc("POST /api/notes/{id}/archive", s.archiveNote)
	s.mux.HandleFunc("POST /api/notes/{id}/unarchive", s.unarchiveNote)

	// Export
	s.mux.HandleFunc("GET /api/notes/{id}/export", s.exportNote)
	s.mux.HandleFunc("GET /api/export", s.exportAll)

	// Meta
	s.mux.HandleFunc("GET /api/tags", s.allTags)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)

	// Personalization
	s.mux.HandleFunc("GET /api/config", s.configHandler)

	// Extras (works for both notebooks and notes via {resource})
	s.mux.HandleFunc("GET /api/extras/{resource}", s.listExtras)
	s.mux.HandleFunc("GET /api/extras/{resource}/{id}", s.getExtras)
	s.mux.HandleFunc("PUT /api/extras/{resource}/{id}", s.putExtras)

	// License activation — accepts a key, validates, persists, hot-reloads tier
	s.mux.HandleFunc("POST /api/license/activate", s.activateLicense)

	// Tier — read-only license info for dashboard banner. Always reachable.
	s.mux.HandleFunc("GET /api/tier", s.tierInfo)

	// Dashboard
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)

	return s
}

// ServeHTTP wraps the underlying mux with a license-gate middleware.
// In trial-required mode, all writes (POST/PUT/DELETE/PATCH) return 402
// EXCEPT POST /api/license/activate (the only way out of trial state).
// Reads are always allowed — the brand promise is that data on disk
// stays accessible even without an active license.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.shouldBlockWrite(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte(`{"error":"Trial required. Start a 14-day free trial at https://stockyard.dev/ — or paste an existing license key in the dashboard under \"Activate License\".","tier":"trial-required"}`))
		return
	}
	s.mux.ServeHTTP(w, r)
}

func (s *Server) shouldBlockWrite(r *http.Request) bool {
	s.limMu.RLock()
	tier := s.limits.Tier
	s.limMu.RUnlock()
	if tier != "trial-required" {
		return false
	}
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	}
	switch r.URL.Path {
	case "/api/license/activate":
		return false
	}
	return true
}

func (s *Server) activateLicense(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 10*1024))
	if err != nil {
		writeErr(w, 400, "could not read request body")
		return
	}
	var req struct {
		LicenseKey string `json:"license_key"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeErr(w, 400, "invalid json: "+err.Error())
		return
	}
	key := strings.TrimSpace(req.LicenseKey)
	if key == "" {
		writeErr(w, 400, "license_key is required")
		return
	}
	if !ValidateLicenseKey(key) {
		writeErr(w, 400, "license key is not valid for this product — make sure you copied the entire key from the welcome email, including the SY- prefix")
		return
	}
	if err := PersistLicense(s.dataDir, key); err != nil {
		log.Printf("notebook: license persist failed: %v", err)
		writeErr(w, 500, "could not save the license key to disk: "+err.Error())
		return
	}
	s.limMu.Lock()
	s.limits = ProLimits()
	s.limMu.Unlock()
	log.Printf("notebook: license activated via dashboard, persisted to %s/%s", s.dataDir, licenseFilename)
	writeJSON(w, 200, map[string]any{
		"ok":   true,
		"tier": "pro",
	})
}

func (s *Server) tierInfo(w http.ResponseWriter, r *http.Request) {
	s.limMu.RLock()
	tier := s.limits.Tier
	s.limMu.RUnlock()
	resp := map[string]any{
		"tier": tier,
	}
	if tier == "trial-required" {
		resp["trial_required"] = true
		resp["start_trial_url"] = "https://stockyard.dev/"
		resp["message"] = "Your trial is not active. Reads work, but you cannot create or change notes until you start a 14-day trial or activate an existing license key."
	} else {
		resp["trial_required"] = false
	}
	writeJSON(w, 200, resp)
}

// ─── helpers ──────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/ui", http.StatusFound)
}

// ─── personalization ──────────────────────────────────────────────

func (s *Server) loadPersonalConfig() {
	path := filepath.Join(s.dataDir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var cfg map[string]json.RawMessage
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("notebook: warning: could not parse config.json: %v", err)
		return
	}
	s.pCfg = cfg
	log.Printf("notebook: loaded personalization from %s", path)
}

func (s *Server) configHandler(w http.ResponseWriter, r *http.Request) {
	if s.pCfg == nil {
		writeJSON(w, 200, map[string]any{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.pCfg)
}

// ─── extras ───────────────────────────────────────────────────────

func (s *Server) listExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	all := s.db.AllExtras(resource)
	out := make(map[string]json.RawMessage, len(all))
	for id, data := range all {
		out[id] = json.RawMessage(data)
	}
	writeJSON(w, 200, out)
}

func (s *Server) getExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	id := r.PathValue("id")
	data := s.db.GetExtras(resource, id)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func (s *Server) putExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	id := r.PathValue("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErr(w, 400, "read body")
		return
	}
	var probe map[string]any
	if err := json.Unmarshal(body, &probe); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if err := s.db.SetExtras(resource, id, string(body)); err != nil {
		writeErr(w, 500, "save failed")
		return
	}
	writeJSON(w, 200, map[string]string{"ok": "saved"})
}

// ─── notebooks ────────────────────────────────────────────────────

func (s *Server) listNotebooks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"notebooks": orEmpty(s.db.ListNotebooks())})
}

func (s *Server) createNotebook(w http.ResponseWriter, r *http.Request) {
	var nb store.Notebook
	if err := json.NewDecoder(r.Body).Decode(&nb); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if nb.Name == "" {
		writeErr(w, 400, "name required")
		return
	}
	if err := s.db.CreateNotebook(&nb); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, nb)
}

func (s *Server) getNotebook(w http.ResponseWriter, r *http.Request) {
	nb := s.db.GetNotebook(r.PathValue("id"))
	if nb == nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, nb)
}

// updateNotebook accepts a partial payload and preserves omitted fields.
// Uses the JSON RawMessage pattern.
func (s *Server) updateNotebook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ex := s.db.GetNotebook(id)
	if ex == nil {
		writeErr(w, 404, "not found")
		return
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}

	patch := *ex
	if v, ok := raw["name"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Name = s
		}
	}
	if v, ok := raw["slug"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Slug = s
		}
	}
	if v, ok := raw["color"]; ok {
		json.Unmarshal(v, &patch.Color)
	}

	if err := s.db.UpdateNotebook(id, &patch); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetNotebook(id))
}

func (s *Server) deleteNotebook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.db.DeleteNotebook(id); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	s.db.DeleteExtras("notebooks", id)
	writeJSON(w, 200, map[string]string{"deleted": "ok"})
}

// ─── notes ────────────────────────────────────────────────────────

func (s *Server) listNotes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	f := store.NoteFilter{
		NotebookID: q.Get("notebook_id"),
		Tag:        q.Get("tag"),
		Search:     q.Get("search"),
		Pinned:     q.Get("pinned"),
		Archived:   q.Get("archived"),
		SortBy:     q.Get("sort"),
		SortDir:    q.Get("dir"),
		Limit:      limit,
		Offset:     offset,
	}
	notes, total := s.db.ListNotes(f)
	writeJSON(w, 200, map[string]any{"notes": orEmpty(notes), "total": total})
}

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	var n store.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if n.Title == "" {
		writeErr(w, 400, "title required")
		return
	}
	if err := s.db.CreateNote(&n); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, n)
}

func (s *Server) getNote(w http.ResponseWriter, r *http.Request) {
	n := s.db.GetNote(r.PathValue("id"))
	if n == nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, n)
}

// updateNote accepts a partial note and properly preserves omitted fields.
// The original implementation only preserved Title, Tags, and NotebookID,
// silently nuking Body, Pinned, and Archived on every partial PUT. PUT
// with {"title":"x"} would blank the body and unpin/unarchive the note.
//
// Uses the JSON RawMessage pattern so:
//   - body can be set to empty string explicitly (clearing a note)
//   - pinned and archived bools can be set to false explicitly
//   - tags can be set to an empty array
//   - notebook_id can be set to empty string (move to inbox)
func (s *Server) updateNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ex := s.db.GetNote(id)
	if ex == nil {
		writeErr(w, 404, "not found")
		return
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}

	patch := *ex
	if v, ok := raw["title"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Title = s
		}
	}
	if v, ok := raw["body"]; ok {
		json.Unmarshal(v, &patch.Body)
	}
	if v, ok := raw["notebook_id"]; ok {
		json.Unmarshal(v, &patch.NotebookID)
	}
	if v, ok := raw["tags"]; ok {
		var tags []string
		if err := json.Unmarshal(v, &tags); err == nil {
			patch.Tags = tags
		}
	}
	if v, ok := raw["pinned"]; ok {
		json.Unmarshal(v, &patch.Pinned)
	}
	if v, ok := raw["archived"]; ok {
		json.Unmarshal(v, &patch.Archived)
	}

	if err := s.db.UpdateNote(id, &patch); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetNote(id))
}

func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.db.DeleteNote(id); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	s.db.DeleteExtras("notes", id)
	writeJSON(w, 200, map[string]string{"deleted": "ok"})
}

func (s *Server) pinNote(w http.ResponseWriter, r *http.Request) {
	s.db.PinNote(r.PathValue("id"), true)
	writeJSON(w, 200, s.db.GetNote(r.PathValue("id")))
}

func (s *Server) unpinNote(w http.ResponseWriter, r *http.Request) {
	s.db.PinNote(r.PathValue("id"), false)
	writeJSON(w, 200, s.db.GetNote(r.PathValue("id")))
}

func (s *Server) archiveNote(w http.ResponseWriter, r *http.Request) {
	s.db.ArchiveNote(r.PathValue("id"), true)
	writeJSON(w, 200, s.db.GetNote(r.PathValue("id")))
}

func (s *Server) unarchiveNote(w http.ResponseWriter, r *http.Request) {
	s.db.ArchiveNote(r.PathValue("id"), false)
	writeJSON(w, 200, s.db.GetNote(r.PathValue("id")))
}

// ─── export ───────────────────────────────────────────────────────

func (s *Server) exportNote(w http.ResponseWriter, r *http.Request) {
	md := s.db.ExportMarkdown(r.PathValue("id"))
	if md == "" {
		writeErr(w, 404, "not found")
		return
	}
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=note.md")
	w.Write([]byte(md))
}

func (s *Server) exportAll(w http.ResponseWriter, r *http.Request) {
	md := s.db.ExportAll()
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=all-notes.md")
	w.Write([]byte(md))
}

// ─── meta ─────────────────────────────────────────────────────────

func (s *Server) allTags(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"tags": orEmpty(s.db.AllTags())})
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.db.Stats())
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats()
	writeJSON(w, 200, map[string]any{
		"status":    "ok",
		"service":   "notebook",
		"notes":     st.Notes,
		"notebooks": st.Notebooks,
	})
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
