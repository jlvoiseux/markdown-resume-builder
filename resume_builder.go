package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func buildResume(ctx context.Context, srv *http.Server, wg *sync.WaitGroup, sourceFile string, mode string, fontAwesomeKitUrl string, photoFile string) {

	header := `<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="style.css">
	<style>
		.markdown-body {
			box-sizing: border-box;
			min-width: 200px;
			max-width: 1200px;
			margin: 0 auto;
			padding: 45px;
		}
	
		@media (max-width: 767px) {
			.markdown-body {
				padding: 15px;
			}
		}
	</style>`

	if fontAwesomeKitUrl != "" {
		header += fmt.Sprintf("\n<script src='%s' crossorigin='anonymous'></script>\n", fontAwesomeKitUrl)
	}

	source, err := os.ReadFile(string(sourceFile))
	if err != nil {
		panic(err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}

	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(header))
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte("<article class='markdown-body'>\n"))
	if err != nil {
		panic(err)
	}

	if photoFile != "" {
		targetHeight := 150
		width, height := getImageDimensions(photoFile)
		ratio := float64(height) / float64(targetHeight)
		targetWidth := int(float64(width) / ratio)

		_, err = f.Write([]byte(fmt.Sprintf("<img style='float: right;' width=%d height=%d src='%s'>\n", targetWidth, targetHeight, photoFile)))
		if err != nil {
			panic(err)
		}
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte("</article>\n"))
	if err != nil {
		panic(err)
	}

	_, err = copyFile("../style.css", "style.css")
	if err != nil {
		panic(err)
	}

	if mode == "PDF" {
		var buf []byte
		if err := chromedp.Run(ctx, printToPDF(`http://localhost:3000`, &buf, 0.7)); err != nil {
			panic(err)
		}
		if err := os.WriteFile("index.pdf", buf, 0o644); err != nil {
			panic(err)
		}
	}
}
