package ginsupport

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

var DefaultDelims = render.Delims{Left: "{{", Right: "}}"}
var _LineSeparator = string(os.PathSeparator)

type compositeHtmlRender struct {
	basepath        string
	ext             string
	store           map[string]*renderEntry
	engine          *gin.Engine
	defaultTemplate *template.Template
	delims          render.Delims

	printDebugInfo bool
}

type renderEntry struct {
	templ *template.Template
	store map[string]*renderEntry
}

func (c *compositeHtmlRender) Instance(s string, a any) render.Render {
	fullpathSplits := strings.Split(filepath.Join(c.basepath, s), _LineSeparator)
	store := c.store
	var templ *template.Template
	for _, lv := range fullpathSplits[:len(fullpathSplits)-1] { // last level should be the file name, not folder
		entry, ok := store[lv]
		if ok {
			templ = entry.templ
			store = entry.store
		} else {
			// not found and return default template
			templ = nil
			break
		}
	}
	if templ == nil {
		templ = c.defaultTemplate
	}

	return wrappedHtmlRender{
		Template: templ,
		Name:     fullpathSplits[len(fullpathSplits)-1], // last path should be the filename
		Data:     a,
	}
}

type wrappedHtmlRender struct {
	Template *template.Template
	Name     string
	Data     any
}

func (r wrappedHtmlRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

func (r wrappedHtmlRender) WriteContentType(w http.ResponseWriter) {
	func(w http.ResponseWriter, value []string) {
		header := w.Header()
		if val := header["Content-Type"]; len(val) == 0 {
			header["Content-Type"] = value
		}
	}(w, []string{"text/html; charset=utf-8"})
}

// addPath will parse the relative path and store the full path containing basePath
// when using the data after addPath method, basePath have to be added as prefix of the template path
// for example: basePath for all templates, the path in the storage should be basePath+templatePath
//
//	basePath=templates, templatePath=default/admin, expected path=templates/default/admin
func (c *compositeHtmlRender) addPath(v string) {
	fullpathSplits := strings.Split(filepath.Join(c.basepath, v), _LineSeparator)
	store := c.store
	var processPath string
	for _, lv := range fullpathSplits {
		entry, ok := store[lv]
		processPath = filepath.Join(processPath, lv)
		if c.printDebugInfo {
			_, _ = fmt.Fprintln(gin.DefaultWriter, "Load template path:", filepath.Join(processPath, "*."+c.ext))
		}
		if !ok {
			entry = &renderEntry{store: map[string]*renderEntry{}}
			store[lv] = entry

			// load folder templates for the first time
			left := c.delims.Left
			right := c.delims.Right
			loadPath := filepath.Join(processPath, "*."+c.ext)
			if filenames, err := filepath.Glob(loadPath); err != nil {
				panic(err)
			} else if len(filenames) == 0 {
				// no files, to avoid panic inside template
				entry.templ = template.New("").Funcs(c.engine.FuncMap)
			} else {
				// has files
				templ, err := template.New("").Delims(left, right).Funcs(c.engine.FuncMap).ParseGlob(loadPath)
				if err != nil {
					panic(err)
				}
				entry.templ = templ.Funcs(c.engine.FuncMap)
				if c.printDebugInfo {
					_, _ = fmt.Fprintln(gin.DefaultWriter, "Loaded files:")
					for _, t := range templ.Templates() {
						_, _ = fmt.Fprintln(gin.DefaultWriter, " - ", t.Name())
					}
				}
			}
		}
		store = entry.store
	}
}

func LoadHTMLTemplateFolder(engine *gin.Engine, folder string, ext string, delims render.Delims) render.HTMLRender {
	return loadHTMLTemplateFolder0(engine, folder, ext, delims, true)
}

func loadHTMLTemplateFolder0(engine *gin.Engine, folder string, ext string, delims render.Delims, printDebugInfo bool) render.HTMLRender {
	if len(folder) <= 0 {
		panic(errors.New("should specify a folder"))
	}
	if len(ext) <= 0 {
		panic(errors.New("should specify a file extension"))
	}

	var scanFolders []string
	if err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			relpath, err := filepath.Rel(folder, path)
			if err != nil {
				panic(err)
			}
			scanFolders = append(scanFolders, relpath)
			return nil
		} else {
			// skip file
			return nil
		}
	}); err != nil {
		panic(err)
	}

	r := new(compositeHtmlRender)
	r.store = map[string]*renderEntry{}
	r.basepath = filepath.Join(folder)
	r.ext = ext
	r.engine = engine
	r.defaultTemplate = template.New("")
	r.delims = delims
	r.printDebugInfo = printDebugInfo
	for _, v := range scanFolders {
		r.addPath(v)
	}

	return r
}

func InitHtmlTemplate(engine *gin.Engine, r render.HTMLRender) {
	engine.HTMLRender = r
}

func InitHtmlTemplateFromFolder(engine *gin.Engine, folder string, ext string) {
	engine.HTMLRender = LoadHTMLTemplateFolder(engine, folder, ext, DefaultDelims)
}

func InitHtmlTemplateFromFolderWithDebugReloading(engine *gin.Engine, folder string, ext string) {
	engine.HTMLRender = LoadHTMLTemplateFolder(engine, folder, ext, DefaultDelims)

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			func() {
				<-ticker.C
				defer func() {
					if e := recover(); e != nil {
						_, _ = fmt.Fprintln(gin.DefaultWriter, "Reload template failed:", e)
					}
				}()
				newRender := loadHTMLTemplateFolder0(engine, folder, ext, DefaultDelims, false)
				<-ticker.C
				engine.HTMLRender = newRender
			}()
		}
	}()
}
