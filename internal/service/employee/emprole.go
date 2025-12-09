package employee

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/dtos/schema"
	"go-transport-hub/internal/daos"
	"go-transport-hub/utils"
)

type EmpRoleObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
	empRoleDao  daos.EmpRoleDao
}

var (
	ErrUnableToPingDB = errors.New("Unable to ping database")
	USER_SUCCESS      = "User logged in successfully"
	ERROR_IN_UPDATE   = "Error; UpdateUserLoginAterCredentialSuccess  - "
	INVALIDE_         = "invalid credentials"
	LOGIN_FAILED      = "login failed  - "
)

const (
	EmployeeMobilePattern = `^((\+)?(\d{2}[-])?(\d{10}){1})?(\d{11}){0,1}?$`
)

// DatePattern
const (
	DatePattern = `^\d{1,2}\/\d{1,2}\/\d{4}$`
)

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *EmpRoleObj {
	return &EmpRoleObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		empRoleDao:  daos.NewEmpRoleObj(l, dbConnMSSQL),
	}
}

func (ul *EmpRoleObj) GetEmployee(orgId int64, limit string, offset string) (*dtos.EmployeeEntries, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}
	employee := &dtos.EmployeeEntries{}

	res, errA := ul.empRoleDao.GetEmployees(orgId, limitI, offsetI)
	if errA != nil {
		ul.l.Error("ERROR: GetEmployee", errA)
		return nil, errA
	}
	//date
	//res
	for i := range *res {
		if (*res)[i].EmpId != 0 {
			empVehicle, err := ul.empRoleDao.GetVehicleAssignedToEmployee((*res)[i].EmpId)
			if err != nil {
				ul.l.Error("ERROR: GetVehicleAssignedToEmployee", err)
				continue
			}
			if empVehicle == nil || len(*empVehicle) == 0 {
				continue
			}

			var vehicleList []string
			vehicleSeen := make(map[string]bool)
			for _, empV := range *empVehicle {
				if empV.VehicleNumber == "" {
					continue
				}
				if !vehicleSeen[empV.VehicleNumber] {
					vehicleSeen[empV.VehicleNumber] = true
					vehicleList = append(vehicleList, empV.VehicleNumber)
				}
			}
			(*res)[i].VehicleAssigned = strings.Join(vehicleList, " ")
		}

		// contactInfo, errA := vh.vendorDao.GetVendorContactInfo((*res)[i].VendorId)
		// if errA != nil {
		// 	vh.l.Error("ERROR: GetVendorContactInfo", errA)
		// 	return nil, errA
		// }
		// (*res)[i].ContactInfo = *contactInfo

		// // Getting vehicles
		// vehicles, errA := vh.vendorDao.GetVehiclesByVendorID((*res)[i].VendorId)
		// if errA != nil {
		// 	vh.l.Error("ERROR: GetVehiclesByVendorID", errA)
		// 	return nil, errA
		// }
		// (*res)[i].Vehicles = *vehicles

		// //
		// doclarations, errA := vh.vendorDao.GetDeclarationDocByVendorID((*res)[i].VendorId)
		// if errA != nil {
		// 	vh.l.Error("ERROR: GetDeclarationDocByVendorID", errA)
		// 	return nil, errA
		// }
		// (*res)[i].DeclarationDocuments = *doclarations
	}

	employee.EmployeeEntiry = res
	employee.Total = ul.empRoleDao.GetTotalCount(orgId)
	employee.Limit = limitI
	employee.OffSet = offsetI

	return employee, nil
}

func (ul *EmpRoleObj) UpdateEmployeeActiveStatus(isA string, employeeId int64) (*dtos.EmpActiveStatusResponse, error) {

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid status update")
	}
	if employeeId == 0 {
		ul.l.Error("unknown employee", employeeId)
		return nil, errors.New("unknown employee")
	}
	empStatusRes := dtos.EmpActiveStatusResponse{}
	emplopleeInfo, errA := ul.empRoleDao.GetEmpployee(employeeId)
	if errA != nil {
		ul.l.Error("ERROR: GetEmployee not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(emplopleeInfo)
	ul.l.Info("GetEmpployee: ******* ", string(jsonBytes))
	emplopleeInfo.Version = emplopleeInfo.Version + 1

	errU := ul.empRoleDao.UpdateEmployeeActiveStatus(employeeId, isActive, emplopleeInfo.Version)
	if errA != nil {
		ul.l.Error("ERROR: UpdateEmployeeActiveStatus", errU)
		return nil, errA
	}
	empStatusRes.Message = "Employee updated successfully"
	empStatusRes.Name = emplopleeInfo.FirstName
	empStatusRes.IsActive = isActive
	empStatusRes.EmpId = employeeId
	return &empStatusRes, nil
}

func (ul *EmpRoleObj) CreateEmployee(employeeReq dtos.CreateEmployeeReq) (*dtos.EmployeeCreateResponse, error) {

	if employeeReq.FirstName == "" {
		ul.l.Error("Error CreateEmployee: first name should not empty")
		return nil, errors.New("first name should not empty")
	}
	if employeeReq.EmployeeCode == "" {
		ul.l.Error("Error CreateEmployee: Employee code should not empty")
		return nil, errors.New("employee code should not empty")
	}
	if employeeReq.MobileNo == "" {
		if mobileNumberLength := len(employeeReq.MobileNo); mobileNumberLength < 10 || mobileNumberLength > 15 {
			return nil, errors.New("mobile mumber should be length of 10 to 15 digits")
		}
		ul.l.Error("Error CreateEmployee: mobile numger should not empty")
		return nil, errors.New("mobile numger should not empty")
	}
	if employeeReq.RoleID == 0 {
		ul.l.Error("Error CreateEmployee: role not set")
		return nil, errors.New("employee role not set properly try login again")
	}

	if employeeReq.JoiningDate == "" {
		ul.l.Error("Error CreateEmployee: joining_date should not empty")
		return nil, errors.New("joining date should not empty")
	}

	if employeeReq.Gender == "" {
		ul.l.Error("Error CreateEmployee: gender should not empty")
		return nil, errors.New("gender should not empty")
	}
	if employeeReq.AadharNo == "" {
		ul.l.Error("Error CreateEmployee: AadharNo should not empty")
		return nil, errors.New("AadharNo should not empty")
	}
	if employeeReq.DOB == "" {
		ul.l.Error("Error CreateEmployee: JoiningDate should not empty")
		return nil, errors.New("JoiningDate should not empty")
	}
	if employeeReq.OrgID == 0 {
		ul.l.Error("Error CreateEmployee: organisation not set properly try login again")
		return nil, errors.New("organisation not set properly try login again")
	}

	employeeId, err := ul.empRoleDao.CreateEmployee(employeeReq)
	if err != nil {
		ul.l.Error("Employee not saved ", employeeReq.FirstName, err)
		return nil, err
	}

	ul.l.Info("Employee created successfully! : ", employeeId, employeeReq.FirstName)
	//Create Login fpoe the Employee
	if employeeReq.MobileNo != "" {
		user := &schema.UserLogin{
			FirstName:    employeeReq.FirstName,
			LastName:     employeeReq.LastName,
			EmailId:      employeeReq.EmailID,
			MobileNo:     employeeReq.MobileNo,
			Password:     utils.SHAEncoding("welcome123"),
			RoleID:       employeeReq.RoleID,
			OrgId:        employeeReq.OrgID,
			IsSuperAdmin: int(employeeReq.IsSuperAdmin),
			IsAdmin:      int(employeeReq.IsAdmin),
			LoginType:    employeeReq.LoginType,
			EmployeeId:   employeeId,
		}
		userId, err := ul.empRoleDao.InsertLoginInfoForEmployee(user)
		if err != nil {
			ul.l.Error("InsertLoginInfoForEmployee ERROR: user insert failed ", employeeReq.FirstName, err)
			return nil, err
		}
		ul.l.Info("User login details successfully! : ", userId, employeeId, employeeReq.FirstName)
	}

	roleResponse := dtos.EmployeeCreateResponse{}
	roleResponse.Message = fmt.Sprintf("Employee saved successfully: %s", employeeReq.FirstName)
	roleResponse.Login = fmt.Sprintf("Enabled the login for : %s", employeeReq.FirstName)
	roleResponse.EmployeeId = employeeId

	return &roleResponse, nil
}

func (ul *EmpRoleObj) UpdateEmployee(employeeId int64, employeeReq dtos.UpdateEmployeeReq) (*dtos.EmployeeUpdateResponse, error) {

	if employeeReq.FirstName == "" {
		ul.l.Error("Error UpdateEmployee: first name should not empty")
		return nil, errors.New("first name should not empty")
	}
	if employeeReq.EmployeeCode == "" {
		ul.l.Error("Error UpdateEmployee: Employee code should not empty")
		return nil, errors.New("employee code should not empty")
	}
	if employeeReq.MobileNo == "" {
		if mobileNumberLength := len(employeeReq.MobileNo); mobileNumberLength < 10 || mobileNumberLength > 15 {
			return nil, errors.New("mobile mumber should be length of 10 to 15 digits")
		}
		ul.l.Error("Error UpdateEmployee: mobile numger should not empty")
		return nil, errors.New("mobile numger should not empty")
	}
	if employeeReq.RoleID == 0 {
		ul.l.Error("Error UpdateEmployee: role not set")
		return nil, errors.New("employee role not set properly try login again")
	}

	if employeeReq.JoiningDate == "" {
		ul.l.Error("Error UpdateEmployee: joining_date should not empty")
		return nil, errors.New("joining date should not empty")
	}

	if employeeReq.Gender == "" {
		ul.l.Error("Error UpdateEmployee: gender should not empty")
		return nil, errors.New("gender should not empty")
	}
	if employeeReq.AadharNo == "" {
		ul.l.Error("Error UpdateEmployee: AadharNo should not empty")
		return nil, errors.New("AadharNo should not empty")
	}
	if employeeReq.DOB == "" {
		ul.l.Error("Error UpdateEmployee: JoiningDate should not empty")
		return nil, errors.New("JoiningDate should not empty")
	}

	emplopleeInfo, errA := ul.empRoleDao.GetEmpployee(employeeId)
	if errA != nil {
		ul.l.Error("ERROR: UpdateEmployee(GetEmployee) not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(emplopleeInfo)
	ul.l.Info("UpdateEmployee(GetEmployee): ******* ", string(jsonBytes))
	employeeReq.Version = emplopleeInfo.Version + 1
	employeeReq.IsActive = emplopleeInfo.IsActive

	err := ul.empRoleDao.UpdateEmployee(employeeId, employeeReq)
	if err != nil {
		ul.l.Error("UpdateEmployee not updated: ", employeeReq.FirstName, err)
		return nil, err
	}

	userId := ul.empRoleDao.GetEmployeeInUser(employeeId)
	if userId != 0 {
		err := ul.empRoleDao.UpdateUserLoginTable(employeeId, employeeReq)
		if err != nil {
			ul.l.Error("UpdateUserLoginTable not updated: ", employeeReq.FirstName, err)
			return nil, err
		}
	}

	ul.l.Info("employee updated successfully!: ", emplopleeInfo.EmpId, employeeReq.FirstName)

	employeeUpdateRes := dtos.EmployeeUpdateResponse{}
	employeeUpdateRes.Message = fmt.Sprintf("Employee updated successfully: %s", employeeReq.FirstName)
	employeeUpdateRes.EmpId = emplopleeInfo.EmpId
	employeeUpdateRes.Name = emplopleeInfo.FirstName

	return &employeeUpdateRes, nil
}

func (ul *EmpRoleObj) GetRoles(orgId int64, limit string, offset string) (*dtos.RoleEmpEntries, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}
	roleEmpEntries := dtos.RoleEmpEntries{}

	res, errA := ul.empRoleDao.GetEmpRoles(orgId, limitI, offsetI)
	if errA != nil {
		ul.l.Error("ERROR: GetRoles", errA)
		return nil, errA
	}

	for i := range *res {
		screensPermission, errA := ul.empRoleDao.GetScreensPermissions((*res)[i].RoleId)
		if errA != nil {
			ul.l.Error("ERROR: GetScreensPermissions", errA)
			return nil, errA
		}
		var overview, tripMgmt, vendors, customers, reports, operations, settings string
		for _, screens := range *screensPermission {
			switch screens.ScreenName {
			case "Customers":
				customers = screens.PermisstionLabel
			case "Operations":
				operations = screens.PermisstionLabel
			case "Overview":
				overview = screens.PermisstionLabel
			case "Settings":
				settings = screens.PermisstionLabel
			case "Trip Management":
				tripMgmt = screens.PermisstionLabel
			case "Vendors":
				vendors = screens.PermisstionLabel
			case "Reports":
				reports = screens.PermisstionLabel
			}
		}
		(*res)[i].Overview = overview
		(*res)[i].Customer = customers
		(*res)[i].Vendors = vendors
		(*res)[i].Operations = operations
		(*res)[i].TripManagement = tripMgmt
		(*res)[i].Settings = settings
		(*res)[i].Reports = reports
		(*res)[i].ScreensPermission = screensPermission
	}

	roleEmpEntries.RoleEmpEntrie = res
	roleEmpEntries.Total = ul.empRoleDao.GetTotalCountRole(orgId)
	roleEmpEntries.Limit = limitI
	roleEmpEntries.OffSet = offsetI
	return &roleEmpEntries, nil
}

func (ul *EmpRoleObj) CreateEmployeeRole(empRole dtos.EmpRole) (*dtos.EmpRoleResponse, error) {
	empRole.RoleCode = strings.ToUpper(empRole.RoleCode)

	if empRole.RoleName == "" {
		ul.l.Error("Error CreateEmployeeRole: Role name should not empty")
		return nil, errors.New("role name should not empty")
	}
	if empRole.RoleCode == "" {
		ul.l.Error("Error CreateEmployeeRole: Role Code should not empty")
		return nil, errors.New("role code should not empty")
	}
	if empRole.OrgId == 0 {
		ul.l.Error("Error CreateEmployeeRole: something wrong, role organization not set properly")
		return nil, errors.New("something went wrong, not created try login again")
	}

	roleId, err := ul.empRoleDao.CreateEmpRole(empRole)
	if err != nil {
		ul.l.Error("Role not saved ", empRole, err)
		return nil, err
	}

	for _, screenP := range empRole.ScreensPermission {
		err := ul.empRoleDao.ScreensPermission(empRole.RoleCode, empRole.RoleName, roleId, screenP)
		if err != nil {
			ul.l.Error("ScreensPermission not saved ", empRole, err)
			return nil, err
		}
	}

	ul.l.Info("ScreensPermission saved successfully : ", len(empRole.ScreensPermission))

	roleResponse := dtos.EmpRoleResponse{}
	roleResponse.Message = fmt.Sprintf("Role saved successfully : %s", empRole.RoleName)
	roleResponse.RoleCode = empRole.RoleCode
	roleResponse.RoleName = empRole.RoleName
	return &roleResponse, nil
}

func (ul *EmpRoleObj) UpdateEmployeeRole(empRole dtos.EmpRoleUpdate, roleId int64) (*dtos.EmpRoleResponse, error) {
	empRole.RoleCode = strings.ToUpper(empRole.RoleCode)

	if empRole.RoleName == "" {
		ul.l.Error("Error CreateEmployeeRole: Role name should not empty")
		return nil, errors.New("role name should not empty")
	}
	if empRole.RoleCode == "" {
		ul.l.Error("Error CreateEmployeeRole: Role Code should not empty")
		return nil, errors.New("role code should not empty")
	}
	if empRole.OrgId == 0 {
		ul.l.Error("Error CreateEmployeeRole: something wrong, role organization not set properly")
		return nil, errors.New("something went wrong, not created try login again")
	}

	empRole.RoleId = roleId

	roleExisting, errrR := ul.empRoleDao.IsRoleExists(empRole.RoleId)
	if errrR != nil {
		ul.l.Error("role not exists ", empRole, errrR)
		return nil, errrR
	}
	empRole.Version = roleExisting.Version + 1

	err := ul.empRoleDao.UpdateEmpRole(empRole)
	if err != nil {
		ul.l.Error("Role not saved ", empRole, err)
		return nil, err
	}

	errD := ul.empRoleDao.DeleteScreensPermission(roleId)
	if errD != nil {
		ul.l.Error("DeleteScreensPermission not deleted ", roleId, empRole.RoleCode, errD)
		return nil, errD
	}

	for _, screenP := range empRole.ScreensPermission {
		err := ul.empRoleDao.ScreensPermission(empRole.RoleCode, empRole.RoleName, roleId, screenP)
		if err != nil {
			ul.l.Error("ScreensPermission not saved ", empRole, err)
			return nil, err
		}
	}

	roleResponse := dtos.EmpRoleResponse{}
	roleResponse.Message = fmt.Sprintf("Role updated successfully : %s", empRole.RoleName)
	roleResponse.RoleName = empRole.RoleName
	roleResponse.RoleCode = empRole.RoleCode
	return &roleResponse, nil
}

func (vo *EmpRoleObj) UploadEmployeeProfile(employeeId int64, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.UploadEmpResponse, error) {

	bool := CheckSizeImage(fileHeader, 10000, vo.l)
	if !bool {
		vo.l.Error("image size issue ")
		return nil, errors.New("image size issue ")
	}

	emplopleeInfo, errA := vo.empRoleDao.GetEmpployee(employeeId)
	if errA != nil {
		vo.l.Error("ERROR: GetEmployee not found", errA)
		return nil, errA
	}
	empVersion := emplopleeInfo.Version + 1

	// imageTypes := [...]string{
	// 	"rc_expiry_doc",
	// 	"insurance_doc",
	// 	"pucc_expiry_doc",
	// 	"np_expire_doc",
	// 	"fitness_expiry_doc",
	// 	"tax_expiry_doc",
	// 	"mp_expire_doc",
	// 	"pancard_img",
	// 	"bank_passbook_or_cheque_img",
	// }
	// exists := false
	// for _, imgType := range imageTypes {
	// 	if imgType == imageFor {
	// 		exists = true
	// 		break
	// 	}
	// }
	// if !exists {
	// 	return nil, errors.New("imageFor should not be empty. imageFor:rc_expiry_doc/insurance_doc/pucc_expiry_doc/np_expire_doc/fitness_expiry_doc/tax_expiry_doc/mp_expire_doc/pancard_img/bank_passbook_or_cheque_img")
	// }

	baseDirectory := os.Getenv("BASE_DIRECTORY")
	uploadPath := os.Getenv("IMAGE_DIRECTORY")
	if uploadPath == "" || baseDirectory == "" {
		vo.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
	}

	imageDirectory := filepath.Join(uploadPath, "employee", fmt.Sprintf("%d", employeeId))
	fullPath := filepath.Join(baseDirectory, imageDirectory)
	vo.l.Infof("employee: %v imageDirectory: %s, fullPath: %s", employeeId, imageDirectory, fullPath)

	err := os.MkdirAll(fullPath, os.ModePerm) // os.ModePerm sets permissions to 0777
	if err != nil {
		vo.l.Error("ERROR: MkdirAll failed for path: ", fullPath, " error: ", err)
		return nil, fmt.Errorf("failed to create directory %s: %w", fullPath, err)
	}
	
	// Verify directory was created
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		vo.l.Error("ERROR: Directory does not exist after MkdirAll: ", fullPath)
		return nil, fmt.Errorf("directory was not created: %s", fullPath)
	}

	imageFor := "employee"

	extension := strings.Split(fileHeader.Filename, ".")
	lengthExt := len(extension)

	imageName := fmt.Sprintf("%s_%v.%v", imageFor, employeeId, extension[lengthExt-1])
	vendorCodeImageFullPath := filepath.Join(fullPath, imageName)
	vo.l.Info("employee ImageFullPath: ", imageFor, vendorCodeImageFullPath)

	out, err := os.Create(vendorCodeImageFullPath)
	if err != nil {
		if utils.CheckFileExists(vendorCodeImageFullPath) {
			vo.l.Error("updating  is already exitis: ", employeeId, uploadPath, err)
		} else {
			vo.l.Error("employee ImageFullPath create error: ", employeeId, uploadPath, err)
			defer out.Close()
			return nil, err
		}
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		vo.l.Error("vendor upload Copy error: ", employeeId, uploadPath, err)
		return nil, err
	}

	imageDirectory = filepath.Join(imageDirectory, imageName)
	vo.l.Info("##### table to be stored path: ", imageFor, vendorCodeImageFullPath, imageDirectory)
	errU := vo.empRoleDao.UpdateEmployeeImagePath(employeeId, empVersion, imageDirectory)
	if errU != nil {
		vo.l.Error("ERROR: UpdateEmployeeImagePath", errU)
		return nil, errU
	}

	vo.l.Info("Image uploaded successfully: ", imageFor, emplopleeInfo.FirstName)
	roleResponse := dtos.UploadEmpResponse{}
	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v,%v", imageFor, emplopleeInfo.FirstName)
	roleResponse.ImagePath = imageDirectory
	roleResponse.Employee = emplopleeInfo.FirstName

	return &roleResponse, nil
}

func (ul *EmpRoleObj) UploadEmployeeProfile_old(employeeId int64, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.Messge, error) {

	bool := CheckSizeImage(fileHeader, 10000, ul.l)
	if !bool {
		ul.l.Error("image size issue ")
		return nil, errors.New("image size issue ")
	}

	emplopleeInfo, errA := ul.empRoleDao.GetEmpployee(employeeId)
	if errA != nil {
		ul.l.Error("ERROR: GetEmployee not found", errA)
		return nil, errA
	}
	empVersion := emplopleeInfo.Version + 1

	uploadPath := os.Getenv("EMPLOYEE_IMAGE_DIRECTORY")
	if uploadPath == "" {
		uploadPath = "/t_hub_document/employee"
	}
	baseDirectory := os.Getenv("BASE_DIRECTORY")

	extension := strings.Split(fileHeader.Filename, ".")
	lengthExt := len(extension)
	employeeUploadPath := filepath.Join(baseDirectory, uploadPath, fmt.Sprintf("%v.%v", employeeId, extension[lengthExt-1]))
	ul.l.Info("Employee image: ", employeeUploadPath)
	// Save the file to the upload directory

	out, err := os.Create(employeeUploadPath)
	if err != nil {
		if checkFileExists(employeeUploadPath) {
			ul.l.Error("updating  is already exitis: ", employeeId, uploadPath, err)
		} else {
			ul.l.Error("employeeUploadPath create error: ", employeeId, uploadPath, err)
			defer out.Close()
			return nil, err
		}
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ul.l.Error("employeeUploadPath Copy error: ", employeeId, uploadPath, err)
		return nil, err
	}

	errU := ul.empRoleDao.UpdateEmployeeImagePath(employeeId, empVersion, employeeUploadPath)
	if errA != nil {
		ul.l.Error("ERROR: UpdateEmployeeImagePath", errU)
		return nil, errA
	}

	ul.l.Info("Profile image uploaded successfully: ", emplopleeInfo.FirstName)
	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Profile image uploaded successfully : %v", emplopleeInfo.FirstName)
	return &roleResponse, nil
}

func CheckSizeImage(file *multipart.FileHeader, limit int64, ul *log.Logger) bool {
	size := file.Size / 1024
	ul.Info("Employee image size kb - ", size)
	return size <= limit
}

func checkFileExists(filePath string) bool {
	_, err := os.Open(filePath)
	return err == nil
}

func (ul *EmpRoleObj) ViweEmployeeProfile(employeeId int64, image string) ([]byte, error) {

	imageByte, err := os.ReadFile(image)
	if err != nil {
		ul.l.Error("Employee image size kb - ", err)
		return nil, err
	}

	return imageByte, nil
}
