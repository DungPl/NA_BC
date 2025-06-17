package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-pdf/fpdf"
)

func drawTableRow(pdf *fpdf.Fpdf, font string, data []string, isHeader bool) {
	pdf.SetX(20)
	pdf.SetFont(font, "", 10)
	if isHeader {
		pdf.SetFont(font, "B", 12)
		pdf.SetFillColor(200, 220, 255)
	} else {
		pdf.SetFont(font, "", 12)
	}

	widths := []float64{20, 50, 60, 30, 100, 80}
	aligns := []string{"C", "L", "L", "C", "L", "L"}

	for i, text := range data {
		style := "1"
		align := aligns[i]
		pdf.CellFormat(widths[i], 10, text, style, 0, align, false, 0, "")
	}
	pdf.Ln(10)
}

func ExportAdvanceRequest(filename string) error {
	pdf := fpdf.New("P", "mm", "A2", "")
	pdf.AddPage()

	fontDir := filepath.Join(os.Getenv("USERPROFILE"), "NA_BC", "fonts")
	chosenFont := "DejaVuSans"
	fontPath := filepath.Join(fontDir, chosenFont+".ttf")
	fontBoldPath := filepath.Join(fontDir, chosenFont+"-Bold.ttf")

	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return fmt.Errorf("missing font: %s", fontPath)
	}
	if _, err := os.Stat(fontBoldPath); os.IsNotExist(err) {
		return fmt.Errorf("missing font: %s", fontBoldPath)
	}

	pdf.AddUTF8Font(chosenFont, "", fontPath)
	pdf.AddUTF8Font(chosenFont, "B", fontBoldPath)

	// Insert logo (optional)
	logoPath := filepath.Join(os.Getenv("USERPROFILE"), "NA_BC", "logo.png")
	if _, err := os.Stat(logoPath); err == nil {
		pdf.ImageOptions(logoPath, 20, 10, 30, 0, false,
			fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	}

	// Company header
	pdf.SetFont(chosenFont, "B", 10)
	pdf.SetXY(15, 15) // Điều chỉnh tọa độ cho A2
	pdf.CellFormat(0, 0, "CTY CP LIÊN KẾT VẬN TẢI THÔNG MINH", "", 1, "C", false, 0, "")
	pdf.SetFont(chosenFont, "", 10)
	pdf.CellFormat(0, 6, "Số 6 tổ 1 Nam Sơn, P.Đằng Giang, Q.Ngô Quyền, TP Hải Phòng, VN", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, "MST: 0201957913", "", 1, "C", false, 0, "")

	pdf.Ln(10)

	// Title
	pdf.SetFont(chosenFont, "B", 14)
	pdf.CellFormat(0, 10, "GIẤY ĐỀ NGHỊ TẠM ỨNG", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Request details (aligned left)
	pdf.SetFont(chosenFont, "", 10)

	pdf.SetXY(20, 50) // Điều chỉnh tọa độ cho A2
	pdf.Write(7, "Kính gửi: ")
	pdf.SetFont(chosenFont, "B", 10)
	pdf.Write(7, "Ban lãnh đạo Công ty CP Liên Kết Vận Tải Thông Minh")
	pdf.Ln(7)

	pdf.SetFont(chosenFont, "", 10)
	pdf.SetXY(20, 57)
	pdf.Write(7, "Tên tôi là: ")
	pdf.SetFont(chosenFont, "B", 10)
	pdf.Write(7, "Trần Minh Hiếu")
	pdf.Ln(7)

	pdf.SetFont(chosenFont, "", 10)
	pdf.SetXY(20, 64)
	pdf.Write(7, "Bộ phận: ")
	pdf.SetFont(chosenFont, "B", 10)
	pdf.Write(7, "Ops")
	pdf.Ln(7)

	pdf.SetFont(chosenFont, "", 10)
	pdf.SetXY(20, 71)
	pdf.Write(7, "Xin đề nghị thanh toán cho bên vận tải: ")
	pdf.SetFont(chosenFont, "B", 10)
	pdf.Write(7, "Công ty cổ phần tiếp vận 3T")
	pdf.Ln(10)

	// Table headers
	headers := []string{"#", "Ngày tháng", "Số đơn", "Số Cont", "Tuyến đường", "Nội dung", "Số tiền", ""}
	widths := []float64{20, 60, 50, 30, 90, 60, 30, 30}

	// Dùng font nhỏ hơn
	pdf.SetFont(chosenFont, "B", 9)
	pdf.SetXY(20, 90)
	for i, h := range headers {
		pdf.CellFormat(widths[i], 9, h, "T,B", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	// Dòng dữ liệu
	pdf.SetFont(chosenFont, "", 9)
	pdf.SetXY(20, 99) // Thẳng hàng với header
	values := []string{"1", "2025-01-21T17:00:00Z", "Test12.05020225.5", "1", "Phú Khê - Bắc Ninh - Ngô Quyền - Hải Phòng", "", "", "31"}
	for i, v := range values {
		pdf.CellFormat(widths[i], 9, v, "T,B", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.Ln(10)
	// Dòng Cộng (ngay sau dữ liệu)
	pdf.SetFont(chosenFont, "B", 9)
	pdf.SetXY(40, 108)
	pdf.CellFormat(widths[1], 9, "Cộng", "T,B", 0, "L", false, 0, "")
	pdf.SetX(20 + widths[0] + widths[1] + widths[2] + widths[3] + widths[4] + widths[5] + widths[6])

	pdf.CellFormat(widths[7], 9, "31", "T,B", 0, "L", false, 0, "")
	pdf.Ln(-1)

	// Dòng Tổng (căn thẳng hàng với Cộng)
	pdf.SetXY(20, 117)
	pdf.CellFormat(widths[0], 9, "Tổng", "T,B", 0, "L", false, 0, "")

	pdf.SetX(20 + widths[0] + widths[1] + widths[2] + widths[3] + widths[4] + widths[5] + widths[6])
	pdf.CellFormat(widths[7], 9, "31", "T,B", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Bank info
	pdf.SetFont(chosenFont, "B", 10)
	pdf.SetXY(20, 140)
	pdf.Cell(0, 6, "Thông tin tài khoản nhận tiền")
	pdf.Ln(6)
	pdf.SetFont(chosenFont, "", 10)
	pdf.Cell(0, 6, "Chủ tài khoản: Trần Minh Hiếu")
	pdf.Ln(6)
	pdf.Cell(0, 6, "Chủ tài khoản: Số TK: 6769796888 tại Ngân hàng: Ngân hàng ACB - Chi nhánh Duyên Hải, Hải Phòng")
	pdf.Ln(15)

	// Signatures
	pdf.SetXY(20, 160)
	signatures := []string{"NGƯỜI ĐỀ NGHỊ", "TRƯỞNG BỘ PHẬN", "KẾ TOÁN", "GIÁM ĐỐC"}
	colWidth := 100.0
	for _, s := range signatures {
		pdf.CellFormat(colWidth, 10, s, "", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	// Save PDF file
	if err := pdf.OutputFileAndClose(filename); err != nil {
		log.Printf("Error: Failed to save PDF to %s: %v", filename, err)
		return fmt.Errorf("failed to save PDF to %s: %v", filename, err)
	}
	return nil
}
