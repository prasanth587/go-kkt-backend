package driver

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
	"go-transport-hub/internal/daos"
	"go-transport-hub/utils"
)

type DriverObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
	driverDao   daos.DriverDao
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

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *DriverObj {
	return &DriverObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		driverDao:   daos.NewDriverObj(l, dbConnMSSQL),
	}
}

func (dr *DriverObj) CreateDriver(driverReg dtos.CreateDriverReq) (*dtos.Messge, error) {

	if driverReg.FirstName == "" {
		dr.l.Error("Error CreateEmployee: first name should not empty")
		return nil, errors.New("first name should not empty")
	}
	if driverReg.LicenseNumber == "" {
		dr.l.Error("Error CreateDriver: License Number code should not empty")
		return nil, errors.New("license number should not empty")
	}
	if driverReg.LicenseExpiryDate == "" {
		dr.l.Error("Error LicenseExpiryDate: license expiry date")
		return nil, errors.New("license expiry date should not empty")
	}
	err := utils.ValidateLicenseExpiryDate(driverReg.LicenseExpiryDate)
	if err != nil {
		dr.l.Error("Error LicenseExpiryDate: license expiry date must be a future date")
		return nil, err
	}

	if driverReg.MobileNo == "" {
		if mobileNumberLength := len(driverReg.MobileNo); mobileNumberLength < 10 || mobileNumberLength > 15 {
			return nil, errors.New("mobile mumber should be length of 10 to 15 digits")
		}
		dr.l.Error("Error CreateDriver: mobile numger should not empty")
		return nil, errors.New("mobile numger should not empty")
	}
	errE := utils.ValidateDriverExperience(driverReg.DriverExperience)
	if errE != nil {
		dr.l.Error("Error driver experience: ", errE)
		return nil, errE
	}

	if driverReg.JoiningDate != "" {
		errJ := utils.ValidateDateStr(driverReg.JoiningDate)
		if errJ != nil {
			dr.l.Error("Error driver joining: ", errJ)
			return nil, errJ
		}
	}

	driverId, err := dr.driverDao.CreateDriver(driverReg)
	if err != nil {
		dr.l.Error("Driver not saved ", driverReg.FirstName, err)
		return nil, err
	}

	dr.l.Info("Driver created successfully! : ", driverReg.FirstName, driverId)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Driver saved successfully: %s", driverReg.FirstName)

	return &roleResponse, nil
}

func (ul *DriverObj) GetDrivers(orgId int64, limit string, offset string) (*dtos.DriverEntries, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	res, errA := ul.driverDao.GetDrivers(orgId, limitI, offsetI)
	if errA != nil {
		ul.l.Error("ERROR: GetDrivers", errA)
		return nil, errA
	}

	driverEntries := dtos.DriverEntries{}
	driverEntries.DriverEntries = res
	driverEntries.Total = ul.driverDao.GetTotalCount(orgId)
	driverEntries.Limit = limitI
	driverEntries.OffSet = offsetI
	return &driverEntries, nil
}

func (ul *DriverObj) UpdateDriverctiveStatus(isA string, driverId int64) (*dtos.Messge, error) {

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid status update")
	}
	if driverId == 0 {
		ul.l.Error("unknown driver", driverId)
		return nil, errors.New("unknown driver")
	}
	statusRes := dtos.Messge{}
	driverInfo, errA := ul.driverDao.GetDriver(driverId)
	if errA != nil {
		ul.l.Error("ERROR: GetDriver not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(driverInfo)
	ul.l.Info("driverInfo: ", string(jsonBytes))

	errU := ul.driverDao.UpdateDriverActiveStatus(driverId, isActive)
	if errU != nil {
		ul.l.Error("ERROR: UpdateDriverActiveStatus ", errU)
		return nil, errU
	}
	statusRes.Message = "driver updated successfully : " + driverInfo.FirstName
	return &statusRes, nil
}

func (dr *DriverObj) UpdateDriver(driverID int64, driverReg dtos.UpdateDriverReq) (*dtos.Messge, error) {

	if driverReg.FirstName == "" {
		dr.l.Error("Error CreateEmployee: first name should not empty")
		return nil, errors.New("first name should not empty")
	}
	if driverReg.LicenseNumber == "" {
		dr.l.Error("Error CreateDriver: License Number code should not empty")
		return nil, errors.New("license number should not empty")
	}
	if driverReg.LicenseExpiryDate == "" {
		dr.l.Error("Error LicenseExpiryDate: license expiry date")
		return nil, errors.New("license expiry date should not empty")
	}
	err := utils.ValidateLicenseExpiryDate(driverReg.LicenseExpiryDate)
	if err != nil {
		dr.l.Error("Error LicenseExpiryDate: license expiry date must be a future date")
		return nil, err
	}

	if driverReg.MobileNo == "" {
		if mobileNumberLength := len(driverReg.MobileNo); mobileNumberLength < 10 || mobileNumberLength > 15 {
			return nil, errors.New("mobile mumber should be length of 10 to 15 digits")
		}
		dr.l.Error("Error CreateDriver: mobile numger should not empty")
		return nil, errors.New("mobile numger should not empty")
	}
	errE := utils.ValidateDriverExperience(driverReg.DriverExperience)
	if errE != nil {
		dr.l.Error("Error driver experience: ", errE)
		return nil, errE
	}

	if driverReg.JoiningDate != "" {
		errJ := utils.ValidateDateStr(driverReg.JoiningDate)
		if errJ != nil {
			dr.l.Error("Error driver joining: ", errJ)
			return nil, errJ
		}
	}

	driverInfo, errA := dr.driverDao.GetDriver(driverID)
	if errA != nil {
		dr.l.Error("ERROR: GetDriver not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(driverInfo)
	dr.l.Info("driverInfo: ", string(jsonBytes))

	errS := dr.driverDao.UpdateDriver(driverID, driverReg)
	if errS != nil {
		dr.l.Error("Driver not updated ", driverReg.FirstName, errS)
		return nil, errS
	}

	dr.l.Info("Driver updated successfully! : ", driverReg.FirstName)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Driver saved successfully: %s", driverReg.FirstName)

	return &roleResponse, nil
}

func (ul *DriverObj) UploadDriverImages(driverId int64, imageFor string, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.Messge, error) {

	if imageFor == "" {
		return nil, errors.New("imageFor should not be empty image/license_front_img/license_back_img")
	}

	imageTypes := [...]string{
		"image",
		"license_front_img",
		"license_back_img",
		"other_document",
	}
	exists := false
	for _, imgType := range imageTypes {
		if imgType == imageFor {
			exists = true
			break
		}
	}
	if !exists {
		return nil, errors.New("imageFor not valid image/license_front_img/license_back_img")
	}

	bool := utils.CheckSizeImage(fileHeader, 10000, ul.l)
	if !bool {
		ul.l.Error("image size issue ")
		return nil, errors.New("image size issue ")
	}

	driverInfo, errA := ul.driverDao.GetDriver(driverId)
	if errA != nil {
		ul.l.Error("ERROR: GetDriver not found", errA)
		return nil, errA
	}

	baseDirectory := os.Getenv("BASE_DIRECTORY")
	uploadPath := os.Getenv("IMAGE_DIRECTORY")
	if uploadPath == "" || baseDirectory == "" {
		ul.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
	}

	imageDirectory := filepath.Join(uploadPath, "driver", strconv.Itoa(int(driverId)))
	fullPath := filepath.Join(baseDirectory, imageDirectory)
	ul.l.Infof("customerId: %v imageDirectory: %s, fullPath: %s", driverId, imageDirectory, fullPath)

	err := os.MkdirAll(fullPath, os.ModePerm) // os.ModePerm sets permissions to 0777
	if err != nil {
		ul.l.Error("ERROR: MkdirAll ", fullPath, err)
		return nil, err
	}

	extension := strings.Split(fileHeader.Filename, ".")
	lengthExt := len(extension)

	imageName := fmt.Sprintf("%v_%v.%v", imageFor, driverId, extension[lengthExt-1])
	driverImageFullPath := filepath.Join(fullPath, imageName)
	ul.l.Info("driverImageFullPath: ", imageFor, driverImageFullPath)

	out, err := os.Create(driverImageFullPath)
	if err != nil {
		if utils.CheckFileExists(driverImageFullPath) {
			ul.l.Error("updating  is already exitis: ", driverId, uploadPath, err)
		} else {
			ul.l.Error("driverImageFullPath create error: ", driverId, uploadPath, err)
			defer out.Close()
			return nil, err
		}
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ul.l.Error("driver upload Copy error: ", driverId, uploadPath, err)
		return nil, err
	}

	imageDirectory = filepath.Join(imageDirectory, imageName)
	ul.l.Info("##### table to be stored path: ", imageFor, imageDirectory)
	updateQuery := fmt.Sprintf(`UPDATE driver SET %v = '%v' WHERE driver_id = '%v'`, imageFor, imageDirectory, driverId)
	errU := ul.driverDao.UpdateDriverImagePath(updateQuery)
	if errU != nil {
		ul.l.Error("ERROR: UpdateVehicleImagePath", driverId, errU)
		return nil, errU
	}

	ul.l.Info("Image uploaded successfully: ", driverInfo.FirstName)
	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", driverInfo.FirstName)
	return &roleResponse, nil
}
