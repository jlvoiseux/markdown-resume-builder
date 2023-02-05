package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
)

func handleCliArgs() {

	if len(os.Args) == 1 {
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	srv, wg := initServer()

	var sourceFile string
	var mode string
	var fontAwesomeKitUrl *url.URL
	var photoFile string
	var err error

	if len(os.Args) > 1 {
		sourceFile = os.Args[1]
		_, err := os.Stat(sourceFile)
		if err != nil {
			fmt.Println("Invalid Markdown path.")
			os.Exit(1)
		}
	}

	if len(os.Args) > 2 {
		switch strings.ToLower(os.Args[2]) {
		case "html":
			mode = "HTML"
			break
		case "pdf":
			mode = "PDF"
			break
		default:
			fmt.Println("Please provide a valid output format (html or pdf).")
			os.Exit(1)
		}
	}

	if len(os.Args) > 3 {
		fontAwesomeKitTemp := os.Args[3]
		fontAwesomeKitUrl, err = url.ParseRequestURI(fontAwesomeKitTemp)
		if err != nil {
			fmt.Println("Invalid Font Awesome Kit URL.")
			os.Exit(1)
		}
	}

	if len(os.Args) > 4 {
		photoFile = os.Args[4]
		_, err := os.Stat(photoFile)
		if err != nil {
			fmt.Println("Invalid picture path.")
			os.Exit(1)
		}
	} else {
		fmt.Println("Please provide a valid photo file.")
		os.Exit(1)
	}

	buildResume(ctx, srv, wg, sourceFile, mode, fontAwesomeKitUrl.String(), photoFile)
	closeServer(ctx, srv, wg)
	os.Exit(0)
}
