package utils

import (
	"ecommerce-test/internal/models"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func GenerateProductPDF(products []models.Product) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Product Report", false)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Product Report")

	pdf.Ln(12)
	pdf.SetFont("Arial", "B", 10)

	// Headers
	headers := []string{"Reference", "Name", "Category", "Status", "Date", "Price", "Supplier", "Qty"}
	for _, header := range headers {
		pdf.CellFormat(25, 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)

	// Content
	for _, p := range products {
		pdf.CellFormat(25, 6, p.ProductReference, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 6, p.ProductName, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 6, p.ProductCategory, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 6, p.Status, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 6, p.CreatedAt.String(), "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("%.2f", p.Price), "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 6, p.Supplier, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("%d", p.Quantity), "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	return pdf, nil
}
