package builder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func HandleCli(ctx context.Context, srv *http.Server, wg *sync.WaitGroup) error {
	var sourceFile string
	var mode string
	var destinationFolder string
	var fontAwesomeKitUrl *url.URL
	var photoFile string

	var err error

	if len(os.Args) > 1 {
		sourceFile = os.Args[1]
		_, err := os.Stat(sourceFile)
		if err != nil {
			fmt.Println("Invalid Markdown path.")
			return err
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
			return err
		}
	}

	if len(os.Args) > 3 {
		destinationFolder = os.Args[3]
		_, err := os.Stat(sourceFile)
		if err != nil {
			fmt.Println("Invalid destination folder.")
			return err
		}
	}

	if len(os.Args) > 4 {
		fontAwesomeKitTemp := os.Args[4]
		fontAwesomeKitUrl, err = url.ParseRequestURI(fontAwesomeKitTemp)
		if err != nil {
			fmt.Println("Invalid Font Awesome Kit URL.")
			return err
		}
	}

	if len(os.Args) > 5 {
		photoFile = os.Args[5]
		_, err := os.Stat(photoFile)
		if err != nil {
			fmt.Println("Invalid picture path.")
			return err
		}
	}

	err = buildResume(ctx, srv, wg, sourceFile, mode, destinationFolder, fontAwesomeKitUrl.String(), photoFile)
	if err != nil {
		return err
	}

	return nil
}
