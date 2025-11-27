package dtos

type VendorRequest struct {
	VendorName               string                `json:"vendor_name"`
	VendorCode               string                `json:"vendor_code"`
	OwnerName                string                `json:"owner_name"`
	GSTNumber                string                `json:"gst_number"`
	PreferredOperatingRoutes string                `json:"preferred_operating_routes"`
	AddressLine1             string                `json:"address_line1"`
	City                     string                `json:"city"`
	State                    string                `json:"state"`
	PANNumber                string                `json:"pan_number"`
	TDSDeclaration           string                `json:"tds_declaration"`
	Remark                   string                `json:"remark"`
	BankAccountHolderName    string                `json:"bank_account_holder_name"`
	BankAccountNumber        string                `json:"bank_account_number"`
	BankName                 string                `json:"bank_name"`
	BankIFSCCode             string                `json:"bank_ifsc_code"`
	PancardImg               string                `json:"pancard_img"`
	BankPassbookOrChequeImg  string                `json:"bank_passbook_or_cheque_img"`
	Status                   string                `json:"status"`
	IsActive                 int                   `json:"is_active"`
	OrgID                    int64                 `json:"org_id"`
	LoginType                string                `json:"login_type"`
	ContactInfo              []VendorContactInfo   `json:"contact_info"`
	Vehicles                 []VehicleModel        `json:"vehicle"`
	DeclarationDocument      []DeclarationDocument `json:"declaration_document"`
}

type VendorContactInfo struct {
	ContactInfoID     int64  `json:"contact_info_id"`
	VendorID          int64  `json:"vendor_id"`
	ContactPersonName string `json:"contact_person_name"`
	Post              string `json:"post"`
	EmailID           string `json:"email_id"`
	ContactNumber1    string `json:"contact_number1"`
	ContactNumber2    string `json:"contact_number2"`
}

type VendorResponse struct {
	VendorId                 int64                    `json:"vendor_id"`
	VendorName               string                   `json:"vendor_name"`
	VendorCode               string                   `json:"vendor_code"`
	OwnerName                string                   `json:"owner_name"`
	GSTNumber                string                   `json:"gst_number"`
	PreferredOperatingRoutes string                   `json:"preferred_operating_routes"`
	AddressLine1             string                   `json:"address_line1"`
	City                     string                   `json:"city"`
	State                    string                   `json:"state"`
	PANNumber                string                   `json:"pan_number"`
	TDSDeclaration           string                   `json:"tds_declaration"`
	Remark                   string                   `json:"remark"`
	BankAccountHolderName    string                   `json:"bank_account_holder_name"`
	BankAccountNumber        string                   `json:"bank_account_number"`
	BankName                 string                   `json:"bank_name"`
	BankIFSCCode             string                   `json:"bank_ifsc_code"`
	PancardImg               string                   `json:"pancard_img"`
	BankPassbookOrChequeImg  string                   `json:"bank_passbook_or_cheque_img"`
	Status                   string                   `json:"status"`
	IsActive                 int64                    `json:"is_active"`
	OrgID                    int64                    `json:"org_id"`
	LoginType                string                   `json:"login_type"`
	ContactInfo              []VendorContactInfo      `json:"contact_info"`
	Vehicles                 []VehicleModel           `json:"vehicle"`
	DeclarationDocuments     []DeclarationDocumentObj `json:"declaration_years"`
}

type VendorV1Entries struct {
	VendorEntiry *[]VendorResponse `json:"vendors"`
	Total        int64             `json:"total"`
	Limit        int64             `json:"limit"`
	OffSet       int64             `json:"offset"`
}

type VendorV1Update struct {
	VendorName               string                `json:"vendor_name"`
	VendorCode               string                `json:"vendor_code"`
	OwnerName                string                `json:"owner_name"`
	GSTNumber                string                `json:"gst_number"`
	PreferredOperatingRoutes string                `json:"preferred_operating_routes"`
	AddressLine1             string                `json:"address_line1"`
	City                     string                `json:"city"`
	State                    string                `json:"state"`
	PANNumber                string                `json:"pan_number"`
	TDSDeclaration           string                `json:"tds_declaration"`
	Remark                   string                `json:"remark"`
	BankAccountHolderName    string                `json:"bank_account_holder_name"`
	BankAccountNumber        string                `json:"bank_account_number"`
	BankName                 string                `json:"bank_name"`
	BankIFSCCode             string                `json:"bank_ifsc_code"`
	PancardImg               string                `json:"pancard_img"`
	BankPassbookOrChequeImg  string                `json:"bank_passbook_or_cheque_img"`
	Status                   string                `json:"status"`
	IsActive                 int                   `json:"is_active"`
	OrgID                    int64                 `json:"org_id"`
	LoginType                string                `json:"login_type"`
	ContactInfo              []VendorContactInfo   `json:"contact_info"`
	Vehicles                 []VehicleModel        `json:"vehicle"`
	DeclarationDocument      []DeclarationDocument `json:"declaration_document"`
}

type DeclarationDocument struct {
	DeclarationYear     string `json:"declaration_year"`
	DeclarationDocImage string `json:"declaration_doc_image"`
	VendorID            int64  `json:"vendor_id"`
}
type DeclarationDocumentObj struct {
	DeclarationYear       string `json:"declaration_year"`
	DeclarationDocImage   string `json:"declaration_doc_image"`
	VendorID              int64  `json:"vendor_id"`
	DeclarationYearInfoID int64  `json:"declaration_year_info_id"`
}

type VehicleModel struct {
	VehicleID           int64  `json:"vehicle_id"`
	VendorID            int64  `json:"vendor_id"`
	VehicleNumber       string `json:"vehicle_number"`
	VehicleType         string `json:"vehicle_type"`
	VehicleMake         string `json:"vehicle_make"`
	VehicleModel        string `json:"vehicle_model"`
	PermitType          string `json:"permit_type"`
	VehicleSize         string `json:"vehicle_size"`
	ClosedOpen          string `json:"closed_open"`
	VehicleCapacityTons string `json:"vehicle_capacity_tons"`
	RCExpiryDoc         string `json:"rc_expiry_doc"`
	InsuranceDoc        string `json:"insurance_doc"`
	PUCCExpiryDoc       string `json:"pucc_expiry_doc"`
	NPExpireDoc         string `json:"np_expire_doc"`
	FitnessExpiryDoc    string `json:"fitness_expiry_doc"`
	TaxExpiryDoc        string `json:"tax_expiry_doc"`
	MPExpireDoc         string `json:"mp_expire_doc"`
}

type VendorAndVehicleUpload struct {
	PancardImg              string `json:"pancard_img"`
	BankPassbookORChequeImg string `json:"bank_passbook_or_cheque_img"`
	RCExpiryDoc             string `json:"rc_expiry_doc"`
	InsuranceDoc            string `json:"insurance_doc"`
	PUCCExpiryDoc           string `json:"pucc_expiry_doc"`
	NPExpireDoc             string `json:"np_expire_doc"`
	FitnessExpiryDoc        string `json:"fitness_expiry_doc"`
	TaxExpiryDoc            string `json:"tax_expiry_doc"`
	MPExpireDoc             string `json:"mp_expire_doc"`
	VehicleID               string `json:"vehicle_id"`
	VendorID                string `json:"vendor_id"`
	TdsDeclaration          string `json:"tdsDeclaration"`
}

type UploadResponse struct {
	//VendorCode string `json:"vendor_code"`
	ImagePath string `json:"image_path"`
	Message   string `json:"message"`
}
