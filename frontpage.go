package frontpage

import (
	"fmt"
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
	binder   *binder
}

func New(html string) *FrontPage {
	fp := &FrontPage{
		r:    fasthttp.NewRouter(),
		vars: NewVars(),
	}
	fp.binder = newBinder(fp.vars)
	fp.WsServer = NewWsServer(fp.r.GetServer().Shutdown)

	// handlers
	fp.HandleFunc("/fp/ws", fp.WsServer.ws)
	fp.HandleFunc("/fp/var.js", func(cx *fasthttp.RequestCtx) {
		cx.SetJsHeader()
		t := template.New("var.js")
		t.Parse(views.Str_var)
		t.Execute(cx, fp.vars)
	})
	fp.HandleFunc("/fp/call/:funcName", fp.binder.handleCall)
	fp.HandleHtml("/", html)

	return fp
}

func (f *FrontPage) HandleHtml(path, html string) {
	f.r.HandleFunc(path, func(cx *fasthttp.RequestCtx) {
		s, e := util.AddHead(html, `<script src="/fp/var.js" type="text/javascript"></script>`)
		if e != nil {
			cx.Error(e.Error(), fasthttp.StatusBadRequest)
			return
		}
		cx.WriteHTML(s)
	})
}

func AddVarsJs(html string) (string, error) {
	s, e := util.AddHead(html, `<script src="/fp/var.js" type="text/javascript"></script>`)
	return s, e
}

func (f *FrontPage) HandleHtmlFunc(path string, fn func(cx *fasthttp.RequestCtx) string) {
	f.r.HandleFunc(path, func(cx *fasthttp.RequestCtx) {
		s, e := util.AddHead(fn(cx), `<script src="/fp/var.js" type="text/javascript"></script>`)
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
	return f.run()
}

func (f *FrontPage) run() error {
	fmt.Println("Listened on http://" + f.vars.Addr)
	return f.r.ListenAndServe(f.vars.Addr)
}

func (f *FrontPage) RunBrowser() error {
	openurl.Open("http://" + f.vars.Addr)
	return f.run()
}

func (f *FrontPage) Bind(v interface{}) {
	f.binder.bind(v)
}

func (f *FrontPage) Eval(s string) {
	f.WsServer.pub(s)
}

func (f *FrontPage) Router() *fasthttp.Router {
	return f.r
}
