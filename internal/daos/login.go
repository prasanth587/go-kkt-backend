package daos

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/dtos/schema"
	"go-transport-hub/utils"
)

type LoginObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewLogin(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *LoginObj {
	return &LoginObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

var (
	UPDATE_ERROR = "ERROR UpdateUserLoginAterCredentialSuccess update: %v "
	VERIFY_ERROR = "ERROR VerifyCredentials: %v "
	INVALIDE_    = "invalid credentials"
	LOGIN_FAILED = "login failed  - "
)

type LoginDao interface {
	VerifyCredentials(reqbody dtos.AdminLogin) (*schema.UserLogin, error)
	UpdateUserLoginAterCredentialSuccess(userLogin *schema.UserLogin) error
	GetOrgDetails(orgId int64) (*dtos.OrganisationResponse, error)
	CreateUser(user schema.UserLogin) error
	UpdateUser(user schema.UserLogin) error
	UpdateUserActiveStatus(userId int64, isActive int) error
	BuildWhereQuery(loginId, searchText string) string
	GetTotalCount(whereQuery string) int64
	GetUsers(limit int64, offset int64, whereQuery string) (*[]dtos.LoginV1Response, error)
	GetScreensPermissions(roleId int64) (*[]dtos.ScreensPermissionRes, error)
	GetAppConfig() (*[]schema.AppConfig, error)
	UpdatePassword(passwordObj dtos.UpdatePassword) error
	GetAttananceForEmployee(employeeID int64, today string) (*schema.EmpAttendance, error)
}

func (rl *LoginObj) GetTotalCount(whereQuery string) int64 {

	countQuery := fmt.Sprintf(`SELECT count(*) FROM user_login as u %v`, whereQuery)
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

func (ul *LoginObj) VerifyCredentials(reqbody dtos.AdminLogin) (*schema.UserLogin, error) {
	userLogin := schema.UserLogin{}

	var firstName, accessToken, loginType, mobileNo, rollName, emailId sql.NullString
	var employeeId, roleID, orgId, version, isSuperAdmin, isAdmin sql.NullInt64
	userId := reqbody.EmailID
	if strings.Contains(reqbody.EmailID, "user") {
		userId = regexp.MustCompile(`[0-9]+`).FindString(userId)
	}

	query := `SELECT u.id,u.first_name, u.access_token, u.role_id, u.employee_id, u.is_super_admin, u.is_admin, u.org_id, u.login_type, u.version, u.mobile_no, r.role_name, u.email_id
				FROM user_login as u
				LEFT JOIN thub_role as r ON r.role_id = u.role_id
				WHERE (u.email_id = ? OR u.mobile_no = ? OR u.id = ?) AND u.password = ?`
	ul.l.Info("Login query for user: ", userId)

	row := ul.dbConnMSSQL.GetQueryer().QueryRow(query, userId, userId, userId, reqbody.Password)
	errS := row.Scan(&userLogin.ID, &firstName, &accessToken, &roleID, &employeeId, &isSuperAdmin, &isAdmin, &orgId, &loginType, &version, &mobileNo, &rollName, &emailId)
	if errS != nil {
		ul.l.Error(VERIFY_ERROR, reqbody.EmailID, errS)
		return nil, errS
	}
	userLogin.FirstName = firstName.String
	userLogin.AccessToken = accessToken.String
	userLogin.LoginType = loginType.String
	userLogin.EmployeeId = employeeId.Int64
	userLogin.EmailId = emailId.String
	userLogin.RoleID = roleID.Int64
	userLogin.OrgId = orgId.Int64
	userLogin.Version = version.Int64
	userLogin.IsSuperAdmin = int(isSuperAdmin.Int64)
	userLogin.IsAdmin = int(isAdmin.Int64)
	userLogin.RoleName = rollName.String
	userLogin.MobileNo = mobileNo.String
	return &userLogin, nil
}

func (ul *LoginObj) UpdateUserLoginAterCredentialSuccess(userLogin *schema.UserLogin) error {
	queryLLU := `UPDATE user_login SET version = ?, last_login = ?  WHERE id = ?`
	_, errup := ul.dbConnMSSQL.GetQueryer().Exec(queryLLU, userLogin.Version, time.Now().In(utils.TimeLoc()), userLogin.ID)
	if errup != nil {
		ul.l.Error(UPDATE_ERROR, userLogin.EmailId, errup)
		return errup
	}
	return nil
}

func (ul *LoginObj) GetOrgDetails(orgId int64) (*dtos.OrganisationResponse, error) {
	ul.l.Info("fetching organisation info ", orgId)
	organization := dtos.OrganisationResponse{}

	var name, displayName, domainName, emailId, contactName, contactNo, logoPath, addressLine1, addressLine2, city sql.NullString
	var orgID, isActive sql.NullInt64

	query := `SELECT org_id, name, display_name, domain_name, email_id, contact_name, contact_no, is_active, logo_path, address_line1, address_line2, city from organisation  WHERE org_id = ?`
	row := ul.dbConnMSSQL.GetQueryer().QueryRow(query, orgId)
	ul.l.Info("dbConnMSSQL ", query, row)
	errS := row.Scan(&orgID,
		&name,
		&displayName,
		&domainName,
		&emailId,
		&contactName,
		&contactNo,
		&isActive,
		&logoPath,
		&addressLine1,
		&addressLine2,
		&city)
	organization.OrgId = orgID.Int64
	organization.Name = name.String
	organization.DisplayName = displayName.String
	organization.DomainName = domainName.String
	organization.EmailId = emailId.String
	organization.ContactName = contactName.String
	organization.ContactNo = contactNo.String
	organization.IsActive = isActive.Int64
	organization.LogoPath = logoPath.String
	organization.AddressLine1 = addressLine1.String
	organization.AddressLine2 = addressLine2.String
	organization.City = city.String

	if errS != nil {
		ul.l.Error("ERROR GetOrgDetails ", orgId, errS)
		return nil, errS
	}

	return &organization, nil
}

func (rl *LoginObj) CreateUser(user schema.UserLogin) error {

	userQuery := fmt.Sprintf(`INSERT INTO 
		user_login (
		first_name, 
		last_name, 
		email_id, 
		mobile_no,
		password,
		role_id, 
		org_id,
		login_type,
		is_active) 
		values('%v', '%v', '%v', '%v', '%v', '%v','%v','%v','%v')`,
		user.FirstName,
		user.LastName,
		user.EmailId,
		user.MobileNo,
		user.Password,
		user.RoleID,
		user.OrgId, user.LoginType, 1)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(userQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(Create user login) user_login: ", err)
		return err
	}
	return nil
}

func (rl *LoginObj) UpdateUser(user schema.UserLogin) error {
	updateQuery := fmt.Sprintf(`UPDATE user_login SET 
		first_name = '%v', 
		last_name = '%v', 
		email_id = '%v', 
		mobile_no = '%v',
		role_id = '%v',
		login_type = '%v'
		WHERE id = %v`,
		user.FirstName,
		user.LastName,
		user.EmailId,
		user.MobileNo,
		user.RoleID,
		user.LoginType,
		user.ID)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(Update user login) user_login: ", err)
		return err
	}
	return nil
}

func (rl *LoginObj) UpdateUserActiveStatus(userId int64, isActive int) error {
	updateQuery := fmt.Sprintf(`UPDATE user_login SET is_active = %v WHERE id = %v`, isActive, userId)
	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateUserActiveStatus) user_login: ", err)
		return err
	}
	return nil
}

func (rl *LoginObj) BuildWhereQuery(loginId, searchText string) string {

	whereQuery := ""
	if loginId != "" || searchText != "" {
		whereQuery = "WHERE u.is_active IN (0,1)"
	}

	if loginId != "" {
		whereQuery = fmt.Sprintf(" %v AND  u.id = '%v'", whereQuery, loginId)
	}

	if searchText != "" {
		whereQuery = fmt.Sprintf(" %v AND (u.id LIKE '%%%v%%' OR u.first_name LIKE '%%%v%%' OR u.last_name LIKE '%%%v%%' OR u.mobile_no LIKE '%%%v%%' OR u.email_id LIKE '%%%v%%' OR u.login_type LIKE '%%%v%%' OR u.last_login LIKE '%%%v%%' ) ", whereQuery, searchText, searchText, searchText, searchText, searchText, searchText, searchText)
	}

	rl.l.Info("user login whereQuery:\n ", whereQuery)

	return whereQuery
}

func (br *LoginObj) GetUsers(limit int64, offset int64, whereQuery string) (*[]dtos.LoginV1Response, error) {
	list := []dtos.LoginV1Response{}

	whereQuery = fmt.Sprintf(" %v ORDER BY u.updated_at DESC LIMIT %v OFFSET %v", whereQuery, limit, offset)

	loadingpointQuery := fmt.Sprintf(`SELECT u.id, u.first_name, u.last_name, u.mobile_no, u.email_id, 
              u.is_active, u.role_id, u.login_type, u.last_login, r.role_code, r.role_name
              FROM user_login as u
              LEFT JOIN thub_role as r ON r.role_id = u.role_id %v;`, whereQuery)

	br.l.Info("loadingpointQuery:\n ", loadingpointQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(loadingpointQuery)
	if err != nil {
		br.l.Error("Error LoadingPoints ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var firstName, lastName, mobileNo, emailId, loginType, lastLogin, roleCode, roleName sql.NullString
		var userId, roleId, isActive sql.NullInt64

		user := &dtos.LoginV1Response{}
		err := rows.Scan(&userId, &firstName, &lastName, &mobileNo, &emailId, &isActive, &roleId, &loginType, &lastLogin, &roleCode, &roleName)
		if err != nil {
			br.l.Error("Error GetLoadingpoints scan: ", err)
			return nil, err
		}

		lastLoginDate := lastLogin.String
		if lastLoginDate != "" {
			lastLoginDate = br.GetReadableTime(lastLoginDate)
		}

		user.UserId = userId.Int64
		user.UserIDText = fmt.Sprintf("USER_%v", userId.Int64)
		user.FirstName = firstName.String
		user.LastName = lastName.String
		user.IsActive = isActive.Int64
		user.MobileNo = mobileNo.String
		user.EmailId = emailId.String
		user.RoleID = roleId.Int64
		user.LoginType = loginType.String
		user.LastLogin = lastLoginDate
		user.RoleCode = roleCode.String
		user.RoleName = roleName.String

		list = append(list, *user)
	}

	return &list, nil
}

func (br *LoginObj) GetReadableTime(lastLogin string) string {
	t, err := time.Parse("2006-01-02 15:04:05", lastLogin)
	if err != nil {
		br.l.Error("ErrError parsing time:: ", err)
		return lastLogin
	}
	return t.Format("2006-01-02 03:04 PM")
}

func (rl *LoginObj) GetScreensPermissions(roleId int64) (*[]dtos.ScreensPermissionRes, error) {
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

func (rl *LoginObj) GetAppConfig() (*[]schema.AppConfig, error) {
	list := []schema.AppConfig{}

	query := "select app_config_id, config_code, config_name, value from app_config where config_code = '%s';"
	query = fmt.Sprintf(query, constant.SESSION_TIMEOUT)
	rl.l.Info("GetAppConfig: ", query)
	rows, err := rl.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		rl.l.Error("Error GetEmpRoles ", err)
		return nil, err
	}
	//app_config_id,config_code,config_name,value
	defer rows.Close()
	for rows.Next() {
		var configCode, configName, value sql.NullString
		var appconfigid sql.NullInt64

		appConfig := &schema.AppConfig{}
		err := rows.Scan(&appconfigid,
			&configCode,
			&configName,
			&value,
		)
		if err != nil {
			rl.l.Error("Error GetEmpRoles scan: ", err)
			return nil, err
		}
		appConfig.ConfigCode = configCode.String
		appConfig.ConfigName = configName.String
		appConfig.ConfigValue = value.String
		appConfig.AppConfigID = appconfigid.Int64
		list = append(list, *appConfig)
	}
	return &list, nil
}

func (rl *LoginObj) UpdatePassword(passwordObj dtos.UpdatePassword) error {
	updateQuery := fmt.Sprintf(`UPDATE user_login SET password = '%s' WHERE mobile_no = '%s' OR email_id = '%s'`, passwordObj.Password, passwordObj.UserLogin, passwordObj.UserLogin)
	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdatePassword) user_login: ", err)
		return err
	}
	//roleResult
	return nil
}

func (rl *LoginObj) GetAttananceForEmployee(employeeID int64, today string) (*schema.EmpAttendance, error) {
	empAttendance := schema.EmpAttendance{}
	empAttQuery := fmt.Sprintf(`SELECT check_in_date_str,in_time,out_time from attendance Where employee_id = "%v" AND check_in_date_str = "%s"`, employeeID, today)

	var checkInDateStr, inIime, outTime sql.NullString

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(empAttQuery)

	errS := row.Scan(&checkInDateStr, &inIime, &outTime)
	if errS != nil {
		rl.l.Error("ERROR GetAttananceForEmployee", employeeID, errS)
		return nil, errS
	}
	empAttendance.CheckINDate = checkInDateStr.String
	empAttendance.InTime = inIime.String
	empAttendance.OutTime = outTime.String

	return &empAttendance, nil
}
