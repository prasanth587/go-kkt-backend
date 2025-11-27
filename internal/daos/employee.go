package daos

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/dtos/schema"
)

type EmpObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewEmpRoleObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *EmpObj {
	return &EmpObj{
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

type EmpRoleDao interface {
	CreateEmpRole(empRole dtos.EmpRole) (int64, error)
	UpdateEmpRole(empRole dtos.EmpRoleUpdate) error
	IsRoleExists(roleId int64) (*dtos.RoleEmp, error)
	GetEmpRoles(orgId int64, limit int64, offset int64) (*[]dtos.RoleEmp, error)
	GetOrgBaseInfo(orgId int64) (*dtos.OrganisationBaseResponse, error)
	GetTotalCountRole(orgId int64) int64
	ScreensPermission(roleCode, roleName string, roleID int64, sp dtos.ScreensPermission) error
	DeleteScreensPermission(roleID int64) error
	GetScreensPermissions(roleId int64) (*[]dtos.ScreensPermissionRes, error)

	//EMPLOYEE
	CreateEmployee(employee dtos.CreateEmployeeReq) (int64, error)
	InsertLoginInfoForEmployee(user *schema.UserLogin) (int64, error)
	GetEmployees(orgId int64, limit int64, offset int64) (*[]dtos.EmployeeEntiry, error)
	GetEmpployee(employeeId int64) (*dtos.EmployeeEntiry, error)
	UpdateEmployeeActiveStatus(employeeId, isActive, version int64) error
	UpdateEmployee(employeeId int64, emp dtos.UpdateEmployeeReq) error
	UpdateEmployeeImagePath(employeeId, version int64, imagePath string) error
	GetTotalCount(orgId int64) int64
	GetVehicleAssignedToEmployee(employeeID int64) (*[]dtos.EmployeeAssigned, error)
	GetEmployeeInUser(employeeId int64) int64
	UpdateUserLoginTable(employeeId int64, emp dtos.UpdateEmployeeReq) error
}

func (rl *EmpObj) GetVehicleAssignedToEmployee(employeeID int64) (*[]dtos.EmployeeAssigned, error) {
	list := []dtos.EmployeeAssigned{}

	statuses := fmt.Sprintf("'%s', '%s'", constant.STATUS_SUBMITTED, constant.STATUS_CREATED)

	selectCountQuery := fmt.Sprintf(`SELECT u.id, t.vehicle_number, t.trip_sheet_id from user_login as u 
	LEFT JOIN trip_sheet_header as t ON ((t.created_by = u.id OR t.managed_by = u.id) AND (t.load_status IN (%s))) 
	where u.employee_id = '%v'`, statuses, employeeID)

	rl.l.Info(" selectCountQuery: ", selectCountQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(selectCountQuery)
	if err != nil {
		rl.l.Error("Error GetEmpRoles ", err)
		return nil, err
	}
	defer rows.Close()

	var userId, tripSheetId sql.NullInt64
	var vehicleNumber sql.NullString
	for rows.Next() {
		empVehicle := dtos.EmployeeAssigned{}
		err := rows.Scan(&userId,
			&vehicleNumber,
			&tripSheetId,
		)
		if err != nil {
			rl.l.Error("Error GetEmpployee scan: ", err)
			return nil, err
		}
		empVehicle.UserID = userId.Int64
		empVehicle.VehicleNumber = vehicleNumber.String
		empVehicle.TripSheetId = tripSheetId.Int64
		list = append(list, empVehicle)
	}
	return &list, nil
}

func (rl *EmpObj) GetTotalCountRole(orgId int64) int64 {

	selectCountQuery := fmt.Sprintf(`SELECT count(*) FROM thub_role WHERE org_id = '%v'`, orgId)
	rl.l.Info(" GetTotalCount select query: ", selectCountQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(selectCountQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return 0
	}

	return count.Int64
}

func (rl *EmpObj) GetTotalCount(orgId int64) int64 {

	empSelectQuery := fmt.Sprintf(`SELECT count(*) FROM employee WHERE org_id = '%v'`, orgId)
	rl.l.Info(" GetTotalCount select query: ", empSelectQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(empSelectQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return 0
	}

	return count.Int64
}

func (rl *EmpObj) GetEmployees(orgId int64, limit int64, offset int64) (*[]dtos.EmployeeEntiry, error) {
	list := []dtos.EmployeeEntiry{}

	// Get today's date in YYYY-MM-DD format for attendance check
	today := fmt.Sprintf("%s", time.Now().Format("2006-01-02"))

	empSelectQuery := fmt.Sprintf(`SELECT e.emp_id, e.first_name, e.last_name, e.employee_code, e.mobile_no, 
	e.email_id, e.role_id, e.dob, e.gender, e.aadhar_no, e.access_no, e.is_active, e.access_token, e.joining_date, 
	e.relieving_date, e.address_line1, e.address_line2, e.city, e.state, e.country, e.is_super_admin, 
	e.is_admin, e.org_id, e.login_type, e.image, e.pin_code, e.version, e.employee_performance, 
	e.vehicle_assigned, e.monthly_salary, e.annual_salary, e.annual_bonus,
	COALESCE(a.status, 'Leave') as attendance_status
	FROM employee e
	LEFT JOIN attendance a ON a.employee_id = e.emp_id AND a.check_in_date_str = '%s'
	WHERE e.org_id = '%v' ORDER BY e.updated_at DESC LIMIT %v OFFSET %v;`, today, orgId, limit, offset)
	rl.l.Info(" GetEmpployee select query: ", empSelectQuery)
	rows, err := rl.dbConnMSSQL.GetQueryer().Query(empSelectQuery)
	if err != nil {
		rl.l.Error("Error GetEmpRoles ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var firstName, lastLame, employeeCode, mobileNo, emailId, dob, gender, aadharNo, accessNo, accessToken, joiningDate, relievingDate sql.NullString
		var addressLine1, addressLine2, city, country, state, loginType, image, vehicleAssigned, employeePerformance sql.NullString
		var empId, roleId, isActive, orgIdN, version, isSuperAdmin, isAdmin, pinCode sql.NullInt64
		var monthlySalary, annualSalary, annualBonus sql.NullInt64
		var attendanceStatus sql.NullString

		employee := &dtos.EmployeeEntiry{}
		err := rows.Scan(&empId,
			&firstName,
			&lastLame,
			&employeeCode,
			&mobileNo,
			&emailId,
			&roleId,
			&dob,
			&gender,
			&aadharNo,
			&accessNo,
			&isActive,
			&accessToken,
			&joiningDate,
			&relievingDate,
			&addressLine1,
			&addressLine2,
			&city,
			&state,
			&country,
			&isSuperAdmin,
			&isAdmin,
			&orgIdN,
			&loginType,
			&image,
			&pinCode,
			&version,
			&employeePerformance,
			&vehicleAssigned, &monthlySalary, &annualSalary, &annualBonus,
			&attendanceStatus,
		)
		if err != nil {
			rl.l.Error("Error GetEmpployee scan: ", err)
			return nil, err
		}
		employee.EmpId = empId.Int64
		employee.EmployeeIDText = fmt.Sprintf("EMP_%v", empId.Int64)
		employee.FirstName = firstName.String
		employee.LastName = lastLame.String
		employee.EmployeeCode = employeeCode.String
		employee.MobileNo = mobileNo.String
		employee.EmailID = emailId.String
		employee.RoleID = roleId.Int64
		employee.DOB = dob.String
		employee.Gender = gender.String

		employee.AadharNo = aadharNo.String
		employee.AccessNo = accessNo.String
		employee.IsActive = isActive.Int64
		//employee.AccessToken = accessToken.String
		employee.JoiningDate = joiningDate.String
		employee.RelievingDate = relievingDate.String
		employee.Gender = gender.String
		employee.AddressLine1 = addressLine1.String
		employee.AddressLine2 = addressLine2.String
		employee.City = city.String
		employee.State = state.String
		employee.Country = country.String
		employee.IsSuperAdmin = isSuperAdmin.Int64
		employee.IsAdmin = isAdmin.Int64
		employee.OrgID = orgIdN.Int64
		employee.LoginType = loginType.String
		employee.Image = image.String
		employee.PinCode = pinCode.Int64
		employee.EmployeePerformance = employeePerformance.String
		employee.VehicleAssigned = vehicleAssigned.String
		employee.MonthlySalary = monthlySalary.Int64
		employee.AnnualSalary = annualSalary.Int64
		employee.AnnualBonus = annualBonus.Int64
		employee.EmployeePerformance = constant.EMPLOYEE_PERFORMANCE_2_NAME
		// Set attendance status - "Present" if checked in, "Absent" if not
		if attendanceStatus.String == "Present" {
			employee.AttendanceStatus = "Present"
		} else {
			employee.AttendanceStatus = "Absent"
		}
		list = append(list, *employee)
	}

	return &list, nil
}

func (rl *EmpObj) GetEmpployee(employeeId int64) (*dtos.EmployeeEntiry, error) {
	employee := dtos.EmployeeEntiry{}

	empSelectQuery := fmt.Sprintf(`select emp_id, first_name, last_name, employee_code, mobile_no, email_id, role_id, dob, gender, aadhar_no, access_no, is_active, access_token, joining_date, relieving_date, address_line1, address_line2, city, state, country, is_super_admin, is_admin, org_id, login_type, image, version, employee_performance, vehicle_assigned, monthly_salary, annual_salary, annual_bonus FROM employee WHERE emp_id = '%v'`, employeeId)
	rl.l.Info(" GetEmpployee select query: ", empSelectQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(empSelectQuery)

	var firstName, lastLame, employeeCode, mobileNo, emailId, dob, gender, aadharNo, accessNo, accessToken, joiningDate, relievingDate, addressLine1, addressLine2, city, country, state, loginType, image sql.NullString
	var empId, roleId, isActive, orgIdN, version, isSuperAdmin, isAdmin sql.NullInt64

	var monthlySalary, annualSalary, annualBonus sql.NullInt64
	var vehicleAssigned, employeePerformance sql.NullString

	errE := row.Scan(&empId,
		&firstName,
		&lastLame,
		&employeeCode,
		&mobileNo,
		&emailId,
		&roleId,
		&dob,
		&gender,
		&aadharNo,
		&accessNo,
		&isActive,
		&accessToken,
		&joiningDate,
		&relievingDate,
		&addressLine1,
		&addressLine2,
		&city,
		&state,
		&country,
		&isSuperAdmin,
		&isAdmin,
		&orgIdN,
		&loginType,
		&image,
		&version,
		&employeePerformance,
		&vehicleAssigned, &monthlySalary, &annualSalary, &annualBonus,
	)
	if errE != nil {
		rl.l.Error("Error GetEmpployee scan: ", errE)
		return nil, errE
	}
	rl.l.Error("isActive: ", isActive)
	employee.EmpId = empId.Int64
	employee.FirstName = firstName.String
	employee.LastName = lastLame.String
	employee.EmployeeCode = employeeCode.String
	employee.MobileNo = mobileNo.String
	employee.EmailID = emailId.String
	employee.RoleID = roleId.Int64
	employee.DOB = dob.String
	employee.Gender = gender.String
	employee.AadharNo = aadharNo.String
	employee.AccessNo = accessNo.String
	employee.IsActive = isActive.Int64
	//employee.AccessToken = accessToken.String
	employee.JoiningDate = joiningDate.String
	employee.RelievingDate = relievingDate.String
	employee.Gender = gender.String
	employee.AddressLine1 = addressLine1.String
	employee.AddressLine2 = addressLine2.String
	employee.City = city.String
	employee.State = state.String
	employee.Country = country.String
	employee.IsSuperAdmin = isSuperAdmin.Int64
	employee.IsAdmin = isAdmin.Int64
	employee.OrgID = orgIdN.Int64
	employee.LoginType = loginType.String
	employee.Image = image.String
	employee.Version = version.Int64
	employee.EmployeePerformance = employeePerformance.String
	employee.VehicleAssigned = vehicleAssigned.String
	employee.MonthlySalary = monthlySalary.Int64
	employee.AnnualSalary = annualSalary.Int64
	employee.AnnualBonus = annualBonus.Int64
	employee.EmployeePerformance = constant.EMPLOYEE_PERFORMANCE_2_NAME
	return &employee, nil
}

func (rl *EmpObj) GetEmpRoles(orgId int64, limit int64, offset int64) (*[]dtos.RoleEmp, error) {
	list := []dtos.RoleEmp{}

	query := "SELECT role_id, role_code, role_name, description, org_id, is_active, updated_at, version FROM thub_role WHERE org_id = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?"

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(query, orgId, limit, offset)
	if err != nil {
		rl.l.Error("Error GetEmpRoles ", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var roleCode, roleName, description sql.NullString
		var roleId, isActive, orgIdN, version sql.NullInt64
		var updatedAt sql.NullString

		role := &dtos.RoleEmp{}
		err := rows.Scan(&roleId,
			&roleCode,
			&roleName,
			&description,
			&orgIdN,
			&isActive,
			&updatedAt,
			&version,
		)
		if err != nil {
			rl.l.Error("Error GetEmpRoles scan: ", err)
			return nil, err
		}
		role.RoleId = roleId.Int64
		role.RoleCode = roleCode.String
		role.RoleName = roleName.String
		role.Description = description.String
		role.OrgId = orgIdN.Int64
		role.IsActive = isActive.Int64
		role.UpdatedAt = updatedAt.String
		role.Version = version.Int64

		list = append(list, *role)
	}

	return &list, nil
}

func (rl *EmpObj) GetScreensPermissions(roleId int64) (*[]dtos.ScreensPermissionRes, error) {
	list := []dtos.ScreensPermissionRes{}

	query := fmt.Sprintf(`
	SELECT t.thub_role_screens_id, t.role_id, t.role_code, t.role_name, t.website_screen_id, w.screen_name, t.permisstion_label_id, p.permisstion_label
	from thub_role_screens_map t
	LEFT JOIN permisstion_label p ON p.permisstion_label_id = t.permisstion_label_id
	LEFT JOIN website_screens w ON w.website_screen_id = t.website_screen_id
	where t.role_id = %v order by t.updated_at`, roleId)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		rl.l.Error("Error GetEmpRoles ", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var roleCode, roleName, screenName, permisstionLabel sql.NullString
		var roleId, thubRoleScreensId, websiteScreenId, permisstionLabelId sql.NullInt64

		screenPermission := &dtos.ScreensPermissionRes{}
		err := rows.Scan(&thubRoleScreensId,
			&roleId,
			&roleCode,
			&roleName,
			&websiteScreenId,
			&screenName,
			&permisstionLabelId,
			&permisstionLabel,
		)
		if err != nil {
			rl.l.Error("Error GetEmpRoles scan: ", err)
			return nil, err
		}
		screenPermission.RoleId = roleId.Int64
		screenPermission.RoleCode = roleCode.String
		screenPermission.RoleName = roleName.String
		screenPermission.WebsiteScreenID = websiteScreenId.Int64
		screenPermission.ScreenName = screenName.String
		screenPermission.ThubRoleScreensID = thubRoleScreensId.Int64
		screenPermission.PermisstionLabel = permisstionLabel.String
		screenPermission.PermisstionLabelID = permisstionLabelId.Int64

		list = append(list, *screenPermission)
	}

	return &list, nil
}

func (rl *EmpObj) UpdateEmpRole(empRole dtos.EmpRoleUpdate) error {

	updateQuery := fmt.Sprintf(`UPDATE thub_role SET 
		role_code = '%v', 
		role_name = '%v', 
		description = '%v', 
		is_active = '%v', 
		version = '%v' 
	WHERE role_id = '%v'`,
		empRole.RoleCode,
		empRole.RoleName,
		empRole.Description,
		empRole.IsActive,
		empRole.Version,
		empRole.RoleId)
	rl.l.Info("Update query ", updateQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateEmpRole) thub_role: ", err)
		return err
	}
	roleId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateEmpRole) thub_role: ", roleId, err)
		//return err
	}
	rl.l.Info("Role updated successfully: ", empRole.RoleId)
	return nil
}

func (rl *EmpObj) UpdateEmployeeActiveStatus(employeeId, isActive, version int64) error {

	updateQuery := fmt.Sprintf(`UPDATE employee SET 
		is_active = '%v', version = '%v'
	WHERE emp_id = '%v'`, isActive, version, employeeId)

	rl.l.Info("UpdateEmployeeActiveStatus Update query ", updateQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateEmployeeActiveStatus) employee: ", err)
		return err
	}
	empId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateEmployeeActiveStatus) employee: ", empId, err)
		return err
	}
	rl.l.Info("Employee updated successfully: ", empId)

	return nil
}

func (rl *EmpObj) CreateEmpRole(empRole dtos.EmpRole) (int64, error) {

	roleQuery := fmt.Sprintf(`INSERT INTO thub_role (
		role_code, role_name, description, org_id, is_active, version) VALUES ('%v', '%v', '%v', %v, %v, %v)`,
		empRole.RoleCode,
		empRole.RoleName,
		empRole.Description,
		empRole.OrgId,
		1,
		1)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(roleQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(CreateEmpRole) thub_role: ", err)
		return 0, err
	}
	roleId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(CreateEmpRole) thub_role: ", roleId, err)
		//return err
	}
	rl.l.Info("Role saved successfully: ", roleId, empRole)
	return roleId, nil
}

func (rl *EmpObj) ScreensPermission(roleCode, roleName string, roleID int64, sp dtos.ScreensPermission) error {

	roleQuery := fmt.Sprintf(`INSERT INTO thub_role_screens_map (
		role_code, role_name, role_id, website_screen_id, permisstion_label_id) VALUES ('%v', '%v', '%v', '%v', '%v')`,
		roleCode,
		roleName,
		roleID,
		sp.WebsiteScreenID,
		sp.PermisstionLabelID)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(roleQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(ScreensPermission) thub_role: ", err)
		return err
	}
	roleId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(ScreensPermission) thub_role: ", roleId, err)
		//return err
	}
	rl.l.Info("ScreensPermission saved successfully: ", roleId)
	return nil
}

func (br *EmpObj) DeleteScreensPermission(roleID int64) error {

	deleteQuery := fmt.Sprintf(`DELETE FROM thub_role_screens_map WHERE role_id = '%v';`, roleID)

	br.l.Info("DeleteScreensPermission query:", deleteQuery)

	_, err := br.dbConnMSSQL.GetQueryer().Exec(deleteQuery)
	if err != nil {
		br.l.Error("Error db.Exec(DeleteScreensPermission): ", err)
		return err
	}
	br.l.Infof("%s thub_role_screens_map deleted successfully \n", roleID)
	return nil
}

func (ul *EmpObj) IsRoleExists(roleId int64) (*dtos.RoleEmp, error) {
	role := dtos.RoleEmp{}

	var roleCode, roleName, description sql.NullString
	var version sql.NullInt64

	query := `SELECT role_code, role_name, description, version FROM thub_role WHERE role_id = ?`
	row := ul.dbConnMSSQL.GetQueryer().QueryRow(query, roleId)
	ul.l.Info("dbConnMSSQL ", query, row)
	errS := row.Scan(
		&roleCode,
		&roleName,
		&description,
		&version)
	role.RoleId = roleId
	role.RoleCode = roleCode.String
	role.RoleName = roleName.String
	role.Description = description.String
	role.Version = version.Int64
	if errS != nil {
		ul.l.Error("ERROR GetOrgBaseInfo ", roleId, errS)
		return nil, errS
	}
	return &role, nil
}

func (ul *EmpObj) GetOrgBaseInfo(orgId int64) (*dtos.OrganisationBaseResponse, error) {
	organization := dtos.OrganisationBaseResponse{}

	var name, displayName sql.NullString
	var isActive sql.NullInt64

	query := `SELECT name, display_name, is_active from organisation  WHERE org_id = ?`
	row := ul.dbConnMSSQL.GetQueryer().QueryRow(query, orgId)
	ul.l.Info("dbConnMSSQL ", query, row)
	errS := row.Scan(
		&name,
		&displayName,
		&isActive)
	organization.OrgId = orgId
	organization.Name = name.String
	organization.DisplayName = displayName.String
	organization.IsActive = isActive.Int64

	if errS != nil {
		ul.l.Error("ERROR GetOrgBaseInfo ", orgId, errS)
		return nil, errS
	}
	return &organization, nil
}

func (rl *EmpObj) CreateEmployee(emp dtos.CreateEmployeeReq) (int64, error) {

	rl.l.Info("JoiningDate : ", emp.JoiningDate)

	ceateEmpQuery := fmt.Sprintf(`INSERT INTO employee (
		first_name, last_name, employee_code, mobile_no, email_id, role_id, dob, gender, aadhar_no, access_no, 
		is_active, access_token, relieving_date, joining_date, 
		address_line1, address_line2, city, state, country, is_super_admin, 
		is_admin, org_id, login_type, image, pin_code, version, monthly_salary, annual_salary, annual_bonus) 
		VALUES 
		( '%v', '%v', '%v', '%v', '%v', '%v',
		   '%v', '%v', '%v', '%v', '%v', '%v',
		   '%v', '%v', '%v', '%v', '%v', '%v',
		   '%v', '%v', '%v', '%v', '%v', '%v','%v', '%v','%v','%v', '%v')`,
		emp.FirstName, emp.LastName, emp.EmployeeCode, emp.MobileNo, emp.EmailID, emp.RoleID,
		emp.DOB, emp.Gender, emp.AadharNo, emp.AccessNo, emp.IsActive, emp.AccessToken, emp.RelievingDate,
		emp.JoiningDate, emp.AddressLine1, emp.AddressLine2, emp.City, emp.State,
		emp.Country, emp.IsSuperAdmin, emp.IsAdmin, emp.OrgID, emp.LoginType, emp.Image, emp.PinCode, 1, emp.MonthlySalary, emp.AnnualSalary, emp.AnnualBonus)

	rl.l.Info("ceateEmpQuery : ", ceateEmpQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(ceateEmpQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(CreateEmployee) employee: ", err)
		return 0, err
	}
	createdEmpId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(CreateEmployee) employee:", createdEmpId, err)
	}
	rl.l.Info("Employee created successfully: ", createdEmpId, emp.FirstName)
	return createdEmpId, nil
}
func (ul *EmpObj) GetEmployeeInUser(employeeId int64) int64 {

	query := `SELECT id FROM user_login WHERE employee_id = ?`

	ul.l.Info("GetEmployeeInUser ", query)

	var userId sql.NullInt64
	row := ul.dbConnMSSQL.GetQueryer().QueryRow(query, employeeId)
	errS := row.Scan(&userId)
	if errS != nil {
		ul.l.Error("No user found for the employee", employeeId, errS)
		return 0
	}
	return userId.Int64
}

func (rl *EmpObj) UpdateUserLoginTable(employeeId int64, emp dtos.UpdateEmployeeReq) error {

	updateUserQuery := fmt.Sprintf(`UPDATE user_login SET 
		first_name = '%v', last_name = '%v', 
	    mobile_no = '%v', role_id = '%v', email_id = '%v'
	    WHERE employee_id = '%v'`,
		emp.FirstName, emp.LastName,
		emp.MobileNo, emp.RoleID, emp.EmailID,
		employeeId)

	rl.l.Info("updateUserQuery query ", updateUserQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(updateUserQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(updateUserQuery) employee: ", err)
		return err
	}
	updatedEmpId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(updateUserQuery) employee:", updatedEmpId, err)
	}
	rl.l.Info("user updated successfully: ", updatedEmpId, emp.FirstName)
	return nil
}

func (rl *EmpObj) UpdateEmployee(employeeId int64, emp dtos.UpdateEmployeeReq) error {

	updateEmpQuery := fmt.Sprintf(`UPDATE employee SET 
		first_name = '%v', last_name = '%v', 
		employee_code = '%v', mobile_no = '%v', role_id = '%v', dob = '%v', gender = '%v',
		aadhar_no = '%v', is_active = '%v', access_token = '%v', joining_date = '%v',
		relieving_date = '%v', address_line1 = '%v', address_line2 = '%v',
		city = '%v', state = '%v', country = '%v',
		is_super_admin = '%v', is_admin = '%v', login_type = '%v', pin_code = '%v', version= '%v', 
		monthly_salary= '%v', annual_salary= '%v', annual_bonus= '%v', email_id = '%v'
	    WHERE emp_id = '%v'`,
		emp.FirstName, emp.LastName,
		emp.EmployeeCode, emp.MobileNo, emp.RoleID, emp.DOB, emp.Gender,
		emp.AadharNo, emp.IsActive, emp.AccessToken, emp.JoiningDate,
		emp.RelievingDate, emp.AddressLine1, emp.AddressLine2,
		emp.City, emp.State, emp.Country,
		emp.IsSuperAdmin, emp.IsAdmin, emp.LoginType, emp.PinCode, emp.Version, emp.MonthlySalary, emp.AnnualSalary, emp.AnnualBonus, emp.EmailID,
		employeeId)

	rl.l.Info("UpdateEmployee query ", updateEmpQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(updateEmpQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateEmployee) employee: ", err)
		return err
	}
	updatedEmpId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateEmployee) employee:", updatedEmpId, err)
	}
	rl.l.Info("employee updated successfully: ", updatedEmpId, emp.FirstName)
	return nil
}

func (rl *EmpObj) InsertLoginInfoForEmployee(user *schema.UserLogin) (int64, error) {

	userQuery := fmt.Sprintf(`INSERT INTO 
		user_login ( first_name, 
		last_name, 
		email_id, 
		mobile_no, 
		password, 
		role_id, 
		org_id,
		is_super_admin,
		is_admin,
		login_type,
		employee_id,
		version) 
		values('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v','%v','%v')`,
		user.FirstName,
		user.LastName,
		user.EmailId,
		user.MobileNo,
		user.Password,
		user.RoleID,
		user.OrgId,
		user.IsSuperAdmin,
		user.IsAdmin,
		user.LoginType,
		user.EmployeeId,
		user.Version)

	userResult, err := rl.dbConnMSSQL.GetQueryer().Exec(userQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(InsertLoginInfoForEmployee) user_login: ", err)
		return 0, err
	}
	userID, err := userResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error InsertLoginInfoForEmployee insert ID: ", err)
	}
	return userID, nil
}

func (rl *EmpObj) UpdateEmployeeImagePath(employeeId, version int64, imagePath string) error {

	updateQuery := fmt.Sprintf(`UPDATE employee SET image = '%v', version = '%v' WHERE emp_id = '%v'`, imagePath, version, employeeId)
	rl.l.Info("UpdateEmployeeImagePath Update query: ", updateQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateEmployeeImagePath) employee: ", err)
		return err
	}
	empId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateEmployeeImagePath) employee: ", empId, err)
		return err
	}

	return nil
}
