package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

var colWidths2 = map[string]float64{
	"A": 8.88,  // ~40 pixels
	"B": 8.88,  // ~77 pixels
	"C": 11.99, // ~141 pixels
	"D": 14.44, // ~101 pixels
	"E": 10.66, // ~246 pixels
	"F": 15.88, // ~153 pixels
	"G": 8.88,  // ~124 pixels
	"H": 8.88,  // ~83 pixels
	"I": 8.88,
	"K": 8.88,
	"L": 8.88,
	"M": 8.88,
	"N": 8.88,
	"J": 8.88,
	"O": 8.88,
	"P": 8.88,
	"Q": 8.88,
	"R": 8.88,
	"S": 8.88,
	"T": 8.88,
	"U": 8.88,
	"V": 8.88,
	"W": 8.88,
}

var colWidths1 = map[string]float64{
	"A": 5.21,  // ~40 pixels
	"B": 9.21,  // ~77 pixels
	"C": 12.77, // ~141 pixels
	"D": 12.77, // ~101 pixels
	"E": 11.21, // ~246 pixels
	"F": 6.55,  // ~153 pixels
	"G": 33.21, // ~124 pixels
	"H": 11.21, // ~83 pixels
	"I": 11.21,
	"K": 11.21,
	"L": 11.21,
	"M": 11.21,
	"N": 11.21,
	"J": 11.21,
	"O": 12.21,
	"P": 11.99,
	"Q": 11.99,
	"R": 10.77,
	"S": 9.44,
	"T": 10.21,
	"U": 9.44,
	"V": 9.44,
	"W": 9.44,
	"X": 15.21,
	"Y": 22.1,
}

var rowHeight = map[int]float64{
	1: 15.00, // ~100 pixels
}

const (
	sheet1            = "BẢNG KÊ THU"
	sheet2            = "BẢNG KÊ CHI"
	defaultFontShip   = "Times New Roman"
	titleRowShip      = 7
	headerRowShip     = 14
	dataStartRowShip  = 15
	dataStartRowShip2 = 13
	//bankInfoStartRow = 8+7*n+2 // Thông tin ngân hàng ở hàng 24
	//signatureRow     = 8+7*n+5 // Chữ ký ở hàng 27

)

var tableHeadersTotal = []string{
	"STT",
	"NGÀY",
	"SỐ ĐH",
	"SỐ BL / BOOK",
	"SỐ CONT",
	"LOẠI CONT",
	"NƠI ĐÓNG / TRẢ",
	"SỐ XE VC",
	"CƯỚC VC",
	"lạch huyện",
	"ca xe",
	"NẰM LẠI",
	"CP KHÁC",
	"PHÍ DICH VỤ",
	"CỘNG TIỀN HÀNG",

	"TIỀN THUẾ GTGT",
	"TỔNG THANH TOÁN",
	"CSHT", "NÂNG", "HẠ", "CP KHÁC",
	"CỘNG CHI HỘ", "SỐ HĐ", "TỔNG THU", "GHI CHÚ (CƯỚC)",
}
var TableHeadersFee = []string{
	"STT", "NGÀY", "SỐ ĐH", "SỐ CONT", "LOẠI CONT",
	"NƠI ĐÓNG / TRẢ", "SỐ XE VC", "CƯỚC VC", "lạch huyện",
	"ca xe", "NẰM LẠI", "CP KHÁC", "CỘNG TIỀN HÀNG",
	"TIỀN THUẾ GTGT", "TỔNG THANH TOÁN", "GHI CHÚ (CƯỚC)",
}

func toFloat(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if trimmed := strings.TrimSpace(v); trimmed != "" && trimmed != "-" {
			if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
				return f
			}
		}
		return 0
	default:
		return 0
	}
}

type StylesShip struct {
	HeaderStyle       int
	MoneyStyle        int
	BoldMoneyStyle    int
	CenterStyle       int
	ElevenBold        int
	Eleven            int
	TwelveBold        int
	Twelve            int
	TitleStyle        int
	Bold              int
	BoldItalic        int
	BoldIatlicRed     int
	NoBorder          int
	Eight             int
	BorderNoLeftRight int
	BoldCenter        int
	BoldItalicCenter  int
	YellowFill        int
	TealFill          int
	SixteenBold       int

	TwentyBoldUnderline int
	LightPinkBold       int
	LightPink           int
}

func initializeStylesShip(f *excelize.File) (*StylesShip, error) {
	styleship := &StylesShip{}

	var err error
	styleship.HeaderStyle, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 8},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#e3edf7"}, Pattern: 1},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %w", err)
	}

	customFmt := "#,##0"
	styleship.MoneyStyle, err = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 8},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		CustomNumFmt: &customFmt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create money style: %w", err)
	}

	styleship.BoldMoneyStyle, err = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true, Size: 8},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		CustomNumFmt: &customFmt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bold money style: %w", err)
	}

	styleship.CenterStyle, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 8},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create center style: %w", err)
	}

	styleship.ElevenBold, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 8},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eleven bold style: %w", err)
	}
	styleship.Eleven, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 11},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eleven bold style: %w", err)
	}
	styleship.TwelveBold, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eleven bold style: %w", err)
	}
	styleship.Twelve, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eleven bold style: %w", err)
	}
	styleship.TitleStyle, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Underline: "single", Size: 20},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eleven bold style: %w", err)
	}
	styleship.Bold, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eleven bold style: %w", err)
	}
	styleship.BoldItalic, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Italic: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eleven bold style: %w", err)
	}
	styleship.BoldIatlicRed, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Italic: true, Color: "FF0000"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eleven bold style: %w", err)
	}
	styleship.NoBorder, err = f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 1},   // No left border
			{Type: "right", Color: "FFFFFF", Style: 1},  // No right border
			{Type: "top", Color: "FFFFFF", Style: 1},    // No top border
			{Type: "bottom", Color: "FFFFFF", Style: 1}, // No bottom border
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.Eight, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 8},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.BorderNoLeftRight, err = f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.BoldCenter, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.BoldItalic, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Italic: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.YellowFill, err = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 8},
		Fill:         excelize.Fill{Type: "pattern", Color: []string{"FFFF00"}, Pattern: 1},
		CustomNumFmt: &customFmt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.TealFill, err = f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#e3edf7"}, Pattern: 1},
		Font: &excelize.Font{Size: 8},
	})
	if err != nil {

		return nil, fmt.Errorf("Failed to create teal fill style: %v\n", err)
	}
	styleship.SixteenBold, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.TwentyBoldUnderline, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 20, Underline: "single"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#e3edf7"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.LightPinkBold, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#f5dedd"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}
	styleship.LightPink, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 8},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#f5dedd"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		CustomNumFmt: &customFmt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create no border style: %w", err)
	}

	return styleship, nil
}
func setHeader1Section(f *excelize.File, sheetName string, styleship *StylesShip) error {
	headerInfo := []struct {
		cell  string
		value string
		merge string
		style int
	}{
		{"F1", "CÔNG TY CỔ PHẦN LIÊN KẾT VẬN TẢI THÔNG MINH", "I1", styleship.Eleven},
		{"F2", "( Smart Link Jsc)", "G2", styleship.TwelveBold},
		{"F3", "Địa chỉ: Số 6 tổ 1 Nam Sơn, Phường Đằng Giang, Quận Ngô Quyền, Thành phố Hải Phòng", "K3", styleship.Twelve},
		{"F4", "VPĐD: Số 19 Bến Láng, Trung Hành 5, Phường Đằng Lâm, Quận Hải An, Thành phố Hải Phòng", "L4", styleship.Twelve},
		{"F5", "ĐT: 0948898868  Email: Smartlink.haiphong@gmail.com", "I5", styleship.Eleven},
	}

	for _, info := range headerInfo {
		if err := f.SetCellValue(sheetName, info.cell, info.value); err != nil {
			return fmt.Errorf("failed to set cell %s: %w", info.cell, err)
		}
		if info.merge != "" {
			if err := f.MergeCell(sheetName, info.cell, info.merge); err != nil {
				return fmt.Errorf("failed to merge cells %s to %s: %w", info.cell, info.merge, err)
			}
		}
		if err := f.SetCellStyle(sheetName, info.cell, info.cell, info.style); err != nil {
			return fmt.Errorf("failed to style cell %s: %w", info.cell, err)
		}
	}
	return nil
}
func setTitle1(f *excelize.File, sheetName string, styleship *StylesShip) error {
	if err := f.SetCellValue(sheetName, "A7", "BẢNG KÊ SẢN LƯỢNG VẬN CHUYỂN"); err != nil {
		return fmt.Errorf("failed to set cell A4: %w", err)
	}
	if err := f.MergeCell(sheetName, "A7", "Y7"); err != nil {
		return fmt.Errorf("failed to merge cells A7 to I7: %w", err)
	}
	if err := f.SetCellStyle(sheetName, "A7", "Y7", styleship.TitleStyle); err != nil {
		return fmt.Errorf("failed to style cell A4: %w", err)
	}
	f.SetCellValue(sheetName, "A9", "Đối tác")
	f.SetCellValue(sheetName, "A10", "Tên viết tắt")
	f.SetCellValue(sheetName, "A11", "Người phụ trách")
	f.SetCellValue(sheetName, "G9", "CHI NHÁNH CÔNG TY CỔ PHẦN GIAO NHẬN VẬN TẢI CON ONG")
	f.SetCellValue(sheetName, "G10", "BEEHAN")
	f.SetCellValue(sheetName, "A12", "Hai bên cùng xác nhận sản lượng vận chuyển từ ngày")
	f.SetCellValue(sheetName, "Y12", "25-12-2024")
	f.SetCellValue(sheetName, "H13", "Cước")
	f.SetCellValue(sheetName, "R13", "CHI HỘ")
	f.SetCellStyle(sheetName, "A8", "W11", styleship.NoBorder)
	f.SetCellStyle(sheetName, "A13", "W13", styleship.BorderNoLeftRight)
	f.MergeCell(sheetName, "A12", "X12")
	f.MergeCell(sheetName, "H13", "O13")
	f.MergeCell(sheetName, "R13", "W13")
	f.SetCellStyle(sheetName, "A12", "A12", styleship.BoldItalic)
	f.SetCellStyle(sheetName, "G9", "G10", styleship.Bold)
	f.SetCellStyle(sheetName, "Y12", "Y12", styleship.BoldIatlicRed)
	f.SetCellStyle(sheetName, "A13", "G13", styleship.NoBorder)
	f.SetCellStyle(sheetName, "H13", "H13", styleship.BoldCenter)
	f.SetCellStyle(sheetName, "R13", "R13", styleship.BoldItalic)

	return nil
}
func setFooter1Section(f *excelize.File, sheetName string, styleship *StylesShip, totalRow int) error {
	// Data rows and summary
	summaryRow := totalRow + 1
	if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", summaryRow), "TỔNG CỘNG PHẢI TRẢ"); err != nil {
		return fmt.Errorf("failed to set summary label in A%d: %w", summaryRow, err)
	}
	if err := f.MergeCell(sheetName, fmt.Sprintf("A%d", summaryRow), fmt.Sprintf("G%d", summaryRow)); err != nil {
		return fmt.Errorf("failed to merge cells A%d to G%d: %w", summaryRow, summaryRow, err)
	}
	if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", summaryRow), fmt.Sprintf("G%d", summaryRow), styleship.ElevenBold); err != nil {
		return fmt.Errorf("failed to apply bold style to A%d:G%d: %w", summaryRow, summaryRow, err)
	}

	// Bank information
	bankRow1 := summaryRow + 1
	if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", bankRow1), "TK 123956865 - ACB - CÔNG TY CỔ PHẦN LIÊN KẾT VẬN TẢI THÔNG MINH"); err != nil {
		return fmt.Errorf("failed to set bank info in A%d: %w", bankRow1, err)
	}
	f.MergeCell(sheetName, fmt.Sprintf("A%d", bankRow1), fmt.Sprintf("X%d", bankRow1))

	if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", bankRow1), fmt.Sprintf("A%d", bankRow1), styleship.TealFill); err != nil {
		return fmt.Errorf("failed to apply style to A%d: %w", bankRow1, err)
	}

	bankRow2 := bankRow1 + 1
	if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", bankRow2), "TK 123956865 - ACB - GIANG ĐỨC HIẾU"); err != nil {
		return fmt.Errorf("failed to set bank info in A%d: %w", bankRow2, err)
	}
	f.MergeCell(sheetName, fmt.Sprintf("A%d", bankRow2), fmt.Sprintf("X%d", bankRow2))
	if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", bankRow2), fmt.Sprintf("A%d", bankRow2), styleship.TealFill); err != nil {
		return fmt.Errorf("failed to apply style to A%d: %w", bankRow2, err)
	}

	// Net total in column Y
	if err := f.SetCellFloat(sheetName, fmt.Sprintf("Y%d", bankRow2), 34776000, 0, 64); err != nil {
		return fmt.Errorf("failed to set net total in Y%d: %w", bankRow2, err)
	}
	if err := f.SetCellStyle(sheetName, fmt.Sprintf("Y%d", bankRow2), fmt.Sprintf("Y%d", bankRow2), styleship.BoldMoneyStyle); err != nil {
		return fmt.Errorf("failed to apply bold money style to Y%d: %w", bankRow2, err)
	}

	// Signature section
	signatureRow := bankRow2 + 2
	if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", signatureRow), "XÁC NHẬN CỦA KHÁCH HÀNG                                    XÁC NHẬN CỦA SMARTLINK"); err != nil {
		return fmt.Errorf("failed to set signature in A%d: %w", signatureRow, err)
	}
	if err := f.MergeCell(sheetName, fmt.Sprintf("A%d", signatureRow), fmt.Sprintf("Y%d", signatureRow)); err != nil {
		return fmt.Errorf("failed to merge cells A%d to Y%d: %w", signatureRow, signatureRow, err)
	}
	if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", signatureRow), fmt.Sprintf("Y%d", signatureRow), styleship.SixteenBold); err != nil {
		return fmt.Errorf("failed to apply bold style to A%d:Y%d: %w", signatureRow, signatureRow, err)
	}

	// Add yellow highlight for signature areas (simplified approximation)
	signatureCol1 := 1  // Column A
	signatureCol2 := 15 // Column O (approximate middle for "Customer Confirmation")
	signatureCol3 := 19 // Column S (approximate for "Smartlink Confirmation")
	if err := f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", 'A'+signatureCol1, signatureRow+1), fmt.Sprintf("%c%d", 'A'+signatureCol1, signatureRow+1), styleship.YellowFill); err != nil {
		return fmt.Errorf("failed to apply yellow fill to signature area 1 in %d: %w", signatureRow+1, err)
	}
	if err := f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", 'A'+signatureCol2, signatureRow+1), fmt.Sprintf("%c%d", 'A'+signatureCol2, signatureRow+1), styleship.YellowFill); err != nil {
		return fmt.Errorf("failed to apply yellow fill to signature area 2 in %d: %w", signatureRow+1, err)
	}
	if err := f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", 'A'+signatureCol3, signatureRow+1), fmt.Sprintf("%c%d", 'A'+signatureCol3, signatureRow+1), styleship.YellowFill); err != nil {
		return fmt.Errorf("failed to apply yellow fill to signature area 3 in %d: %w", signatureRow+1, err)
	}

	return nil
}
func SetHeader2Section(f *excelize.File, sheetName string, styleship *StylesShip) error {
	f.SetCellValue(sheet2, "A3", "BẢNG KÊ SẢN LƯỢNG VẬN CHUYỂN")

	f.MergeCell(sheet2, "A3", "W3")

	f.SetCellStyle(sheet2, "A7", "Y7", styleship.TwentyBoldUnderline)

	f.SetCellValue(sheet2, "A5", "Nhà xe")
	f.SetCellValue(sheet2, "A6", "Tên viết tắt")
	f.SetCellValue(sheet2, "A7", "Người phụ trách")
	f.SetCellValue(sheet2, "F5", "CHI NHÁNH CÔNG TY CỔ PHẦN GIAO NHẬN VẬN TẢI CON ONG")
	f.SetCellValue(sheet2, "F6", "BEEHAN")
	f.SetCellValue(sheet2, "A8", "Hai bên cùng xác nhận sản lượng vận chuyển từ ngày")
	f.SetCellValue(sheet2, "W8", "25-12-2024")
	f.SetCellValue(sheet2, "G11", "CƯỚC")

	f.MergeCell(sheet2, "A6", "B6")
	f.MergeCell(sheet2, "A7", "B7")
	f.MergeCell(sheet2, "F5", "P5")
	f.MergeCell(sheet2, "F6", "V6")
	f.MergeCell(sheet2, "A8", "V8")
	f.MergeCell(sheet2, "G11", "M11")
	f.SetCellStyle(sheet2, "A3", "W3", styleship.TwentyBoldUnderline)
	f.SetCellStyle(sheet2, "A8", "A8", styleship.BoldItalic)
	f.SetCellStyle(sheet2, "F5", "P5", styleship.Bold)
	f.SetCellStyle(sheet2, "F6", "V6", styleship.Bold)
	f.SetCellStyle(sheet2, "W8", "W8", styleship.BoldIatlicRed)
	f.SetCellStyle(sheet2, "A4", "W7", styleship.NoBorder)
	f.SetCellStyle(sheet2, "A11", "F11", styleship.BorderNoLeftRight)
	f.SetCellStyle(sheet2, "G11", "M11", styleship.LightPinkBold)
	f.SetCellStyle(sheet2, "N11", "N11", styleship.LightPink)
	f.SetCellStyle(sheet2, "O11", "O11", styleship.YellowFill)

	return nil
}
func ExportPaymentRequestExcelv3(filename string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()
	sheet1 := "Bảng Kê Thu"
	sheet2 := "Bảng Kê Chi"

	// Tạo sheet mới cho Sheet2
	f.NewSheet(sheet2)
	f.SetDefaultFont(defaultFontShip)
	// Đặt tên cho Sheet1 (mặc định là Sheet1)
	f.SetSheetName("Sheet1", sheet1)

	// Định nghĩa các style
	styleship, err := initializeStylesShip(f)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for col, width := range colWidths1 {
		if err := f.SetColWidth(sheet1, col, col, width); err != nil {
			return fmt.Errorf("failed to set width for column %s: %w", col, err)
		}
	}
	for col, width := range colWidths2 {
		if err := f.SetColWidth(sheet2, col, col, width); err != nil {
			return fmt.Errorf("failed to set width for column %s: %w", col, err)
		}
	}
	if err := setHeader1Section(f, sheet1, styleship); err != nil {
		fmt.Println("Error setting header section:", err)
		return nil
	}

	if err := setTitle1(f, sheet1, styleship); err != nil {
		fmt.Println("Error setting title:", err)
		return nil
	}
	if err := SetHeader2Section(f, sheet2, styleship); err != nil {
		fmt.Println("Error setting title:", err)
		return nil
	}
	for colIdx, header := range tableHeadersTotal {
		cell := fmt.Sprintf("%c%d", 'A'+colIdx, 14)
		f.SetCellValue(sheet1, cell, header)
		f.SetCellStyle(sheet1, cell, cell, styleship.HeaderStyle)
	}

	// Điền header cho Sheet2
	for colIdx2, header := range TableHeadersFee {
		cell := fmt.Sprintf("%c%d", 'A'+colIdx2, 12)
		f.SetCellValue(sheet2, cell, header)
		f.SetCellStyle(sheet2, cell, cell, styleship.HeaderStyle)
	}

	data := [][]interface{}{
		{1, "26/11/2024", "KH.26112024.13", "", "GAOU6765765", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H 07303", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{2, "27/11/2024", "KH.27112024.41", "", "TXGU6839401", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C17786", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{3, "27/11/2024", "KH.27112024.42", "", "TLLU8616128", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "29C71582", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{4, "28/11/2024", "KH.28112024.50", "", "TXGU6554190", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H08479", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{5, "29/11/2024", "KH.29112024.31", "", "TSSU5176234", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C22535", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{6, "29/11/2024", "KH.29112024.32", "", "FFAU4382558", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C03270", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{7, "03/12/2024", "KH.03122024.31", "", "TLLU8699441", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C13956", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{8, "03/12/2024", "KH.03122024.32", "", "CAAU7626926", "40HC", "CANON THĂNG LONG, HÀ NỘI - HẢI PHÒNG", "15H03856", 2500000, "", "", "", "", "", 2500000, 200000, 2700000, "", "", "", "", "", "", ""},
		{9, "04/12/2024", "KH.04122024.64", "", "TRHU8855824", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H05864", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{10, "05/12/2024", "KH.06122024.14", "", "TSSU5061847", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C15654", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{11, "05/12/2024", "KH.06122024.16", "", "TSSU5154385", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H05147", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{12, "05/12/2024", "KH.06122024.19", "", "TSSU5187917", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C20649", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{13, "06/12/2024", "KH.06122024.37", "", "GAOU7078193", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15F-01759", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{14, "06/12/2024", "KH.06122024.38", "", "BHCU5015964", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15F-01721", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{15, "06/12/2024", "KH.06122024.39", "", "GAOU6728421", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "29H-77080", 2700000, "", "", 200000, "", "", 2900000, 232000, 3132000, "", "", "", "", "", "", ""},
		{16, "06/12/2024", "KH.06122024.40", "", "TSSU5210581", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H07388", 2700000, "", "", 200000, "", "", 2900000, 232000, 3132000, "", "", "", "", "", "", ""},
		{17, "06/12/2024", "KH.06122024.41", "", "BEAU6454029", "40HC", "PHỐ NỐI, HƯNG YÊN - HẢI PHÒNG", "15C31373", 2200000, "", "", "", "", "", 2200000, 176000, 2376000, "", "", "", "", "", "", ""},
		{18, "07/12/2024", "KH.07122024.22", "", "TSSU5097855", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C 14626", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{19, "07/12/2024", "KH.07122024.24", "", "TXGU8264060", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H07778", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{20, "08/12/2024", "KH.09122024.13", "", "TXGU5812767", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "99E-01867", 2700000, "", "", 200000, "", "", 2900000, 232000, 3132000, "", "", "", "", "", "", ""},
		{21, "09/12/2024", "KH.09122024.18", "", "TXGU6906346", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C31806", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{22, "09/12/2024", "CK.10122024.3", "", "BMOU5006575", "40HC", "HẢI PHÒNG - QUẾ VÕ, BẮC NINH", "15C32924", 3000000, "", "", "", 200000, "", 3200000, 256000, 3456000, "", "", "", "", "", "", ""},
		{23, "10/12/2024", "KH.10122024.38", "", "TLLU8591610", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "29H76185", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{24, "10/12/2024", "KH.10122024.39", "", "GCXU6517842", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C14207", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{25, "10/12/2024", "KH.10122024.40", "", "TXGU7360706", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C04028", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{26, "11/12/2024", "KH.11122024.35", "", "GAOU6961979", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C14996", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{27, "11/12/2024", "KH.11122024.36", "", "TSSU5159685", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H06828", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{28, "17/12/2024", "KH.17122024.43", "", "GAOU6732653", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C31169", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{29, "17/12/2024", "KH.17122024.44", "", "TSSU5184883", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H00442", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{30, "17/12/2024", "CK.18122024.11-1", "", "NLLU4149610", "40HC", "ĐẠI TỪ, THÁI NGUYÊN - HẢI PHÒNG", "15H07954", 4500000, "", "", "", 500000, "", 5000000, 400000, 5400000, "", "", "", "", "", "", ""},
		{31, "18/12/2024", "KH.18122024.52", "", "GAOU6769919", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H01822", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{32, "18/12/2024", "KH.18122024.53", "", "CAAU8038444", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H08421", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{33, "18/12/2024", "KH.18122024.54", "", "FFAU4394380", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H00521", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{34, "19/12/2024", "KH.19122024.27", "", "GAOU7010721", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C-14420", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
		{35, "19/12/2024", "CK.19122024.53", "", "CSNU8501462", "40HC", "ĐẠI TỪ, THÁI NGUYÊN - HẢI PHÒNG", "29K041.53", 4500000, "", "", "", 652000, "", 5152000, 412160, 5564160, "", "", "", "", "", "", ""},
		{36, "20/12/2024", "KH.20122024.29", "", "SEKU4083291", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C21249", 2700000, "", "", "", "", "", 2700000, 216000, 2916000, "", "", "", "", "", "", ""},
	}

	// Tính tổng cho Sheet1
	var totalCucVC, totalLachHuyen, totalCaXe, totalNamLai, totalCPKhac, totalPhiDichVu, totalCongTienHang, totalTienThueGTGT, totalTongThanhToan, totalCSHT, totalNang, totalHa, totalCPKhac2, totalCongChiHo, totalTongThu float64
	for _, row := range data {
		totalCucVC += toFloat(row[8])
		totalLachHuyen += toFloat(row[9])
		totalCaXe += toFloat(row[10])
		totalNamLai += toFloat(row[11])
		totalCPKhac += toFloat(row[12])
		totalPhiDichVu += toFloat(row[13])
		totalCongTienHang += toFloat(row[14])
		totalTienThueGTGT += toFloat(row[15])
		totalTongThanhToan += toFloat(row[16])
		totalCSHT += toFloat(row[17])
		totalNang += toFloat(row[18])
		totalHa += toFloat(row[19])
		totalCPKhac2 += toFloat(row[20])
		totalCongChiHo += toFloat(row[21])
		totalTongThu += toFloat(row[22])
	}

	// Thêm dòng tổng cho Sheet1

	// Populate data for Sheet1
	for rowIdx, row := range data {
		for colIdx, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, dataStartRowShip+rowIdx)
			if colIdx >= 8 && colIdx <= 16 { // Numeric columns (CƯỚC VC to TỔNG THANH TOÁN)
				if floatVal, ok := value.(float64); ok {
					if err := f.SetCellFloat(sheet1, cell, floatVal, 0, 64); err != nil {
						return fmt.Errorf("failed to set float value in cell %s: %w", cell, err)
					}
				} else if strVal, ok := value.(string); ok {
					if floatVal, err := strconv.ParseFloat(strings.TrimSpace(strVal), 64); err == nil {
						if err := f.SetCellFloat(sheet1, cell, floatVal, 0, 64); err != nil {
							return fmt.Errorf("failed to parse and set float value in cell %s: %w", cell, err)
						}
					} else {
						if err := f.SetCellValue(sheet1, cell, value); err != nil {
							return fmt.Errorf("failed to set value in cell %s: %w", cell, err)
						}
					}
				} else {
					if err := f.SetCellValue(sheet1, cell, value); err != nil {
						return fmt.Errorf("failed to set value in cell %s: %w", cell, err)
					}
				}
			} else {
				if err := f.SetCellValue(sheet1, cell, value); err != nil {
					return fmt.Errorf("failed to set value in cell %s: %w", cell, err)
				}
			}
			// Apply styles
			if colIdx == 1 || colIdx == 6 { // NGÀY and NƠI ĐÓNG / TRẢ
				if err := f.SetCellStyle(sheet1, cell, cell, styleship.CenterStyle); err != nil {
					return fmt.Errorf("failed to apply center style to cell %s: %w", cell, err)
				}
			} else if colIdx == 14 { // CỘNG TIỀN HÀNG
				if err := f.SetCellStyle(sheet1, cell, cell, styleship.LightPink); err != nil {
					return fmt.Errorf("failed to apply peach money style to cell %s: %w", cell, err)
				}
			} else if colIdx == 15 { // TIỀN THUẾ GTGT
				if err := f.SetCellStyle(sheet1, cell, cell, styleship.LightPink); err != nil {
					return fmt.Errorf("failed to apply peach money style to cell %s: %w", cell, err)
				}
			} else if colIdx == 16 { // TỔNG THANH TOÁN
				if err := f.SetCellStyle(sheet1, cell, cell, styleship.YellowFill); err != nil {
					return fmt.Errorf("failed to apply yellow money style to cell %s: %w", cell, err)
				}
			} else if colIdx >= 8 && colIdx <= 16 { // Numeric columns
				if err := f.SetCellStyle(sheet1, cell, cell, styleship.MoneyStyle); err != nil {
					return fmt.Errorf("failed to apply money style to cell %s: %w", cell, err)
				}
			} else {
				if err := f.SetCellStyle(sheet1, cell, cell, styleship.CenterStyle); err != nil {
					return fmt.Errorf("failed to apply center style to cell %s: %w", cell, err)
				}
			}
		}
	}

	// Add total row for Sheet1
	totalRow := dataStartRowShip + len(data)
	if err := f.SetCellValue(sheet1, fmt.Sprintf("A%d", totalRow), "TỔNG"); err != nil {
		return fmt.Errorf("failed to set total label in A%d: %w", totalRow, err)
	}
	if err := f.MergeCell(sheet1, fmt.Sprintf("A%d", totalRow), fmt.Sprintf("H%d", totalRow)); err != nil {
		return fmt.Errorf("failed to merge cells A%d to N%d: %w", totalRow, totalRow, err)
	}
	if err := f.SetCellStyle(sheet1, fmt.Sprintf("A%d", totalRow), fmt.Sprintf("H%d", totalRow), styleship.ElevenBold); err != nil {
		return fmt.Errorf("failed to apply bold style to A%d:N%d: %w", totalRow, totalRow, err)
	}

	f.SetCellFloat(sheet1, fmt.Sprintf("I%d", totalRow), totalCucVC, 0, 64)
	f.SetCellFloat(sheet1, fmt.Sprintf("L%d", totalRow), totalNamLai, 0, 64)
	f.SetCellFloat(sheet1, fmt.Sprintf("M%d", totalRow), totalCPKhac, 0, 64)
	if err := f.SetCellStyle(sheet1, fmt.Sprintf("O%d", totalRow), fmt.Sprintf("O%d", totalRow), styleship.BoldMoneyStyle); err != nil {
		return fmt.Errorf("failed to apply bold money style to O%d: %w", totalRow, err)
	}
	f.SetCellStyle(sheet1, fmt.Sprintf("I%d", totalRow), fmt.Sprintf("I%d", totalRow), styleship.BoldMoneyStyle)
	f.SetCellStyle(sheet1, fmt.Sprintf("L%d", totalRow), fmt.Sprintf("L%d", totalRow), styleship.BoldMoneyStyle)
	f.SetCellStyle(sheet1, fmt.Sprintf("M%d", totalRow), fmt.Sprintf("M%d", totalRow), styleship.BoldMoneyStyle)
	f.SetCellStyle(sheet1, fmt.Sprintf("X%d", totalRow), fmt.Sprintf("X%d", totalRow), styleship.BoldMoneyStyle)
	if err := f.SetCellFloat(sheet1, fmt.Sprintf("P%d", totalRow), totalTienThueGTGT, 0, 64); err != nil {
		return fmt.Errorf("failed to set Tien Thue GTGT total in P%d: %w", totalRow, err)
	}
	if err := f.SetCellStyle(sheet1, fmt.Sprintf("P%d", totalRow), fmt.Sprintf("P%d", totalRow), styleship.BoldMoneyStyle); err != nil {
		return fmt.Errorf("failed to apply bold money style to P%d: %w", totalRow, err)
	}

	if err := f.SetCellStyle(sheet1, fmt.Sprintf("Q%d", totalRow), fmt.Sprintf("Q%d", totalRow), styleship.BoldMoneyStyle); err != nil {
		return fmt.Errorf("failed to apply bold money style to Q%d: %w", totalRow, err)
	}
	setFooter1Section(f, sheet1, styleship, totalRow)

	dataSheet2 := [][]interface{}{
		{1, "26/11/2024", "KH.26112024.13", "GAOU6765765", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H 07303", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{2, "27/11/2024", "KH.27112024.41", "TXGU6839401", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C17786", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{3, "27/11/2024", "KH.27112024.42", "TLLU8616128", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "29C71582", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{4, "28/11/2024", "KH.28112024.50", "TXGU6554190", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H08479", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{5, "29/11/2024", "KH.29112024.31", "TSSU5176234", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C22535", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{6, "29/11/2024", "KH.29112024.32", "FFAU4382558", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C03270", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{7, "03/12/2024", "KH.03122024.31", "TLLU8699441", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C13956", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{8, "03/12/2024", "KH.03122024.32", "CAAU7626926", "40HC", "CANON THĂNG LONG, HÀ NỘI - HẢI PHÒNG", "15H03856", 2500000, "", "", "", "", 2500000, 200000, 2700000, ""},
		{9, "04/12/2024", "KH.04122024.64", "TRHU8855824", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H05864", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{10, "05/12/2024", "KH.06122024.14", "TSSU5061847", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C15654", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{11, "05/12/2024", "KH.06122024.16", "TSSU5154385", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H05147", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{12, "05/12/2024", "KH.06122024.19", "TSSU5187917", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C20649", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{13, "06/12/2024", "KH.06122024.37", "GAOU7078193", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15F-01759", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{14, "06/12/2024", "KH.06122024.38", "BHCU5015964", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15F-01721", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{15, "06/12/2024", "KH.06122024.39", "GAOU6728421", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "29H-77080", 2700000, "", "", 200000, "", 2900000, 232000, 3132000, ""},
		{16, "06/12/2024", "KH.06122024.40", "TSSU5210581", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H07388", 2700000, "", "", 200000, "", 2900000, 232000, 3132000, ""},
		{17, "06/12/2024", "KH.06122024.41", "BEAU6454029", "40HC", "PHỐ NỐI, HƯNG YÊN - HẢI PHÒNG", "15C31373", 2200000, "", "", "", "", 2200000, 176000, 2376000, ""},
		{18, "07/12/2024", "KH.07122024.22", "TSSU5097855", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C 14626", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{19, "07/12/2024", "KH.07122024.24", "TXGU8264060", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H07778", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{20, "08/12/2024", "KH.09122024.13", "TXGU5812767", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "99E-01867", 2700000, "", "", 200000, "", 2900000, 232000, 3132000, ""},
		{21, "09/12/2024", "KH.09122024.18", "TXGU6906346", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C31806", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{22, "09/12/2024", "CK.10122024.3", "BMOU5006575", "40HC", "HẢI PHÒNG - QUẾ VÕ, BẮC NINH", "15C32924", 3000000, "", "", "", 200000, 3200000, 256000, 3456000, ""},
		{23, "10/12/2024", "KH.10122024.38", "TLLU8591610", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "29H76185", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{24, "10/12/2024", "KH.10122024.39", "GCXU6517842", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C14207", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{25, "10/12/2024", "KH.10122024.40", "TXGU7360706", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C04028", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{26, "11/12/2024", "KH.11122024.35", "GAOU6961979", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C14996", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{27, "11/12/2024", "KH.11122024.36", "TSSU5159685", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H06828", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{28, "17/12/2024", "KH.17122024.43", "GAOU6732653", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C31169", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{29, "17/12/2024", "KH.17122024.44", "TSSU5184883", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H00442", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{30, "17/12/2024", "CK.18122024.11-1", "NLLU4149610", "40HC", "ĐẠI TỪ, THÁI NGUYÊN - HẢI PHÒNG", "15H07954", 4500000, "", "", "", 500000, 5000000, 400000, 5400000, ""},
		{31, "18/12/2024", "KH.18122024.52", "GAOU6769919", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H01822", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{32, "18/12/2024", "KH.18122024.53", "CAAU8038444", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H08421", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{33, "18/12/2024", "KH.18122024.54", "FFAU4394380", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15H00521", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{34, "19/12/2024", "KH.19122024.27", "GAOU7010721", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C-14420", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
		{35, "19/12/2024", "CK.19122024.53", "CSNU8501462", "40HC", "ĐẠI TỪ, THÁI NGUYÊN - HẢI PHÒNG", "29K041.53", 4500000, "", "", "", 652000, 5152000, 412160, 5564160, ""},
		{36, "20/12/2024", "KH.20122024.29", "SEKU4083291", "40HC", "PHONG KHÊ, BẮC NINH - HẢI PHÒNG", "15C21249", 2700000, "", "", "", "", 2700000, 216000, 2916000, ""},
	}
	// Điền dữ liệu vào Sheet2
	for rowIdx, row := range dataSheet2 {
		for colIdx, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, dataStartRowShip+rowIdx-2)
			if colIdx >= 7 && colIdx <= 14 {
				if floatVal, ok := value.(float64); ok {
					if err := f.SetCellFloat(sheet2, cell, floatVal, 0, 64); err != nil {
						return fmt.Errorf("failed to set float value in cell %s: %w", cell, err)
					}
				} else if strVal, ok := value.(string); ok {
					if floatVal, err := strconv.ParseFloat(strings.TrimSpace(strVal), 64); err == nil {
						if err := f.SetCellFloat(sheet2, cell, floatVal, 0, 64); err != nil {
							return fmt.Errorf("failed to parse and set float value in cell %s: %w", cell, err)
						}
					} else {
						if err := f.SetCellValue(sheet2, cell, value); err != nil {
							return fmt.Errorf("failed to set value in cell %s: %w", cell, err)
						}
					}
				} else {
					if err := f.SetCellValue(sheet2, cell, value); err != nil {
						return fmt.Errorf("failed to set value in cell %s: %w", cell, err)
					}
				}
			} else {
				if err := f.SetCellValue(sheet2, cell, value); err != nil {
					return fmt.Errorf("failed to set value in cell %s: %w", cell, err)
				}
			}
			if colIdx == 1 || colIdx == 5 {
				if err := f.SetCellStyle(sheet2, cell, cell, styleship.CenterStyle); err != nil {
					return fmt.Errorf("failed to apply center style to cell %s: %w", cell, err)
				}
			} else if colIdx == 12 { // CỘNG TIỀN HÀNG
				if err := f.SetCellStyle(sheet2, cell, cell, styleship.LightPink); err != nil {
					return fmt.Errorf("failed to apply peach money style to cell %s: %w", cell, err)
				}
			} else if colIdx == 13 { // TIỀN THUẾ GTGT
				if err := f.SetCellStyle(sheet2, cell, cell, styleship.LightPink); err != nil {
					return fmt.Errorf("failed to apply peach money style to cell %s: %w", cell, err)
				}
			} else if colIdx == 14 { // TỔNG THANH TOÁN
				if err := f.SetCellStyle(sheet2, cell, cell, styleship.YellowFill); err != nil {
					return fmt.Errorf("failed to apply yellow money style to cell %s: %w", cell, err)
				}
			} else if colIdx >= 7 && colIdx <= 14 {
				if err := f.SetCellStyle(sheet2, cell, cell, styleship.MoneyStyle); err != nil {
					return fmt.Errorf("failed to apply money style to cell %s: %w", cell, err)
				}
			} else {
				if err := f.SetCellStyle(sheet2, cell, cell, styleship.CenterStyle); err != nil {
					return fmt.Errorf("failed to apply center style to cell %s: %w", cell, err)
				}
			}
		}
	}
	var totalCucVCSheet2, totalCPKhacSheet2, totalCongTienHangSheet2, totalTienThueGTGTSheet2, totalTongThanhToanSheet2, totalNamLaiSheet2 float64
	for _, row := range dataSheet2 {
		totalCucVCSheet2 += toFloat(row[7])
		totalCPKhacSheet2 += toFloat(row[11])
		totalCongTienHangSheet2 += toFloat(row[12])
		totalTienThueGTGTSheet2 += toFloat(row[13])
		totalTongThanhToanSheet2 += toFloat(row[14])
		totalNamLaiSheet2 += toFloat(row[10])
	}
	totalRowSheet2 := dataStartRowShip2 + len(dataSheet2)
	if err := f.SetCellValue(sheet2, fmt.Sprintf("A%d", totalRowSheet2), "TỔNG"); err != nil {
		return fmt.Errorf("failed to set total label in A%d: %w", totalRowSheet2, err)
	}
	if err := f.MergeCell(sheet2, fmt.Sprintf("A%d", totalRowSheet2), fmt.Sprintf("G%d", totalRowSheet2)); err != nil {
		return fmt.Errorf("failed to merge cells A%d to N%d: %w", totalRowSheet2, totalRowSheet2, err)
	}
	if err := f.SetCellStyle(sheet2, fmt.Sprintf("A%d", totalRowSheet2), fmt.Sprintf("G%d", totalRowSheet2), styleship.ElevenBold); err != nil {
		return fmt.Errorf("failed to apply bold style to A%d:N%d: %w", totalRowSheet2, totalRowSheet2, err)
	}
	f.SetCellFloat(sheet2, fmt.Sprintf("H%d", totalRowSheet2), totalCucVCSheet2, 0, 64)

	f.SetCellStyle(sheet2, fmt.Sprintf("H%d", totalRowSheet2), fmt.Sprintf("H%d", totalRowSheet2), styleship.BoldMoneyStyle)

	f.SetCellFloat(sheet2, fmt.Sprintf("K%d", totalRowSheet2), totalNamLaiSheet2, 0, 64)
	f.SetCellStyle(sheet2, fmt.Sprintf("K%d", totalRowSheet2), fmt.Sprintf("K%d", totalRowSheet2), styleship.BoldMoneyStyle)

	f.SetCellFloat(sheet2, fmt.Sprintf("L%d", totalRowSheet2), totalCPKhacSheet2, 0, 64)

	f.SetCellStyle(sheet2, fmt.Sprintf("L%d", totalRowSheet2), fmt.Sprintf("L%d", totalRowSheet2), styleship.BoldMoneyStyle)

	f.SetCellFloat(sheet2, fmt.Sprintf("N%d", totalRowSheet2), totalTienThueGTGTSheet2, 0, 64)
	f.SetCellStyle(sheet2, fmt.Sprintf("N%d", totalRowSheet2), fmt.Sprintf("N%d", totalRowSheet2), styleship.BoldMoneyStyle)

	//Thêm dòng tổng cho Sheet2
	// totalRowSheet2 := len(data) + 2
	// f.SetCellValue(sheet2, fmt.Sprintf("A%d", totalRowSheet2), "TỔNG")
	// f.MergeCell(sheet2, fmt.Sprintf("A%d", totalRowSheet2), fmt.Sprintf("M%d", totalRowSheet2))
	// f.SetCellStyle(sheet2, fmt.Sprintf("A%d", totalRowSheet2), fmt.Sprintf("M%d", totalRowSheet2), styleship.ElevenBold)
	// f.SetCellFloat(sheet2, fmt.Sprintf("N%d", totalRowSheet2), totalCongTienHang, 0, 64)
	// f.SetCellStyle(sheet2, fmt.Sprintf("N%d", totalRowSheet2), fmt.Sprintf("N%d", totalRowSheet2), styleship.BoldMoneyStyle)
	// f.SetCellFloat(sheet2, fmt.Sprintf("O%d", totalRowSheet2), totalTienThueGTGT, 0, 64)
	// f.SetCellStyle(sheet2, fmt.Sprintf("O%d", totalRowSheet2), fmt.Sprintf("O%d", totalRowSheet2), styleship.BoldMoneyStyle)
	// f.SetCellFloat(sheet2, fmt.Sprintf("P%d", totalRowSheet2), totalTongThanhToan, 0, 64)
	// f.SetCellStyle(sheet2, fmt.Sprintf("P%d", totalRowSheet2), fmt.Sprintf("P%d", totalRowSheet2), styleship.BoldMoneyStyle)

	// Đặt Sheet1 là sheet mặc định
	f.SetActiveSheet(0)

	// Lưu file
	if err := f.SaveAs("Bảng kê 1.xlsx"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("File Bảng kê.xlsx đã được tạo thành công!")
	}
	return nil

}
