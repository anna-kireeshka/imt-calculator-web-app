package server

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"unicode"

	"github.com/a-h/templ"
)

type PageData struct {
	Title    string
	Tab      string
	Name     string
	Initials string
	Detail   string
	Centered bool
}

//go:embed templates
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

var pages = parsePages()

// Каждой странице — свой набор шаблонов: layout, partials и она сама.
// Один общий набор не подошёл бы: страницы определяют одно и то же имя
// "content", и в общем наборе выжило бы только последнее определение.
func parsePages() map[string]*template.Template {
	files, err := fs.Glob(templatesFS, "templates/pages/*.html")
	if err != nil {
		panic(err)
	}

	var result = make(map[string]*template.Template, len(files))

	for _, page := range files {
		result[path.Base(page)] = template.Must(template.ParseFS(templatesFS,
			"templates/layouts/*.html",
			"templates/partials/*.html",
			page,
		))
	}

	return result
}

// renderComponent отдаёт templ-компонент. Как и render, сначала пишет в буфер:
// Render может упасть на середине, уже отправив статус 200 и часть разметки.
func (s *Server) renderComponent(ctx context.Context, w http.ResponseWriter, component templ.Component) {
	var buf bytes.Buffer
	if err := component.Render(ctx, &buf); err != nil {
		log.Printf("render component: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func (s *Server) render(w http.ResponseWriter, page string, data PageData) {
	s.renderStatus(w, page, data, http.StatusOK)
}

func (s *Server) renderStatus(w http.ResponseWriter, page string, data PageData, status int) {
	tmpl, ok := pages[page]
	if !ok {
		log.Printf("шаблон не найден: %s", page)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if data.Initials == "" {
		data.Initials = initials(data.Name)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base.html", data); err != nil {
		log.Printf("render %s: %v", page, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func initials(name string) string {
	var result []rune

	var newWord = true
	for _, r := range name {
		switch {
		case unicode.IsSpace(r):
			newWord = true
		case newWord:
			result = append(result, unicode.ToUpper(r))
			newWord = false
			if len(result) == 2 {
				return string(result)
			}
		}
	}

	if len(result) == 0 {
		return "?"
	}
	return string(result)
}
