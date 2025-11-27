package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type VendorObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewVendorObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *VendorObj {
	return &VendorObj{
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

type VendorDao interface {
	CreateVendor(vendorReg dtos.VendorReq) error
	GetVendors(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.VendorRes, error)
	UpdateVendor(vendorId int64, veh dtos.VendorUpdate) error
	GetVendor(vendorId int64) (*dtos.VendorRes, error)
	UpdateVendorActiveStatus(vendorId, isActive int64) error
	UpdateVendorStatus(vendorId int64, status string) error
	UpdateVendorImagePath(updateQuery string) error
	GetTotalCount(whereQuery string) int64
	BuildWhereQuery(orgId int64, vendorId, searchText string) string

	CreateVendorV1(vendorReg dtos.VendorRequest) (int64, error)
	CreateVendorContactInfo(vendorId int64, contactReq dtos.VendorContactInfo) error
	GetVendorsV1(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.VendorResponse, error)
	GetTotalCountV1(whereQuery string) int64
	GetVendorContactInfo(vendorID int64) (*[]dtos.VendorContactInfo, error)
	GetVendorV1(vendorId int64) (*dtos.VendorResponse, error)
	UpdateVendorStatusV1(vendorId int64, status string) error
	UpdateVendorActiveStatusV1(vendorId, isActive int64) error
	UpdateVendorV1(vendorId int64, veh dtos.VendorV1Update) error
	DeleteVendorContactInfo(vendorId int64) error
	DeleteVendorVehiclesInfo(vendorId int64) error
	DeleteVendorDeclarations(vendorId int64) error

	/// VehicleModel
	CreateVehicle(vehicle dtos.VehicleModel) error
	GetVehiclesByVendorID(vendorID int64) (*[]dtos.VehicleModel, error)
	GetDeclarationDocByVendorID(vendorID int64) (*[]dtos.DeclarationDocumentObj, error)

	UpdateVendorImage(vendorId string, imageFor, imagePath string) error
	UpdateVehicleImage(vehicleId string, imageFor, imagePath string) error

	CreateVehicleDoclarationDoc(vehicle dtos.DeclarationDocument) error
}

func (rl *VendorObj) BuildWhereQuery(orgId int64, vendorId, searchText string) string {

	whereQuery := fmt.Sprintf("WHERE org_id = '%v'", orgId)

	if vendorId != "" {
		whereQuery = fmt.Sprintf(" %v AND vendor_id = '%v'", whereQuery, vendorId)
	}

	if searchText != "" {
		whereQuery = fmt.Sprintf(" %v AND (vendor_name LIKE '%%%v%%' OR vendor_code LIKE '%%%v%%' OR owner_name LIKE '%%%v%%' OR gst_number LIKE '%%%v%%' OR preferred_operating_routes LIKE '%%%v%%' OR address_line1 LIKE '%%%v%%' OR city LIKE '%%%v%%' OR remark LIKE '%%%v%%' OR bank_account_holder_name LIKE '%%%v%%' OR bank_name LIKE '%%%v%%' ) ", whereQuery, searchText, searchText, searchText, searchText, searchText, searchText, searchText, searchText, searchText, searchText)
	}

	rl.l.Info("vendor whereQuery:\n ", whereQuery)
	return whereQuery
}

func (rl *VendorObj) GetTotalCount(whereQuery string) int64 {
	countQuery := fmt.Sprintf(`SELECT count(*) FROM vendors %v`, whereQuery)
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

func (rl *VendorObj) GetTotalCountV1(whereQuery string) int64 {
	countQuery := fmt.Sprintf(`SELECT count(*) FROM vendors %v`, whereQuery)
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

func (rl *VendorObj) CreateVendor(vendorReg dtos.VendorReq) error {

	rl.l.Info("CreateVendor : ", vendorReg.VendorCode)

	createVendorQuery := fmt.Sprintf(`INSERT INTO vendor (
        vendor_name, vendor_code, mobile_number, contact_person, 
		alternative_number, address_line1, address_line2, city, 
		state, status, org_id, is_active)
        VALUES
        ('%v', '%v', '%v', '%v', 
		'%v', '%v', '%v', '%v', 
		'%v', '%v', '%v','%v')`,
		vendorReg.VendorName, vendorReg.VendorCode, vendorReg.MobileNumber, vendorReg.ContactPerson,
		vendorReg.AlternativeNumber, vendorReg.AddressLine1, vendorReg.AddressLine2, vendorReg.City,
		vendorReg.State, vendorReg.Status, vendorReg.OrgId, vendorReg.IsActive)

	rl.l.Info("ceateVendorQuery : ", createVendorQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(createVendorQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(CreateVendor): ", err)
		return err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(CreateVendor):", createdId, err)
		return err
	}
	rl.l.Info("Vendor created successfully: ", createdId, vendorReg.VendorCode)
	return nil
}

func (rl *VendorObj) CreateVendorV1(vendorReg dtos.VendorRequest) (int64, error) {
	rl.l.Info("CreateVendor : ", vendorReg.VendorCode)

	createVendorQuery := fmt.Sprintf(`INSERT INTO vendors (
        vendor_name, vendor_code, owner_name, gst_number, 
        preferred_operating_routes, address_line1, city, state, 
        pan_number, tds_declaration, remark, bank_account_holder_name, 
        bank_account_number, bank_name, bank_ifsc_code, pancard_img, 
        bank_passbook_or_cheque_img, status, is_active, org_id, login_type
    ) VALUES (
        '%v', '%v', '%v', '%v', 
        '%v', '%v', '%v', '%v', 
        '%v', '%v', '%v', '%v', 
        '%v', '%v', '%v', '%v', 
        '%v', '%v', '%v', '%v', '%v'
    )`,
		vendorReg.VendorName, vendorReg.VendorCode, vendorReg.OwnerName, vendorReg.GSTNumber,
		vendorReg.PreferredOperatingRoutes, vendorReg.AddressLine1, vendorReg.City, vendorReg.State,
		vendorReg.PANNumber, vendorReg.TDSDeclaration, vendorReg.Remark, vendorReg.BankAccountHolderName,
		vendorReg.BankAccountNumber, vendorReg.BankName, vendorReg.BankIFSCCode, vendorReg.PancardImg,
		vendorReg.BankPassbookOrChequeImg, vendorReg.Status, vendorReg.IsActive, vendorReg.OrgID, vendorReg.LoginType)

	rl.l.Info("CreateVendor query: ", createVendorQuery)

	result, err := rl.dbConnMSSQL.GetQueryer().Exec(createVendorQuery)
	if err != nil {
		rl.l.Error("Error creating vendor: ", err)
		return 0, err
	}

	createdId, err := result.LastInsertId()
	if err != nil {
		rl.l.Error("Error getting last insert ID: ", err)
		return 0, err
	}

	rl.l.Info("Vendor created successfully. ID: ", createdId)
	return createdId, nil
}

func (cb *VendorObj) CreateVendorContactInfo(vendorId int64, contactReq dtos.VendorContactInfo) error {

	createContactInfoQuery := fmt.Sprintf(`INSERT INTO vendor_contact_info (
    vendor_id, contact_person_name, post, email_id, contact_nummber1, contact_nummber2
	) VALUES (
		'%v', '%v', '%v', '%v', '%v', '%v'
	)`,
		vendorId,
		contactReq.ContactPersonName,
		contactReq.Post,
		contactReq.EmailID,
		contactReq.ContactNumber1,
		contactReq.ContactNumber2,
	)
	cb.l.Info("createContactInfoQuery:\n ", createContactInfoQuery)

	result, err := cb.dbConnMSSQL.GetQueryer().Exec(createContactInfoQuery)
	if err != nil {
		cb.l.Error("Error db.Exec(CreateVendorContactInfo): ", err)
		return err
	}
	createdId, err := result.LastInsertId()
	if err != nil {
		cb.l.Error("Error db.Exec(CreateVendorContactInfo):", createdId, err)
		return err
	}
	cb.l.Info("Vendor Contact Into created successfully: ", createdId, contactReq.ContactPersonName)
	return nil
}

func (rl *VendorObj) GetVendors(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.VendorRes, error) {
	list := []dtos.VendorRes{}

	whereQuery = fmt.Sprintf(" %v ORDER BY updated_at DESC LIMIT %v OFFSET %v;", whereQuery, limit, offset)
	rl.l.Info("vendor whereQuery:\n ", whereQuery)

	vendorQuery := fmt.Sprintf(`SELECT vendor_id, vendor_name, vendor_code, mobile_number, contact_person,
    alternative_number, address_line1, address_line2, is_active, status, city, state, visiting_card_image, pancard_img, aadhar_card_img, cancelled_check_book_img, bank_passbook_img, gst_document_img FROM vendors %v;`, whereQuery)

	rl.l.Info("vendorQuery:\n ", vendorQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(vendorQuery)
	if err != nil {
		rl.l.Error("Error Vendors ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vendorName, vendorCode, mobileNumber, contactPerson, alternativeNumber, addressLine1, addressLine2, status, city, state, visitingCardImage, pancardImg, aadharCardImg, cancelledCheckBookImg, bankPassbookImg,
			gstDocumentImg sql.NullString

		var vendorId, isActive sql.NullInt64

		vendorRes := &dtos.VendorRes{}
		err := rows.Scan(&vendorId, &vendorName, &vendorCode, &mobileNumber, &contactPerson, &alternativeNumber,
			&addressLine1, &addressLine2, &isActive, &status, &city, &state,
			&visitingCardImage, &pancardImg, &aadharCardImg, &cancelledCheckBookImg, &bankPassbookImg, &gstDocumentImg)
		if err != nil {
			rl.l.Error("Error GetVendors scan: ", err)
			return nil, err
		}
		vendorRes.VendorId = vendorId.Int64
		vendorRes.VendorName = vendorName.String
		vendorRes.VendorCode = vendorCode.String
		vendorRes.MobileNumber = mobileNumber.String
		vendorRes.ContactPerson = contactPerson.String
		vendorRes.AlternativeNumber = alternativeNumber.String
		vendorRes.AddressLine1 = addressLine1.String
		vendorRes.AddressLine2 = addressLine2.String
		vendorRes.IsActive = isActive.Int64
		vendorRes.Status = status.String
		vendorRes.State = state.String
		vendorRes.City = city.String
		vendorRes.VisitingCardImage = visitingCardImage.String
		vendorRes.PancardImg = pancardImg.String
		vendorRes.AadharCardImg = aadharCardImg.String
		vendorRes.CancelledCheckBookImg = cancelledCheckBookImg.String
		vendorRes.BankPassbookImg = bankPassbookImg.String
		vendorRes.GstDocumentImg = gstDocumentImg.String
		list = append(list, *vendorRes)
	}

	return &list, nil
}

func (rl *VendorObj) GetVendorsV1(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.VendorResponse, error) {
	list := []dtos.VendorResponse{}

	whereQuery = fmt.Sprintf(" %v ORDER BY updated_at DESC LIMIT %v OFFSET %v", whereQuery, limit, offset)
	vendorQuery := fmt.Sprintf(`SELECT 
        vendor_id, vendor_name, vendor_code, owner_name, gst_number,
        preferred_operating_routes, address_line1, city, state, pan_number,
        tds_declaration, remark, bank_account_holder_name, bank_account_number,
        bank_name, bank_ifsc_code, pancard_img, bank_passbook_or_cheque_img,
        status, is_active, org_id, login_type
        FROM vendors %v;`, whereQuery)

	rl.l.Info("whereQuery:\n", whereQuery)
	rl.l.Info("Executing vendor query:\n", vendorQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(vendorQuery)
	if err != nil {
		rl.l.Error("Error querying vendors: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			vendorName, vendorCode, ownerName, gstNumber, preferredOperatingRoutes, addressLine1, city,
			state, panNumber, tdsDeclaration, remark, bankAccountHolderName, bankAccountNumber,
			bankName, bankIFSCCode, pancardImg, bankPassbookChequeImg, status, loginType sql.NullString
			vendorID, orgID, isActive sql.NullInt64
		)

		err := rows.Scan(
			&vendorID, &vendorName, &vendorCode, &ownerName, &gstNumber,
			&preferredOperatingRoutes, &addressLine1, &city, &state, &panNumber,
			&tdsDeclaration, &remark, &bankAccountHolderName, &bankAccountNumber,
			&bankName, &bankIFSCCode, &pancardImg, &bankPassbookChequeImg,
			&status, &isActive, &orgID, &loginType,
		)

		if err != nil {
			rl.l.Error("Error scanning vendor row: ", err)
			return nil, err
		}

		vendorRes := dtos.VendorResponse{
			VendorId:                 vendorID.Int64,
			VendorName:               vendorName.String,
			VendorCode:               vendorCode.String,
			OwnerName:                ownerName.String,
			GSTNumber:                gstNumber.String,
			PreferredOperatingRoutes: preferredOperatingRoutes.String,
			AddressLine1:             addressLine1.String,
			City:                     city.String,
			State:                    state.String,
			PANNumber:                panNumber.String,
			TDSDeclaration:           tdsDeclaration.String,
			Remark:                   remark.String,
			BankAccountHolderName:    bankAccountHolderName.String,
			BankAccountNumber:        bankAccountNumber.String,
			BankName:                 bankName.String,
			BankIFSCCode:             bankIFSCCode.String,
			PancardImg:               pancardImg.String,
			BankPassbookOrChequeImg:  bankPassbookChequeImg.String,
			Status:                   status.String,
			IsActive:                 isActive.Int64,
			OrgID:                    orgID.Int64,
			LoginType:                loginType.String,
		}
		list = append(list, vendorRes)
	}
	return &list, nil
}

func (rl *VendorObj) GetVendorV1(vendorId int64) (*dtos.VendorResponse, error) {

	vendorQuery := fmt.Sprintf(`SELECT 
        vendor_id, vendor_name, vendor_code, owner_name, gst_number,
        preferred_operating_routes, address_line1, city, state, pan_number,
        tds_declaration, remark, bank_account_holder_name, bank_account_number,
        bank_name, bank_ifsc_code, pancard_img, bank_passbook_or_cheque_img,
        status, is_active, org_id, login_type
        FROM vendors WHERE vendor_id = '%v';`, vendorId)

	rl.l.Info("Executing vendor query:\n", vendorQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(vendorQuery)
	var (
		vendorName, vendorCode, ownerName, gstNumber, preferredOperatingRoutes, addressLine1, city,
		state, panNumber, tdsDeclaration, remark, bankAccountHolderName, bankAccountNumber,
		bankName, bankIFSCCode, pancardImg, bankPassbookChequeImg, status, loginType sql.NullString
		vendorID, orgID, isActive sql.NullInt64
	)
	err := row.Scan(
		&vendorID, &vendorName, &vendorCode, &ownerName, &gstNumber,
		&preferredOperatingRoutes, &addressLine1, &city, &state, &panNumber,
		&tdsDeclaration, &remark, &bankAccountHolderName, &bankAccountNumber,
		&bankName, &bankIFSCCode, &pancardImg, &bankPassbookChequeImg,
		&status, &isActive, &orgID, &loginType,
	)

	if err != nil {
		rl.l.Error("Error scanning vendor row: ", err)
		return nil, err
	}
	vendorRes := dtos.VendorResponse{
		VendorId:                 vendorID.Int64,
		VendorName:               vendorName.String,
		VendorCode:               vendorCode.String,
		OwnerName:                ownerName.String,
		GSTNumber:                gstNumber.String,
		PreferredOperatingRoutes: preferredOperatingRoutes.String,
		AddressLine1:             addressLine1.String,
		City:                     city.String,
		State:                    state.String,
		PANNumber:                panNumber.String,
		TDSDeclaration:           tdsDeclaration.String,
		Remark:                   remark.String,
		BankAccountHolderName:    bankAccountHolderName.String,
		BankAccountNumber:        bankAccountNumber.String,
		BankName:                 bankName.String,
		BankIFSCCode:             bankIFSCCode.String,
		PancardImg:               pancardImg.String,
		BankPassbookOrChequeImg:  bankPassbookChequeImg.String,
		Status:                   status.String,
		IsActive:                 isActive.Int64,
		OrgID:                    orgID.Int64,
		LoginType:                loginType.String,
	}
	return &vendorRes, nil
}

func (rl *VendorObj) GetVendorContactInfo(vendorID int64) (*[]dtos.VendorContactInfo, error) {
	list := []dtos.VendorContactInfo{}

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(`
        SELECT contact_info_id, vendor_id, contact_person_name, 
            post, email_id, contact_nummber1, contact_nummber2, contact_nummber3
        FROM vendor_contact_info
        WHERE vendor_id = ?
    `, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			contactInfoID                                                                    sql.NullInt64
			vendorID                                                                         sql.NullInt64
			contactPersonName, post, emailID, contactNumber1, contactNumber2, contactNumber3 sql.NullString
		)
		err := rows.Scan(&contactInfoID, &vendorID, &contactPersonName,
			&post, &emailID, &contactNumber1, &contactNumber2, &contactNumber3,
		)
		if err != nil {
			return nil, err
		}

		vendorC := dtos.VendorContactInfo{
			ContactInfoID:     contactInfoID.Int64,
			VendorID:          vendorID.Int64,
			ContactPersonName: contactPersonName.String,
			Post:              post.String,
			EmailID:           emailID.String,
			ContactNumber1:    contactNumber1.String,
			ContactNumber2:    contactNumber2.String,
		}

		list = append(list, vendorC)
	}

	return &list, nil
}

func (rl *VendorObj) UpdateVendorV1(vendorId int64, veh dtos.VendorV1Update) error {
	updateVendorQuery := fmt.Sprintf(`
        UPDATE vendors SET
            vendor_name = '%v',
            vendor_code = '%v',
            owner_name = '%v',
            gst_number = '%v',
            preferred_operating_routes = '%v',
            address_line1 = '%v',
            city = '%v',
            state = '%v',
            pan_number = '%v',
            tds_declaration = '%v',
            remark = '%v',
            bank_account_holder_name = '%v',
            bank_account_number = '%v',
            bank_name = '%v',
            bank_ifsc_code = '%v',
            pancard_img = '%v',
            bank_passbook_or_cheque_img = '%v',
            status = '%v',
            is_active = %v,
            org_id = %v,
            login_type = '%v'
        WHERE vendor_id = '%v';`,
		veh.VendorName, veh.VendorCode, veh.OwnerName, veh.GSTNumber, veh.PreferredOperatingRoutes,
		veh.AddressLine1, veh.City, veh.State, veh.PANNumber, veh.TDSDeclaration, veh.Remark,
		veh.BankAccountHolderName, veh.BankAccountNumber, veh.BankName, veh.BankIFSCCode, veh.PancardImg,
		veh.BankPassbookOrChequeImg, veh.Status, veh.IsActive, veh.OrgID, veh.LoginType, vendorId)

	rl.l.Info("UpdateVendorV1 query: ", updateVendorQuery)
	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateVendorQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendorV1): ", err)
		return err
	}
	rl.l.Info("Vendor updated successfully: ", vendorId)
	return nil
}

func (rl *VendorObj) UpdateVendor(vendorId int64, veh dtos.VendorUpdate) error {

	updateVendorQuery := fmt.Sprintf(`
    UPDATE vendor SET
        vendor_name = '%v',
        vendor_code = '%v',
        mobile_number = '%v',
        contact_person = '%v',
        alternative_number = '%v',
        address_line1 = '%v',
        address_line2 = '%v',
        city = '%v',
        state = '%v',
        status = '%v',
        org_id = '%v',
		is_active = '%v' 
    	WHERE vendor_id = '%v';`,
		veh.VendorName, veh.VendorCode, veh.MobileNumber, veh.ContactPerson, veh.AlternativeNumber,
		veh.AddressLine1, veh.AddressLine2, veh.City, veh.State, veh.Status, veh.OrgId, veh.IsActive, vendorId)

	rl.l.Info("UpdateVendor query: ", updateVendorQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateVendorQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendor): ", err)
		return err
	}
	rl.l.Info("Vendor updated successfully: ", vendorId)
	return nil
}

func (rl *VendorObj) GetVendor(vendorId int64) (*dtos.VendorRes, error) {

	vendorQuery := fmt.Sprintf(`SELECT vendor_id, vendor_name, vendor_code, mobile_number, contact_person,
    alternative_number, address_line1, address_line2, is_active, status, city, state, visiting_card_image, pancard_img, aadhar_card_img, cancelled_check_book_img, bank_passbook_img, gst_document_img FROM vendors WHERE vendor_id = '%v' `, vendorId)

	rl.l.Info("vendorQuery:\n ", vendorQuery)

	var vendorName, vendorCode, mobileNumber, contactPerson, alternativeNumber, addressLine1, addressLine2, status, city, state, visitingCardImage, pancardImg, aadharCardImg, cancelledCheckBookImg, bankPassbookImg,
		gstDocumentImg sql.NullString

	var vendorIdN, isActive sql.NullInt64

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(vendorQuery)

	vendorRes := dtos.VendorRes{}

	err := row.Scan(&vendorIdN, &vendorName, &vendorCode, &mobileNumber, &contactPerson, &alternativeNumber,
		&addressLine1, &addressLine2, &isActive, &status, &city, &state,
		&visitingCardImage, &pancardImg, &aadharCardImg, &cancelledCheckBookImg, &bankPassbookImg, &gstDocumentImg)
	if err != nil {
		rl.l.Error("Error GetVendors scan: ", err)
		return nil, err
	}
	vendorRes.VendorId = vendorIdN.Int64
	vendorRes.VendorName = vendorName.String
	vendorRes.VendorCode = vendorCode.String
	vendorRes.MobileNumber = mobileNumber.String
	vendorRes.ContactPerson = contactPerson.String
	vendorRes.AlternativeNumber = alternativeNumber.String
	vendorRes.AddressLine1 = addressLine1.String
	vendorRes.AddressLine2 = addressLine2.String
	vendorRes.IsActive = isActive.Int64
	vendorRes.Status = status.String
	vendorRes.State = state.String
	vendorRes.City = city.String
	vendorRes.VisitingCardImage = visitingCardImage.String
	vendorRes.PancardImg = pancardImg.String
	vendorRes.AadharCardImg = aadharCardImg.String
	vendorRes.CancelledCheckBookImg = cancelledCheckBookImg.String
	vendorRes.BankPassbookImg = bankPassbookImg.String
	vendorRes.GstDocumentImg = gstDocumentImg.String

	return &vendorRes, nil
}

func (rl *VendorObj) UpdateVendorActiveStatus(vendorId, isActive int64) error {

	updateQuery := fmt.Sprintf(`UPDATE vendor SET is_active = '%v' WHERE vendor_id = '%v'`, isActive, vendorId)

	rl.l.Info("UpdateVendorActiveStatus Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendorActiveStatus): ", err)
		return err
	}

	rl.l.Info("Vendor status updated successfully: ", vendorId)

	return nil
}

func (rl *VendorObj) UpdateVendorActiveStatusV1(vendorId, isActive int64) error {

	updateQuery := fmt.Sprintf(`UPDATE vendors SET is_active = '%v' WHERE vendor_id = '%v'`, isActive, vendorId)

	rl.l.Info("UpdateVendorActiveStatus Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendorActiveStatus): ", err)
		return err
	}

	rl.l.Info("Vendor status updated successfully: ", vendorId)

	return nil
}

func (rl *VendorObj) UpdateVendorStatus(vendorId int64, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE vendor SET status = '%v' WHERE vendor_id = '%v'`, status, vendorId)

	rl.l.Info("UpdateVendorStatus Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendorStatus): ", err)
		return err
	}

	rl.l.Info("Vendor status updated successfully: ", vendorId)

	return nil
}

func (rl *VendorObj) UpdateVendorStatusV1(vendorId int64, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE vendors SET status = '%v' WHERE vendor_id = '%v'`, status, vendorId)

	rl.l.Info("UpdateVendorStatus Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendorStatus): ", err)
		return err
	}

	rl.l.Info("Vendor status updated successfully: ", vendorId)

	return nil
}

func (rl *VendorObj) UpdateVendorImagePath(updateQuery string) error {

	rl.l.Info("UpdateVendorImagePath Update query: ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendorImagePath) VendorObj: ", err)
		return err
	}

	return nil
}

func (vh *VendorObj) DeleteVendorContactInfo(vendorId int64) error {

	deleteQuery := fmt.Sprintf(`DELETE FROM vendor_contact_info WHERE vendor_id = '%v';`, vendorId)

	vh.l.Info("DeleteVendorContactInfo query:", deleteQuery)

	roleResult, err := vh.dbConnMSSQL.GetQueryer().Exec(deleteQuery)
	if err != nil {
		vh.l.Error("Error db.Exec(DeleteVendorContactInfo): ", err)
		return err
	}

	vh.l.Info("vendor_contact_info deleted successfully: ", vendorId, roleResult)

	return nil
}

func (vh *VendorObj) DeleteVendorVehiclesInfo(vendorId int64) error {

	deleteQuery := fmt.Sprintf(`DELETE FROM vehicles WHERE vendor_id = '%v';`, vendorId)

	vh.l.Info("DeleteVendorVehiclesInfo query:", deleteQuery)

	roleResult, err := vh.dbConnMSSQL.GetQueryer().Exec(deleteQuery)
	if err != nil {
		vh.l.Error("Error db.Exec(DeleteVendorVehiclesInfo): ", err)
		return err
	}

	vh.l.Info("vehicles deleted successfully: ", vendorId, roleResult)

	return nil
}

func (vh *VendorObj) DeleteVendorDeclarations(vendorId int64) error {

	deleteQuery := fmt.Sprintf(`DELETE FROM declaration_year_info WHERE vendor_id = '%v';`, vendorId)

	vh.l.Info("DeleteVendorDeclarations query:", deleteQuery)

	roleResult, err := vh.dbConnMSSQL.GetQueryer().Exec(deleteQuery)
	if err != nil {
		vh.l.Error("Error db.Exec(DeleteVendorDeclarations): ", err)
		return err
	}

	vh.l.Info("vendor declarations are deleted successfully: ", vendorId, roleResult)

	return nil
}

func (rl *VendorObj) CreateVehicle(vehicle dtos.VehicleModel) error {
	rl.l.Info("CreateVehicle : ", vehicle.VehicleNumber)

	createVehicleQuery := fmt.Sprintf(`
        INSERT INTO vehicles (
            vendor_id, vehicle_number, vehicle_type, vehicle_make, vehicle_model,
            permit_type, vehicle_size, closed_open, vehicle_capacity_tons,
            rc_expiry_doc, insurance_doc, pucc_expiry_doc, np_expire_doc,
            fitness_expiry_doc, tax_expiry_doc, mp_expire_doc
        )
        VALUES (
            %v, '%v', '%v', '%v', '%v',
            '%v', '%v', '%v', '%v',
            '%v', '%v', '%v', '%v',
            '%v', '%v', '%v'
        )`,
		vehicle.VendorID, vehicle.VehicleNumber, vehicle.VehicleType,
		vehicle.VehicleMake, vehicle.VehicleModel, vehicle.PermitType,
		vehicle.VehicleSize, vehicle.ClosedOpen, vehicle.VehicleCapacityTons,
		vehicle.RCExpiryDoc, vehicle.InsuranceDoc, vehicle.PUCCExpiryDoc,
		vehicle.NPExpireDoc, vehicle.FitnessExpiryDoc, vehicle.TaxExpiryDoc,
		vehicle.MPExpireDoc,
	)

	rl.l.Info("createVehicleQuery : ", createVehicleQuery)

	result, err := rl.dbConnMSSQL.GetQueryer().Exec(createVehicleQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(CreateVehicle): ", err)
		return err
	}

	createdId, err := result.LastInsertId()
	if err != nil {
		rl.l.Error("Error getting last insert ID: ", err)
		return err
	}

	rl.l.Info("Vehicle created successfully: ", createdId, vehicle.VehicleNumber)
	return nil
}

//CreateVehicleDoclarationDoc(vehicle dtos.DeclarationDocument) error

func (rl *VendorObj) CreateVehicleDoclarationDoc(vehicle dtos.DeclarationDocument) error {
	rl.l.Info("CreateVehicleDoclarationDoc : ", vehicle.DeclarationYear)

	createVehicleDecQuery := fmt.Sprintf(`
        INSERT INTO declaration_year_info (
            vendor_id, declaration_year, declaration_doc
        )
        VALUES (
            %v, '%v', '%v'
        )`,
		vehicle.VendorID, vehicle.DeclarationYear, vehicle.DeclarationDocImage,
	)

	rl.l.Info("createVehicleDecQuery : ", createVehicleDecQuery)
	result, err := rl.dbConnMSSQL.GetQueryer().Exec(createVehicleDecQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(createVehicleDecQuery): ", err)
		return err
	}

	createdId, err := result.LastInsertId()
	if err != nil {
		rl.l.Error("Error getting last insert ID: ", err)
		return err
	}

	rl.l.Info("Vehicle created successfully: ", createdId, vehicle.VendorID, vehicle.DeclarationYear)
	return nil
}

func (rl *VendorObj) GetVehiclesByVendorID(vendorID int64) (*[]dtos.VehicleModel, error) {
	list := []dtos.VehicleModel{}

	query := fmt.Sprintf(`
    SELECT vehicle_id, vendor_id, vehicle_number, vehicle_type, vehicle_make, vehicle_model,
           permit_type, vehicle_size, closed_open, vehicle_capacity_tons,
           rc_expiry_doc, insurance_doc, pucc_expiry_doc, np_expire_doc,
           fitness_expiry_doc, tax_expiry_doc, mp_expire_doc
    FROM vehicles
    WHERE vendor_id = '%v'
    ORDER BY updated_at DESC`, vendorID)

	rl.l.Info("GetVehiclesByVendorID query : ", query)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			vehicleID, vendorID sql.NullInt64
			vehicleNumber, vehicleType, vehicleMake, vehicleModel, permitType, vehicleSize,
			closedOpen, vehicleCapacityTons sql.NullString
			rcExpiryDoc, insuranceDoc, puccExpiryDoc, npExpireDoc,
			fitnessExpiryDoc, taxExpiryDoc, mpExpireDoc sql.NullString
		)
		err := rows.Scan(
			&vehicleID, &vendorID, &vehicleNumber, &vehicleType, &vehicleMake, &vehicleModel,
			&permitType, &vehicleSize, &closedOpen, &vehicleCapacityTons,
			&rcExpiryDoc, &insuranceDoc, &puccExpiryDoc, &npExpireDoc,
			&fitnessExpiryDoc, &taxExpiryDoc, &mpExpireDoc,
		)
		if err != nil {
			return nil, err
		}

		vehicle := dtos.VehicleModel{
			VehicleID:           vehicleID.Int64,
			VendorID:            vendorID.Int64,
			VehicleNumber:       vehicleNumber.String,
			VehicleType:         vehicleType.String,
			VehicleMake:         vehicleMake.String,
			VehicleModel:        vehicleModel.String,
			PermitType:          permitType.String,
			VehicleSize:         vehicleSize.String,
			ClosedOpen:          closedOpen.String,
			VehicleCapacityTons: vehicleCapacityTons.String,
			RCExpiryDoc:         rcExpiryDoc.String,
			InsuranceDoc:        insuranceDoc.String,
			PUCCExpiryDoc:       puccExpiryDoc.String,
			NPExpireDoc:         npExpireDoc.String,
			FitnessExpiryDoc:    fitnessExpiryDoc.String,
			TaxExpiryDoc:        taxExpiryDoc.String,
			MPExpireDoc:         mpExpireDoc.String,
		}

		list = append(list, vehicle)
	}

	return &list, nil
}

func (rl *VendorObj) GetDeclarationDocByVendorID(vendorID int64) (*[]dtos.DeclarationDocumentObj, error) {
	list := []dtos.DeclarationDocumentObj{}

	query := fmt.Sprintf(`SELECT declaration_year_info_id, vendor_id, declaration_doc, declaration_year 
    FROM declaration_year_info 
    WHERE vendor_id = '%v' ORDER BY updated_at DESC`, vendorID)

	rl.l.Info("GetDeclarationDocByVendorID query : ", query)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			declarationYearInfoID, vendorID sql.NullInt64
			declarationDoc, declarationYear sql.NullString
		)
		err := rows.Scan(&declarationYearInfoID, &vendorID, &declarationDoc, &declarationYear)
		if err != nil {
			return nil, err
		}

		declarations := dtos.DeclarationDocumentObj{
			DeclarationYearInfoID: declarationYearInfoID.Int64,
			VendorID:              vendorID.Int64,
			DeclarationDocImage:   declarationDoc.String,
			DeclarationYear:       declarationYear.String,
		}

		list = append(list, declarations)
	}

	return &list, nil
}

func (rl *VendorObj) UpdateVendorImage(vendorId string, imageFor, imagePath string) error {

	updateQuery := fmt.Sprintf(`UPDATE vendors SET %v = '%v' WHERE vendor_id = '%v'`, imageFor, imagePath, vendorId)

	rl.l.Info("UpdateVendorImage Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendorImage): ", err)
		return err
	}

	rl.l.Info("Vendor image updated successfully: ", vendorId)

	return nil
}

func (rl *VendorObj) UpdateVehicleImage(vehicleId string, imageFor, imagePath string) error {

	updateQuery := fmt.Sprintf(`UPDATE vehicles SET %v = '%v' WHERE vehicle_id = '%v'`, imageFor, imagePath, vehicleId)

	rl.l.Info("UpdateVendorImage Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVendorImage): ", err)
		return err
	}

	rl.l.Info("Vendor image updated successfully: ", vehicleId)

	return nil
}
