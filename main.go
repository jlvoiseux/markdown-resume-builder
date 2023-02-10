package main

import (
	"context"
	"markdown-resume-builder/builder"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
)

func main() {

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cwd := filepath.Dir(exe)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	srv, wg, err := builder.InitServer(cwd)
	if err != nil {
		panic(err)
	}
	defer builder.CloseServer(ctx, srv, wg)

	if len(os.Args) == 1 {
		err = builder.HandleGui(ctx, srv, wg, cwd)
	} else {
		err = builder.HandleCli(ctx, srv, wg, cwd)
	}

	if err != nil {
		panic(err)
	}
}
