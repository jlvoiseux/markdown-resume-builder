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
	srv, wg, err := initServer()
	if err != nil {
		panic(err)
	}

	var sourceFile string
	var mode string
	var destinationFolder string
	var fontAwesomeKitUrl *url.URL
	var photoFile string

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
		case "pdf":
			mode = "PDF"
		default:
			fmt.Println("Please provide a valid output format (html or pdf).")
			os.Exit(1)
		}
	}

	if len(os.Args) > 3 {
		destinationFolder = os.Args[3]
		_, err := os.Stat(sourceFile)
		if err != nil {
			fmt.Println("Invalid destination folder.")
			os.Exit(1)
		}
	}

	if len(os.Args) > 4 {
		fontAwesomeKitTemp := os.Args[4]
		fontAwesomeKitUrl, err = url.ParseRequestURI(fontAwesomeKitTemp)
		if err != nil {
			fmt.Println("Invalid Font Awesome Kit URL.")
			os.Exit(1)
		}
	}

	if len(os.Args) > 5 {
		photoFile = os.Args[5]
		_, err := os.Stat(photoFile)
		if err != nil {
			fmt.Println("Invalid picture path.")
			os.Exit(1)
		}
	} else {
		fmt.Println("Please provide a valid photo file.")
		os.Exit(1)
	}

	err = buildResume(ctx, srv, wg, sourceFile, mode, destinationFolder, fontAwesomeKitUrl.String(), photoFile)
	if err != nil {
		panic(err)
	}

	err = closeServer(ctx, srv, wg)
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
