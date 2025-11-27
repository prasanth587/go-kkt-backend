package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type CustomerObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewCustomerObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *CustomerObj {
	return &CustomerObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

var (
// UPDATE_ERROR = "ERROR UpdateUserLoginAterCredentialSuccess update: %v "
// VERIFY_ERROR = "ERROR VerifyCredentials: %v "
// INVALIDE_    = "invalide credentials"
// LOGIN_FAILED = "login failed  - "
)

type CustomerDao interface {
	// CreateCustomer(customerReq dtos.CustomerReq) error
	// GetCustomers(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.CustomerRes, error)
	// UpdateCustomer(customerId int64, customerUpdate dtos.CustomerUpdate) error
	// GetCustomer(customerId int64) (*dtos.CustomerRes, error)
	// UpdateCustomerActiveStatus(customerId, isActive int64) error
	// UpdateCustomerStatus(customerId int64, status string) error
	// UpdateCustomerImagePath(updateQuery string) error

	//V2
	GetTotalCount(whereQuery string) int64
	BuildWhereQuery(orgId int64, customerId, searchText string) string
	CreateCustomerV1(customerReq dtos.CustomersReq) (int64, error)
	CreateCustomerContactInfo(customerId int64, contactInfo dtos.CustomersContactReq) error
	GetCustomersV1(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.CustomersResV1, error)
	GetCustomerContactInfo(customerId int64) (*[]dtos.CustomersContact, error)
	UpdateCustomerV1(customerId int64, customerReq dtos.CustomersUpdateReq) error
	GetCustomerV1(customerId int64) (*dtos.CustomersResV1, error)
	DeleteCustomerContactInfo(customerId int64) error
	UpdateCustomerStatusV1(customerId int64, status string) error
	UpdateCustomerActiveStatusV1(customerId, isActive int64) error
	SaveCustomerAggrement(customerReq dtos.CustomersAggrementUpload, aggrementPath string) error
	UpdateCustomerAggrement(customerReq dtos.CustomersAggrementUpload, aggrementPath string) error
	GetCustomerAggrements(customerId int64) (*[]dtos.CustomersAggrement, error)
}

func (br *CustomerObj) DeleteCustomerContactInfo(customerId int64) error {

	deleteQuery := fmt.Sprintf(`DELETE FROM contact_info WHERE customer_id = '%v';`, customerId)

	br.l.Info("DeleteCustomerContactInfo query:", deleteQuery)

	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(deleteQuery)
	if err != nil {
		br.l.Error("Error db.Exec(DeleteCustomerContactInfo): ", err)
		return err
	}

	br.l.Info("deleted successfully: ", roleResult)

	return nil
}

func (rl *CustomerObj) BuildWhereQuery(orgId int64, customerId, searchText string) string {

	whereQuery := fmt.Sprintf("WHERE org_id = '%v'", orgId)

	if customerId != "" {
		whereQuery = fmt.Sprintf(" %v AND customer_id = '%v'", whereQuery, customerId)
	}

	if searchText != "" {
		whereQuery = fmt.Sprintf(" %v AND (customer_code LIKE '%%%v%%' OR customer_name LIKE '%%%v%%' OR branch_name LIKE '%%%v%%' OR payment_terms LIKE '%%%v%%' OR gst_number LIKE '%%%v%%' OR kkt_responsible_emp LIKE '%%%v%%' OR city LIKE '%%%v%%' OR address LIKE '%%%v%%' OR company_type LIKE '%%%v%%' ) ", whereQuery, searchText, searchText, searchText, searchText, searchText, searchText, searchText, searchText, searchText)
	}

	rl.l.Info("customer whereQuery:\n ", whereQuery)

	return whereQuery
}

func (rl *CustomerObj) GetTotalCount(whereQuery string) int64 {

	countQuery := fmt.Sprintf(`SELECT count(*) FROM customers %v`, whereQuery)
	rl.l.Info(" GetTotalCount select query: ", countQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return 0
	}

	return count.Int64
}

// func (rl *CustomerObj) CreateCustomer(customerReq dtos.CustomerReq) error {

// 	rl.l.Info("CreateCustomer: ", customerReq.CustomerCode)

// 	createCustomerQuery := fmt.Sprintf(`INSERT INTO customer (
//         branch_id, customer_code, customer_name, address, city,
// 		state, pincode, country, gstin_type, gstin_no,
// 		pan_number, contact_person, mobile_number, alternative_number, email_id,
// 		status, branch_sales_person_name, branch_responsible_person, employee_code, business_start_date,
// 		is_active, org_id, credibility_days, company_type)
//         VALUES
//         ('%v', '%v', '%v', '%v', '%v',
// 		'%v', '%v', '%v', '%v', '%v',
// 		'%v', '%v', '%v', '%v', '%v',
// 		'%v', '%v', '%v', '%v', '%v',
// 		'%v', '%v', '%v','%v')`,
// 		customerReq.BranchId, customerReq.CustomerCode, customerReq.CustomerName, customerReq.Address, customerReq.City,
// 		customerReq.State, customerReq.PinCode, customerReq.Country, customerReq.GSTNType, customerReq.GSTNNo,
// 		customerReq.PanNumber, customerReq.ContactPerson, customerReq.MobileNumber, customerReq.AlternativeNumber, customerReq.EmailID,
// 		customerReq.Status, customerReq.BranchSalesPersonName, customerReq.BranchResponsiblePerson, customerReq.EmployeeCode, customerReq.BusinessStartDate,
// 		customerReq.IsActive, customerReq.OrgId, customerReq.CredibilityDays, customerReq.CompanyType)

// 	rl.l.Info("ceateCustomerQuery:\n ", createCustomerQuery)

// 	result, err := rl.dbConnMSSQL.GetQueryer().Exec(createCustomerQuery)
// 	if err != nil {
// 		rl.l.Error("Error db.Exec(CreateCustomer): ", err)
// 		return err
// 	}
// 	createdId, err := result.LastInsertId()
// 	if err != nil {
// 		rl.l.Error("Error db.Exec(CreateCustomer):", createdId, err)
// 		return err
// 	}
// 	rl.l.Info("Customer created successfully: ", createdId, customerReq.CustomerCode)
// 	return nil
// }

func (rl *CustomerObj) CreateCustomerV1(customerReq dtos.CustomersReq) (int64, error) {

	createCustomerQuery := fmt.Sprintf(`INSERT INTO customers (
    customer_code, customer_name, branch_name, nick_name, payment_terms,
    gst_number, remark, kkt_responsible_emp_id, fassi, pan_number,
    address, city, state, pincode, country,
    status, is_active, org_id, company_type, branch_id
	) VALUES (
		'%v', '%v', '%v', '%v', '%v',
		'%v', '%v', '%v', '%v', '%v',
		'%v', '%v', '%v', '%v', '%v',
		'%v', '%v', '%v', '%v', '%v'
	)`,
		customerReq.CustomerCode, customerReq.CustomerName, customerReq.BranchName, customerReq.NickName, customerReq.PaymentTerms,
		customerReq.GSTNumber, customerReq.Remark, customerReq.KKTResponsibleEmpID, customerReq.Fassi, customerReq.PanNumber,
		customerReq.Address, customerReq.City, customerReq.State, customerReq.PinCode, customerReq.Country,
		customerReq.Status, customerReq.IsActive, customerReq.OrgId, customerReq.CompanyType, customerReq.BranchId,
	)

	rl.l.Info("ceateCustomerQuery:\n ", createCustomerQuery)

	result, err := rl.dbConnMSSQL.GetQueryer().Exec(createCustomerQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(CreateCustomer): ", err)
		return 0, err
	}
	createdId, err := result.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(CreateCustomer):", createdId, err)
		return 0, err
	}
	rl.l.Info("Customer created successfully: ", createdId, customerReq.CustomerCode)

	return createdId, nil
}

func (cb *CustomerObj) CreateCustomerContactInfo(customerId int64, contactReq dtos.CustomersContactReq) error {

	createContactInfoQuery := fmt.Sprintf(`INSERT INTO contact_info (
    customer_id, contact_person_name, post, email_id, contact_nummber1, contact_nummber2
	) VALUES (
		'%v', '%v', '%v', '%v', '%v', '%v'
	)`,
		customerId,
		contactReq.ContactPerson,
		contactReq.Post,
		contactReq.EmailID,
		contactReq.ContactNumber1,
		contactReq.ContactNumber2,
	)
	cb.l.Info("createContactInfoQuery:\n ", createContactInfoQuery)

	result, err := cb.dbConnMSSQL.GetQueryer().Exec(createContactInfoQuery)
	if err != nil {
		cb.l.Error("Error db.Exec(CreateCustomerContactInfo): ", err)
		return err
	}
	createdId, err := result.LastInsertId()
	if err != nil {
		cb.l.Error("Error db.Exec(CreateCustomerContactInfo):", createdId, err)
		return err
	}
	cb.l.Info("Customer Contact Into created successfully: ", createdId, contactReq.ContactPerson)
	return nil
}

// func (rl *CustomerObj) GetCustomers(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.CustomerRes, error) {
// 	list := []dtos.CustomerRes{}

// 	whereQuery = fmt.Sprintf(" %v ORDER BY updated_at DESC LIMIT %v OFFSET %v;", whereQuery, limit, offset)

// 	rl.l.Info("customer whereQuery:\n ", whereQuery)

// 	customerQuery := fmt.Sprintf(`SELECT customer_id, branch_id, customer_code, customer_name, address, city, state, pincode, country, gstin_type, gstin_no, pan_number, contact_person, mobile_number, alternative_number, email_id, status, branch_sales_person_name, branch_responsible_person, employee_code, business_start_date, is_active, credibility_days, company_type, agreement_doc_image FROM customer %v `, whereQuery)

// 	rl.l.Info("get customers Query:\n ", customerQuery)

// 	rows, err := rl.dbConnMSSQL.GetQueryer().Query(customerQuery)
// 	if err != nil {
// 		rl.l.Error("Error Customers ", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var customerName, customerCode, mobileNumber, contactPerson, alternativeNumber, address, status, city, state,
// 			companyType, emailId, agreementDocImage, country, gstinType, gstinNumber, panNumber,
// 			branchSalesPersonName, branchResponsiblePerson, employeeCode, businessStartDate sql.NullString

// 		var branchID, cusId, isActive, credibilityDays, pinCode sql.NullInt64

// 		customerRes := &dtos.CustomerRes{}
// 		err := rows.Scan(&cusId, &branchID, &customerCode, &customerName, &address, &city,
// 			&state, &pinCode, &country, &gstinType, &gstinNumber, &panNumber,
// 			&contactPerson, &mobileNumber, &alternativeNumber, &emailId, &status, &branchSalesPersonName,
// 			&branchResponsiblePerson, &employeeCode, &businessStartDate, &isActive, &credibilityDays,
// 			&companyType, &agreementDocImage)
// 		if err != nil {
// 			rl.l.Error("Error GetCustomers scan: ", err)
// 			return nil, err
// 		}
// 		customerRes.CustomerId = cusId.Int64
// 		customerRes.BranchId = branchID.Int64
// 		customerRes.CustomerCode = customerCode.String
// 		customerRes.CustomerName = customerName.String
// 		customerRes.Address = address.String
// 		customerRes.City = city.String

// 		customerRes.State = state.String
// 		customerRes.PinCode = pinCode.Int64
// 		customerRes.Country = country.String
// 		customerRes.GSTNType = gstinType.String
// 		customerRes.GSTNNo = gstinNumber.String
// 		customerRes.PanNumber = panNumber.String

// 		customerRes.ContactPerson = contactPerson.String
// 		customerRes.MobileNumber = mobileNumber.String
// 		customerRes.AlternativeNumber = alternativeNumber.String
// 		customerRes.EmailID = emailId.String
// 		customerRes.Status = status.String
// 		customerRes.BranchSalesPersonName = branchSalesPersonName.String

// 		customerRes.BranchResponsiblePerson = branchResponsiblePerson.String
// 		customerRes.EmployeeCode = employeeCode.String
// 		customerRes.BusinessStartDate = businessStartDate.String
// 		customerRes.IsActive = isActive.Int64
// 		customerRes.CredibilityDays = credibilityDays.Int64

// 		customerRes.CompanyType = companyType.String
// 		customerRes.AgreementDocImage = agreementDocImage.String
// 		customerRes.OrgId = orgId
// 		list = append(list, *customerRes)
// 	}
// 	return &list, nil
// }

func (rl *CustomerObj) GetCustomerV1(customerId int64) (*dtos.CustomersResV1, error) {

	customerQuery := fmt.Sprintf(`
    SELECT customer_id, customer_code, customer_name, branch_name, branch_id, nick_name, 
	payment_terms, gst_number,remark, kkt_responsible_emp_id,fassi, pan_number, 
	address,city,state,pincode,country,status,
	is_active, org_id,company_type
    FROM customers where customer_id = '%v' `, customerId)

	rl.l.Info("get customers Query:\n ", customerQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(customerQuery)

	var customerName, customerCode, address, status, city, state,
		companyType, country, fassi, panNumber, gstNumber, remark,
		nickName, paymentTerms, branchName sql.NullString

	var branchID, cusId, isActive, pinCode, orgId, kktResponsibleEmpId sql.NullInt64

	customerRes := &dtos.CustomersResV1{}
	err := row.Scan(&cusId, &customerCode, &customerName, &branchName, &branchID, &nickName,
		&paymentTerms, &gstNumber, &remark, &kktResponsibleEmpId, &fassi, &panNumber,
		&address, &city, &state, &pinCode, &country, &status,
		&isActive, &orgId, &companyType)
	if err != nil {
		rl.l.Error("Error GetCustomers scan: ", err)
		return nil, err
	}
	customerRes.CustomerId = cusId.Int64
	customerRes.CustomerCode = customerCode.String
	customerRes.CustomerName = customerName.String
	customerRes.BranchName = branchName.String
	customerRes.BranchId = branchID.Int64
	customerRes.NickName = nickName.String
	customerRes.PaymentTerms = paymentTerms.String
	customerRes.GSTNumber = gstNumber.String
	customerRes.Remark = remark.String
	customerRes.KKTResponsibleEmpID = kktResponsibleEmpId.Int64
	customerRes.Fassi = fassi.String
	customerRes.PanNumber = panNumber.String
	customerRes.Address = address.String
	customerRes.City = city.String
	customerRes.State = state.String
	customerRes.PinCode = pinCode.Int64
	customerRes.Country = country.String
	customerRes.Status = status.String
	customerRes.IsActive = isActive.Int64
	customerRes.OrgId = orgId.Int64
	customerRes.CompanyType = companyType.String
	return customerRes, nil
}

func (rl *CustomerObj) GetCustomersV1(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.CustomersResV1, error) {

	list := []dtos.CustomersResV1{}

	whereQuery = fmt.Sprintf(" %v ORDER BY updated_at DESC LIMIT %v OFFSET %v;", whereQuery, limit, offset)

	rl.l.Info("customer whereQuery:\n ", whereQuery)

	customerQuery := fmt.Sprintf(`
    SELECT customer_id, customer_code, customer_name, branch_name, branch_id, nick_name, 
	payment_terms, gst_number,remark, kkt_responsible_emp_id,fassi, pan_number, 
	address,city,state,pincode,country,status,
	is_active, org_id,company_type
    FROM customers %v`, whereQuery)

	rl.l.Info("get customers Query:\n ", customerQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(customerQuery)
	if err != nil {
		rl.l.Error("Error Customers ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customerName, customerCode, address, status, city, state,
			companyType, country, fassi, panNumber, gstNumber, remark,
			nickName, paymentTerms, branchName sql.NullString

		var branchID, cusId, isActive, pinCode, kktResponsibleEmpId sql.NullInt64

		customerRes := &dtos.CustomersResV1{}
		err := rows.Scan(&cusId, &customerCode, &customerName, &branchName, &branchID, &nickName,
			&paymentTerms, &gstNumber, &remark, &kktResponsibleEmpId, &fassi, &panNumber,
			&address, &city, &state, &pinCode, &country, &status,
			&isActive, &orgId, &companyType)
		if err != nil {
			rl.l.Error("Error GetCustomers scan: ", err)
			return nil, err
		}
		customerRes.CustomerId = cusId.Int64
		customerRes.CustomerCode = customerCode.String
		customerRes.CustomerName = customerName.String
		customerRes.BranchName = branchName.String
		customerRes.BranchId = branchID.Int64
		customerRes.NickName = nickName.String
		customerRes.PaymentTerms = paymentTerms.String
		customerRes.GSTNumber = gstNumber.String
		customerRes.Remark = remark.String
		customerRes.KKTResponsibleEmpID = kktResponsibleEmpId.Int64
		customerRes.Fassi = fassi.String
		customerRes.PanNumber = panNumber.String
		customerRes.Address = address.String
		customerRes.City = city.String
		customerRes.State = state.String
		customerRes.PinCode = pinCode.Int64
		customerRes.Country = country.String
		customerRes.Status = status.String
		customerRes.IsActive = isActive.Int64
		customerRes.OrgId = orgId
		customerRes.CompanyType = companyType.String
		list = append(list, *customerRes)
	}
	return &list, nil
}

func (cs *CustomerObj) GetCustomerContactInfo(customerId int64) (*[]dtos.CustomersContact, error) {

	list := []dtos.CustomersContact{}

	contactQuery := fmt.Sprintf(`
    SELECT contact_info_id, contact_person_name, post, email_id, contact_nummber1, contact_nummber2
    FROM contact_info WHERE customer_id = '%v'`, customerId)

	cs.l.Info("contactQuery:\n ", contactQuery)

	rows, err := cs.dbConnMSSQL.GetQueryer().Query(contactQuery)
	if err != nil {
		cs.l.Error("Error Customers ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contactPersonName, post, emailId, contactNumber1, contactNumber2 sql.NullString

		var contactInfoId sql.NullInt64

		contact := &dtos.CustomersContact{}
		err := rows.Scan(&contactInfoId, &contactPersonName, &post, &emailId, &contactNumber1, &contactNumber2)
		if err != nil {
			cs.l.Error("Error GetCustomers scan: ", err)
			return nil, err
		}
		contact.ContactInfoId = contactInfoId.Int64
		contact.ContactPerson = contactPersonName.String
		contact.Post = post.String
		contact.EmailID = emailId.String
		contact.ContactNumber1 = contactNumber1.String
		contact.ContactNumber2 = contactNumber2.String

		list = append(list, *contact)
	}

	return &list, nil
}

func (rl *CustomerObj) UpdateCustomerV1(customerId int64, customerReq dtos.CustomersUpdateReq) error {

	updateCustomerQuery := fmt.Sprintf(`
    UPDATE customers SET
        customer_code = '%v', customer_name = '%v',
        branch_name = '%v', nick_name = '%v',
        payment_terms = '%v', gst_number = '%v',
        remark = '%v', kkt_responsible_emp_id = '%v',
        fassi = '%v', pan_number = '%v', address = '%v', city = '%v', state = '%v',
        pincode = %v, country = '%v', status = '%v', is_active = %v,
        org_id = %v, company_type = '%v', branch_id = '%v'
    WHERE customer_id = %v;`,
		customerReq.CustomerCode, customerReq.CustomerName, customerReq.BranchName,
		customerReq.NickName, customerReq.PaymentTerms, customerReq.GSTNumber,
		customerReq.Remark, customerReq.KKTResponsibleEmpID, customerReq.Fassi,
		customerReq.PanNumber, customerReq.Address, customerReq.City,
		customerReq.State, customerReq.PinCode, customerReq.Country, customerReq.Status,
		customerReq.IsActive, customerReq.OrgId, customerReq.CompanyType, customerReq.BranchId,
		customerId,
	)

	rl.l.Info("UpdateCustomer query: ", updateCustomerQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateCustomerQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateCustomer): ", err)
		return err
	}
	rl.l.Info("Customer updated successfully: ", customerId)
	return nil
}

// func (rl *CustomerObj) UpdateCustomer(customerId int64, customerReq dtos.CustomerUpdate) error {

// 	updateCustomerQuery := fmt.Sprintf(`
//     UPDATE customer SET
//         branch_id = '%v',
//         customer_code = '%v',
//         customer_name = '%v',
//         address = '%v',
//         city = '%v',
//         state = '%v',
//         pincode = '%v',
//         country = '%v',
//         gstin_type = '%v',
//         gstin_no = '%v',
//         pan_number = '%v',
//         contact_person = '%v',
//         mobile_number = '%v',
//         alternative_number = '%v',
//         email_id = '%v',
//         status = '%v',
//         branch_sales_person_name = '%v',
//         branch_responsible_person = '%v',
//         employee_code = '%v',
//         business_start_date = '%v',
//         is_active = %v,
//         org_id = '%v',
//         credibility_days = %v,
//         company_type = '%v'
//     WHERE customer_id = '%v';`,
// 		customerReq.BranchId, customerReq.CustomerCode, customerReq.CustomerName, customerReq.Address, customerReq.City,
// 		customerReq.State, customerReq.PinCode, customerReq.Country, customerReq.GSTNType, customerReq.GSTNNo,
// 		customerReq.PanNumber, customerReq.ContactPerson, customerReq.MobileNumber, customerReq.AlternativeNumber, customerReq.EmailID,
// 		customerReq.Status, customerReq.BranchSalesPersonName, customerReq.BranchResponsiblePerson, customerReq.EmployeeCode, customerReq.BusinessStartDate,
// 		customerReq.IsActive, customerReq.OrgId, customerReq.CredibilityDays, customerReq.CompanyType, customerId)

// 	rl.l.Info("UpdateCustomer query: ", updateCustomerQuery)

// 	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateCustomerQuery)
// 	if err != nil {
// 		rl.l.Error("Error db.Exec(UpdateCustomer): ", err)
// 		return err
// 	}
// 	rl.l.Info("Customer updated successfully: ", customerId)
// 	return nil
// }

// func (rl *CustomerObj) GetCustomer(customerId int64) (*dtos.CustomerRes, error) {

// 	customerQuery := fmt.Sprintf(`SELECT customer_id, branch_id, customer_code, customer_name, address, city, state,
// 	pincode, country, gstin_type, gstin_no, pan_number, contact_person,
// 	mobile_number, alternative_number, email_id, status, branch_sales_person_name,
// 	branch_responsible_person, employee_code, business_start_date, is_active,
// 	credibility_days, company_type, agreement_doc_image,
// 	org_id FROM customers WHERE customer_id = '%v';`, customerId)

// 	rl.l.Info("GET customer Query:\n ", customerQuery)

// 	var customerName, customerCode, mobileNumber, contactPerson, alternativeNumber, address, status, city, state,
// 		companyType, emailId, agreementDocImage, country, gstinType, gstinNumber, panNumber,
// 		branchSalesPersonName, branchResponsiblePerson, employeeCode, businessStartDate sql.NullString

// 	var branchID, isActive, credibilityDays, pinCode, orgId sql.NullInt64

// 	row := rl.dbConnMSSQL.GetQueryer().QueryRow(customerQuery)

// 	customerRes := dtos.CustomerRes{}
// 	err := row.Scan(&customerId, &branchID, &customerCode, &customerName, &address, &city,
// 		&state, &pinCode, &country, &gstinType, &gstinNumber, &panNumber,
// 		&contactPerson, &mobileNumber, &alternativeNumber, &emailId, &status, &branchSalesPersonName,
// 		&branchResponsiblePerson, &employeeCode, &businessStartDate, &isActive, &credibilityDays,
// 		&companyType, &agreementDocImage, &orgId)
// 	if err != nil {
// 		rl.l.Error("Error GetCustomers scan: ", err)
// 		return nil, err
// 	}
// 	customerRes.CustomerId = customerId
// 	customerRes.BranchId = branchID.Int64
// 	customerRes.CustomerCode = customerCode.String
// 	customerRes.CustomerName = customerName.String
// 	customerRes.Address = address.String
// 	customerRes.City = city.String

// 	customerRes.State = state.String
// 	customerRes.PinCode = pinCode.Int64
// 	customerRes.Country = country.String
// 	customerRes.GSTNType = gstinType.String
// 	customerRes.GSTNNo = gstinNumber.String
// 	customerRes.PanNumber = panNumber.String

// 	customerRes.ContactPerson = contactPerson.String
// 	customerRes.MobileNumber = mobileNumber.String
// 	customerRes.AlternativeNumber = alternativeNumber.String
// 	customerRes.EmailID = emailId.String
// 	customerRes.Status = status.String
// 	customerRes.BranchSalesPersonName = branchSalesPersonName.String

// 	customerRes.BranchResponsiblePerson = branchResponsiblePerson.String
// 	customerRes.EmployeeCode = employeeCode.String
// 	customerRes.BusinessStartDate = businessStartDate.String
// 	customerRes.IsActive = isActive.Int64
// 	customerRes.CredibilityDays = credibilityDays.Int64

// 	customerRes.CompanyType = companyType.String
// 	customerRes.AgreementDocImage = agreementDocImage.String
// 	customerRes.OrgId = orgId.Int64
// 	return &customerRes, nil
// }

// func (rl *CustomerObj) UpdateCustomerActiveStatus(customerId, isActive int64) error {

// 	updateQuery := fmt.Sprintf(`UPDATE customer SET is_active = '%v' WHERE customer_id = '%v'`, isActive, customerId)

// 	rl.l.Info("UpdateCustomerActiveStatus Update query ", updateQuery)

// 	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
// 	if err != nil {
// 		rl.l.Error("Error db.Exec(UpdateCustomerActiveStatus): ", err)
// 		return err
// 	}

// 	rl.l.Info("Customer status updated successfully: ", customerId)

// 	return nil
// }

func (rl *CustomerObj) UpdateCustomerActiveStatusV1(customerId, isActive int64) error {

	updateQuery := fmt.Sprintf(`UPDATE customers SET is_active = '%v' WHERE customer_id = '%v'`, isActive, customerId)

	rl.l.Info("UpdateCustomerActiveStatusV1 Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateCustomerActiveStatusV1): ", err)
		return err
	}

	rl.l.Info("Customer status updated successfully: ", customerId)

	return nil
}

// func (rl *CustomerObj) UpdateCustomerStatus(customerId int64, status string) error {

// 	updateQuery := fmt.Sprintf(`UPDATE customer SET status = '%v' WHERE customer_id = '%v'`, status, customerId)

// 	rl.l.Info("UpdateCustomerStatus Update query ", updateQuery)

// 	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
// 	if err != nil {
// 		rl.l.Error("Error db.Exec(UpdateCustomerStatus): ", err)
// 		return err
// 	}

// 	rl.l.Info("Customer status updated successfully: ", customerId)

// 	return nil
// }

func (rl *CustomerObj) UpdateCustomerStatusV1(customerId int64, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE customers SET status = '%v' WHERE customer_id = '%v'`, status, customerId)

	rl.l.Info("UpdateCustomerStatus Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateCustomerStatusV1): ", err)
		return err
	}

	rl.l.Info("Customer status updated successfully: ", customerId)

	return nil
}

// func (rl *CustomerObj) UpdateCustomerImagePath(updateQuery string) error {

// 	rl.l.Info("UpdateCustomerImagePath Update query: ", updateQuery)

// 	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
// 	if err != nil {
// 		rl.l.Error("Error db.Exec(UpdateCustomerImagePath) CustomerObj: ", err)
// 		return err
// 	}

// 	return nil
// }

func (rl *CustomerObj) SaveCustomerAggrement(customerReq dtos.CustomersAggrementUpload, aggrementPath string) error {

	rl.l.Info("SaveCustomerAggrement: ", customerReq.AggrementName)

	saveCustomerAggrement := fmt.Sprintf(`INSERT INTO customer_aggrement (
        customer_id, aggrement_number, aggrement_period, aggrement_name, aggrement_type, 
		aggrement_doc, remark)
        VALUES
        ('%v', '%v', '%v', '%v', '%v', '%v', '%v')`,
		customerReq.CustomerId, customerReq.AggrementNo, customerReq.Period, customerReq.AggrementName, customerReq.AggrementType,
		aggrementPath, customerReq.Remark)

	rl.l.Info("SaveCustomerAggrement:\n ", saveCustomerAggrement)

	result, err := rl.dbConnMSSQL.GetQueryer().Exec(saveCustomerAggrement)
	if err != nil {
		rl.l.Error("Error db.Exec(saveCustomerAggrement): ", err)
		return err
	}
	createdId, err := result.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(saveCustomerAggrement):", createdId, err)
		return err
	}
	rl.l.Info("Customer Aggrement successfully: ", createdId, customerReq.AggrementName)
	return nil
}

func (rl *CustomerObj) UpdateCustomerAggrement(customerReq dtos.CustomersAggrementUpload, aggrementPath string) error {

	updateCustomerAggrement := fmt.Sprintf(`UPDATE customer_aggrement
    SET customer_id = '%v', aggrement_number = '%v', aggrement_period = '%v',
        aggrement_name = '%v', aggrement_type = '%v', aggrement_doc = '%v', remark = '%v'
    WHERE
        aggrement_id = '%v'`,
		customerReq.CustomerId,
		customerReq.AggrementNo,
		customerReq.Period,
		customerReq.AggrementName,
		customerReq.AggrementType,
		aggrementPath,
		customerReq.Remark,
		customerReq.AggrementId)

	rl.l.Info("updateCustomerAggrement Update query ", updateCustomerAggrement)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateCustomerAggrement)
	if err != nil {
		rl.l.Error("Error db.Exec(updateCustomerAggrement): ", err)
		return err
	}

	rl.l.Info("customer aggrement update successfully: ", customerReq.CustomerId, customerReq.AggrementName)

	return nil
}

func (rl *CustomerObj) GetCustomerAggrements(customerId int64) (*[]dtos.CustomersAggrement, error) {

	aggrementQuery := fmt.Sprintf(`SELECT aggrement_id, customer_id, aggrement_number, aggrement_period, aggrement_name, aggrement_type, aggrement_doc, remark FROM customer_aggrement WHERE customer_id = '%v' ORDER BY aggrement_id ASC;`, customerId)

	list := []dtos.CustomersAggrement{}

	rl.l.Info("aggrementQuery:\n ", aggrementQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(aggrementQuery)
	if err != nil {
		rl.l.Error("Error Customers ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var aggrementNumber, aggrementPeriod, aggrementName, aggrementType, aggrementDoc, remark sql.NullString

		var aggrementId, customerId sql.NullInt64

		aggrement := &dtos.CustomersAggrement{}
		err := rows.Scan(&aggrementId, &customerId, &aggrementNumber, &aggrementPeriod, &aggrementName, &aggrementType, &aggrementDoc, &remark)
		if err != nil {
			rl.l.Error("Error GetCustomers scan: ", err)
			return nil, err
		}
		aggrement.AggrementId = aggrementId.Int64
		aggrement.CustomerId = customerId.Int64
		aggrement.AggrementNo = aggrementNumber.String
		aggrement.Period = aggrementPeriod.String
		aggrement.AggrementName = aggrementName.String
		aggrement.AggrementType = aggrementType.String
		aggrement.AggrementDoc = aggrementDoc.String
		aggrement.Remark = remark.String
		list = append(list, *aggrement)
	}

	return &list, nil

}
