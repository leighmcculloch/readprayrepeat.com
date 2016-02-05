package static

type Page struct {
	Path string
	Func PageFunc
}

type PageFunc func(path string) (data interface{}, tmpls []string, tmpl string)
