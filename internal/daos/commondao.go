package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type PreRequisiteObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewPreRequisiteObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *PreRequisiteObj {
	return &PreRequisiteObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

type PreRequisiteDao interface {
	GetLoadingPoints(orgId int64) (*[]dtos.LoadingPoints, error)
	GetCustomers(orgId int64) (*[]dtos.Customers, error)
	GetBranches(orgId int64) (*[]dtos.BranchT, error)
	GetVendors(orgId int64) (*[]dtos.VendorT, error)
	InsertAndGetTripNumber(tripSheeetNumber string) (int64, error)
	GetPodTripInfo(podRequired int) (*[]dtos.TripSheetInfo, error)
	GetLRTripInfo() (*[]dtos.TripSheetInfo, error)
	GetLocationNameById(loadingPointId int64) (*dtos.LoadUnLoadLoc, error)
	GetEmpRoles(orgId int64) (*[]dtos.RoleEmpPre, error)
	GetLastTripSheetRow() (int64, string, error)
	GetLastCustomerRow() (int64, string, error)
	GetLastVendorRow() (int64, string, error)
	GetVehicleSizeTypes() (*[]dtos.VehicleSizeTypePre, error)
	GetWebsiteScreen() (*[]dtos.WebsiteScreen, error)
	GetPermissionLabel() (*[]dtos.PermissionLabel, error)
	GetEmployeePre() (*[]dtos.EmployeesPre, error)
}

func (br *PreRequisiteObj) GetLoadingPoints(orgId int64) (*[]dtos.LoadingPoints, error) {
	list := []dtos.LoadingPoints{}

	loadingpointQuery := fmt.Sprintf(`SELECT loading_point_id, city_code, city_name FROM loading_point WHERE org_id = '%v' AND is_active = '1' ORDER BY city_code ASC;`, orgId)

	br.l.Info("loadingpointQuery:\n ", loadingpointQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(loadingpointQuery)
	if err != nil {
		br.l.Error("Error LoadingPoints ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cityCode, cityName sql.NullString
		var loadingpointId sql.NullInt64

		loadingPointE := &dtos.LoadingPoints{}
		err := rows.Scan(&loadingpointId, &cityCode, &cityName)
		if err != nil {
			br.l.Error("Error GetLoadingpoints scan: ", err)
			return nil, err
		}
		loadingPointE.LoadingPointID = loadingpointId.Int64
		loadingPointE.CityCode = cityCode.String
		loadingPointE.CityName = cityName.String
		list = append(list, *loadingPointE)
	}

	return &list, nil
}

// GetLastTripSheetRow() (int64, string, error)
func (br *PreRequisiteObj) GetLastTripSheetRow() (int64, string, error) {

	lastTripQuery := `SELECT trip_sheet_id, trip_sheet_num 
	FROM trip_sheet_header ORDER BY trip_sheet_id DESC LIMIT 1;`

	br.l.Info("lastTripQuery:\n ", lastTripQuery)

	row := br.dbConnMSSQL.GetQueryer().QueryRow(lastTripQuery)
	var tripSheetNum sql.NullString
	var tripSheetId sql.NullInt64
	errE := row.Scan(&tripSheetId, &tripSheetNum)
	if errE != nil {
		br.l.Error("Error lastTripQuery scan: ", errE)
		return 0, "", errE
	}
	return tripSheetId.Int64, tripSheetNum.String, nil
}

func (br *PreRequisiteObj) GetLastCustomerRow() (int64, string, error) {

	lastCustomerQuery := `SELECT customer_id, customer_code 
	FROM customers ORDER BY customer_id DESC LIMIT 1;`

	br.l.Info("lastCustomerQuery:\n ", lastCustomerQuery)

	row := br.dbConnMSSQL.GetQueryer().QueryRow(lastCustomerQuery)
	var customerCode sql.NullString
	var customerId sql.NullInt64
	errE := row.Scan(&customerId, &customerCode)
	if errE != nil {
		br.l.Error("Error GetLastCustomerRow scan: ", errE)
		return 0, "", errE
	}
	return customerId.Int64, customerCode.String, nil
}

func (br *PreRequisiteObj) GetLastVendorRow() (int64, string, error) {

	lastVendorQuery := `SELECT vendor_id, vendor_code 
	FROM vendors ORDER BY vendor_id DESC LIMIT 1;`

	br.l.Info("lastVendorQuery:\n ", lastVendorQuery)

	row := br.dbConnMSSQL.GetQueryer().QueryRow(lastVendorQuery)
	var vendorCode sql.NullString
	var vendorId sql.NullInt64
	errE := row.Scan(&vendorId, &vendorCode)
	if errE != nil {
		br.l.Error("Error GetLastCustomerRow scan: ", errE)
		return 0, "", errE
	}
	return vendorId.Int64, vendorCode.String, nil
}

func (br *PreRequisiteObj) GetCustomers(orgId int64) (*[]dtos.Customers, error) {
	list := []dtos.Customers{}

	loadingpointQuery := fmt.Sprintf(`SELECT customer_id, customer_name, customer_code 
	FROM customers WHERE org_id = '%v' AND is_active = '1' AND status IN ('%v', '%v') ORDER BY customer_code ASC;`, orgId, constant.STATUS_APPROVED, constant.STATUS_CREATED)

	br.l.Info("loadingpointQuery:\n ", loadingpointQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(loadingpointQuery)
	if err != nil {
		br.l.Error("Error GetCustomers ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customerName, customerCode sql.NullString

		var customerId sql.NullInt64

		customerE := &dtos.Customers{}
		err := rows.Scan(&customerId, &customerName, &customerCode)
		if err != nil {
			br.l.Error("Error GetCustomers scan: ", err)
			return nil, err
		}
		customerE.CustomerId = customerId.Int64
		customerE.CustomerName = fmt.Sprintf("%v - %v", customerName.String, customerCode.String)
		list = append(list, *customerE)
	}

	return &list, nil
}

func (br *PreRequisiteObj) GetBranches(orgId int64) (*[]dtos.BranchT, error) {
	list := []dtos.BranchT{}

	branchQuery := fmt.Sprintf(`SELECT branch_id, branch_name, branch_code FROM branch WHERE org_id = '%v' AND is_active='1' ORDER BY branch_code ASC;`, orgId)

	br.l.Info("branchQuery:\n ", branchQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(branchQuery)
	if err != nil {
		br.l.Error("Error Branches ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branchName, branchCode sql.NullString
		var branchId sql.NullInt64

		branch := &dtos.BranchT{}
		err := rows.Scan(&branchId, &branchName, &branchCode)
		if err != nil {
			br.l.Error("Error GetBranchs scan: ", err)
			return nil, err
		}
		branch.BranchId = branchId.Int64
		branch.BranchName = fmt.Sprintf("%v - %v", branchCode.String, branchName.String)
		list = append(list, *branch)
	}

	return &list, nil
}

func (br *PreRequisiteObj) GetVendors(orgId int64) (*[]dtos.VendorT, error) {
	list := []dtos.VendorT{}

	vendorQuery := fmt.Sprintf(`SELECT vendor_id, vendor_name, vendor_code FROM vendors 
	WHERE org_id = '%v' AND is_active='1' AND status = '%v' ORDER BY vendor_code ASC;`, orgId, constant.STATUS_APPROVED)

	br.l.Info("vendorQuery:\n ", vendorQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(vendorQuery)
	if err != nil {
		br.l.Error("Error GetVendors ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vendorName, vendorCode sql.NullString
		var vendorID sql.NullInt64

		vendor := &dtos.VendorT{}
		err := rows.Scan(&vendorID, &vendorName, &vendorCode)
		if err != nil {
			br.l.Error("Error GetVendors scan: ", err)
			return nil, err
		}
		vendor.VendorId = vendorID.Int64
		vendor.VendorName = fmt.Sprintf("%v - %v", vendorName.String, vendorCode.String)
		list = append(list, *vendor)
	}
	return &list, nil
}

func (br *PreRequisiteObj) InsertAndGetTripNumber(tripSheeetNumber string) (int64, error) {

	tripSheetNum := fmt.Sprintf(`INSERT INTO trip_sheet_num (year) VALUES ( '%v' )`, tripSheeetNumber)

	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(tripSheetNum)
	if err != nil {
		br.l.Error("Error db.Exec(InsertAndGetTripNumber): ", err)
		return 0, err
	}
	createdId, _ := roleResult.LastInsertId()

	return createdId, nil

}

func (br *PreRequisiteObj) GetPodTripInfo(podRequired int) (*[]dtos.TripSheetInfo, error) {
	list := []dtos.TripSheetInfo{}

	podQuery := fmt.Sprintf(`SELECT trip_sheet_id, trip_sheet_num, open_trip_date_time, lr_number, 
	customer_id, load_status, customer_invoice_no  
	FROM trip_sheet_header WHERE pod_required = '%v' AND load_status IN ('Intransit') ORDER BY trip_sheet_num ASC;`, podRequired)

	br.l.Info("podQuery:\n ", podQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(podQuery)
	if err != nil {
		br.l.Error("Error GetPodTripInfo ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tripSheetNum, openTripDateTime, lrNumber, loadStatus, invoiceNumber sql.NullString
		var tripSheetId, customerId sql.NullInt64

		tripInfo := &dtos.TripSheetInfo{}
		err := rows.Scan(&tripSheetId, &tripSheetNum, &openTripDateTime, &lrNumber, &customerId, &loadStatus, &invoiceNumber)
		if err != nil {
			br.l.Error("Error GetPodTripInfo scan: ", err)
			return nil, err
		}
		tripInfo.TripSheetID = tripSheetId.Int64
		tripInfo.TripSheetNum = tripSheetNum.String
		tripInfo.LRNumber = lrNumber.String
		tripInfo.CustomerId = customerId.Int64
		tripInfo.LoadStatus = loadStatus.String
		tripInfo.TripDate = openTripDateTime.String
		tripInfo.TripDate = openTripDateTime.String
		tripInfo.PodRequired = podRequired
		tripInfo.InvoiceNumber = invoiceNumber.String

		list = append(list, *tripInfo)
	}

	return &list, nil
}

func (br *PreRequisiteObj) GetLRTripInfo() (*[]dtos.TripSheetInfo, error) {
	list := []dtos.TripSheetInfo{}

	order := "DESC"
	isLRGenerated := 0

	lrQuery := fmt.Sprintf(`SELECT trip_sheet_id, trip_sheet_num, open_trip_date_time,
	lr_number, customer_id, load_status, customer_invoice_no, driver_name,
	mobile_number, vehicle_number, vehicle_size
	FROM trip_sheet_header WHERE is_lr_generated = '%v' ORDER BY trip_sheet_num %s;`, isLRGenerated, order)

	br.l.Info("lrQuery:\n ", lrQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(lrQuery)
	if err != nil {
		br.l.Error("Error GetLRTripInfo ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tripSheetNum, openTripDateTime, lrNumber, loadStatus, invoiceNumber, driverName, mobileNumber, vehicleNumber, vehicleSize sql.NullString
		var tripSheetId, customerId sql.NullInt64

		tripInfo := &dtos.TripSheetInfo{}
		err := rows.Scan(&tripSheetId, &tripSheetNum, &openTripDateTime, &lrNumber, &customerId, &loadStatus, &invoiceNumber, &driverName, &mobileNumber, &vehicleNumber, &vehicleSize)
		if err != nil {
			br.l.Error("Error GetLRTripInfo scan: ", err)
			return nil, err
		}
		tripInfo.TripSheetID = tripSheetId.Int64
		tripInfo.TripSheetNum = tripSheetNum.String
		tripInfo.LRNumber = lrNumber.String
		tripInfo.CustomerId = customerId.Int64
		tripInfo.LoadStatus = loadStatus.String
		tripInfo.TripDate = openTripDateTime.String
		tripInfo.TripDate = openTripDateTime.String
		tripInfo.InvoiceNumber = invoiceNumber.String
		tripInfo.DriverName = driverName.String
		tripInfo.MobileNumber = mobileNumber.String
		tripInfo.VehicleNumber = vehicleNumber.String
		tripInfo.VehicleSize = vehicleSize.String
		list = append(list, *tripInfo)
	}

	return &list, nil
}

// TODO: Above commanted has to work
// func (br *PreRequisiteObj) GetLRTripInfo() (*[]dtos.TripSheetInfo, error) {
// 	list := []dtos.TripSheetInfo{}

// 	order := "DESC"
// 	//isLRGenerated := 0

// 	lrQuery := fmt.Sprintf(`SELECT trip_sheet_id, trip_sheet_num, open_trip_date_time,
// 	lr_number, customer_id, load_status, customer_invoice_no, driver_name,
// 	mobile_number, vehicle_number, vehicle_size
// 	FROM trip_sheet_header
// 	WHERE load_status NOT IN ('%v', '%v')
// 	ORDER BY trip_sheet_num %s;`, constant.STATUS_CLOSED, constant.STATUS_COMPLETED, order)

// 	br.l.Info("lrQuery:\n ", lrQuery)

// 	rows, err := br.dbConnMSSQL.GetQueryer().Query(lrQuery)
// 	if err != nil {
// 		br.l.Error("Error GetLRTripInfo ", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var tripSheetNum, openTripDateTime, lrNumber, loadStatus, invoiceNumber, driverName, mobileNumber, vehicleNumber, vehicleSize sql.NullString
// 		var tripSheetId, customerId sql.NullInt64

// 		tripInfo := &dtos.TripSheetInfo{}
// 		err := rows.Scan(&tripSheetId, &tripSheetNum, &openTripDateTime, &lrNumber, &customerId, &loadStatus, &invoiceNumber, &driverName, &mobileNumber, &vehicleNumber, &vehicleSize)
// 		if err != nil {
// 			br.l.Error("Error GetLRTripInfo scan: ", err)
// 			return nil, err
// 		}
// 		tripInfo.TripSheetID = tripSheetId.Int64
// 		tripInfo.TripSheetNum = tripSheetNum.String
// 		tripInfo.LRNumber = lrNumber.String
// 		tripInfo.CustomerId = customerId.Int64
// 		tripInfo.LoadStatus = loadStatus.String
// 		tripInfo.TripDate = openTripDateTime.String
// 		tripInfo.TripDate = openTripDateTime.String
// 		tripInfo.InvoiceNumber = invoiceNumber.String
// 		tripInfo.DriverName = driverName.String
// 		tripInfo.MobileNumber = mobileNumber.String
// 		tripInfo.VehicleNumber = vehicleNumber.String
// 		tripInfo.VehicleSize = vehicleSize.String
// 		list = append(list, *tripInfo)
// 	}

// 	return &list, nil
// }

func (mp *PreRequisiteObj) GetLocationNameById(loadingPointId int64) (*dtos.LoadUnLoadLoc, error) {

	locQuery := fmt.Sprintf(`SELECT city_code, city_name FROM loading_point WHERE loading_point_id = '%v' `, loadingPointId)
	mp.l.Info("GetLocationNameById whereQuery:\n ", locQuery)

	row := mp.dbConnMSSQL.GetQueryer().QueryRow(locQuery)
	var cityCode, cityName sql.NullString
	errE := row.Scan(&cityCode, &cityName)
	if errE != nil {
		mp.l.Error("Error GetCount scan: ", errE)
		return nil, errE
	}
	loadUnLoad := &dtos.LoadUnLoadLoc{}
	loadUnLoad.CityCode = cityCode.String
	loadUnLoad.CityName = cityName.String
	loadUnLoad.LoadingPointId = loadingPointId

	return loadUnLoad, nil
}

func (rl *PreRequisiteObj) GetEmpRoles(orgId int64) (*[]dtos.RoleEmpPre, error) {
	list := []dtos.RoleEmpPre{}

	query := `SELECT role_id, role_code, role_name, description FROM thub_role WHERE org_id = ? AND is_active = 1 ORDER BY updated_at DESC`
	rl.l.Info("GetEmpRoles whereQuery: ", query)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(query, orgId)
	if err != nil {
		rl.l.Error("Error GetEmpRoles ", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var roleCode, roleName, description sql.NullString
		var roleId sql.NullInt64

		role := &dtos.RoleEmpPre{}
		err := rows.Scan(&roleId,
			&roleCode,
			&roleName,
			&description,
		)
		if err != nil {
			rl.l.Error("Error GetEmpRoles scan: ", err)
			return nil, err
		}
		role.RoleId = roleId.Int64
		role.RoleCode = roleCode.String
		role.RoleName = roleName.String
		role.Description = description.String

		list = append(list, *role)
	}

	return &list, nil
}

func (rl *PreRequisiteObj) GetVehicleSizeTypes() (*[]dtos.VehicleSizeTypePre, error) {
	list := []dtos.VehicleSizeTypePre{}

	vehicleSizeTypesQuery := "SELECT vehicle_size_id, vehicle_size, vehicle_type  FROM vehicle_size_type WHERE is_active = ? ORDER BY vehicle_size"
	rl.l.Info("vehicleSizeTypesQuery: ", vehicleSizeTypesQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(vehicleSizeTypesQuery, 1)
	if err != nil {
		rl.l.Error("Error Vehicles ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vehicleType, vehicleSize sql.NullString
		var vehicleSizeId sql.NullInt64

		vehicleRes := &dtos.VehicleSizeTypePre{}
		err := rows.Scan(&vehicleSizeId, &vehicleSize, &vehicleType)
		if err != nil {
			rl.l.Error("Error GetVehicles scan: ", err)
			return nil, err
		}
		vehicleRes.VehicleSizeId = vehicleSizeId.Int64
		vehicleRes.VehicleType = vehicleType.String
		vehicleRes.VehicleSize = vehicleSize.String

		list = append(list, *vehicleRes)
	}

	return &list, nil
}

func (rl *PreRequisiteObj) GetWebsiteScreen() (*[]dtos.WebsiteScreen, error) {
	list := []dtos.WebsiteScreen{}

	websiteScreenQuery := "SELECT website_screen_id, screen_name, description from website_screens order by screen_name"
	rl.l.Info("websiteScreenQuery: ", websiteScreenQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(websiteScreenQuery)
	if err != nil {
		rl.l.Error("Error websiteScreenQuery ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var screenName, description sql.NullString
		var websiteScreenID sql.NullInt64

		websiteScreenRes := &dtos.WebsiteScreen{}
		err := rows.Scan(&websiteScreenID, &screenName, &description)
		if err != nil {
			rl.l.Error("Error GetVehicles scan: ", err)
			return nil, err
		}
		websiteScreenRes.WebsiteScreenID = websiteScreenID.Int64
		websiteScreenRes.ScreenName = screenName.String
		websiteScreenRes.Description = description.String

		list = append(list, *websiteScreenRes)
	}

	return &list, nil
}

//

func (rl *PreRequisiteObj) GetPermissionLabel() (*[]dtos.PermissionLabel, error) {
	list := []dtos.PermissionLabel{}

	permisstionLabelQuery := "select permisstion_label_id, permisstion_label, description from permisstion_label order by permisstion_label"
	rl.l.Info("permisstionLabelQuery: ", permisstionLabelQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(permisstionLabelQuery)
	if err != nil {
		rl.l.Error("Error permisstionLabelQuery ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permisstionLabel, description sql.NullString
		var permisstionLabelID sql.NullInt64

		labelRes := &dtos.PermissionLabel{}
		err := rows.Scan(&permisstionLabelID, &permisstionLabel, &description)
		if err != nil {
			rl.l.Error("Error GetPermissionLabel scan: ", err)
			return nil, err
		}
		labelRes.PermisstionLabelID = permisstionLabelID.Int64
		labelRes.PermisstionLabel = permisstionLabel.String
		labelRes.Description = description.String

		list = append(list, *labelRes)
	}

	return &list, nil
}

//GetEmployee() (*[]dtos.Employee, error)

func (rl *PreRequisiteObj) GetEmployeePre() (*[]dtos.EmployeesPre, error) {
	list := []dtos.EmployeesPre{}

	employeeQuery := "SELECT emp_id, first_name, last_name from employee where is_active = 1 order by first_name"
	rl.l.Info("employeeQuery: ", employeeQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(employeeQuery)
	if err != nil {
		rl.l.Error("Error employeeQuery", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var firstName, lastName sql.NullString
		var empId sql.NullInt64

		empRes := &dtos.EmployeesPre{}
		err := rows.Scan(&empId, &firstName, &lastName)
		if err != nil {
			rl.l.Error("Error GetEmployee scan: ", err)
			return nil, err
		}
		empRes.EmployeeId = empId.Int64
		empRes.EmployeeName = fmt.Sprintf("%s %s", firstName.String, lastName.String)

		list = append(list, *empRes)
	}

	return &list, nil
}
