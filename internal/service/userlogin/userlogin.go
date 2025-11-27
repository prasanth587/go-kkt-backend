package userlogin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/dtos/schema"
	"go-transport-hub/internal/daos"
	"go-transport-hub/utils"
)

type UserLoginObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
	loginDao    daos.LoginDao
}

var (
	ErrUnableToPingDB = errors.New("Unable to ping database")
	USER_SUCCESS      = "User logged in successfully"
	ERROR_IN_UPDATE   = "Error; UpdateUserLoginAterCredentialSuccess  - "
	INVALIDE_         = "invalid credentials"
	LOGIN_FAILED      = "login failed  - "
)

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *UserLoginObj {
	return &UserLoginObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		loginDao:    daos.NewLogin(l, dbConnMSSQL),
	}
}

func (ul *UserLoginObj) UserLogin(reqbody dtos.AdminLogin) (*dtos.LoginResponse, error) {
	reqbody.EmailID = strings.ToLower(reqbody.EmailID)
	reqbody.Password = utils.SHAEncoding(reqbody.Password)

	loginSchema, err := ul.loginDao.VerifyCredentials(reqbody)
	if err != nil {
		ul.l.Error(LOGIN_FAILED, reqbody.EmailID, err)
		return nil, errors.New(INVALIDE_)
	}

	loginSchema.Version = loginSchema.Version + 1
	errU := ul.loginDao.UpdateUserLoginAterCredentialSuccess(loginSchema)
	if errU != nil {
		ul.l.Error(ERROR_IN_UPDATE, reqbody.EmailID, errU)
	}
	orgInfo, errO := ul.loginDao.GetOrgDetails(loginSchema.OrgId)
	if errO != nil {
		ul.l.Error("GetOrgDetails ERROR", reqbody.EmailID, errO)
	}

	screensPermission, errA := ul.loginDao.GetScreensPermissions(loginSchema.RoleID)
	if errA != nil {
		ul.l.Error("ERROR: GetScreensPermissions", errA)
		//return nil, errA
	}
	empAttendance := schema.EmpAttendance{}
	if loginSchema.EmployeeId != 0 {
		empAttendanc, errA := ul.loginDao.GetAttananceForEmployee(loginSchema.EmployeeId, utils.GetCurrentDateStr())
		if errA != nil {
			if !strings.Contains(errA.Error(), "no rows in result set") {
				ul.l.Error("ERROR: GetAttendanceForEmployee", errA)
			}
		}
		if empAttendanc != nil {
			empAttendance = *empAttendanc
		}
	}

	sessionTimeoutMS := 1800000 // 30 minutes in milliseconds

	appConfigs, errG := ul.loginDao.GetAppConfig()
	if errG != nil {
		ul.l.Error("ERROR: GetAppConfig", errG)
	}

	for _, config := range *appConfigs {
		switch config.ConfigCode {
		case constant.SESSION_TIMEOUT:
			sessionTimeoutMS, _ = strconv.Atoi(config.ConfigValue)
		}
	}

	userRes := &dtos.LoginResponse{}
	userRes.Message = USER_SUCCESS
	userRes.EmailId = loginSchema.EmailId
	userRes.LoginId = loginSchema.ID
	userRes.FirstName = loginSchema.FirstName
	userRes.MobileNo = loginSchema.MobileNo
	userRes.RoleID = loginSchema.RoleID
	userRes.RoleName = loginSchema.RoleName
	userRes.EmployeeId = loginSchema.EmployeeId
	userRes.LoginType = loginSchema.LoginType
	userRes.LastLogin = time.Now().In(utils.TimeLoc())
	userRes.OrganisationResponse = orgInfo
	userRes.SessionTimeoutMS = sessionTimeoutMS
	userRes.EmpAttendance = empAttendance
	ul.l.Info("loginSchema.IsSuperAdmin", loginSchema.IsSuperAdmin)
	if loginSchema.IsSuperAdmin == 1 {
		userRes.Settings.IsEdit = true
		userRes.Settings.IsView = true
		userRes.Employee.IsEdit = true
		userRes.Employee.IsView = true
		getScreenPermissionsTrue(screensPermission, userRes)
	} else {
		getScreenPermissions(screensPermission, userRes)
	}

	ul.l.Info("isCustomerEdit ", userRes.Customer.IsEdit, "isCustomerView", userRes.Customer.IsView, "isCustomerNoAccess", userRes.Customer.IsNoAccess,
		"isOperationsEdit ", userRes.Operations.IsEdit, "isOperationsView", userRes.Operations.IsView, "isOperationsNoAccess", userRes.Operations.IsNoAccess,
		"isOverviewEdit ", userRes.Overview.IsEdit, "isOverviewView", userRes.Overview.IsView, "isOverviewNoAccess", userRes.Overview.IsNoAccess,
		"isSettingsEdit ", userRes.Settings.IsEdit, "isSettingsView", userRes.Settings.IsView, "isSettingsNoAccess", userRes.Settings.IsNoAccess,
		"istripMgmtEdit ", userRes.TripMgmt.IsEdit, "istripMgmtView", userRes.TripMgmt.IsView, "istripMgmtNoAccess", userRes.TripMgmt.IsNoAccess,
		"isVendorsEdit ", userRes.Vendors.IsEdit, "isVendorsView", userRes.Vendors.IsView, "isVendorsNoAccess", userRes.Vendors.IsNoAccess,
		"isReportsEdit ", userRes.Reports.IsEdit, "isReportsView", userRes.Reports.IsView, "isReportsNoAccess", userRes.Reports.IsNoAccess)

	return userRes, nil
}

func getScreenPermissions(screensPermission *[]dtos.ScreensPermissionRes, userRes *dtos.LoginResponse) {
	for _, screen := range *screensPermission {
		switch screen.ScreenName {
		case "Customers":
			userRes.Customer = mapPermissions(screen.PermisstionLabel)
		case "Operations":
			userRes.Operations = mapPermissions(screen.PermisstionLabel)
		case "Overview":
			userRes.Overview = mapPermissions(screen.PermisstionLabel)
		case "Settings":
			userRes.Settings = mapPermissions(screen.PermisstionLabel)
		case "Trip Management":
			userRes.TripMgmt = mapPermissions(screen.PermisstionLabel)
		case "Vendors":
			userRes.Vendors = mapPermissions(screen.PermisstionLabel)
		case "Reports":
			userRes.Reports = mapPermissions(screen.PermisstionLabel)
		case "Employee":
			userRes.Employee = mapPermissions(screen.PermisstionLabel)
		}
	}
}

func getScreenPermissionsTrue(screensPermission *[]dtos.ScreensPermissionRes, userRes *dtos.LoginResponse) {
	for _, screen := range *screensPermission {
		switch screen.ScreenName {
		case "Customers":
			userRes.Customer = dtos.ScreenPermission{IsEdit: true, IsView: true, IsNoAccess: false}
		case "Operations":
			userRes.Operations = dtos.ScreenPermission{IsEdit: true, IsView: true, IsNoAccess: false}
		case "Overview":
			userRes.Overview = dtos.ScreenPermission{IsEdit: true, IsView: true, IsNoAccess: false}
		case "Settings":
			userRes.Settings = dtos.ScreenPermission{IsEdit: true, IsView: true, IsNoAccess: false}
		case "Trip Management":
			userRes.TripMgmt = dtos.ScreenPermission{IsEdit: true, IsView: true, IsNoAccess: false}
		case "Vendors":
			userRes.Vendors = dtos.ScreenPermission{IsEdit: true, IsView: true, IsNoAccess: false}
		case "Reports":
			userRes.Reports = dtos.ScreenPermission{IsEdit: true, IsView: true, IsNoAccess: false}
		case "Employee":
			userRes.Employee = dtos.ScreenPermission{IsEdit: true, IsView: true, IsNoAccess: false}
		}
	}
}

func (ul *UserLoginObj) CreateUser(reqbody schema.UserLogin) (*dtos.Messge, error) {

	reqbody.Password = utils.SHAEncoding("welcome123")

	errC := ul.loginDao.CreateUser(reqbody)
	if errC != nil {
		ul.l.Error(ERROR_IN_UPDATE, reqbody.FirstName, errC)
		return nil, errC
	}

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("User created saved successfully: %s", reqbody.FirstName)
	return &response, nil
}

func (ul *UserLoginObj) UpdateUser(reqbody schema.UserLogin) (*dtos.Messge, error) {
	errC := ul.loginDao.UpdateUser(reqbody)
	if errC != nil {
		ul.l.Error("UpdateUser error:", reqbody.FirstName, errC)
		return nil, errC
	}

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("User updated successfully: %s", reqbody.FirstName)
	return &response, nil
}

func (ul *UserLoginObj) UpdateUserActiveStatus(userId int64, isActive int) (*dtos.Messge, error) {
	errC := ul.loginDao.UpdateUserActiveStatus(userId, isActive)
	if errC != nil {
		ul.l.Error("UpdateUserActiveStatus error:", userId, errC)
		return nil, errC
	}

	status := "activated"
	if isActive == 0 {
		status = "deactivated"
	}
	response := dtos.Messge{}
	response.Message = fmt.Sprintf("User %s successfully", status)
	return &response, nil
}

func (br *UserLoginObj) GetUserList(limit, offset, searchText, loginId string) (*dtos.LoginRes, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	whereQuery := br.loginDao.BuildWhereQuery(loginId, searchText)

	res, errA := br.loginDao.GetUsers(limitI, offsetI, whereQuery)
	if errA != nil {
		br.l.Error("ERROR: GetLoadingPoints", errA)
		return nil, errA
	}

	users := dtos.LoginRes{}
	users.Users = res
	users.Total = br.loginDao.GetTotalCount(whereQuery)
	users.Limit = limitI
	users.OffSet = offsetI
	return &users, nil
}

func mapPermissions(label string) dtos.ScreenPermission {
	edit := label == "EDIT"
	view := label == "VIEW"
	if edit {
		view = true
	}
	noAccess := !edit && !view
	return dtos.ScreenPermission{
		IsEdit:     edit,
		IsView:     view,
		IsNoAccess: noAccess,
	}
}

func (ul *UserLoginObj) UpdatePassword(userLogin dtos.UpdatePassword) (*dtos.Messge, error) {

	if userLogin.Password == "" {
		ul.l.Error("Error Password: Password should not empty")
		return nil, errors.New("password should not empty")
	}
	if userLogin.UserLogin == "" {
		ul.l.Error("Error UserLogin: UserLogin should not empty")
		return nil, errors.New("userLogin should not empty")
	}

	userLogin.Password = utils.SHAEncoding(userLogin.Password)
	err := ul.loginDao.UpdatePassword(userLogin)
	if err != nil {
		ul.l.Error("Password not saved ", userLogin.UserLogin, err)
		return nil, err
	}

	ul.l.Info("Password updated successfully! : ")

	resMsg := dtos.Messge{}
	resMsg.Message = fmt.Sprintf("Password saved successfully %s, login now!", userLogin.UserLogin)
	return &resMsg, nil
}
