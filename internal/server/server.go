package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/stockyard-dev/stockyard-notebook/internal/store"
)

type Server struct {
	db     *store.DB
	mux    *http.ServeMux
	limits Limits
}

func New(db *store.DB, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), limits: limits}

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

	// Dashboard
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/notebook/"})})

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/ui", http.StatusFound)
}

// ── Notebooks ──

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
func (s *Server) updateNotebook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ex := s.db.GetNotebook(id)
	if ex == nil {
		writeErr(w, 404, "not found")
		return
	}
	var nb store.Notebook
	if err := json.NewDecoder(r.Body).Decode(&nb); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if nb.Name == "" {
		nb.Name = ex.Name
	}
	if nb.Slug == "" {
		nb.Slug = ex.Slug
	}
	if nb.Color == "" {
		nb.Color = ex.Color
	}
	if err := s.db.UpdateNotebook(id, &nb); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetNotebook(id))
}
func (s *Server) deleteNotebook(w http.ResponseWriter, r *http.Request) {
	if err := s.db.DeleteNotebook(r.PathValue("id")); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"deleted": "ok"})
}

// ── Notes ──

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
func (s *Server) updateNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ex := s.db.GetNote(id)
	if ex == nil {
		writeErr(w, 404, "not found")
		return
	}
	var n store.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if n.Title == "" {
		n.Title = ex.Title
	}
	if n.Tags == nil {
		n.Tags = ex.Tags
	}
	if n.NotebookID == "" {
		n.NotebookID = ex.NotebookID
	}
	if err := s.db.UpdateNote(id, &n); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetNote(id))
}
func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	if err := s.db.DeleteNote(r.PathValue("id")); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
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

// ── Export ──

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

// ── Meta ──

func (s *Server) allTags(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"tags": orEmpty(s.db.AllTags())})
}
func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.db.Stats())
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats()
	writeJSON(w, 200, map[string]any{"status": "ok", "service": "notebook", "notes": st.Notes, "notebooks": st.Notebooks})
}

func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
func init() { log.SetFlags(log.LstdFlags | log.Lshortfile) }
