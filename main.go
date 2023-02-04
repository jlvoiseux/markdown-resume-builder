package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	HTML string = "html"
	PDF         = "pdf"
)

func main() {

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

	var sourceFile []byte
	mode := HTML

	if len(os.Args) > 1 {
		sourceFile = []byte(os.Args[1])
	} else {
		fmt.Println("Please provide a source markdown file.")
		fmt.Println("Ex: .\\markdown-resume-builder.exe resume.md")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		switch os.Args[2] {
		case HTML:
			break
		case PDF:
			mode = PDF
			break
		default:
			fmt.Println("Please provide a valid output format (html or pdf).")
			os.Exit(1)
		}
	}

	if len(os.Args) > 3 {
		fontAwesomeKitTemp := []byte(os.Args[3])
		fontAwesomeKitUrl, err := url.ParseRequestURI(string(fontAwesomeKitTemp))
		if err != nil {
			fmt.Println("Invalid Font Awesome Kit URL.")
			os.Exit(1)
		}
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

	_, err = f.Write([]byte(`<article class="markdown-body">`))
	if err != nil {
		panic(err)
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(`</article>`))
	if err != nil {
		panic(err)
	}

	_, err = copyFile("../style.css", "style.css")
	if err != nil {
		panic(err)
	}

	if mode == PDF {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)
		srv := startHttpServer(httpServerExitDone)

		var buf []byte
		if err := chromedp.Run(ctx, printToPDF(`http://localhost:3000`, &buf, 0.7)); err != nil {
			panic(err)
		}
		if err := os.WriteFile("index.pdf", buf, 0o644); err != nil {
			panic(err)
		}

		if err := srv.Shutdown(ctx); err != nil {
			panic(err)
		}

		httpServerExitDone.Wait()
	}

}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func startHttpServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":3000"}
	http.Handle("/", http.FileServer(http.Dir(".")))

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("ListenAndServe(): %v", err))
		}
	}()
	return srv
}

func printToPDF(urlstr string, res *[]byte, scale float64) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(1 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPaperHeight(11.69).
				WithPaperWidth(8.27).
				WithPrintBackground(true).
				WithScale(scale).
				WithMarginTop(0).
				WithMarginRight(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
