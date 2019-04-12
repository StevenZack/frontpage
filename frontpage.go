package frontpage

import (
	"fmt"

	"github.com/StevenZack/fasthttp"
	"github.com/StevenZack/tools/strToolkit"
)

func Run(str string) error {
	port := strToolkit.RandomPort()
	r := fasthttp.NewRouter()
	r.HandleFunc("/", func(cx *fasthttp.RequestCtx) {
		cx.SetHtmlHeader()
		cx.WriteString(str)
	})
	fmt.Println("listened on http://localhost:" + port)
	return r.ListenAndServe(":" + port)
}
