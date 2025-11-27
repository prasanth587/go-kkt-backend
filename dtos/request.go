package dtos

type AdminLogin struct {
	EmailID  string `json:"email_id"`
	MobileNo string `json:"mobile_mo"`
	Password string `json:"password"`
	UserName string `json:"user_name"`
}

type EmpRole struct {
	RoleCode          string              `json:"role_code"`
	RoleName          string              `json:"role_name"`
	Description       string              `json:"description"`
	OrgId             int64               `json:"org_id"`
	Version           int64               `json:"version"`
	ScreensPermission []ScreensPermission `json:"screens_permission"`
}

type ScreensPermission struct {
	WebsiteScreenID    int64 `json:"website_screen_id"`
	PermisstionLabelID int64 `json:"permisstion_label_id"`
}

type EmpRoleUpdate struct {
	RoleId            int64               `json:"role_id"`
	RoleCode          string              `json:"role_code"`
	RoleName          string              `json:"role_name"`
	Description       string              `json:"description"`
	OrgId             int64               `json:"org_id"`
	IsActive          int64               `json:"is_active"`
	Version           int64               `json:"version"`
	ScreensPermission []ScreensPermission `json:"screens_permission"`
}

type CreateEmployeeReq struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	EmployeeCode  string `json:"employee_code"`
	MobileNo      string `json:"mobile_no"`
	EmailID       string `json:"email_id"`
	RoleID        int64  `json:"role_id"`
	DOB           string `json:"dob"`
	Gender        string `json:"gender"`
	AadharNo      string `json:"aadhar_no"`
	AccessNo      string `json:"access_no"`
	AccessToken   string `json:"access_token"`
	IsActive      int64  `json:"is_active"`
	JoiningDate   string `json:"joining_date"`
	RelievingDate string `json:"relieving_date"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	IsSuperAdmin  int64  `json:"is_super_admin" `
	IsAdmin       int64  `json:"is_admin"`
	OrgID         int64  `json:"org_id"`
	LoginType     string `json:"login_type" `
	Image         string `json:"image"`
	Version       int64  `json:"version"`
	PinCode       int64  `json:"pin_code"`
	MonthlySalary int64  `json:"monthly_salary"`
	AnnualSalary  int64  `json:"annual_salary"`
	AnnualBonus   int64  `json:"annual_bonus"`
}

type UpdatePassword struct {
	Password  string `json:"new_password"`
	UserLogin string `json:"user_login"`
}
type UpdateEmployeeReq struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	EmployeeCode  string `json:"employee_code"`
	MobileNo      string `json:"mobile_no"`
	EmailID       string `json:"email_id"`
	RoleID        int64  `json:"role_id"`
	DOB           string `json:"dob"`
	Gender        string `json:"gender"`
	AadharNo      string `json:"aadhar_no"`
	AccessNo      string `json:"access_no"`
	AccessToken   string `json:"access_token"`
	IsActive      int64  `json:"is_active"`
	JoiningDate   string `json:"joining_date"`
	RelievingDate string `json:"relieving_date"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	IsSuperAdmin  int64  `json:"is_super_admin" `
	IsAdmin       int64  `json:"is_admin"`
	LoginType     string `json:"login_type" `
	Version       int64  `json:"version"`
	PinCode       int64  `json:"pin_code"`
	MonthlySalary int64  `json:"monthly_salary"`
	AnnualSalary  int64  `json:"annual_salary"`
	AnnualBonus   int64  `json:"annual_bonus"`
}

type EmployeeEntiry struct {
	EmpId          int64  `json:"emp_id"`
	EmployeeIDText string `json:"emp_id_text"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	EmployeeCode   string `json:"employee_code"`
	MobileNo       string `json:"mobile_no"`
	EmailID        string `json:"email_id"`
	RoleID         int64  `json:"role_id"`
	DOB            string `json:"dob"`
	Gender         string `json:"gender"`
	AadharNo       string `json:"aadhar_no"`
	AccessNo       string `json:"access_no"`
	//AccessToken   string `json:"access_token"`
	IsActive            int64  `json:"is_active"`
	JoiningDate         string `json:"joining_date"`
	RelievingDate       string `json:"relieving_date"`
	AddressLine1        string `json:"address_line1"`
	AddressLine2        string `json:"address_line2"`
	City                string `json:"city"`
	State               string `json:"state"`
	Country             string `json:"country"`
	IsSuperAdmin        int64  `json:"is_super_admin" `
	IsAdmin             int64  `json:"is_admin"`
	OrgID               int64  `json:"org_id"`
	LoginType           string `json:"login_type" `
	Image               string `json:"image"`
	PinCode             int64  `json:"pin_code"`
	Version             int64  `json:"version"`
	EmployeePerformance string `json:"employee_performance"`
	VehicleAssigned     string `json:"vehicle_assigned"`
	MonthlySalary       int64  `json:"monthly_salary"`
	AnnualSalary        int64  `json:"annual_salary"`
	AnnualBonus         int64  `json:"annual_bonus"`
	AttendanceStatus    string `json:"attendance_status"` // Present, Leave, Absent
}
type EmployeeEntries struct {
	EmployeeEntiry *[]EmployeeEntiry `json:"employees"`
	Total          int64             `json:"total"`
	Limit          int64             `json:"limit"`
	OffSet         int64             `json:"offset"`
}

type EmpActiveStatusResponse struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	IsActive int64  `json:"is_active"`
	EmpId    int64  `json:"emp_id"`
}

type EmployeeUpdateResponse struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	EmpId   int64  `json:"emp_id"`
}

type CreateDriverReq struct {
	FirstName              string  `json:"first_name"`
	LastName               string  `json:"last_name"`
	LicenseNumber          string  `json:"license_number"`
	LicenseExpiryDate      string  `json:"license_expiry_date"`
	MobileNo               string  `json:"mobile_no"`
	AlternateContactNumber string  `json:"alternate_contact_number"`
	EmailID                string  `json:"email_id"`
	JoiningDate            string  `json:"joining_date"`
	RelievingDate          string  `json:"relieving_date"`
	IsActive               int64   `json:"is_active"`
	VehicleID              int64   `json:"vehicle_id"`
	DriverExperience       float64 `json:"driver_experience"`
	AddressLine1           string  `json:"address_line1"`
	AddressLine2           string  `json:"address_line2"`
	City                   string  `json:"city"`
	State                  string  `json:"state"`
	Country                string  `json:"country"`
	OrgID                  int     `json:"org_id"`
	LoginType              string  `json:"login_type"`
}

type DriverRes struct {
	DriverId               int64   `json:"driver_id"`
	FirstName              string  `json:"first_name"`
	LastName               string  `json:"last_name"`
	LicenseNumber          string  `json:"license_number"`
	LicenseExpiryDate      string  `json:"license_expiry_date"`
	MobileNo               string  `json:"mobile_no"`
	AlternateContactNumber string  `json:"alternate_contact_number"`
	EmailID                string  `json:"email_id"`
	JoiningDate            string  `json:"joining_date"`
	RelievingDate          string  `json:"relieving_date"`
	IsActive               int64   `json:"is_active"`
	VehicleID              int64   `json:"vehicle_id"`
	DriverExperience       float64 `json:"driver_experience"`
	AddressLine1           string  `json:"address_line1"`
	AddressLine2           string  `json:"address_line2"`
	City                   string  `json:"city"`
	State                  string  `json:"state"`
	Country                string  `json:"country"`
	OrgID                  int64   `json:"org_id"`
	LoginType              string  `json:"login_type"`
	LicenseFrontImg        string  `json:"license_front_img"`
	LicenseBackImg         string  `json:"license_back_img"`
	OtherDocument          string  `json:"other_document"`
	ProfileImage           string  `json:"profile_image"`
}

type DriverEntries struct {
	DriverEntries *[]DriverRes `json:"drivers"`
	Total         int64        `json:"total"`
	Limit         int64        `json:"limit"`
	OffSet        int64        `json:"offset"`
}

type UpdateDriverReq struct {
	FirstName              string  `json:"first_name"`
	LastName               string  `json:"last_name"`
	LicenseNumber          string  `json:"license_number"`
	LicenseExpiryDate      string  `json:"license_expiry_date"`
	MobileNo               string  `json:"mobile_no"`
	AlternateContactNumber string  `json:"alternate_contact_number"`
	EmailID                string  `json:"email_id"`
	JoiningDate            string  `json:"joining_date"`
	RelievingDate          string  `json:"relieving_date"`
	IsActive               int64   `json:"is_active"`
	VehicleID              int64   `json:"vehicle_id"`
	DriverExperience       float64 `json:"driver_experience"`
	AddressLine1           string  `json:"address_line1"`
	AddressLine2           string  `json:"address_line2"`
	City                   string  `json:"city"`
	State                  string  `json:"state"`
	Country                string  `json:"country"`
	OrgID                  int     `json:"org_id"`
	LoginType              string  `json:"login_type"`
}

type VehicleReq struct {
	VehicleType             string `json:"vehicle_type"`
	VehicleNumber           string `json:"vehicle_number"`
	VehicleModel            string `json:"vehicle_model"`
	VehicleYear             int64  `json:"vehicle_year"`
	VehicleCapacity         string `json:"vehicle_capacity"`
	VehicleInsuranceNumber  string `json:"vehicle_insurance_number"`
	InsuranceExpiryDate     string `json:"insurance_expiry_date"`
	VehicleRegistrationDate string `json:"vehicle_registration_date"`
	VehicleRenewalDate      string `json:"vehicle_renewal_date"`
	IsActive                int64  `json:"is_active"`
	DriverID                int64  `json:"driver_id"`
	OrgID                   int64  `json:"org_id"`
	Status                  string `json:"status"`
	// VehicleImage            string `json:"vehicle_image"`
	// InsuranceDocument       string `json:"fitness_certificate"`
	// RegistrationDocument    string `json:"registration_document"`
}
type VehicleUpdate struct {
	VehicleType             string `json:"vehicle_type"`
	VehicleNumber           string `json:"vehicle_number"`
	VehicleModel            string `json:"vehicle_model"`
	VehicleYear             int64  `json:"vehicle_year"`
	VehicleCapacity         string `json:"vehicle_capacity"`
	VehicleInsuranceNumber  string `json:"vehicle_insurance_number"`
	InsuranceExpiryDate     string `json:"insurance_expiry_date"`
	VehicleRegistrationDate string `json:"vehicle_registration_date"`
	VehicleRenewalDate      string `json:"vehicle_renewal_date"`
	IsActive                int64  `json:"is_active"`
	DriverID                int64  `json:"driver_id"`
	OrgID                   int64  `json:"org_id"`
	Status                  string `json:"status"`
	// VehicleImage            string `json:"vehicle_image"`
	// InsuranceDocument       string `json:"insurance_document"`
	// RegistrationDocument    string `json:"registration_document"`
}

type VehicleRes struct {
	VehicleId                    int64  `json:"vehicle_id"`
	VehicleType                  string `json:"vehicle_type"`
	VehicleNumber                string `json:"vehicle_number"`
	VehicleModel                 string `json:"vehicle_model"`
	VehicleYear                  int64  `json:"vehicle_year"`
	VehicleCapacity              string `json:"vehicle_capacity"`
	VehicleInsuranceNumber       string `json:"vehicle_insurance_number"`
	InsuranceExpiryDate          string `json:"insurance_expiry_date"`
	VehicleRegistrationDate      string `json:"vehicle_registration_date"`
	VehicleRenewalDate           string `json:"vehicle_renewal_date"`
	IsActive                     int64  `json:"is_active"`
	DriverID                     int64  `json:"driver_id"`
	OrgID                        int64  `json:"org_id"`
	VehicleImage                 string `json:"vehicle_image"`
	FitnessCertificate           string `json:"fitness_certificate"`
	InsuranceCertificate         string `json:"insurance_certificate"`
	PollutionCertificate         string `json:"pollution_certificate"`
	NationalPermitsCertificate   string `json:"national_permits_certificate"`
	RegistrationCertificate      string `json:"registration_certificate"`
	AnnualMaintenanceCertificate string `json:"annual_maintenance_certificate"`
	Status                       string `json:"status"`
}

type VehicleEntries struct {
	VendorEntiry *[]VehicleRes `json:"vehicles"`
	Total        int64         `json:"total"`
	Limit        int64         `json:"limit"`
	OffSet       int64         `json:"offset"`
}

type VendorReq struct {
	VendorName        string `json:"vendor_name"`
	VendorCode        string `json:"vendor_code"`
	MobileNumber      string `json:"mobile_number"`
	ContactPerson     string `json:"contact_person"`
	AlternativeNumber string `json:"alternative_number"`
	AddressLine1      string `json:"address_line1"`
	AddressLine2      string `json:"address_line2"`
	City              string `json:"city"`
	State             string `json:"state"`
	Status            string `json:"status"`
	OrgId             int64  `json:"org_id"`
	IsActive          int64  `json:"is_active"`
}
type VendorUpdate struct {
	VendorName        string `json:"vendor_name"`
	VendorCode        string `json:"vendor_code"`
	MobileNumber      string `json:"mobile_number"`
	ContactPerson     string `json:"contact_person"`
	AlternativeNumber string `json:"alternative_number"`
	AddressLine1      string `json:"address_line1"`
	AddressLine2      string `json:"address_line2"`
	City              string `json:"city"`
	State             string `json:"state"`
	Status            string `json:"status"`
	OrgId             int64  `json:"org_id"`
	IsActive          int64  `json:"is_active"`
}

type VendorRes struct {
	VendorId              int64  `json:"vendor_id"`
	VendorName            string `json:"vendor_name"`
	VendorCode            string `json:"vendor_code"`
	MobileNumber          string `json:"mobile_number"`
	ContactPerson         string `json:"contact_person"`
	AlternativeNumber     string `json:"alternative_number"`
	AddressLine1          string `json:"address_line1"`
	AddressLine2          string `json:"address_line2"`
	IsActive              int64  `json:"is_active"`
	Status                string `json:"status"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	VisitingCardImage     string `json:"visiting_card_image"`
	PancardImg            string `json:"pancard_img"`
	AadharCardImg         string `json:"aadhar_card_img"`
	CancelledCheckBookImg string `json:"cancelled_check_book_img"`
	BankPassbookImg       string `json:"bank_passbook_img"`
	GstDocumentImg        string `json:"gst_document_img"`
}
type VendorEntries struct {
	VendorEntiry *[]VendorRes `json:"vendors"`
	Total        int64        `json:"total"`
	Limit        int64        `json:"limit"`
	OffSet       int64        `json:"offset"`
}

// Customer
type CustomerReq struct {
	BranchId                int64  `json:"branch_id"`
	CustomerCode            string `json:"customer_code"`
	CustomerName            string `json:"customer_name"`
	Address                 string `json:"address"`
	City                    string `json:"city"`
	State                   string `json:"state"`
	PinCode                 int64  `json:"pincode"`
	Country                 string `json:"country"`
	GSTNType                string `json:"gstin_type"`
	GSTNNo                  string `json:"gstin_no"`
	PanNumber               string `json:"pan_number"`
	ContactPerson           string `json:"contact_person"`
	MobileNumber            string `json:"mobile_number"`
	AlternativeNumber       string `json:"alternative_number"`
	EmailID                 string `json:"email_id"`
	Status                  string `json:"status"`
	IsActive                int64  `json:"is_active"`
	OrgId                   int64  `json:"org_id"`
	BranchSalesPersonName   string `json:"branch_sales_person_name"`
	BranchResponsiblePerson string `json:"branch_responsible_person"`
	EmployeeCode            string `json:"employee_code"`
	BusinessStartDate       string `json:"business_start_date"`
	CredibilityDays         string `json:"credibility_days"`
	CompanyType             string `json:"company_type"`
}

type CustomerUpdate struct {
	BranchId                int64  `json:"branch_id"`
	CustomerCode            string `json:"customer_code"`
	CustomerName            string `json:"customer_name"`
	Address                 string `json:"address"`
	City                    string `json:"city"`
	State                   string `json:"state"`
	PinCode                 int64  `json:"pincode"`
	Country                 string `json:"country"`
	GSTNType                string `json:"gstin_type"`
	GSTNNo                  string `json:"gstin_no"`
	PanNumber               string `json:"pan_number"`
	ContactPerson           string `json:"contact_person"`
	MobileNumber            string `json:"mobile_number"`
	AlternativeNumber       string `json:"alternative_number"`
	EmailID                 string `json:"email_id"`
	Status                  string `json:"status"`
	IsActive                int64  `json:"is_active"`
	BranchSalesPersonName   string `json:"branch_sales_person_name"`
	BranchResponsiblePerson string `json:"branch_responsible_person"`
	EmployeeCode            string `json:"employee_code"`
	BusinessStartDate       string `json:"business_start_date"`
	CredibilityDays         string `json:"credibility_days"`
	CompanyType             string `json:"company_type"`
	OrgId                   int64  `json:"org_id"`
}

type CustomerRes struct {
	CustomerId              int64  `json:"customer_id"`
	BranchId                int64  `json:"branch_id"`
	CustomerCode            string `json:"customer_code"`
	CustomerName            string `json:"customer_name"`
	Address                 string `json:"address"`
	City                    string `json:"city"`
	State                   string `json:"state"`
	PinCode                 int64  `json:"pincode"`
	Country                 string `json:"country"`
	GSTNType                string `json:"gstin_type"`
	GSTNNo                  string `json:"gstin_no"`
	PanNumber               string `json:"pan_number"`
	ContactPerson           string `json:"contact_person"`
	MobileNumber            string `json:"mobile_number"`
	AlternativeNumber       string `json:"alternative_number"`
	EmailID                 string `json:"email_id"`
	Status                  string `json:"status"`
	IsActive                int64  `json:"is_active"`
	OrgId                   int64  `json:"org_id"`
	BranchSalesPersonName   string `json:"branch_sales_person_name"`
	BranchResponsiblePerson string `json:"branch_responsible_person"`
	EmployeeCode            string `json:"employee_code"`
	BusinessStartDate       string `json:"business_start_date"`
	CredibilityDays         int64  `json:"credibility_days"`
	CompanyType             string `json:"company_type"`
	AgreementDocImage       string `json:"agreement_doc_image"`
}
type CustomerEntries struct {
	CustomerEntiry *[]CustomerRes `json:"customers"`
	Total          int64          `json:"total"`
	Limit          int64          `json:"limit"`
	OffSet         int64          `json:"offset"`
}

type BranchReq struct {
	BranchName   string `json:"branch_name"`
	BranchCode   string `json:"branch_code"`
	OrgID        int64  `json:"org_id"`
	IsActive     int64  `json:"is_active"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}

type BranchUpdate struct {
	BranchName   string `json:"branch_name"`
	BranchCode   string `json:"branch_code"`
	OrgID        int64  `json:"org_id"`
	IsActive     int64  `json:"is_active"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}

type Branches struct {
	Branch *[]Branch `json:"branches"`
	Total  int64     `json:"total"`
	Limit  int64     `json:"limit"`
	OffSet int64     `json:"offset"`
}

type Branch struct {
	BranchId     int64  `json:"branch_id"`
	BranchName   string `json:"branch_name"`
	BranchCode   string `json:"branch_code"`
	OrgID        int64  `json:"org_id"`
	IsActive     int64  `json:"is_active"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}

type LoadingPointReq struct {
	BranchID    int64  `json:"branch_id"`
	OrgID       int64  `json:"org_id"`
	IsActive    int64  `json:"is_active"`
	CityCode    string `json:"city_code"`
	CityName    string `json:"city_name"`
	AddressLine string `json:"address_line"`
	MapLink     string `json:"map_link"`
	State       string `json:"state"`
	Country     string `json:"country"`
}
type LoadingPointUpdate struct {
	BranchID    int64  `json:"branch_id"`
	OrgID       int64  `json:"org_id"`
	IsActive    int64  `json:"is_active"`
	CityCode    string `json:"city_code"`
	CityName    string `json:"city_name"`
	AddressLine string `json:"address_line"`
	MapLink     string `json:"map_link"`
	State       string `json:"state"`
	Country     string `json:"country"`
}

type LoadingPointEntries struct {
	LoadingPointID int64  `json:"loading_point_id"`
	BranchID       int64  `json:"branch_id"`
	OrgID          int64  `json:"org_id"`
	IsActive       int64  `json:"is_active"`
	CityCode       string `json:"city_code"`
	CityName       string `json:"city_name"`
	AddressLine    string `json:"address_line"`
	MapLink        string `json:"map_link"`
	State          string `json:"state"`
	Country        string `json:"country"`
}

type LoadingPointes struct {
	LoadingPointEntries *[]LoadingPointEntries `json:"loading_unloading_points"`
	Total               int64                  `json:"total"`
	Limit               int64                  `json:"limit"`
	OffSet              int64                  `json:"offset"`
}

type CreateTripSheetHeader struct {
	TripSheetNum     string `json:"trip_sheet_num"`
	TripType         string `json:"trip_type"`
	TripSheetType    string `json:"trip_sheet_type"`
	LoadHoursType    string `json:"load_hours_type"`
	OpenTripDateTime string `json:"open_trip_date_time"`
	// Foreign Keys
	BranchID          uint64   `json:"branch_id"`
	CustomerID        int64    `json:"customer_id"`
	LoadingPointIDs   []uint64 `json:"loading_point_ids"`
	UnLoadingPointIDs []uint64 `json:"un_loading_point_ids"`
	VendorID          uint64   `json:"vendor_id"`
	// Vehicle Details
	VehicleCapacityTon string `json:"vehicle_capacity_ton"`
	VehicleNumber      string `json:"vehicle_number"`
	VehicleSize        string `json:"vehicle_size"`
	MobileNumber       string `json:"mobile_number"`
	DriverName         string `json:"driver_name"`
	DriverLicenseImage string `json:"driver_license_image"`
	OrgId              uint64 `json:"org_id"`

	//Line Trip
	ZonalName string `json:"zonal_name"`

	LoadStatus  string `json:"load_status"`
	UserLoginId int64  `json:"user_login_id"`
}

type UpdateTripSheetHeader struct {
	LoadingPointIDs   []int64 `json:"loading_point_ids"`
	UnLoadingPointIDs []int64 `json:"un_loading_point_ids"`
	TripSheetType     string  `json:"trip_sheet_type"`
	TripType          string  `json:"trip_type"`
	VehicleNumber     string  `json:"vehicle_number"`
	VehicleSize       string  `json:"vehicle_size"`
	TripSheetNum      string  `json:"trip_sheet_num"`
	OpenTripDateTime  string  `json:"open_trip_date_time"`
	LRNUmber          string  `json:"lr_number"`
	LoadHoursType     string  `json:"load_hours_type"`
	CustomerID        int64   `json:"customer_id"`
	VehicleSizeID     int64   `json:"vehicle_size_id"`

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
	CustomerPerLoadHire                    float64 `json:"customer_per_load_hire"`
	CustomerRunningKM                      float64 `json:"customer_running_km"`
	CustomerPerKMPrice                     float64 `json:"customer_per_km_price"`
	CustomerPlacedVehicleSize              string  `json:"customer_placed_vehicle_size"`
	CustomerLoadCancelled                  string  `json:"customer_load_cancelled"`
	CustomerReportedDateTimeForHaltingCalc string  `json:"customer_reported_date_time_for_halting_calc"`
	PodReceived                            int     `json:"pod_received"`

	VendorCommission       float64 `json:"vendor_commission"`
	VendorPaidDate         string  `json:"vendor_paid_date"`
	VendorCrossDock        string  `json:"vendor_cross_dock"`
	VendorRemark           string  `json:"vendor_remark"`
	VendorID               uint64  `json:"vendor_id"`
	VendorBaseRate         float64 `json:"vendor_base_rate"`
	VendorKMCost           float64 `json:"vendor_km_cost"`
	VendorToll             float64 `json:"vendor_toll"`
	VendorTotalHire        float64 `json:"vendor_total_hire"`
	VendorAdvance          float64 `json:"vendor_advance"`
	VendorDebitAmount      float64 `json:"vendor_debit_amount"`
	VendorBalanceAmount    float64 `json:"vendor_balance_amount"`
	VendorBreakDown        string  `json:"vendor_break_down"`
	VendorPaidBy           string  `json:"vendor_paid_by"`
	VendorLoadUnLoadAmount float64 `json:"vendor_load_unload_amount"`
	VendorHaltingDays      float64 `json:"vendor_halting_days"`
	VendorHaltingPaid      float64 `json:"vendor_halting_paid"`
	VendorExtraDelivery    float64 `json:"vendor_extra_delivery"`
	VendorMonul            float64 `json:"vendor_monul"`
	VendorTotalAmount      float64 `json:"vendor_total_amount"`
	PodRequired            int     `json:"pod_required"`
	LoadStatus             string  `json:"load_status"`
	VehicleCapacityTon     string  `json:"vehicle_capacity_ton"`
	MobileNumber           string  `json:"mobile_number"`
	DriverName             string  `json:"driver_name"`
	ZonalName              string  `json:"zonal_name"`

	TripSubmittedDate string `json:"trip_submitted_date"`

	LRGateImage string `json:"lr_gate_image"`
	UserLoginId int64  `json:"user_login_id"`
}

type TripSheet struct {
	TripSheetID      int64  `json:"trip_sheet_id"`
	TripSheetNum     string `json:"trip_sheet_num"`
	TripType         string `json:"trip_type"`
	TripSheetType    string `json:"trip_sheet_type"`
	LoadHoursType    string `json:"load_hours_type"`
	OpenTripDateTime string `json:"open_trip_date_time"`
	// Foreign Keys
	BranchID          int64                 `json:"branch_id"`
	Customer          *Customers            `json:"customer"`
	Vendor            *VendorT              `json:"vendor"`
	VendorID          int64                 `json:"vendor_id"`
	LoadingPointIDs   *[]TripSheetLoading   `json:"loading_points"`
	UnLoadingPointIDs *[]TripSheetUnLoading `json:"un_loading_point"`
	// Vehicle Details
	VehicleCapacityTon string `json:"vehicle_capacity_ton"`
	VehicleNumber      string `json:"vehicle_number"`
	VehicleSize        string `json:"vehicle_size"`
	MobileNumber       string `json:"mobile_number"`
	DriverName         string `json:"driver_name"`
	DriverLicenseImage string `json:"driver_license_image"`
	LRGateImage        string `json:"lr_gate_image"`
	LRNumber           string `json:"lr_number"`

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
	VendorBreakDown        string  `json:"vendor_break_down"`
	PodReceived            int64   `json:"pod_received"`
	LoadStatus             string  `json:"load_status"`
	VendorMonul            float64 `json:"vendor_monul"`
	VendorTotalAmount      float64 `json:"vendor_total_amount"`
	ZonalName              string  `json:"zonal_name"`

	TripSubmittedDate string `json:"trip_submitted_date"`
	TripClosedDate    string `json:"trip_closed_date"`
	TripDeliveredDate string `json:"trip_delivered_date"`
	TripCompletedDate string `json:"trip_completed_date"`

	VehicleSizeID    int64   `json:"vehicle_size_id"`
	VendorCommission float64 `json:"vendor_commission"`
}

type TripSheetRes struct {
	TripSheet *[]TripSheet `json:"trip_sheets"`
	Total     int64        `json:"total"`
	Limit     int64        `json:"limit"`
	OffSet    int64        `json:"offset"`
}

type TripSheetLoading struct {
	LoadingPointId int64  `json:"loading_point_id"`
	CityCode       string `json:"city_code"`
}

type TripSheetUnLoading struct {
	UnLoadingPointId int64  `json:"un_loading_point_id"`
	CityCode         string `json:"city_code"`
}

type TripSheetLoadUnLoadPoints struct {
	LoadingPointID int64  `json:"loading_point_id"`
	Type           string `json:"type"`
}

type TripStats struct {
	TripSheetID   int64  `json:"trip_sheet_id"`
	TripSheetNum  string `json:"trip_sheet_num"`
	TripType      string `json:"trip_type"`
	TripSheetType string `json:"trip_sheet_type"`
	PodRequired   int64  `json:"pod_required"`
	PodReceived   int64  `json:"pod_received"`
	LoadStatus    string `json:"load_status"`
}

type TripStatsRes struct {
	TripStats *[]TripStats `json:"trip_stats"`
	Total     int          `json:"total"`
}

type TripStatsChallanInfo struct {
	AddressTop    string `json:"address_top"`
	Title         string `json:"title"`
	Month         string `json:"month"`
	SerialNumber  string `json:"serial_number"`
	LoadingDate   string `json:"loading_date"`
	LRNumber      string `json:"lr_number"`
	VehicleNumber string `json:"vehicle_number"`
	VehicleSize   string `json:"vehicle_size"`
	FromLocations string `json:"from_locations"`
	ToLocations   string `json:"to_locations"`
	DriverName    string `json:"driver_name"`
	ContactNumber string `json:"contact_number"`
	TransportName string `json:"transport_name"`

	VehicleHire    float64 `json:"vehicle_hire"`
	VehicleAdvance float64 `json:"vehicle_advance"`
	Mamul          float64 `json:"mamul"`
	Balance1       float64 `json:"balance1"`

	LoadingUNCharges float64 `json:"loading_un_charges"`
	HaltingCharges   float64 `json:"halting_charges"`
	Balance2         float64 `json:"balance2"`

	IssuedBy     string `json:"issued_by"`
	ReceivedDate string `json:"received_date"`
	CompanyName  string `json:"company_name"`

	OfficeBalance3     float64 `json:"office_balance3"`
	OfficeCommission   float64 `json:"office_commission"`
	OfficeTotal        float64 `json:"office_total"`
	OfficeApprovedBy   string  `json:"office_approved_by"`
	OfficeAuthorizedBy string  `json:"office_authorized_by"`
	Notes              string  `json:"notes"`
}

type EmployeeAssigned struct {
	UserID        int64  `json:"user_id"`
	VehicleNumber string `json:"vehicle_number"`
	TripSheetId   int64  `json:"trip_sheet_id"`
}

type InvoiceObj struct {
	ID            int64  `json:"id"`
	InvoiceNumber string `json:"invoice_number"`
	InvoiceRef    string `json:"invoice_ref"`
	//WorkType          string  `json:"work_type"`
	WorkStartDate     string  `json:"work_start_date"`
	WorkEndDate       string  `json:"work_end_date"`
	WorkPeriod        string  `json:"work_period"`
	DocumentDateStr   string  `json:"document_date_str"`
	DocumentDate      string  `json:"document_date"`
	InvoiceStatus     string  `json:"invoice_status"`
	InvoiceStatusDesc string  `json:"invoice_status_desc"`
	InvoiceAmount     float64 `json:"invoice_amount"`
	InvoiceDate       string  `json:"invoice_date"`
	PaymentDate       string  `json:"payment_date"`
	CustomerName      string  `json:"customer_name"`
	CustomerId        int64   `json:"customer_id"`
	CustomerCode      string  `json:"customer_code"`
	PaymentTerms      string  `json:"payment_terms"`
	IsOverDue         bool    `json:"is_over_due"`
}

type InvoiceBreach struct {
	ID            int64  `json:"id"`
	InvoiceNumber string `json:"invoice_number"`
	InvoiceStatus string `json:"invoice_status"`
	InvoiceDate   string `json:"invoice_date"`
	PaymentTerms  string `json:"payment_terms"`
}

type InvoiceBreachEntries struct {
	TotalInvoices    int `json:"total_invoices"`
	BreachedInvoices int `json:"breached_invoices"`
}

type InvoiceEntries struct {
	InvoiceObj *[]InvoiceObj `json:"invoices"`
	Total      int64         `json:"total"`
	Limit      int64         `json:"limit"`
	OffSet     int64         `json:"offset"`
}

type InvoicePDFInfo struct {
	Header                  string            `json:"header"`
	BilledFromText          string            `json:"billed_from_text"`
	BilledFromOfficeName    string            `json:"billed_from_office_name"`
	BilledFromOfficeAddress string            `json:"billed_from_office_address"`
	BilledFromOfficeEmail   string            `json:"billed_from_office_email"`
	BilledFromOfficeGSTIN   string            `json:"billed_from_office_gstin"`
	BilledFromOfficePan     string            `json:"billed_from_office_pan"`
	BilledToText            string            `json:"billed_to_text"`
	BilledToOffice          string            `json:"billed_to_office"`
	BilledToAddress         string            `json:"billed_to_address"`
	BilledToState           string            `json:"billed_to_state"`
	BilledToGSTINUIN        string            `json:"billed_to_gstin_uin"`
	BillingPeriodText       string            `json:"billing_period_text"`
	BillingPeriod           string            `json:"billing_period"`
	PaymentTermsText        string            `json:"payment_terms_text"`
	InFavourOf              string            `json:"in_favour_of"`
	BankBranch              string            `json:"bank_branch"`
	AccountNo               string            `json:"account_no"`
	IFSCCode                string            `json:"ifsc_code"`
	InvoiceDetails          InvoiceDetails    `json:"invoice_details"`
	BillSupplyDetails       BillSupplyDetails `json:"bill_supply_details"`
	TermsAndConditions      string            `json:"terms_and_conditions"`
	ForCompany              string            `json:"for_company"`
	AuthorizedSignatureText string            `json:"authorized_signature_text"`
}

type InvoiceDetails struct {
	SLNO                  string `json:"sl_no"`
	DescriptionOfServices string `json:"description_of_services"`
	HSNSACCode            string `json:"hsn_sac_code"`
	Amount                string `json:"amount"`
	TotalAmount           string `json:"total_amount"`
	AmountInWords         string `json:"amount_in_words"`
}
type BillSupplyDetails struct {
	BillSupplyDate             string `json:"bill_supply_date"`
	BillSupplyNo               string `json:"bill_supply_no"`
	BillSupplyTypeOfEnterprise string `json:"bill_supply_type_of_enterprise:"`
}
