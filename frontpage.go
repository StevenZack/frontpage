package frontpage

import (
	"text/template"

	"github.com/StevenZack/openurl"

	"github.com/StevenZack/frontpage/views"

	"github.com/StevenZack/fasthttp"
	"github.com/StevenZack/frontpage/util"
)

type FrontPage struct {
	r        *fasthttp.Router
	vars     *Vars
	WsServer *WsServer
}

func New(html string) *FrontPage {
	fp := &FrontPage{
		r:    fasthttp.NewRouter(),
		vars: NewVars(),
	}
	fp.WsServer = NewWsServer(fp.r.GetServer().Shutdown)
	fp.HandleFunc("/var.js", func(cx *fasthttp.RequestCtx) {
		cx.SetJsHeader()
		t := template.New("var.js")
		t.Parse(views.Str_var)
		t.Execute(cx, fp.vars)
	})
	fp.HandleHtml("/", html)
	fp.HandleFunc("/ws", fp.WsServer.ws)
	return fp
}

func (f *FrontPage) HandleHtml(path, html string) {
	f.r.HandleFunc(path, func(cx *fasthttp.RequestCtx) {
		s, e := util.AddHead(html, `<script src="/var.js" type="text/javascript"></script>`)
		if e != nil {
			cx.Error(e.Error(), fasthttp.StatusBadRequest)
			return
		}
		cx.WriteHTML(s)
	})
}

func (f *FrontPage) HandleJs(path, js string) {
	f.r.HandleFunc(path, func(cx *fasthttp.RequestCtx) {
		cx.SetJsHeader()
		cx.WriteString(js)
	})
}

func (f *FrontPage) HandleCss(path, css string) {
	f.r.HandleFunc(path, func(cx *fasthttp.RequestCtx) {
		cx.SetCssHeader()
		cx.WriteString(css)
	})
}

func (f *FrontPage) HandleFunc(path string, handler func(cx *fasthttp.RequestCtx)) {
	f.r.HandleFunc(path, handler)
}

func (f *FrontPage) Run() error {
	openurl.OpenApp("http://" + f.vars.Addr)
	return f.r.ListenAndServe(f.vars.Addr)
}

func (f *FrontPage) RunBrowser() error {
	openurl.Open("http://" + f.vars.Addr)
	return f.r.ListenAndServe(f.vars.Addr)
}
