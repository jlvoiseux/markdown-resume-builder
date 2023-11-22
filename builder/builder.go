package builder

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func buildResume(ctx context.Context, srv *http.Server, wg *sync.WaitGroup, cwd string, sourceFile string, mode string, destinationFolder string, fontAwesomeKitUrl string, photoFile string) error {

	header := `<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="style.css">
	<style>
		.markdown-body {
			box-sizing: border-box;
			min-width: 200px;
			max-width: 1200px;
			margin: 0 auto;
			padding: 25px;
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
		return err
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
		return err
	}

	f, err := os.Create(filepath.Join(cwd, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(header))
	if err != nil {
		return err
	}

	_, err = f.Write([]byte("<article class='markdown-body'>\n"))
	if err != nil {
		return err
	}

	if photoFile != "" {
		targetHeight := 150
		width, height, err := getImageDimensions(photoFile)
		if err != nil {
			return err
		}

		ratio := float64(height) / float64(targetHeight)
		targetWidth := int(float64(width) / ratio)

		photoExt := filepath.Ext(photoFile)

		_, err = copyFile(photoFile, filepath.Join(cwd, "photo"+photoExt))
		if err != nil {
			return err
		}

		_, err = f.Write([]byte(fmt.Sprintf("<img style='float: right;' width=%d height=%d src='%s'>\n", targetWidth, targetHeight, "photo"+photoExt)))
		if err != nil {
			return err
		}
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = f.Write([]byte("</article>\n"))
	if err != nil {
		return err
	}

	if mode == "PDF" {
		var buf []byte
		if err := chromedp.Run(ctx, printToPDF(`http://localhost:3000`, &buf, 0.7)); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(cwd, "index.pdf"), buf, 0o644); err != nil {
			return err
		}
	}

	_, err = os.Stat(filepath.Join(destinationFolder, "md-resume-builder-output"))
	if os.IsNotExist(err) {
		if err = os.Mkdir(filepath.Join(destinationFolder, "md-resume-builder-output"), os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	switch mode {
	case "HTML":
		_, err = copyFile(filepath.Join(cwd, "index.html"), filepath.Join(destinationFolder, "md-resume-builder-output", "resume.html"))
		if err != nil {
			return err
		}
		_, err = copyFile(filepath.Join(cwd, "style.css"), filepath.Join(destinationFolder, "md-resume-builder-output", "style.css"))
		if err != nil {
			return err
		}
		if photoFile != "" {
			_, err = copyFile(filepath.Join(cwd, filepath.Join("photo", filepath.Ext(photoFile))), filepath.Join(destinationFolder, "md-resume-builder-output", filepath.Join("photo", filepath.Ext(photoFile))))
			if err != nil {
				return err
			}
		}
	case "PDF":
		_, err = copyFile(filepath.Join(cwd, "index.pdf"), filepath.Join(destinationFolder, "md-resume-builder-output", "resume.pdf"))
		if err != nil {
			return err
		}
	}

	return nil
}
