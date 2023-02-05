package main

import (
	"context"

	"github.com/chromedp/chromedp"
)

func main() {

	handleCliArgs()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	srv, wg := initServer()

	handleGui(ctx, srv, wg)
	closeServer(ctx, srv, wg)
}
