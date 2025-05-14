package constants

const (
	MISSING_LOGIN_INPUT  = "Thiếu thông tin đăng nhập"
	SUCCESS_LOGIN        = "Đăng nhập thành công"
	ERROR_INTERNAL_ERROR = "Lỗi nội bộ máy chủ"
	INVALID_USERNAME     = "Username không tìm thấy"
	// ERROR_CREATE                                  = "Không thể thêm mới dữ liệu này"
	ERROR_EDIT                 = "Chỉnh sửa thất bại"
	INVALID_PASSWORD           = "Mật khẩu không đúng"
	ACCOUNT_NOT_ACTIVE         = "Tài khoản không được phép hoạt động"
	NOT_ADMIN                  = "Chỉ tài khoản Admin được lấy dữ liệu này"
	ERROR_INPUT                = "Thiếu dữ liệu input yêu cầu"
	ERROR_PARSE_DATA_TO_LOCALS = "Dữ liệu truyền vào không đúng yêu cầu"
	// ERROR_DATA_IS_ALSO                            = "Lỗi dữ liệu đồng thời tồn tại"
	DATA_INPUT_IS_NOT_NUMBER              = "Dữ liệu truyền vào phải là số"
	DATA_INPUT_IS_NOT_BOOL                = "Dữ liệu truyền vào phải là boolean"
	NEW_PASSWORD_NOT_SAME_REPEAT_PASSWORD = "Mật khẩu mới và mật khẩu nhắc lại không trùng khớp"
	NEW_PASSWORD_SAME_CURRENT_PASSWORD    = "Mật khẩu mới không được trùng với mật khẩu hiện tại"
	CAN_NOT_HASH_PASSWORD                 = "Không thể mã hoá được mật khẩu"
	NEW_PASSWORD_LENGTH_INVALID           = "Mật khẩu mới không phù hợp với yêu cầu"
	NOT_FOUND_RECORDS                     = "Không tìm thấy dữ liệu"
	IDENTIFICATION_CARD_EXISTS            = "Số CCCD này đã được đăng kí"
	EMAIL_EXISTS                          = "Email này đã được đăng kí "
	ERROR_PERMISSION_DENIED               = "Bạn không có quyền truy cập vào dữ liệu này"
	// TAX_CODE_EXISTS                               = "Mã số thuế này đã được đăng kí"
	// TAX_CODE_INVALID                              = "Mã số thuế phải gồm các số và có 10 hoặc 14 kí tự"
	PHONE_NUMBER_EXISTS = "Số điện thoại này đã được đăng kí"
	// STAFF_CANNOT_CREATE_WAITING_ACCEPT_ON_REFUND  = "Nhân viên không thể tạo đơn hoàn lương với trạng thái khác chờ duyệt"
	// STAFF_CANNOT_REFUND                           = "Nhân viên không thể hoàn lương"
	// TOTAL_MONEY_REFUND_CANNOT_GREATER_THAN_TOTAL  = "Số tiền hoàn ứng lương không được lớn hơn số tiền đã ứng còn lại"
	// CAN_NOT_EDIT_REFUND_CONFIRM                   = "Không thể chỉnh sửa đơn hoàn ứng lương khi đã được duyệt"
	// STATUS_ADVANCE_EXISTS                         = "Trạng thái ứng lương này đã tồn tại"
	// STATUS_ADVANCE_INVALID                        = "Trạng thái ứng lương không đúng"
	// STATUS_REFUND_INVALID                         = "Trạng thái hoàn ứng không đúng"
	// STAFF_CANNOT_CREATE_WAITING_ACCEPT_ON_ADVANCE = "Nhân viên không thể tạo trạng thái ứng lương khác chờ duyệt"
	// CAN_NOT_EDIT_ADVANCE_CONFIRM                  = "Không thể chỉnh sửa đơn ứng lương khi đã được duyệt"
	// CAN_NOT_ADVANCE_FOR_STAFF_AND_EMPLOYEE        = "Không được ứng cho nhân viên hoặc công nhân trong 1 tiến trình"
	// MISSING_STAFF_OR_EMPLOYEE_IN_ADVANCE          = "Phải chọn đối tượng là nhân viên hoặc công nhân trong 1 tiến trình"
	// STATUS_WORKING_INVALID                        = "Trạng thái làm việc không tồn tại"
	TYPE_WORKING_INVALID  = "Hình thức làm việc không tồn tại"
	TYPE_POSITION_INVALID = "Vị trí làm việc không tồn tại"
	GENDER_INVALID        = "Giới tính không tồn tại"
	// CAN_NOT_EDIT_SALARY_IS_PAYMENT                = "Không thể sửa bảng lương đã thanh toán"
	// PERMISSION_INVALID                            = "Không có quyền thực hiện thao tác này"
	// KEY_INVALID                                   = "Key tìm kiếm không được hỗ trợ"
	// FILE_IMPORT_INVALID                           = "File import không đúng định dạng"
	// MONEY_REFUND_MUST_EQUAL_MONEY_ADVANCE         = "Số tiền hoàn ứng phải bằng số tiền đã ứng"
	// DATA_INPUT_IS_NOT_ARRAY                       = "Dữ liệu truyền vào phải là mảng"
	// SALARY_ADVANCE_IS_REFUND                      = "Đơn ứng lương đã được hoàn đủ"
	REFRESH_TOKEN_NOT_FOUND = "Phiên đăng nhập bị lỗi "
	//
)

var ROLE = []string{"ADMIN", "QUANLY", "KETOAN", "SALE"}
var ROLE_ADMIN = "ADMIN"
var ROLE_QUANLY = "QUANLY"
var ROLE_KETOAN = "KETOAN"
var ROLE_SALE = "SALE"

var TYPE_WORKING = []string{"đang làm", "nghỉ phép ", "công tác ", "nghỉ việc"}

var TYPE_POSITION = []string{"Nhân viên", "Công nhân", "Giám đốc", "QUANLY", "KETOAN", "SALE", "Thư ký ", "Trưởng phòng", "Phó giám đốc", "Phó trưởng phòng"}

var GENDER = []string{"Nam", "Nữ "}

// var COLUMN_OF_KEY_FILTER_EMPLOYEE = map[string]string{
// 	"code":               "code",
// 	"name":               "name",
// 	"identificationCard": "identification_card",
// 	"phoneNumber":        "phone_number",
// }

// var COLUMN_OF_KEY_FILTER_SALARY_ADVANCE = map[string]string{
// 	// "codeStaff":                  "staffs.code",
// 	"codeEmployee":               "employees.code",
// 	"nameStaff":                  "staffs.name",
// 	"nameEmployee":               "employees.name",
// 	"identificationCardStaff":    "staffs.identification_card",
// 	"identificationCardEmployee": "employees.identification_card",
// 	"phoneNumberStaff":           "staffs.phone_number",
// 	"phoneNumberEmployee":        "employees.phone_number",
// }

// var COLUMN_OF_KEY_FILTER_SALARY_REFUND = map[string]string{
// 	// "codeStaff":                  "salary_advances.staffs.code",
// 	"codeEmployee":               "salary_advances.employees.code",
// 	"nameStaff":                  "salary_advances.staffs.name",
// 	"nameEmployee":               "salary_advances.employees.name",
// 	"identificationCardStaff":    "salary_advances.staffs.identification_card",
// 	"identificationCardEmployee": "salary_advances.employees.identification_card",
// 	"phoneNumberStaff":           "salary_advances.staffs.phone_number",
// 	"phoneNumberEmployee":        "salary_advances.employees.phone_number",
// }

// var COLUMN_OF_KEY_FILTER_SALARY = map[string]string{
// 	"codeEmployee":            "employees.code",
// 	"nameStaff":               "staffs.name",
// 	"nameEmployee":            "employees.name",
// 	"identificationCardStaff": "staffs.identification_card",
// }
