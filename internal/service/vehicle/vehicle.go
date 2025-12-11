package vehicle

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

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
	"go-transport-hub/utils"
)

type VehicleObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
	vehicleDao  daos.VehicleDao
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

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *VehicleObj {
	return &VehicleObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		vehicleDao:  daos.NewVehicleObj(l, dbConnMSSQL),
	}
}

func (vh *VehicleObj) CreateVehicle(vehicleReg dtos.VehicleReq) (*dtos.Messge, error) {

	if vehicleReg.VehicleType == "" {
		vh.l.Error("Error CreateVehicle: first name should not empty")
		return nil, errors.New("vehicle type should not empty")
	}
	if vehicleReg.VehicleNumber == "" {
		vh.l.Error("Error CreateVehicle: vehicle number should not empty")
		return nil, errors.New("vehicle number should not empty")
	}
	vehicleReg.VehicleNumber = strings.ToUpper(vehicleReg.VehicleNumber)

	if vehicleReg.VehicleModel == "" {
		vh.l.Error("Error CreateVehicle: vehicle model should not empty")
		return nil, errors.New("vehicle model should not empty")
	}
	if vehicleReg.VehicleYear == 0 {
		vh.l.Error("Error VehicleYear: vehicle year should not empty")
		return nil, errors.New("vehicle year should not empty")
	}
	if vehicleReg.InsuranceExpiryDate != "" {
		err := utils.ValidateDateStr(vehicleReg.InsuranceExpiryDate)
		if err != nil {
			vh.l.Error("Error InsuranceExpiryDate: invalid date format")
			return nil, errors.New("InsuranceExpiryDate: invalid date format")
		}
	}
	if vehicleReg.VehicleRegistrationDate != "" {
		err := utils.ValidateDateStr(vehicleReg.VehicleRegistrationDate)
		if err != nil {
			vh.l.Error("Error VehicleRegistrationDate: invalid date format")
			return nil, errors.New("VehicleRegistrationDate: invalid date format")
		}
	}
	if vehicleReg.VehicleRenewalDate != "" {
		err := utils.ValidateDateStr(vehicleReg.VehicleRenewalDate)
		if err != nil {
			vh.l.Error("Error VehicleRenewalDate: invalid date format")
			return nil, errors.New("VehicleRenewalDate: invalid date format")
		}
	}
	if vehicleReg.VehicleInsuranceNumber != "" {
		vehicleReg.VehicleInsuranceNumber = strings.ToUpper(vehicleReg.VehicleInsuranceNumber)
	}
	vehicleReg.Status = constant.STATUS_CREATED

	err1 := vh.vehicleDao.CreateVehicle(vehicleReg)
	if err1 != nil {
		vh.l.Error("vehicle not saved ", vehicleReg.VehicleNumber, err1)
		return nil, err1
	}

	vh.l.Info("Vehicle created successfully! : ", vehicleReg.VehicleNumber)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vehicle saved successfully: %s", vehicleReg.VehicleNumber)
	return &roleResponse, nil
}

func (vh *VehicleObj) GetVehicles(orgId int64, limit string, offset string) (*dtos.VehicleEntries, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	res, errA := vh.vehicleDao.GetVehicles(orgId, limitI, offsetI)
	if errA != nil {
		vh.l.Error("ERROR: GetVehicles", errA)
		return nil, errA
	}

	vehicleEntries := dtos.VehicleEntries{}
	vehicleEntries.VendorEntiry = res
	vehicleEntries.Total = vh.vehicleDao.GetTotalCount(orgId)
	vehicleEntries.Limit = limitI
	vehicleEntries.OffSet = offsetI

	return &vehicleEntries, nil
}

func (vh *VehicleObj) UpdateVehicle(vehicleId int64, vehicleReg dtos.VehicleUpdate) (*dtos.Messge, error) {

	if vehicleReg.VehicleType == "" {
		vh.l.Error("Error CreateVehicle: first name should not empty")
		return nil, errors.New("vehicle type should not empty")
	}
	if vehicleReg.VehicleNumber == "" {
		vh.l.Error("Error CreateVehicle: vehicle number should not empty")
		return nil, errors.New("vehicle number should not empty")
	}
	vehicleReg.VehicleNumber = strings.ToUpper(vehicleReg.VehicleNumber)

	if vehicleReg.VehicleModel == "" {
		vh.l.Error("Error CreateVehicle: vehicle model should not empty")
		return nil, errors.New("vehicle model should not empty")
	}
	if vehicleReg.VehicleYear == 0 {
		vh.l.Error("Error VehicleYear: vehicle year should not empty")
		return nil, errors.New("vehicle year should not empty")
	}
	if vehicleReg.InsuranceExpiryDate != "" {
		err := utils.ValidateDateStr(vehicleReg.InsuranceExpiryDate)
		if err != nil {
			vh.l.Error("Error InsuranceExpiryDate: invalid date format")
			return nil, errors.New("InsuranceExpiryDate: invalid date format")
		}
	}
	if vehicleReg.VehicleRegistrationDate != "" {
		err := utils.ValidateDateStr(vehicleReg.VehicleRegistrationDate)
		if err != nil {
			vh.l.Error("Error VehicleRegistrationDate: invalid date format")
			return nil, errors.New("VehicleRegistrationDate: invalid date format")
		}
	}
	if vehicleReg.VehicleRenewalDate != "" {
		err := utils.ValidateDateStr(vehicleReg.VehicleRenewalDate)
		if err != nil {
			vh.l.Error("Error VehicleRenewalDate: invalid date format")
			return nil, errors.New("VehicleRenewalDate: invalid date format")
		}
	}
	if vehicleReg.VehicleInsuranceNumber != "" {
		vehicleReg.VehicleInsuranceNumber = strings.ToUpper(vehicleReg.VehicleInsuranceNumber)
	}

	vehicleInfo, errA := vh.vehicleDao.GetVehicle(vehicleId)
	if errA != nil {
		vh.l.Error("ERROR: Vehicle not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(vehicleInfo)
	vh.l.Info("vehicleInfo: ", string(jsonBytes))
	vehicleReg.Status = vehicleInfo.Status

	err1 := vh.vehicleDao.UpdateVehicle(vehicleId, vehicleReg)
	if err1 != nil {
		vh.l.Error("vehicle not updated ", vehicleReg.VehicleNumber, err1)
		return nil, err1
	}

	vh.l.Info("Vehicle updated successfully! : ", vehicleReg.VehicleNumber)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vehicle saved successfully: %s", vehicleReg.VehicleNumber)
	return &roleResponse, nil
}

func (ul *VehicleObj) UpdateVehicleActiveStatus(isA string, vehicleID int64) (*dtos.Messge, error) {

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid status update")
	}
	if vehicleID == 0 {
		ul.l.Error("unknown vehicle", vehicleID)
		return nil, errors.New("unknown vehicle")
	}
	statusRes := dtos.Messge{}
	vehicleInfo, errA := ul.vehicleDao.GetVehicle(vehicleID)
	if errA != nil {
		ul.l.Error("ERROR: Vehicle not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(vehicleInfo)
	ul.l.Info("vehicleInfo: ", string(jsonBytes))

	errU := ul.vehicleDao.UpdateVehicleActiveStatus(vehicleID, isActive)
	if errU != nil {
		ul.l.Error("ERROR: UpdateVehicleActiveStatus ", errU)
		return nil, errU
	}
	statusRes.Message = "vehicle updated successfully : " + vehicleInfo.VehicleNumber
	return &statusRes, nil
}

func (ul *VehicleObj) UploadVehicleImages(vehicleId int64, imageFor string, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.Messge, error) {

	if imageFor == "" {
		return nil, errors.New("imageFor should not be empty.\nfitness_certificate/insurance_certificate/pollution_certificate/national_permits_certificate/registration_certificate/annual_maintenance_certificate")
	}

	imageTypes := [...]string{
		"fitness_certificate",
		"insurance_certificate",
		"pollution_certificate",
		"national_permits_certificate",
		"registration_certificate",
		"annual_maintenance_certificate",
	}

	exists := false
	for _, imgType := range imageTypes {
		if imgType == imageFor {
			exists = true
			break
		}
	}
	if !exists {
		return nil, errors.New("imageFor not valid fitness_certificate/insurance_certificate/pollution_certificate/national_permits_certificate/registration_certificate/annual_maintenance_certificate")
	}

	bool := utils.CheckSizeImage(fileHeader, 10000, ul.l)
	if !bool {
		ul.l.Error("image size issue ")
		return nil, errors.New("image size issue ")
	}

	vehicleInfo, errA := ul.vehicleDao.GetVehicle(vehicleId)
	if errA != nil {
		ul.l.Error("error: vehicle not found", errA)
		return nil, errA
	}

	baseDirectory := os.Getenv("BASE_DIRECTORY")
	uploadPath := os.Getenv("IMAGE_DIRECTORY")
	if uploadPath == "" || baseDirectory == "" {
		ul.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
	}

	imageDirectory := filepath.Join(uploadPath, "vehicle", strconv.Itoa(int(vehicleId)))
	fullPath := filepath.Join(baseDirectory, imageDirectory)
	ul.l.Infof("vehicle: %v imageDirectory: %s, fullPath: %s", vehicleId, imageDirectory, fullPath)

	err := os.MkdirAll(fullPath, os.ModePerm) // os.ModePerm sets permissions to 0777
	if err != nil {
		ul.l.Error("ERROR: MkdirAll failed for path: ", fullPath, " error: ", err)
		return nil, fmt.Errorf("failed to create directory %s: %w", fullPath, err)
	}

	// Verify directory was created
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		ul.l.Error("ERROR: Directory does not exist after MkdirAll: ", fullPath)
		return nil, fmt.Errorf("directory was not created: %s", fullPath)
	}

	extension := strings.Split(fileHeader.Filename, ".")
	lengthExt := len(extension)

	imageName := fmt.Sprintf("%v_%v.%v", imageFor, vehicleId, extension[lengthExt-1])
	vehicleImageFullPath := filepath.Join(fullPath, imageName)
	ul.l.Info("vehicleImageFullPath: ", imageFor, vehicleImageFullPath)

	out, err := os.Create(vehicleImageFullPath)
	if err != nil {
		if utils.CheckFileExists(vehicleImageFullPath) {
			ul.l.Error("updating  is already exitis: ", vehicleId, uploadPath, err)
		} else {
			ul.l.Error("vehicleImagePath create error: ", vehicleId, uploadPath, err)
			defer out.Close()
			return nil, err
		}
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ul.l.Error("vehicle upload Copy error: ", vehicleId, uploadPath, err)
		return nil, err
	}

	imageDirectory = filepath.Join(imageDirectory, imageName)
	ul.l.Info("##### table to be stored path: ", imageFor, imageDirectory)
	updateQuery := fmt.Sprintf(`UPDATE vehicle SET %v = '%v' WHERE vehicle_id = '%v'`, imageFor, imageDirectory, vehicleId)
	errU := ul.vehicleDao.UpdateVehicleImagePath(updateQuery)
	if errU != nil {
		ul.l.Error("ERROR: UpdateVehicleImagePath", vehicleId, errU)
		return nil, errU
	}

	ul.l.Info("Image uploaded successfully: ", vehicleInfo.VehicleNumber)
	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", vehicleInfo.VehicleNumber)
	return &roleResponse, nil
}

func (vh *VehicleObj) CreateVehicleSizeTypes(vehicleObj dtos.VehicleSizeType) (*dtos.Messge, error) {

	if vehicleObj.VehicleSize == "" {
		vh.l.Error("Error VehicleSize: vehicle size not empty")
		return nil, errors.New("vehicle size not empty")
	}

	if vehicleObj.VehicleSize == "" {
		vh.l.Error("Error VehicleSize: vehicle size not empty")
		return nil, errors.New("vehicle size not empty")
	}

	vehicleObj.Status = constant.STATUS_CREATED

	err1 := vh.vehicleDao.CreateVehicleSizeTypes(vehicleObj)
	if err1 != nil {
		vh.l.Error("CreateVehicleSizeTypes not saved ", vehicleObj.VehicleSize, err1)
		return nil, err1
	}

	// Note: Vehicle size notifications skipped as orgId is not available in VehicleSizeType DTO
	// Vehicle size types appear to be global, not org-specific

	vh.l.Info("Vehicle size type created successfully! : ", vehicleObj.VehicleSize, vehicleObj.VehicleType)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vehicle size type created successfully!: %s %s", vehicleObj.VehicleSize, vehicleObj.VehicleType)
	return &roleResponse, nil
}

func (vh *VehicleObj) GetVehicleSizeTypes(limit string, offset string) (*dtos.VehicleSizeTypeEntries, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	res, errA := vh.vehicleDao.GetVehicleSizeTypes(limitI, offsetI)
	if errA != nil {
		vh.l.Error("ERROR: GetVehicleSizeTypes", errA)
		return nil, errA
	}

	vehicleTypeEntries := dtos.VehicleSizeTypeEntries{}
	vehicleTypeEntries.VehicleSizeType = res
	vehicleTypeEntries.Total = vh.vehicleDao.VehicleSizeTypesGetTotalCount()
	vehicleTypeEntries.Limit = limitI
	vehicleTypeEntries.OffSet = offsetI

	return &vehicleTypeEntries, nil
}

func (vh *VehicleObj) UpdateVehicleSizeType(vehicleSizeId int64, vehicleObj dtos.VehicleSizeTypeUpdate) (*dtos.Messge, error) {

	if vehicleObj.VehicleSize == "" {
		vh.l.Error("Error VehicleSize: vehicle size not empty")
		return nil, errors.New("vehicle size not empty")
	}

	if vehicleObj.VehicleSize == "" {
		vh.l.Error("Error VehicleSize: vehicle size not empty")
		return nil, errors.New("vehicle size not empty")
	}

	vehicleSizeInfo, errA := vh.vehicleDao.GetVehicleSizeType(vehicleSizeId)
	if errA != nil {
		vh.l.Error("ERROR: Vehicle not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(vehicleSizeInfo)
	vh.l.Info("vehicleSizeInfo: ", string(jsonBytes))

	vehicleObj.Status = constant.STATUS_CREATED

	err1 := vh.vehicleDao.UpdateVehicleSizeType(vehicleSizeId, vehicleObj)
	if err1 != nil {
		vh.l.Error("UpdateVehicleSizeType not updated ", vehicleObj.VehicleSize, err1)
		return nil, err1
	}

	// Note: Vehicle size notifications skipped as orgId is not available
	// Vehicle size types appear to be global, not org-specific

	vh.l.Info("Vehicle size type updated successfully! : ", vehicleSizeId, vehicleObj.VehicleSize)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vehicle updated successfully: %s %s", vehicleObj.VehicleSize, vehicleObj.VehicleType)
	return &roleResponse, nil
}

func (ul *VehicleObj) UpdateVehicleSizeTypeActiveStatus(isA string, vehicleSizeId int64) (*dtos.Messge, error) {

	if isA != "0" && isA != "1" {
		ul.l.Error("unknown status", isA)
		return nil, errors.New("unknown status")
	}

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid status update")
	}
	if vehicleSizeId == 0 {
		ul.l.Error("unknown vehicle", vehicleSizeId)
		return nil, errors.New("unknown vehicle")
	}

	vehicleInfo, errA := ul.vehicleDao.GetVehicleSizeType(vehicleSizeId)
	if errA != nil {
		ul.l.Error("ERROR: Vehicle not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(vehicleInfo)
	ul.l.Info("vehicleInfo: ", string(jsonBytes))

	errU := ul.vehicleDao.UpdateVehicleSizeTypeActiveStatus(vehicleSizeId, isActive)
	if errU != nil {
		ul.l.Error("ERROR: UpdateVehicleSizeTypeActiveStatus ", errU)
		return nil, errU
	}

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vehicle updated successfully: %s %s", vehicleInfo.VehicleSize, vehicleInfo.VehicleType)
	return &roleResponse, nil
}
