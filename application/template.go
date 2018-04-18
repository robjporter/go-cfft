package application

import (
  "io"
  "../data"
  "../packages/template"
  "html/template"
  "github.com/labstack/echo"
  assetfs "github.com/elazarl/go-bindata-assetfs"
)

// Template is custom renderer for Echo, to render html from bindata
type Template struct {
	templates *template.Template
}

// Render renders template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// NewTemplate creates a new template
func NewTemplate() *Template {
	return &Template{
		templates: binhtml.New(data.Asset, data.AssetDir).MustLoadDirectory("public2"),
	}
}

func NewAssets(root string) *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     data.Asset,
		AssetDir:  data.AssetDir,
		AssetInfo: data.AssetInfo,
		Prefix:    root,
	}
}
