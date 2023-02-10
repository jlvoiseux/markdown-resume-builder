package builder

import (
	"context"
	"fmt"
	"image"
	"io"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func printToPDF(urlstr string, res *[]byte, scale float64) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(3 * time.Second),
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

func getImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
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
