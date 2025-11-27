package dtos

import (
	"time"

	"go-transport-hub/dtos/schema"
)

type LoginResponse struct {
	Message              string                `json:"message"`
	EmailId              string                `json:"email_id"`
	FirstName            string                `json:"first_name"`
	LoginId              int64                 `json:"login_id"`
	MobileNo             string                `json:"mobile_no"`
	RoleID               int64                 `json:"role_id"`
	RoleName             string                `json:"role_name"`
	EmployeeId           int64                 `json:"employee_id"`
	LoginType            string                `json:"login_type"`
	LastLogin            time.Time             `json:"last_login"`
	SessionTimeoutMS     int                   `json:"session_timeout_ms"`
	OrganisationResponse *OrganisationResponse `json:"organisation"`
	Customer             ScreenPermission      `json:"customer"`
	Operations           ScreenPermission      `json:"operations"`
	Overview             ScreenPermission      `json:"overview"`
	Settings             ScreenPermission      `json:"settings"`
	TripMgmt             ScreenPermission      `json:"trip_mgmt"`
	Vendors              ScreenPermission      `json:"vendors"`
	Reports              ScreenPermission      `json:"reports"`
	Employee             ScreenPermission      `json:"employee"`
	EmpAttendance        schema.EmpAttendance  `json:"employee_attendance"`
}

type ScreenPermission struct {
	IsEdit     bool `json:"is_edit"`
	IsView     bool `json:"is_view"`
	IsNoAccess bool `json:"is_no_access"`
}

type EmpRoleResponse struct {
	Message  string `json:"message"`
	RoleCode string `json:"role_code"`
	RoleName string `json:"role_name"`
}

type EmpRoles struct {
	Roles *[]schema.Role `json:"roles"`
}

type RoleEmp struct {
	RoleId            int64                   `json:"role_id"`
	RoleCode          string                  `json:"role_code"`
	RoleName          string                  `json:"role_name"`
	Description       string                  `json:"description"`
	OrgId             int64                   `json:"org_id"`
	IsActive          int64                   `json:"is_active"`
	Version           int64                   `json:"version"`
	UpdatedAt         string                  `json:"updated_at"`
	Overview          string                  `json:"overview"`
	TripManagement    string                  `json:"trip_management"`
	Vendors           string                  `json:"vendors"`
	Settings          string                  `json:"settings"`
	Operations        string                  `json:"operations"`
	Customer          string                  `json:"customer"`
	Reports           string                  `json:"reports"`
	ScreensPermission *[]ScreensPermissionRes `json:"screens_permission"`
}

type ScreensPermissionRes struct {
	ThubRoleScreensID  int64  `json:"thub_role_screens_id"`
	RoleId             int64  `json:"role_id"`
	RoleCode           string `json:"role_code"`
	RoleName           string `json:"role_name"`
	WebsiteScreenID    int64  `json:"website_screen_id"`
	PermisstionLabelID int64  `json:"permisstion_label_id"`
	ScreenName         string `json:"screen_name"`
	PermisstionLabel   string `json:"permisstion_label"`
}

type RoleEmpPre struct {
	RoleId      int64  `json:"role_id"`
	RoleCode    string `json:"role_code"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

type RoleEmpEntries struct {
	RoleEmpEntrie *[]RoleEmp `json:"roles"`
	Total         int64      `json:"total"`
	Limit         int64      `json:"limit"`
	OffSet        int64      `json:"offset"`
}

type OrganisationResponse struct {
	OrgId        int64  `json:"org_id"`
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	DomainName   string `json:"domain_name"`
	EmailId      string `json:"email_id"`
	ContactName  string `json:"contact_name"`
	ContactNo    string `json:"contact_no"`
	IsActive     int64  `json:"is_active"`
	LogoPath     string `json:"logo_path"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
}
type OrganisationBaseResponse struct {
	OrgId       int64  `json:"org_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	IsActive    int64  `json:"is_active"`
}

type EmployeeCreateResponse struct {
	Message    string `json:"message"`
	Login      string `json:"login"`
	EmployeeId int64  `json:"employee_id"`
}

type Messge struct {
	Message string `json:"message"`
}

type TripPrerequisite struct {
	TripSheetNumber string                `json:"trip_sheet_number"`
	CustomerCode    string                `json:"customer_code"`
	VendorCode      string                `json:"vendor_code"`
	TripType        []string              `json:"trip_type"`
	TripSheetType   []string              `json:"trip_sheet_type"`
	Branch          *[]BranchT            `json:"branches"`
	LoadingPoints   *[]LoadingPoints      `json:"loading_points"`
	Customers       *[]Customers          `json:"customers"`
	Vendor          *[]VendorT            `json:"vendors"`
	TripStatus      []string              `json:"trip_status"`
	PodTrips        *[]TripSheetInfo      `json:"pod_trips"`
	RegularTrips    *[]TripSheetInfo      `json:"regular_trips"`
	LRTrips         *[]TripSheetInfo      `json:"lr_trips"`
	Roles           *[]RoleEmpPre         `json:"roles"`
	VehicleSizeType *[]VehicleSizeTypePre `json:"vehicle_sizes"`
	DeclarationYear []string              `json:"declaration_year"`
	WebsiteScreen   []WebsiteScreen       `json:"website_screens"`
	PermissionLabel []PermissionLabel     `json:"permission_labels"`
	EmployeesPre    []EmployeesPre        `json:"employees"`
	InvoiceStatus   []string              `json:"invoice_status"`
}

type LoadingPoints struct {
	LoadingPointID int64  `json:"loading_point_id"`
	CityCode       string `json:"city_code"`
	CityName       string `json:"city_name"`
}

type Customers struct {
	CustomerId   int64  `json:"customer_id"`
	CustomerName string `json:"customer_name"`
}

type BranchT struct {
	BranchId   int64  `json:"branch_id"`
	BranchName string `json:"branch_name"`
}
type VendorT struct {
	VendorId   int64  `json:"vendor_id"`
	VendorName string `json:"vendor_name"`
}

type UploadTripSheetResponse struct {
	TripSheetNumber string `json:"trip_sheet_number"`
	ImagePath       string `json:"imagePath"`
	Message         string `json:"message"`
}

type TripSheetInfo struct {
	TripSheetID   int64  `json:"trip_sheet_id"`
	TripSheetNum  string `json:"trip_sheet_num"`
	LRNumber      string `json:"lr_number"`
	TripDate      string `json:"trip_date"`
	CustomerId    int64  `json:"customer_id"`
	PodRequired   int    `json:"pod_required"`
	LoadStatus    string `json:"trip_status"`
	InvoiceNumber string `json:"invoice_number"`
	DriverName    string `json:"driver_name"`
	MobileNumber  string `json:"mobile_number"`
	VehicleNumber string `json:"vehicle_number"`
	VehicleSize   string `json:"vehicle_size"`
}

type UploadManagePODResponse struct {
	TripSheetNumber int64  `json:"trip_sheet_number"`
	ImagePath       string `json:"imagePath"`
	Message         string `json:"message"`
}

type UploadEmpResponse struct {
	Employee  string `json:"employee"`
	ImagePath string `json:"image_path"`
	Message   string `json:"message"`
}

type LoginV1Response struct {
	UserId     int64  `json:"user_id"`
	UserIDText string `json:"user_id_text"`
	EmailId    string `json:"email_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	LoginId    int64  `json:"login_id"`
	IsActive   int64  `json:"is_active"`
	MobileNo   string `json:"mobile_no"`
	RoleID     int64  `json:"role_id"`
	EmployeeId int64  `json:"employee_id"`
	LoginType  string `json:"login_type"`
	LastLogin  string `json:"last_login"`
	RoleCode   string `json:"role_code"`
	RoleName   string `json:"role_name"`
}

type LoginRes struct {
	Users  *[]LoginV1Response `json:"users"`
	Total  int64              `json:"total"`
	Limit  int64              `json:"limit"`
	OffSet int64              `json:"offset"`
}

type InvoiceInfo struct {
	ID                   int64                  `json:"id"`
	InvoiceNumber        string                 `json:"invoice_number"`
	InvoiceRef           string                 `json:"invoice_ref"`
	WorkType             string                 `json:"work_type"`
	WorkStartDate        string                 `json:"work_start_date"`
	WorkEndDate          string                 `json:"work_end_date"`
	WorkPeriod           string                 `json:"work_period"`
	DocumentDateStr      string                 `json:"document_date_str"`
	DocumentDate         string                 `json:"document_date"`
	InvoiceStatus        string                 `json:"invoice_status"`
	InvoiceStatusDesc    string                 `json:"invoice_status_desc"`
	InvoiceAmount        float64                `json:"invoice_amount"`
	InvoiceDate          string                 `json:"invoice_date"`
	PaymentDate          string                 `json:"payment_date"`
	TotalTrips           string                 `json:"total_trips"`
	TripSheetsForInvoice []TripSheetsForInvoice `json:"trips"`
}
type InvoiceInfoResponse struct {
	InvoiceInfo          InvoiceObj             `json:"invoice_info"`
	TripSheetsForInvoice []TripSheetsForInvoice `json:"trips"`
}

type TripSheetsForInvoice struct {
	TripSheetID        int64  `json:"trip_sheet_id"`
	TripSheetNum       string `json:"trip_sheet_num"`
	TripType           string `json:"trip_type"`
	TripSheetType      string `json:"trip_sheet_type"`
	LoadHoursType      string `json:"load_hours_type"`
	OpenTripDateTime   string `json:"open_trip_date_time"`
	CustomerName       string `json:"customer_name"`
	CustomerCode       string `json:"customer_code"`
	VendorID           int64  `json:"vendor_id"`
	LoadingPointIDs    string `json:"loading_points"`
	UnLoadingPointIDs  string `json:"un_loading_point"`
	VehicleCapacityTon string `json:"vehicle_capacity_ton"`
	VehicleNumber      string `json:"vehicle_number"`
	VehicleSize        string `json:"vehicle_size"`
	MobileNumber       string `json:"mobile_number"`
	DriverName         string `json:"driver_name"`
	DriverLicenseImage string `json:"driver_license_image"`
	LRGateImage        string `json:"lr_gate_image"`
	LRNumber           string `json:"lr_number"`
	FromLocaion        string `json:"from_location"`
	ToLocation         string `json:"to_location"`

	// Customer Section
	CustomerCloseTripDateTime              string  `json:"customer_close_trip_date_time"`
	CustomerInvoiceNo                      string  `json:"customer_invoice_no"`
	CustomerPaymentReceivedDate            string  `json:"customer_payment_received_date"`
	CustomerRemark                         string  `json:"customer_remark"`
	CustomerBillingRaisedDate              string  `json:"customer_billing_raised_date"`
	CustomerBaseRate                       float64 `json:"customer_base_rate"`
	CustomerKMCost                         float64 `json:"customer_km_cost"`
	CustomerToll                           float64 `json:"customer_toll"`
	CustomerExtraHours                     float64 `json:"customer_extra_hours"`
	CustomerExtraKM                        float64 `json:"customer_extra_km"`
	CustomerTotalHire                      float64 `json:"customer_total_hire"`
	CustomerDebitAmount                    float64 `json:"customer_debit_amount"`
	PodRequired                            int64   `json:"pod_required"`
	CustomerPerLoadHire                    float64 `json:"customer_per_load_hire"`
	CustomerRunningKM                      float64 `json:"customer_running_km"`
	CustomerPerKMPrice                     float64 `json:"customer_per_km_price"`
	CustomerPlacedVehicleSize              string  `json:"customer_placed_vehicle_size"`
	CustomerLoadCancelled                  string  `json:"customer_load_cancelled"`
	CustomerReportedDateTimeForHaltingCalc string  `json:"customer_reported_date_time_for_halting_calc"`

	// Vendor Section
	VendorPaidDate         string  `json:"vendor_paid_date"`
	VendorCrossDock        string  `json:"vendor_cross_dock"`
	VendorRemark           string  `json:"vendor_remark"`
	VendorBaseRate         float64 `json:"vendor_base_rate"`
	VendorKMCost           float64 `json:"vendor_km_cost"`
	VendorToll             float64 `json:"vendor_toll"`
	VendorTotalHire        float64 `json:"vendor_total_hire"`
	VendorAdvance          float64 `json:"vendor_advance"`
	VendorDebitAmount      float64 `json:"vendor_debit_amount"`
	VendorBalanceAmount    float64 `json:"vendor_balance_amount"`
	VendorPaidBy           string  `json:"vendor_paid_by"`
	VendorLoadUnLoadAmount float64 `json:"vendor_load_unload_amount"`
	VendorHaltingDays      float64 `json:"vendor_halting_days"`
	VendorHaltingPaid      float64 `json:"vendor_halting_paid"`
	VendorExtraDelivery    float64 `json:"vendor_extra_delivery"`
	VendorBreakDown        float64 `json:"vendor_break_down"`
	PodReceived            int64   `json:"pod_received"`
	LoadStatus             string  `json:"load_status"`
	ZonalName              string  `json:"zonal_name"`
}

type InvStatusResponse struct {
	DraftCount      int64 `json:"draft"`
	InvcoiceRaised  int   `json:"raised"`
	InvcoiceOverDue int   `json:"over_due"`
}
