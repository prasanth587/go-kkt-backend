package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type DriverObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewDriverObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *DriverObj {
	return &DriverObj{
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

type DriverDao interface {
	CreateDriver(driver dtos.CreateDriverReq) (int64, error)
	GetDrivers(orgId int64, limit int64, offset int64) (*[]dtos.DriverRes, error)
	UpdateDriverActiveStatus(driverId, isActive int64) error
	GetDriver(driverId int64) (*dtos.DriverRes, error)
	UpdateDriver(driverId int64, driver dtos.UpdateDriverReq) error
	UpdateDriverImagePath(query string) error
	GetTotalCount(orgId int64) int64
}

func (rl *DriverObj) CreateDriver(driver dtos.CreateDriverReq) (int64, error) {

	rl.l.Info("CreateDriver : ", driver.FirstName)

	ceateDriverQuery := fmt.Sprintf(`INSERT INTO driver (
		first_name, last_name, license_number, 
		license_expiry_date, mobile_no, alternate_contact_number, email_id, 
		joining_date, relieving_date, is_active, 
		vehicle_id, driver_experience, address_line1, 
		address_line2, city, state, 
		login_type, country, org_id)
		VALUES
		( '%v', '%v', '%v', 
		  '%v', '%v', '%v', '%v',
		   '%v', '%v', '%v', 
		   '%v', '%v', '%v',
		   '%v', '%v', '%v',
		   '%v', '%v', '%v')`,
		driver.FirstName, driver.LastName, driver.LicenseNumber,
		driver.LicenseExpiryDate, driver.MobileNo, driver.AlternateContactNumber, driver.EmailID,
		driver.JoiningDate, driver.RelievingDate, driver.IsActive,
		driver.VehicleID, driver.DriverExperience, driver.AddressLine1,
		driver.AddressLine2, driver.City, driver.State,
		driver.LoginType, driver.Country, driver.OrgID,
	)

	rl.l.Info("ceateDriverQuery : ", ceateDriverQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(ceateDriverQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(CreateEmployee) employee: ", err)
		return 0, err
	}
	createdEmpId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(CreateDriver) driver:", createdEmpId, err)
		return 0, err
	}
	rl.l.Info("CreateDriver created successfully: ", createdEmpId, driver.FirstName)
	return createdEmpId, nil
}

func (rl *DriverObj) GetTotalCount(orgId int64) int64 {

	countQuery := fmt.Sprintf(`SELECT count(*) FROM driver WHERE org_id = '%v'`, orgId)
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

func (rl *DriverObj) GetDrivers(orgId int64, limit int64, offset int64) (*[]dtos.DriverRes, error) {
	list := []dtos.DriverRes{}

	query := "select driver_id, first_name, last_name, license_number, license_expiry_date, mobile_no,alternate_contact_number, email_id, joining_date, relieving_date, is_active, vehicle_id,driver_experience, address_line1, address_line2, city, state, country, org_id, login_type,license_front_img, license_back_img, other_document, image FROM driver WHERE org_id=? ORDER BY updated_at DESC LIMIT ? OFFSET ?"
	rl.l.Info("Driver: ", query)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(query, orgId, limit, offset)
	if err != nil {
		rl.l.Error("Error GetDrivers ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var firstName, lastName, licenseNumber, licenseDxpiryDate, mobileNo, alternateContactNumber, emailId, joiningDate, relievingDate, addressLine1, addressLine2, city, state, country, loginLype, licenseFrontImg, licenseBackImg, otherDocument, image sql.NullString

		var driverId, vehicleId, isActive, orgIdN sql.NullInt64

		var driverDxperience sql.NullFloat64

		driverRes := &dtos.DriverRes{}
		err := rows.Scan(&driverId, &firstName, &lastName, &licenseNumber, &licenseDxpiryDate, &mobileNo,
			&alternateContactNumber, &emailId, &joiningDate, &relievingDate, &isActive, &vehicleId,
			&driverDxperience, &addressLine1, &addressLine2, &city, &state, &country, &orgIdN,
			&loginLype, &licenseFrontImg, &licenseBackImg,
			&otherDocument, &image)
		if err != nil {
			rl.l.Error("Error GetDrivers scan: ", err)
			return nil, err
		}
		driverRes.DriverId = driverId.Int64
		driverRes.FirstName = firstName.String
		driverRes.LastName = lastName.String
		driverRes.LicenseNumber = licenseNumber.String
		driverRes.LicenseExpiryDate = licenseDxpiryDate.String
		driverRes.MobileNo = mobileNo.String
		driverRes.AlternateContactNumber = alternateContactNumber.String
		driverRes.EmailID = emailId.String
		driverRes.JoiningDate = joiningDate.String
		driverRes.RelievingDate = relievingDate.String
		driverRes.IsActive = isActive.Int64
		driverRes.VehicleID = vehicleId.Int64
		driverRes.DriverExperience = driverDxperience.Float64
		driverRes.AddressLine1 = addressLine1.String
		driverRes.AddressLine2 = addressLine2.String
		driverRes.City = city.String
		driverRes.State = state.String
		driverRes.Country = country.String
		driverRes.OrgID = orgIdN.Int64
		driverRes.LoginType = loginLype.String
		driverRes.LicenseFrontImg = licenseFrontImg.String
		driverRes.LicenseBackImg = licenseBackImg.String
		driverRes.OtherDocument = otherDocument.String
		driverRes.ProfileImage = image.String
		list = append(list, *driverRes)
	}

	return &list, nil
}

func (rl *DriverObj) GetDriver(driverId int64) (*dtos.DriverRes, error) {
	driverRes := dtos.DriverRes{}

	query := "select driver_id, first_name, last_name, license_number, license_expiry_date, mobile_no,alternate_contact_number, email_id, joining_date, relieving_date, is_active, vehicle_id,driver_experience, address_line1, address_line2, city, state, country, org_id, login_type,license_front_img, license_back_img, other_document, image FROM driver WHERE driver_id=? "
	rl.l.Info("Driver: ", query)

	var firstName, lastName, licenseNumber, licenseDxpiryDate, mobileNo, alternateContactNumber, emailId, joiningDate, relievingDate, addressLine1, addressLine2, city, state, country, loginLype, licenseFrontImg, licenseBackImg, otherDocument, image sql.NullString

	var driverIDn, vehicleId, isActive, orgIdN sql.NullInt64

	var driverDxperience sql.NullFloat64

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(query, driverId)

	err := row.Scan(&driverIDn, &firstName, &lastName, &licenseNumber, &licenseDxpiryDate, &mobileNo,
		&alternateContactNumber, &emailId, &joiningDate, &relievingDate, &isActive, &vehicleId,
		&driverDxperience, &addressLine1, &addressLine2, &city, &state, &country, &orgIdN,
		&loginLype, &licenseFrontImg, &licenseBackImg,
		&otherDocument, &image)
	if err != nil {
		rl.l.Error("Error GetDrivers scan: ", err)
		return nil, err
	}

	driverRes.DriverId = driverIDn.Int64
	driverRes.FirstName = firstName.String
	driverRes.LastName = lastName.String
	driverRes.LicenseNumber = licenseNumber.String
	driverRes.LicenseExpiryDate = licenseDxpiryDate.String
	driverRes.MobileNo = mobileNo.String
	driverRes.AlternateContactNumber = alternateContactNumber.String
	driverRes.EmailID = emailId.String
	driverRes.JoiningDate = joiningDate.String
	driverRes.RelievingDate = relievingDate.String
	driverRes.IsActive = isActive.Int64
	driverRes.VehicleID = vehicleId.Int64
	driverRes.DriverExperience = driverDxperience.Float64
	driverRes.AddressLine1 = addressLine1.String
	driverRes.AddressLine2 = addressLine2.String
	driverRes.City = city.String
	driverRes.State = state.String
	driverRes.Country = country.String
	driverRes.OrgID = orgIdN.Int64
	driverRes.LoginType = loginLype.String
	driverRes.LicenseFrontImg = licenseFrontImg.String
	driverRes.LicenseBackImg = licenseBackImg.String
	driverRes.OtherDocument = otherDocument.String
	driverRes.ProfileImage = image.String

	return &driverRes, nil
}

func (rl *DriverObj) UpdateDriverActiveStatus(driverId, isActive int64) error {

	updateQuery := fmt.Sprintf(`UPDATE driver SET is_active = '%v' WHERE driver_id = '%v'`, isActive, driverId)

	rl.l.Info("UpdateDriverActiveStatus Update query ", updateQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateDriverActiveStatus): ", err)
		return err
	}
	dID, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateDriverActiveStatus): ", dID, err)
		return err
	}
	rl.l.Info("Driver updated successfully: ", driverId)

	return nil
}

func (rl *DriverObj) UpdateDriver(driverId int64, driver dtos.UpdateDriverReq) error {

	updateQuery := fmt.Sprintf(`
    UPDATE driver SET 
        first_name = '%s',
        last_name = '%s',
        license_number = '%s',
        license_expiry_date = '%s',
        mobile_no = '%s',
        alternate_contact_number = '%s',
        email_id = '%s',
        joining_date = '%s',
        relieving_date = '%s',
        is_active = '%d',
        vehicle_id = '%d',
        driver_experience = '%f',
        address_line1 = '%s',
        address_line2 = '%s',
        city = '%s',
        state = '%s',
        country = '%s',
        org_id = '%d',
        login_type = '%s'
    WHERE driver_id = '%d'`,
		driver.FirstName, driver.LastName, driver.LicenseNumber, driver.LicenseExpiryDate,
		driver.MobileNo, driver.AlternateContactNumber, driver.EmailID, driver.JoiningDate,
		driver.RelievingDate, driver.IsActive, driver.VehicleID, driver.DriverExperience,
		driver.AddressLine1, driver.AddressLine2, driver.City, driver.State,
		driver.Country, driver.OrgID, driver.LoginType, driverId)

	rl.l.Info("UpdateDriver Update query: ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateDriver): ", err)
		return err
	}
	rl.l.Info("Driver updated successfully: ", driverId)
	return nil
}

func (rl *DriverObj) UpdateDriverImagePath(updateQuery string) error {

	rl.l.Info("UpdateDriverImagePath Update query: ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateDriverImagePath) driver: ", err)
		return err
	}

	return nil
}
