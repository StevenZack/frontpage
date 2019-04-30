package frontpage

import (
	"fmt"
	"text/template"

	"github.com/StevenZack/frontpage/views"

	"github.com/StevenZack/fasthttp"
)

func (fp *FrontPage)varjs(cx *fasthttp.RequestCtx) {
	cx.SetJsHeader()
	t, e := template.New("varjs").Parse(views.Str_var)
	if e != nil {
		fmt.Println(`parse error :`, e)
		return
	}
	t.Execute(cx, fp.fnMap)
}
