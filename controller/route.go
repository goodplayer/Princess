package controller

import (
	"html/template"
	"path"
)

import (
	"github.com/gin-gonic/gin"

	"moetang.info/prod/Princess/config"
)

type TemplateFunc func(*gin.Context) string
type TemplateCommonFunc func(interface{}) interface{}
type TemplateFunc2 func(*gin.Context) (string, error)

func NewTemplateModel(c *gin.Context) gin.H {
	// build in template obj
	r := gin.H{
		"CONTEXT": c,
		"SITE":    g_TEMPLATE_COMMON_MODEL,
	}
	return r
}

//================================

func RegisterRoute(r *gin.Engine) {
	loadHTMLGlob(r, path.Join(config.GLOBAL_CONFIG.TemplatePath, "*"), registerDefaultFunctions)

	r.GET("/", IndexAction)

	r.Static("/statics", config.GLOBAL_CONFIG.StaticPath)

	r.NoRoute(NoRouteHandler)
}

// register custom template functions

func registerTemplateCommonModel(m map[string]interface{}) {
	m["default_site_name"] = config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName
}

func registerTemplateFunc(m map[string]TemplateFunc) {
	m["princess_requri"] = ShowReqUriFunction
}

func registerTemplateCommonFunc(m map[string]TemplateCommonFunc) {
	m["StdDate"] = StdDate
}

func registerTemplateFunc2(m map[string]TemplateFunc2) {
}

//================================

var (
	g_TEMPLATE_COMMON_MODEL = make(map[string]interface{})
)

func registerDefaultFunctions(tmpl *template.Template) {
	m := make(map[string]TemplateFunc)
	registerTemplateFunc(m)
	tmpl.Funcs(convertTemplateFunc(m))

	m2 := make(map[string]TemplateFunc2)
	registerTemplateFunc2(m2)
	tmpl.Funcs(convertTemplateFunc2(m2))

	mc := make(map[string]TemplateCommonFunc)
	registerTemplateCommonFunc(mc)
	tmpl.Funcs(convertTemplateCommonFunc(mc))

	registerTemplateCommonModel(g_TEMPLATE_COMMON_MODEL)
}

func convertTemplateFunc(m map[string]TemplateFunc) template.FuncMap {
	rm := make(template.FuncMap)
	for k, v := range m {
		rm[k] = v
	}
	return rm
}

func convertTemplateCommonFunc(m map[string]TemplateCommonFunc) template.FuncMap {
	rm := make(template.FuncMap)
	for k, v := range m {
		rm[k] = v
	}
	return rm
}

func convertTemplateFunc2(m map[string]TemplateFunc2) template.FuncMap {
	rm := make(template.FuncMap)
	for k, v := range m {
		rm[k] = v
	}
	return rm
}

//TODO customize templates
func loadHTMLGlob(r *gin.Engine, pattern string, templFunc func(*template.Template)) {
	if gin.IsDebugging() {
		if len(pattern) <= 0 {
			panic("the HTML debug render was created without files or glob pattern")
		}
	}
	templ := template.New("")
	templFunc(templ)
	templ = template.Must(templ.ParseGlob(pattern))
	r.SetHTMLTemplate(templ)
}
