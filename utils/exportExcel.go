package utils

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ExportPaymentRequestExcel(filename string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Set default font
	f.SetDefaultFont("Times New Roman")

	// Create a new sheet
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// Set column widths
	f.SetColWidth(sheet, "A", "H", 15)
	f.SetColWidth(sheet, "F", "F", 20) // Wider for Nội Dung
	f.SetColWidth(sheet, "G", "G", 20) // Wider for Số Tiền

	// Define styles
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 2, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 2, Color: "000000"},
		},
	})

	centerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	moneyStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		NumFmt:    3, // Thousands separator (e.g., 6,000,000)
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	boldMoneyStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11}, // In đậm
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		NumFmt:    3, // Dấu phân cách hàng nghìn (e.g., 6,000,000)
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	elevenBoldCenterStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	elevenBoldStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
	})
	centerNoBoldStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	tweleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 12,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	tweleBoldStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "top", Style: 2, Color: "000000"},
		},
	})
	tenStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 7, Color: "000000"},
			{Type: "bottom", Style: 7, Color: "000000"},
		},
	})
	elevenBottomStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		NumFmt: 3,
		Border: []excelize.Border{
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 7, Color: "000000"},
		},
	})
	elevenRightBottomStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},

		Border: []excelize.Border{
			{Type: "top", Style: 7, Color: "000000"},
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 7, Color: "000000"},
			{Type: "right", Style: 2, Color: "000000"},
		},
	})
	elevenBorderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		NumFmt: 3,
		Border: []excelize.Border{
			{Type: "top", Style: 7, Color: "000000"},
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 7, Color: "000000"},
		},
	})
	elevenTopStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		NumFmt: 3,
		Border: []excelize.Border{
			{Type: "top", Style: 7, Color: "000000"},
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	elevenRightTopStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "top", Style: 7, Color: "000000"},
			{Type: "right", Style: 2, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
			{Type: "left", Style: 1, Color: "000000"},
		},
	})
	elevenRightStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},

		Border: []excelize.Border{
			{Type: "right", Style: 2, Color: "000000"},
			{Type: "bottom", Style: 7, Color: "000000"},
			{Type: "top", Style: 7, Color: "000000"},
			{Type: "left", Style: 1, Color: "000000"},
		},
	})

	// Insert the logo
	// logoPath := `C:\Users\lenovo\NA_BC\test_logo.png`
	// if _, err := os.Stat(logoPath); os.IsNotExist(err) {
	// 	fmt.Println("Logo file not found:", logoPath)
	// } else {
	// 	err := f.AddPicture(sheet, "A1", logoPath, &excelize.GraphicOptions{
	// 		OffsetX: 5,
	// 		OffsetY: 5,
	// 		ScaleX:  0.5,
	// 		ScaleY:  0.5,
	// 	})
	// 	if err != nil {
	// 		fmt.Println("Failed to add logo:", err)
	// 		// Fallback: Continue without logo if it fails
	// 	}
	// }
	// Adjust company information position (shift down to avoid overlapping with logo)
	f.SetCellValue(sheet, "C1", "CTY CP LIÊN KẾT VẬN TẢI THÔNG MINH")
	f.SetCellValue(sheet, "C2", "Số 6 tổ 1 Nam Sơn, P. Đằng Giang, Q. Ngô Quyền, TP Hải Phòng, VN")
	f.SetCellValue(sheet, "C3", "MST: 0201957913")
	f.MergeCell(sheet, "C1", "H1")
	f.MergeCell(sheet, "C2", "H2")
	f.MergeCell(sheet, "C3", "H3")
	f.SetCellStyle(sheet, "C1", "C1", elevenBoldCenterStyle)
	f.SetCellStyle(sheet, "C2", "C2", centerNoBoldStyle)
	f.SetCellStyle(sheet, "C3", "C3", centerNoBoldStyle)

	// Title
	f.SetCellValue(sheet, "A4", "GIẤY ĐỀ NGHỊ THANH TOÁN")
	f.MergeCell(sheet, "A4", "H4")
	f.SetCellStyle(sheet, "A4", "A4", titleStyle)

	// Requester Information
	f.SetCellValue(sheet, "A5", "Kính gửi: Ban lãnh đạo Công ty CP Liên Kết Vận Tải Thông Minh")
	f.SetCellValue(sheet, "A6", "Tên tôi là: ")
	f.SetCellValue(sheet, "C6", "Hoàng Thị Phương Ly")
	f.SetCellValue(sheet, "F6", "Bộ phận:")
	f.SetCellValue(sheet, "G6", "Kế hoạch")
	f.SetCellValue(sheet, "A7", "Xin đề nghị thanh toán cho bên vận tải:")
	f.SetCellValue(sheet, "D7", "CÔNG TY CỔ PHẦN TIẾP VẬN 3T")
	f.MergeCell(sheet, "A5", "H5")
	f.MergeCell(sheet, "A6", "B6")
	f.MergeCell(sheet, "A7", "C7")
	f.MergeCell(sheet, "D7", "H7")
	f.SetCellStyle(sheet, "A5", "H5", elevenBoldCenterStyle)

	f.SetCellStyle(sheet, "A6", "B6", elevenBoldStyle)
	f.SetCellStyle(sheet, "F6", "F6", elevenBoldStyle)
	f.SetCellStyle(sheet, "D7", "D7", titleStyle)

	// Table Headers
	headers := []string{"STT", "NGÀY THÁNG", "SỐ ĐƠN", "SỐ CONT", "TUYẾN ĐƯỜNG", "NỘI DUNG", "SỐ TIỀN (VNĐ)", "GHI CHÚ"}
	for col, header := range headers {
		cell := fmt.Sprintf("%c8", 'A'+col)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Table Data
	data := [][]interface{}{
		{1, "13/1/2025", "KD.12012025.28", "MSNU6788986", "HẢI PHÒNG - NGỌC LẶC, THANH HOÁ", "Cước", 6000000, ""},
		{"", "", "", "", "", "Nằm lại", "-", ""},
		{"", "", "", "", "", "Ca xe", "-", ""},
		{"", "", "", "", "", "Lạch Huyện", "-", ""},
		{"", "", "", "", "", "Luật", "-", ""},
		{"", "", "", "", "", "Thuế VAT (8%)", 480000, ""},
		{"", "CỘNG", "", "", "", "", 6480000, ""},
		{2, "13/1/2025", "KD.12012025.29", "MSNU6788137", "HẢI PHÒNG - NGỌC LẶC, THANH HOÁ", "Cước", 6000000, ""},
		{"", "", "", "", "", "Nằm lại", "", ""},
		{"", "", "", "", "", "Ca xe", "", ""},
		{"", "", "", "", "", "Lạch Huyện", "", ""},
		{"", "", "", "", "", "Kiểm Hàng", "", ""},
		{"", "", "", "", "", "Thuế VAT (8%)", 480000, ""},
		{"", "CỘNG", "", "", "", "", 6480000, ""},
		{"TỔNG", "", "", "", "", "", 12960000, ""},
	}

	for rowIdx, row := range data { // Changed from data.Items to data
		for colIdx, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowIdx+9) // Start at row 9 to align with headers at row 14
			if err := f.SetCellValue(sheet, cell, value); err != nil {
				return fmt.Errorf("failed to set cell %s: %w", cell, err)
			}
			// Áp dụng kiểu cho cột "NỘI DUNG" (cột F)
			if colIdx == 6 {
				// Kiểm tra giá trị trong cột "NỘI DUNG" (cột F) của cùng hàng
				if valInt, ok := value.(int); ok && valInt == 6480000 || valInt == 12960000 {
					if err := f.SetCellStyle(sheet, cell, cell, boldMoneyStyle); err != nil {
						return fmt.Errorf("failed to style cell %s: %w", cell, err)
					}
				} else {
					if err := f.SetCellStyle(sheet, cell, cell, moneyStyle); err != nil {
						return fmt.Errorf("failed to style cell %s: %w", cell, err)
					}
				}
				// Áp dụng kiểu cho cột "SỐ ĐƠN" (cột C) khi là "CỘNG"
			} else if strings.ToUpper(fmt.Sprintf("%v", value)) == "CỘNG" {
				cellName, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
				f.SetCellStyle(sheet, cellName, cellName, elevenBoldStyle)
			} else {
				if err := f.SetCellStyle(sheet, cell, cell, centerStyle); err != nil {
					return fmt.Errorf("failed to style cell %s: %w", cell, err)
				}
			}

		}
	}

	f.MergeCell(sheet, "A23", "F23")
	// Apply bold style to the merged cell A23
	f.SetCellStyle(sheet, "A23", "A23", elevenBoldStyle)
	f.SetCellStyle(sheet, "E8", "E8", tweleBoldStyle)
	f.SetCellStyle(sheet, "E9", "E22", tenStyle)

	f.SetCellStyle(sheet, "A9", "D9", elevenBottomStyle)

	f.SetCellStyle(sheet, "A10", "D14", elevenBorderStyle)
	f.SetCellStyle(sheet, "A15", "D15", elevenTopStyle)

	f.SetCellStyle(sheet, "F9", "G9", elevenBottomStyle)
	f.SetCellStyle(sheet, "H9", "H9", elevenRightBottomStyle)
	f.SetCellStyle(sheet, "F10", "G14", elevenBorderStyle)
	f.SetCellStyle(sheet, "H10", "H14", elevenRightStyle)
	f.SetCellStyle(sheet, "F15", "G15", elevenTopStyle)
	f.SetCellStyle(sheet, "H15", "H15", elevenRightTopStyle)

	f.SetCellStyle(sheet, "A16", "D16", elevenBottomStyle)

	f.SetCellStyle(sheet, "F16", "G16", elevenBottomStyle)
	f.SetCellStyle(sheet, "F17", "G21", elevenBorderStyle)
	f.SetCellStyle(sheet, "H16", "H16", elevenRightBottomStyle)
	f.SetCellStyle(sheet, "H17", "H21", elevenRightStyle)
	f.SetCellStyle(sheet, "A17", "D21", elevenBorderStyle)
	f.SetCellStyle(sheet, "A22", "D22", elevenTopStyle)

	f.SetColWidth(sheet, "A", "A", 5.66)
	f.SetColWidth(sheet, "B", "B", 10.99)
	f.SetColWidth(sheet, "C", "C", 20.21)
	f.SetColWidth(sheet, "D", "D", 14.44)
	f.SetColWidth(sheet, "E", "E", 35.10)
	f.SetColWidth(sheet, "F", "F", 21.88)
	f.SetColWidth(sheet, "G", "G", 17.77)
	f.SetColWidth(sheet, "H", "H", 11.88)
	for rowIdx := 1; rowIdx <= 4; rowIdx++ {
		if err := f.SetRowHeight(sheet, rowIdx, 18.6); err != nil {
			fmt.Printf("failed to set height for row %d: %v\n", rowIdx, err)
			return nil
		}
	}
	if err := f.SetRowHeight(sheet, 4, 18); err != nil {
		fmt.Printf("failed to set height for row 4: %v\n", err)
		return nil
	}
	for rowIdx := 5; rowIdx <= 6; rowIdx++ {
		if err := f.SetRowHeight(sheet, rowIdx, 20.1); err != nil {
			fmt.Printf("failed to set height for row %d: %v\n", rowIdx, err)
			return nil
		}
	}
	if err := f.SetRowHeight(sheet, 7, 28.1); err != nil {
		fmt.Printf("failed to set height for row 4: %v\n", err)
		return nil
	}
	if err := f.SetRowHeight(sheet, 8, 25.1); err != nil {
		fmt.Printf("failed to set height for row 4: %v\n", err)
		return nil
	}
	for rowIdx := 9; rowIdx <= 23; rowIdx++ {
		if err := f.SetRowHeight(sheet, rowIdx, 17); err != nil {
			fmt.Printf("failed to set height for row %d: %v\n", rowIdx, err)
			return nil
		}
	}
	for rowIdx := 24; rowIdx <= 25; rowIdx++ {
		if err := f.SetRowHeight(sheet, rowIdx, 20.1); err != nil {
			fmt.Printf("failed to set height for row %d: %v\n", rowIdx, err)
			return nil
		}
	}
	for rowIdx := 26; rowIdx <= 27; rowIdx++ {
		if err := f.SetRowHeight(sheet, rowIdx, 15.6); err != nil {
			fmt.Printf("failed to set height for row %d: %v\n", rowIdx, err)
			return nil
		}
	}
	// Bank Information
	f.SetCellValue(sheet, "A24", "THÔNG TIN TÀI KHOẢN NHẬN TIỀN")
	f.SetCellValue(sheet, "A25", "Chủ tài khoản: ")
	f.SetCellValue(sheet, "C25", "CÔNG TY CỔ PHẦN TIẾP VẬN 3T")
	f.SetCellValue(sheet, "A26", "Số TK: ")
	f.SetCellValue(sheet, "C26", 6796796888)
	f.SetCellValue(sheet, "E26", "tại Ngân hàng: Ngân hàng ACB - Chi nhánh Duyên Hải, Hải Phòng")
	f.MergeCell(sheet, "A24", "C24")
	f.MergeCell(sheet, "A25", "B25")
	f.MergeCell(sheet, "C25", "E25")
	f.MergeCell(sheet, "A26", "B26")
	f.SetCellStyle(sheet, "A24", "A24", elevenBoldStyle)
	f.SetCellStyle(sheet, "E26", "E26", tweleStyle)

	// Signatures
	f.SetCellValue(sheet, "A27", "NGƯỜI ĐỀ NGHỊ")
	f.SetCellValue(sheet, "D27", "TRƯỞNG BỘ PHẬN")
	f.SetCellValue(sheet, "F27", "KẾ TOÁN")
	f.SetCellValue(sheet, "G27", "GIÁM ĐỐC")
	f.MergeCell(sheet, "A27", "C27")
	f.MergeCell(sheet, "D27", "E27")

	f.SetCellStyle(sheet, "A27", "G27", elevenBoldStyle)

	// Save the Excel file with the provided filename
	if err := f.SaveAs(filename); err != nil {
		fmt.Println("Error saving Excel file:", err)
		return err
	}
	fmt.Println("Excel file generated successfully:", filename)
	return nil
}
