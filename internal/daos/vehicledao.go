package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type VehicleObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewVehicleObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *VehicleObj {
	return &VehicleObj{
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

type VehicleDao interface {
	CreateVehicle(vehicleReg dtos.VehicleReq) error
	GetVehicles(orgId int64, limit int64, offset int64) (*[]dtos.VehicleRes, error)
	UpdateVehicle(vehicleId int64, veh dtos.VehicleUpdate) error
	GetVehicle(vehicleId int64) (*dtos.VehicleRes, error)
	UpdateVehicleActiveStatus(vehicleId, isActive int64) error
	UpdateVehicleImagePath(updateQuery string) error
	GetTotalCount(orgId int64) int64

	// Vehicle size type
	CreateVehicleSizeTypes(vehicleReg dtos.VehicleSizeType) error
	VehicleSizeTypesGetTotalCount() int64
	GetVehicleSizeTypes(limit int64, offset int64) (*[]dtos.VehicleSizeTypeObj, error)
	GetVehicleSizeType(vehicleSizeId int64) (*dtos.VehicleSizeTypeObj, error)
	UpdateVehicleSizeType(vehicleSizeId int64, veh dtos.VehicleSizeTypeUpdate) error
	UpdateVehicleSizeTypeActiveStatus(vehicleSizeId, isActive int64) error
}

func (rl *VehicleObj) GetTotalCount(orgId int64) int64 {

	countQuery := fmt.Sprintf(`SELECT count(*) FROM vehicle WHERE org_id = '%v'`, orgId)
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

func (rl *VehicleObj) CreateVehicle(vehicleReg dtos.VehicleReq) error {

	rl.l.Info("CreateDriver : ", vehicleReg.VehicleNumber)

	ceateVehicleQuery := fmt.Sprintf(`INSERT INTO vehicle (
		vehicle_type, vehicle_number, vehicle_model, 
		vehicle_year, vehicle_capacity, vehicle_insurance_number, insurance_expiry_date, 
		vehicle_registration_date, vehicle_renewal_date, is_active, 
		org_id, status)
		VALUES
		( '%v', '%v', '%v', 
		  '%v', '%v', '%v', '%v',
		   '%v', '%v', '%v', 
		   '%v','%v')`,
		vehicleReg.VehicleType, vehicleReg.VehicleNumber, vehicleReg.VehicleModel,
		vehicleReg.VehicleYear, vehicleReg.VehicleCapacity, vehicleReg.VehicleInsuranceNumber, vehicleReg.InsuranceExpiryDate,
		vehicleReg.VehicleRegistrationDate, vehicleReg.VehicleRenewalDate, vehicleReg.IsActive,
		vehicleReg.OrgID, vehicleReg.Status)

	rl.l.Info("ceateVehicleQuery : ", ceateVehicleQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(ceateVehicleQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(CreateVehicle): ", err)
		return err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(CreateVehicle):", createdId, err)
		return err
	}
	rl.l.Info("Vehicle created successfully: ", createdId, vehicleReg.VehicleNumber)
	return nil
}

func (rl *VehicleObj) GetVehicles(orgId int64, limit int64, offset int64) (*[]dtos.VehicleRes, error) {
	list := []dtos.VehicleRes{}

	vehicleQuery := "SELECT vehicle_id, vehicle_type, vehicle_number, vehicle_model, vehicle_year, vehicle_capacity, vehicle_insurance_number, insurance_expiry_date, vehicle_registration_date, vehicle_renewal_date, is_active, driver_id, org_id, vehicle_image, insurance_certificate, registration_certificate, fitness_certificate, pollution_certificate, national_permits_certificate, annual_maintenance_certificate, status FROM vehicle WHERE org_id=? ORDER BY updated_at DESC LIMIT ? OFFSET ?"
	rl.l.Info("Vehicles: ", vehicleQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(vehicleQuery, orgId, limit, offset)
	if err != nil {
		rl.l.Error("Error Vehicles ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vehicleType, vehicleNumber, vehicleModel, vehicleCapacity, vehicleInsuranceNumber, insuranceExpiryDate, vehicleRegistrationDate, vehicleRenewalDate, vehicleImage, insuranceCert, registrationCert, fitnessCert, pollutionCert, nationalPermitsCert, annualMaintenanceCert, status sql.NullString

		var vehicleId, isActive, orgIdN, vehicleYear, driverId sql.NullInt64

		vehicleRes := &dtos.VehicleRes{}
		err := rows.Scan(&vehicleId, &vehicleType, &vehicleNumber, &vehicleModel, &vehicleYear, &vehicleCapacity,
			&vehicleInsuranceNumber, &insuranceExpiryDate, &vehicleRegistrationDate, &vehicleRenewalDate, &isActive, &driverId,
			&orgIdN, &vehicleImage, &insuranceCert, &registrationCert, &fitnessCert, &pollutionCert, &nationalPermitsCert, &annualMaintenanceCert, &status)
		if err != nil {
			rl.l.Error("Error GetVehicles scan: ", err)
			return nil, err
		}
		vehicleRes.VehicleId = vehicleId.Int64
		vehicleRes.VehicleType = vehicleType.String
		vehicleRes.VehicleNumber = vehicleNumber.String
		vehicleRes.VehicleModel = vehicleModel.String
		vehicleRes.VehicleYear = vehicleYear.Int64
		vehicleRes.VehicleCapacity = vehicleCapacity.String
		vehicleRes.VehicleInsuranceNumber = vehicleInsuranceNumber.String
		vehicleRes.InsuranceExpiryDate = insuranceExpiryDate.String
		vehicleRes.VehicleRegistrationDate = vehicleRegistrationDate.String
		vehicleRes.VehicleRenewalDate = vehicleRenewalDate.String
		vehicleRes.IsActive = isActive.Int64
		vehicleRes.OrgID = orgIdN.Int64
		vehicleRes.VehicleImage = vehicleImage.String
		vehicleRes.InsuranceCertificate = insuranceCert.String
		vehicleRes.RegistrationCertificate = registrationCert.String
		vehicleRes.FitnessCertificate = fitnessCert.String
		vehicleRes.PollutionCertificate = pollutionCert.String
		vehicleRes.NationalPermitsCertificate = nationalPermitsCert.String
		vehicleRes.AnnualMaintenanceCertificate = annualMaintenanceCert.String
		vehicleRes.Status = status.String

		list = append(list, *vehicleRes)
	}

	return &list, nil
}

func (rl *VehicleObj) UpdateVehicle(vehicleId int64, veh dtos.VehicleUpdate) error {

	updateVehicleQuery := fmt.Sprintf(`
    UPDATE vehicle SET 
        vehicle_type = '%v',
        vehicle_number = '%v',
        vehicle_model = '%v',
        vehicle_year = '%v',
        vehicle_capacity = '%v',
        vehicle_insurance_number = '%v',
        insurance_expiry_date = '%v',
        vehicle_registration_date = '%v',
        vehicle_renewal_date = '%v',
        is_active = '%v',
        org_id = '%v',
		status = '%v'
    WHERE vehicle_id = '%v'`,
		veh.VehicleType, veh.VehicleNumber, veh.VehicleModel, veh.VehicleYear, veh.VehicleCapacity,
		veh.VehicleInsuranceNumber, veh.InsuranceExpiryDate, veh.VehicleRegistrationDate, veh.VehicleRenewalDate,
		veh.IsActive, veh.OrgID, veh.Status, vehicleId)

	rl.l.Info("UpdateVehicle query: ", updateVehicleQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateVehicleQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVehicle): ", err)
		return err
	}
	rl.l.Info("Vehicle updated successfully: ", vehicleId)
	return nil
}

func (rl *VehicleObj) GetVehicle(vehicleId int64) (*dtos.VehicleRes, error) {

	vehicleQuery := "SELECT vehicle_id, vehicle_type, vehicle_number, vehicle_model, vehicle_year, vehicle_capacity, vehicle_insurance_number, insurance_expiry_date, vehicle_registration_date, vehicle_renewal_date, is_active, driver_id, org_id, vehicle_image, insurance_certificate, registration_certificate, fitness_certificate, pollution_certificate, national_permits_certificate, annual_maintenance_certificate, status FROM vehicle WHERE vehicle_id=? "
	rl.l.Info("GetVehicle: ", vehicleQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(vehicleQuery, vehicleId)

	var vehicleType, vehicleNumber, vehicleModel, vehicleCapacity, vehicleInsuranceNumber, insuranceExpiryDate, vehicleRegistrationDate, vehicleRenewalDate, vehicleImage, insuranceCert, registrationCert, fitnessCert, pollutionCert, nationalPermitsCert, annualMaintenanceCert, status sql.NullString

	var vehicleID, isActive, orgIdN, vehicleYear, driverId sql.NullInt64

	vehicleRes := dtos.VehicleRes{}
	err := row.Scan(&vehicleID, &vehicleType, &vehicleNumber, &vehicleModel, &vehicleYear, &vehicleCapacity,
		&vehicleInsuranceNumber, &insuranceExpiryDate, &vehicleRegistrationDate, &vehicleRenewalDate, &isActive, &driverId,
		&orgIdN, &vehicleImage, &insuranceCert, &registrationCert, &fitnessCert, &pollutionCert, &nationalPermitsCert, &annualMaintenanceCert, &status)
	if err != nil {
		rl.l.Error("Error GetVehicles scan: ", err)
		return nil, err
	}

	vehicleRes.VehicleId = vehicleID.Int64
	vehicleRes.VehicleType = vehicleType.String
	vehicleRes.VehicleNumber = vehicleNumber.String
	vehicleRes.VehicleModel = vehicleModel.String
	vehicleRes.VehicleYear = vehicleYear.Int64
	vehicleRes.VehicleCapacity = vehicleCapacity.String
	vehicleRes.VehicleInsuranceNumber = vehicleInsuranceNumber.String
	vehicleRes.InsuranceExpiryDate = insuranceExpiryDate.String
	vehicleRes.VehicleRegistrationDate = vehicleRegistrationDate.String
	vehicleRes.VehicleRenewalDate = vehicleRenewalDate.String
	vehicleRes.IsActive = isActive.Int64
	vehicleRes.OrgID = orgIdN.Int64
	vehicleRes.VehicleImage = vehicleImage.String
	vehicleRes.InsuranceCertificate = insuranceCert.String
	vehicleRes.RegistrationCertificate = registrationCert.String
	vehicleRes.FitnessCertificate = fitnessCert.String
	vehicleRes.PollutionCertificate = pollutionCert.String
	vehicleRes.NationalPermitsCertificate = nationalPermitsCert.String
	vehicleRes.AnnualMaintenanceCertificate = annualMaintenanceCert.String
	vehicleRes.Status = status.String

	return &vehicleRes, nil
}

func (rl *VehicleObj) UpdateVehicleActiveStatus(vehicleId, isActive int64) error {

	updateQuery := fmt.Sprintf(`UPDATE vehicle SET is_active = '%v' WHERE vehicle_id = '%v'`, isActive, vehicleId)

	rl.l.Info("UpdateVehicleActiveStatus Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVehicleActiveStatus): ", err)
		return err
	}

	rl.l.Info("Vehicle status updated successfully: ", vehicleId)

	return nil
}

func (rl *VehicleObj) UpdateVehicleImagePath(updateQuery string) error {

	rl.l.Info("UpdateVehicleImagePath Update query: ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVehicleImagePath) VehicleObj: ", err)
		return err
	}

	return nil
}
