package views

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"bytes"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func addTemplatePath(files[]string){
		for i, f := range files{
			files[i] = TemplateDir + f
		}
}

func addTemplateExt(files []string){
	for i, f := range files{
		files[i] = f + TemplateExt
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		log.Panicln(err)
	}
	return files
}

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Panicln(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type){
	case Data:
	default:
		data = Data{
			Yield: data,
		}
	}
	//Create a buffer to load the layout to recover from rendering errors.
	var buffer bytes.Buffer
	err := v.Template.ExecuteTemplate(&buffer, v.Layout, data)
	if err != nil{
		http.Error(w, "Oops! Something went wrong. If the problem persists, please contact us.")
		return
	}
	//Use io package to copy buffer contents into ResponseWriter
	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request){
	v.Render(w,nil)
}

type View struct {
	Template *template.Template
	Layout   string
}
