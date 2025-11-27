package constant

const (
	EmployeeMobilePattern = `^((\+)?(\d{2}[-])?(\d{10}){1})?(\d{11}){0,1}?$`

	VENDOR   = "VEN"
	CUSTOMER = "CUS"

	MAX_ROW_UPDATE_TRIPSHEET = 200

	// Trip Status
	STATUS_CREATED        = "Created"
	STATUS_SUBMITTED      = "Intransit"
	STATUS_CANCELLED      = "Cancelled"
	STATUS_CLOSED         = "Closed"
	STATUS_DELIVERED      = "Delivered"
	STATUS_DELETED        = "Deleted"
	STATUS_APPROVED       = "Approved"
	STATUS_COMPLETED      = "Completed"
	STATUS_INVOICE_RAISED = "InvoiceRaised"
	STATUS_PAID           = "Paid"

	STATUS_INVOICE_DRAFT          = "Draft"
	STATUS_INVOICE_DRAFT_DESC     = "Pending invoice details"
	STATUS_INVOICE_RAISED_DESC    = "Payment Pending"
	STATUS_INVOICE_CANCELLED_DESC = "Cancelled thr draft invoice"
	STATUS_PAID_DESC              = "Payment Paid"

	TRIP_TYPE_POD     = "pod"
	TRIP_TYPE_REGULAR = "regular"

	//
	UN_LOADING_POINT         = "un_loading_point"
	LOADING_POINT            = "loading_point"
	LOCAL_SCHEDULED_TRIP     = "Local Scheduled Trip"
	LOCAL_ADHOC_TRIP         = "Local Adhoc Trip"
	LINE_HUAL_SCHEDULED_TRIP = "Line haul Scheduled Trip"
	LINE_HUAL_ADHOC_TRIP     = "Line haul Adhoc Trip"

	ONE_WAY_TRIP = "One way"
	ROUND_TRIP   = "Round Trip"

	// Directroy
	POD_DIRECTORY = "pod_all"
)

var STATUS_TO_COLUMN = map[string]string{
	STATUS_CLOSED:    "trip_closed_date",
	STATUS_COMPLETED: "trip_completed_date",
	STATUS_DELIVERED: "trip_delivered_date",
	STATUS_SUBMITTED: "trip_submitted_date",
}

const (
	SESSION_TIMEOUT       = "session_timeout"
	SESSION_TIMEOUT_MS    = "session_timeout_ms"
	SESSION_TIMEOUT_VALUE = "1800000" // 30 minutes in milliseconds

	EMPLOYEE_PERFORMANCE   = "employee_performance"
	EMPLOYEE_PERFORMANCE_1 = "1"
	EMPLOYEE_PERFORMANCE_2 = "2"
	EMPLOYEE_PERFORMANCE_3 = "3"
	EMPLOYEE_PERFORMANCE_4 = "4"
	EMPLOYEE_PERFORMANCE_5 = "5"

	EMPLOYEE_PERFORMANCE_1_NAME = "Excellent"
	EMPLOYEE_PERFORMANCE_2_NAME = "Good"
	EMPLOYEE_PERFORMANCE_3_NAME = "Average"
	EMPLOYEE_PERFORMANCE_4_NAME = "Needs Improvement"
	EMPLOYEE_PERFORMANCE_5_NAME = "Unsatisfactory"
)

const (
	// Menu Names
	MENU_1_OVERVIEW        = "Overview"
	MENU_2_TRIP_MANAGEMENT = "Trip Management"
	MENU_3_VENDORS         = "Vendors"
	MENU_4_CUSTOMERS       = "Customers"
	MENU_5_OPERATIONS      = "Operations"
	MENU_6_SETTINGS        = "Settings"
	MENU_7_REPORTS         = "Reports"
	MENU_8_EMPLOYEE        = "Employee"

	// Menu Descriptions
	MENU_1_OVERVIEW_DESC        = "trip and payment overview home"
	MENU_2_TRIP_MANAGEMENT_DESC = "Trip details"
	MENU_3_VENDORS_DESC         = "Vendor and vehicle screen"
	MENU_4_CUSTOMERS_DESC       = "Customer"
	MENU_5_OPERATIONS_DESC      = "Operations screens"
	MENU_6_SETTINGS_DESC        = "there are menus - branch, profile, vehicle size"
	MENU_7_REPORTS_DESC         = "Payment related info"
	MENU_8_EMPLOYEE_DESC        = "Employee info"
)
