package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/topdf", chromedpconverttopdf)
	log.Println("Serving on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var chromedpconverttopdf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create PDF
	var pdfBuffer []byte
	err := chromedp.Run(ctx, printToPDF(&pdfBuffer))
	if err != nil {
		log.Fatalf("Failed to create PDF: %v", err)
	}

	// Write the PDF to a file
	err = os.WriteFile("resume.pdf", pdfBuffer, 0644)
	if err != nil {
		log.Fatalf("Failed to write PDF file: %v", err)
	}

	log.Println("PDF created successfully: resume.pdf")
	w.Write(pdfBuffer)
})

// printToPDF converts HTML content to PDF
func printToPDF( res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("http://localhost:8080/resume_pdf.html"),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
                WithPrintBackground(true).
                WithPaperWidth(10). // Set page size to A4 or Letter A4 > 8.5
                WithPaperHeight(19.2).   // A4 > 11
                WithMarginTop(0).  // Set margins
                WithMarginBottom(0).
                WithMarginLeft(0).
                WithMarginRight(0).
                Do(ctx)
            if err != nil {
                return err
            }
            *res = buf
			return nil
		}),
	}
}
