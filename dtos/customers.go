package dtos

type CustomersReq struct {
	CustomerCode        string                `json:"customer_code"`
	CustomerName        string                `json:"customer_name"`
	BranchName          string                `json:"branch_name"`
	BranchId            string                `json:"branch_id"`
	NickName            string                `json:"nick_name"`
	PaymentTerms        string                `json:"payment_terms"`
	GSTNumber           string                `json:"gst_number"`
	Remark              string                `json:"remark"`
	KKTResponsibleEmpID int64                 `json:"kkt_responsible_emp_id"`
	Fassi               string                `json:"fassi"`
	PanNumber           string                `json:"pan_number"`
	Address             string                `json:"address"`
	City                string                `json:"city"`
	State               string                `json:"state"`
	PinCode             int64                 `json:"pincode"`
	Country             string                `json:"country"`
	Status              string                `json:"status"`
	IsActive            int64                 `json:"is_active"`
	OrgId               int64                 `json:"org_id"`
	CompanyType         string                `json:"company_type"`
	EmployeeCode        string                `json:"employee_code"`
	ContactInfo         []CustomersContactReq `json:"contact_info"`
}

type CustomersUpdateReq struct {
	CustomerCode        string                `json:"customer_code"`
	CustomerName        string                `json:"customer_name"`
	BranchName          string                `json:"branch_name"`
	BranchId            string                `json:"branch_id"`
	NickName            string                `json:"nick_name"`
	PaymentTerms        string                `json:"payment_terms"`
	GSTNumber           string                `json:"gst_number"`
	Remark              string                `json:"remark"`
	KKTResponsibleEmpID int64                 `json:"kkt_responsible_emp_id"`
	Fassi               string                `json:"fassi"`
	PanNumber           string                `json:"pan_number"`
	Address             string                `json:"address"`
	City                string                `json:"city"`
	State               string                `json:"state"`
	PinCode             int64                 `json:"pincode"`
	Country             string                `json:"country"`
	Status              string                `json:"status"`
	IsActive            int64                 `json:"is_active"`
	OrgId               int64                 `json:"org_id"`
	CompanyType         string                `json:"company_type"`
	ContactInfo         []CustomersContactReq `json:"contact_info"`
}

type CustomersContactReq struct {
	ContactPerson  string `json:"contact_person_name"`
	Post           string `json:"post"`
	EmailID        string `json:"email_id"`
	ContactNumber1 string `json:"contact_nummber1"`
	ContactNumber2 string `json:"contact_nummber2"`
}

type CustomersResV1 struct {
	CustomerId          int64                 `json:"customer_id"`
	CustomerCode        string                `json:"customer_code"`
	CustomerName        string                `json:"customer_name"`
	BranchName          string                `json:"branch_name"`
	BranchId            int64                 `json:"branch_id"`
	NickName            string                `json:"nick_name"`
	PaymentTerms        string                `json:"payment_terms"`
	GSTNumber           string                `json:"gst_number"`
	Remark              string                `json:"remark"`
	KKTResponsibleEmpID int64                 `json:"kkt_responsible_emp_id"`
	Fassi               string                `json:"fassi"`
	PanNumber           string                `json:"pan_number"`
	Address             string                `json:"address"`
	City                string                `json:"city"`
	State               string                `json:"state"`
	PinCode             int64                 `json:"pincode"`
	Country             string                `json:"country"`
	Status              string                `json:"status"`
	IsActive            int64                 `json:"is_active"`
	OrgId               int64                 `json:"org_id"`
	CompanyType         string                `json:"company_type"`
	EmployeeCode        string                `json:"employee_code"`
	ContactInfo         *[]CustomersContact   `json:"contact_info"`
	CustomersAggrement  *[]CustomersAggrement `json:"aggrement_info"`
}
type CustomersContact struct {
	ContactInfoId  int64  `json:"contact_info_id"`
	ContactPerson  string `json:"contact_person_name"`
	Post           string `json:"post"`
	EmailID        string `json:"email_id"`
	ContactNumber1 string `json:"contact_nummber1"`
	ContactNumber2 string `json:"contact_nummber2"`
}

type CustomersResponse struct {
	CustomerEntiry *[]CustomersResV1 `json:"customers"`
	Total          int64             `json:"total"`
	Limit          int64             `json:"limit"`
	OffSet         int64             `json:"offset"`
}

type CustomersAggrementUpload struct {
	AggrementId   string `json:"aggrement_id"`
	Period        string `json:"period"`
	AggrementName string `json:"aggrement_name"`
	AggrementType string `json:"aggrement_type"`
	Remark        string `json:"remark"`
	AggrementNo   string `json:"aggrementNo"`
	CustomerId    int64  `json:"customer_id"`
}

type CustomersAggrement struct {
	AggrementId   int64  `json:"aggrement_id"`
	Period        string `json:"period"`
	AggrementName string `json:"aggrement_name"`
	AggrementType string `json:"aggrement_type"`
	Remark        string `json:"remark"`
	AggrementNo   string `json:"aggrementNo"`
	CustomerId    int64  `json:"customer_id"`
	AggrementDoc  string `json:"aggrement_doc"`
}
