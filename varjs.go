package frontpage

import (
	"github.com/StevenZack/fasthttp"
	"github.com/StevenZack/frontpage/views"
)

func varjs(cx *fasthttp.RequestCtx) {
	cx.SetJsHeader()
	cx.WriteString(views.Str_var)
}
