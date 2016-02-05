package static

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Static struct {
	SourceDir  string
	BuildDir   string
	ServerPort int

	TemplateFuncs template.FuncMap

	templates map[string]*template.Template
	pages     map[string]*Page
}

func NewStatic() *Static {
	return &Static{
		SourceDir:  "source",
		BuildDir:   "build",
		ServerPort: 4567,
		TemplateFuncs: template.FuncMap{
			"UnsafeHTML": unsafeHTML,
		},
		templates: make(map[string]*template.Template),
		pages:     make(map[string]*Page),
	}
}

func (s *Static) Handle(path string, PageFunc PageFunc) {
	s.pages[path] = &Page{
		Path: path,
		Func: PageFunc,
	}
}

const (
	RunCommandServer = "server"
	RunCommandBuild  = "build"
)

func (s *Static) Run() {
	command := RunCommandServer
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	var err error
	switch command {
	case RunCommandBuild:
		err = s.Build()
	case RunCommandServer:
		err = s.ListenAndServe(":4567")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Static) Build() error {
	for path := range s.pages {
		err := s.BuildPage(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Static) BuildPage(path string) error {
	fp := fmt.Sprintf("%s%s", s.BuildDir, path)
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	err = s.handleRequest(f, path, false)
	return err
}

func (s *Static) ListenAndServe(addr string) error {
	fileServer := http.FileServer(http.Dir(s.BuildDir))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/":
			path = "/index.html"
		case !strings.HasSuffix(path, ".html"):
			path = fmt.Sprintf("%s.html", path)
		}
		err := s.handleRequest(w, path, true)
		if err != nil {
			if err == errNotFound {
				fileServer.ServeHTTP(w, r)
			} else {
				log.Fatal(err)
			}
		}
	})
	return http.ListenAndServe(addr, mux)
}

var errNotFound = errors.New("No handler for path")

func (s *Static) handleRequest(w io.Writer, path string, ignoreCache bool) error {
	log.Printf(" => %s", path)
	page := s.getPageForPath(path)
	if page == nil {
		return errNotFound
	}
	data, tmpls, tmpl := page.Func(path)
	err := s.writeResponse(w, data, tmpls, tmpl, ignoreCache)
	return err
}

func (s *Static) getPageForPath(path string) *Page {
	return s.pages[path]
}

func (s *Static) expandTemplatePaths(tmpl []string) []string {
	expandedTmpl := make([]string, len(tmpl))
	for i, t := range tmpl {
		expandedTmpl[i] = s.templatePath(t)
	}
	return expandedTmpl
}

func getTemplateCacheHash(tmpl []string) string {
	h := sha1.New()
	for _, t := range tmpl {
		io.WriteString(h, t)
	}
	return fmt.Sprintf("% x", h.Sum(nil))
}

func (s *Static) getTemplates(tmpl []string, ignoreCache bool) (*template.Template, error) {
	var err error
	tmpl = s.expandTemplatePaths(tmpl)
	h := getTemplateCacheHash(tmpl)
	if s.templates[h] == nil || ignoreCache {
		s.templates[h], err = template.New("all").Funcs(s.TemplateFuncs).ParseFiles(tmpl...)
	}
	return s.templates[h], err
}

func (s *Static) writeResponse(w io.Writer, data interface{}, tmpls []string, tmpl string, ignoreCache bool) error {
	templates, err := s.getTemplates(tmpls, ignoreCache)
	if err != nil {
		return err
	}
	return templates.ExecuteTemplate(w, tmpl, data)
}

func (s *Static) templatePath(filename string) string {
	return path.Join(s.SourceDir, filename)
}

func unsafeHTML(s string) template.HTML {
	return template.HTML(s)
}
