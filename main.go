package main

import (
	"context"
	"markdown-resume-builder/builder"
	"os"

	"github.com/chromedp/chromedp"
)

func main() {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	srv, wg, err := builder.InitServer()
	if err != nil {
		panic(err)
	}
	defer builder.CloseServer(ctx, srv, wg)

	if len(os.Args) == 1 {
		err = builder.HandleGui(ctx, srv, wg)
	} else {
		err = builder.HandleCli(ctx, srv, wg)
	}

	if err != nil {
		panic(err)
	}
}
