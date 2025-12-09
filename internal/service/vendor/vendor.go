package vendor

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

type VendorObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
	vendorDao   daos.VendorDao
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

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *VendorObj {
	return &VendorObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		vendorDao:   daos.NewVendorObj(l, dbConnMSSQL),
	}
}

func (vh *VendorObj) CreateVendor(vendorReg dtos.VendorReq) (*dtos.Messge, error) {

	if vendorReg.VendorName == "" {
		vh.l.Error("Error VendorName: vendor name should not empty")
		return nil, errors.New("vendor name should not empty")
	}
	if vendorReg.VendorCode == "" {
		vh.l.Error("Error VendorCode: vendor code should not empty")
		return nil, errors.New("vendor code should not empty")
	}
	vendorReg.VendorCode = strings.ToUpper(vendorReg.VendorCode)

	if vendorReg.MobileNumber == "" {
		vh.l.Error("Error MobileNumber: Mobile number should not empty")
		return nil, errors.New("mobile number should not empty")
	}
	if vendorReg.ContactPerson == "" {
		vh.l.Error("Error ContactPerson: contact person should not empty")
		return nil, errors.New("contact person should not empty")
	}
	if vendorReg.City == "" {
		vh.l.Error("Error City: city should not empty")
		return nil, errors.New("city should not empty")
	}
	vendorReg.Status = constant.STATUS_CREATED

	err1 := vh.vendorDao.CreateVendor(vendorReg)
	if err1 != nil {
		vh.l.Error("vendor not saved ", vendorReg.VendorCode, err1)
		return nil, err1
	}

	vh.l.Info("Vendor created successfully! : ", vendorReg.VendorCode)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vendor saved successfully: %s", vendorReg.VendorCode)
	return &roleResponse, nil
}

func (vh *VendorObj) CreateVendorV1(vendorReg dtos.VendorRequest) (*dtos.Messge, error) {

	err := vh.validateVendor(vendorReg)
	if err != nil {
		vh.l.Error("ERROR: CreateVendorV1", err)
		return nil, err
	}
	//trp.l.Info("trip sheet type: ", tripSheetReq.TripSheetType)

	vendorReg.VendorCode = strings.ToUpper(vendorReg.VendorCode)

	vendorReg.Status = constant.STATUS_CREATED

	vendorId, err1 := vh.vendorDao.CreateVendorV1(vendorReg)
	if err1 != nil {
		vh.l.Error("vendor not saved ", vendorReg.VendorCode, err1)
		return nil, err1
	}

	if len(vendorReg.ContactInfo) != 0 {
		for _, contactInfo := range vendorReg.ContactInfo {
			err1 := vh.vendorDao.CreateVendorContactInfo(vendorId, contactInfo)
			if err1 != nil {
				vh.l.Error("CreateVendorContactInfo not saved ", contactInfo.ContactPersonName, err1)
				return nil, err1
			}
		}
	}

	if len(vendorReg.Vehicles) != 0 {
		for _, vehicleObj := range vendorReg.Vehicles {
			vehicleObj.VendorID = vendorId
			err1 := vh.vendorDao.CreateVehicle(vehicleObj)
			if err1 != nil {
				vh.l.Error("CreateVendorContactInfo not saved ", vehicleObj.VehicleNumber, err1)
				return nil, err1
			}
		}
	}

	if len(vendorReg.DeclarationDocument) != 0 {
		for _, declarationDoc := range vendorReg.DeclarationDocument {
			declarationDoc.VendorID = vendorId
			err1 := vh.vendorDao.CreateVehicleDoclarationDoc(declarationDoc)
			if err1 != nil {
				vh.l.Error("DeclarationDocument not saved ", declarationDoc.DeclarationYear, err1)
				return nil, err1
			}
		}
	}

	vh.l.Info("Vendor created successfully! : ", vendorReg.VendorCode)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vendor saved successfully: %s", vendorReg.VendorCode)
	return &roleResponse, nil
}

func (vh *VendorObj) GetVendors(orgId int64, vendorId string, limit string, offset, searchText string) (*dtos.VendorEntries, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	whereQuery := vh.vendorDao.BuildWhereQuery(orgId, vendorId, searchText)

	res, errA := vh.vendorDao.GetVendors(orgId, whereQuery, limitI, offsetI)
	if errA != nil {
		vh.l.Error("ERROR: GetVendors", errA)
		return nil, errA
	}

	vehicleEntries := dtos.VendorEntries{}
	vehicleEntries.VendorEntiry = res
	vehicleEntries.Total = vh.vendorDao.GetTotalCount(whereQuery)
	vehicleEntries.Limit = limitI
	vehicleEntries.OffSet = offsetI
	return &vehicleEntries, nil
}
func (vh *VendorObj) GetVendorsV1(orgId int64, vendorId string, limit string, offset, searchText string) (*dtos.VendorV1Entries, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	whereQuery := vh.vendorDao.BuildWhereQuery(orgId, vendorId, searchText)

	res, errA := vh.vendorDao.GetVendorsV1(orgId, whereQuery, limitI, offsetI)
	if errA != nil {
		vh.l.Error("ERROR: GetVendors", errA)
		return nil, errA
	}

	for i := range *res {
		contactInfo, errA := vh.vendorDao.GetVendorContactInfo((*res)[i].VendorId)
		if errA != nil {
			vh.l.Error("ERROR: GetVendorContactInfo", errA)
			return nil, errA
		}
		(*res)[i].ContactInfo = *contactInfo

		// Getting vehicles
		vehicles, errA := vh.vendorDao.GetVehiclesByVendorID((*res)[i].VendorId)
		if errA != nil {
			vh.l.Error("ERROR: GetVehiclesByVendorID", errA)
			return nil, errA
		}
		(*res)[i].Vehicles = *vehicles

		//
		doclarations, errA := vh.vendorDao.GetDeclarationDocByVendorID((*res)[i].VendorId)
		if errA != nil {
			vh.l.Error("ERROR: GetDeclarationDocByVendorID", errA)
			return nil, errA
		}
		(*res)[i].DeclarationDocuments = *doclarations
	}

	vehicleEntries := dtos.VendorV1Entries{}
	vehicleEntries.VendorEntiry = res
	vehicleEntries.Total = vh.vendorDao.GetTotalCountV1(whereQuery)
	vehicleEntries.Limit = limitI
	vehicleEntries.OffSet = offsetI
	return &vehicleEntries, nil
}

func (vh *VendorObj) UpdateVendor(vendorId int64, vendorReg dtos.VendorUpdate) (*dtos.Messge, error) {

	if vendorReg.VendorName == "" {
		vh.l.Error("Error VendorName: vendor name should not empty")
		return nil, errors.New("vendor name should not empty")
	}
	if vendorReg.VendorCode == "" {
		vh.l.Error("Error VendorCode: vendor code should not empty")
		return nil, errors.New("vendor code should not empty")
	}
	vendorReg.VendorCode = strings.ToUpper(vendorReg.VendorCode)

	if vendorReg.MobileNumber == "" {
		vh.l.Error("Error MobileNumber: Mobile number should not empty")
		return nil, errors.New("mobile number should not empty")
	}
	if vendorReg.ContactPerson == "" {
		vh.l.Error("Error ContactPerson: contact person should not empty")
		return nil, errors.New("contact person should not empty")
	}
	if vendorReg.City == "" {
		vh.l.Error("Error City: city should not empty")
		return nil, errors.New("city should not empty")
	}

	vendorInfo, errV := vh.vendorDao.GetVendor(vendorId)
	if errV != nil {
		vh.l.Error("ERROR: GetEmployee not found", errV)
		return nil, errV
	}
	jsonBytes, _ := json.Marshal(vendorInfo)
	vh.l.Info("GetVendor: ******* ", string(jsonBytes))

	vendorReg.Status = vendorInfo.Status

	err1 := vh.vendorDao.UpdateVendor(vendorId, vendorReg)
	if err1 != nil {
		vh.l.Error("vendor not updated ", vendorReg.VendorName, err1)
		return nil, err1
	}

	vh.l.Info("Vendor updated successfully! : ", vendorReg.VendorName)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vendor saved successfully: %s", vendorReg.VendorName)
	return &roleResponse, nil
}

func (vh *VendorObj) UpdateVendorV1(vendorId int64, vendorReg dtos.VendorV1Update) (*dtos.Messge, error) {

	if vendorReg.VendorName == "" {
		vh.l.Error("Error VendorName: vendor name should not empty")
		return nil, errors.New("vendor name should not empty")
	}
	if vendorReg.VendorCode == "" {
		vh.l.Error("Error VendorCode: vendor code should not empty")
		return nil, errors.New("vendor code should not empty")
	}
	vendorReg.VendorCode = strings.ToUpper(vendorReg.VendorCode)

	vendorInfo, errV := vh.vendorDao.GetVendorV1(vendorId)
	if errV != nil {
		vh.l.Error("ERROR: GetEmployee not found", errV)
		return nil, errV
	}
	jsonBytes, _ := json.Marshal(vendorInfo)
	vh.l.Info("GetVendor: ******* ", string(jsonBytes))

	vendorReg.Status = vendorInfo.Status

	err1 := vh.vendorDao.UpdateVendorV1(vendorId, vendorReg)
	if err1 != nil {
		vh.l.Error("vendor not updated ", vendorReg.VendorName, err1)
		return nil, err1
	}

	// Deleting existing contact info for the vendor
	errD := vh.vendorDao.DeleteVendorContactInfo(vendorId)
	if errD != nil {
		vh.l.Error("DeleteVendorContactInfo not deleted ", vendorInfo.VendorName, vendorId, errD)
		return nil, errD
	}

	if len(vendorReg.ContactInfo) != 0 {
		for _, contactInfo := range vendorReg.ContactInfo {
			err1 := vh.vendorDao.CreateVendorContactInfo(vendorId, contactInfo)
			if err1 != nil {
				vh.l.Error("CreateVendorContactInfo not saved ", contactInfo.ContactPersonName, err1)
				return nil, err1
			}
		}
	}

	// Deleting existing vehicles for the vendor
	errM := vh.vendorDao.DeleteVendorVehiclesInfo(vendorId)
	if errM != nil {
		vh.l.Error("DeleteVendorContactInfo not deleted ", vendorInfo.VendorName, vendorId, errM)
		return nil, errM
	}

	if len(vendorReg.Vehicles) != 0 {
		for _, vehicleObj := range vendorReg.Vehicles {
			vehicleObj.VendorID = vendorId
			err1 := vh.vendorDao.CreateVehicle(vehicleObj)
			if err1 != nil {
				vh.l.Error("CreateVehicle not saved ", vehicleObj.VehicleNumber, err1)
				//return nil, err1
				continue
			}
		}
	}

	// Deleting existing vehicles for the vendor
	errDec := vh.vendorDao.DeleteVendorDeclarations(vendorId)
	if errDec != nil {
		vh.l.Error("DeleteVendorDeclarations not deleted ", vendorInfo.VendorName, vendorId, errDec)
		return nil, errDec
	}
	if len(vendorReg.DeclarationDocument) != 0 {
		for _, declarationDoc := range vendorReg.DeclarationDocument {
			declarationDoc.VendorID = vendorId
			err1 := vh.vendorDao.CreateVehicleDoclarationDoc(declarationDoc)
			if err1 != nil {
				vh.l.Error("CreateVehicleDoclarationDoc not saved ", declarationDoc.DeclarationYear, err1)
				return nil, err1
			}
		}
	}

	vh.l.Info("Vendor updated successfully! : ", vendorReg.VendorName)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Vendor saved successfully: %s", vendorReg.VendorName)
	return &roleResponse, nil
}

func (ul *VendorObj) UpdateVendorActiveStatus(isA, status string, vendorID int64) (*dtos.Messge, error) {

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil && status == "" {
		return nil, errors.New("invalid status update")
	}
	if vendorID == 0 {
		ul.l.Error("unknown vendor", vendorID)
		return nil, errors.New("unknown vendor")
	}
	statusRes := dtos.Messge{}
	vendorInfo, errA := ul.vendorDao.GetVendor(vendorID)
	if errA != nil {
		ul.l.Error("ERROR: Vendor not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(vendorInfo)
	ul.l.Info("vendorInfo: ", status, isActive, string(jsonBytes))

	var errU error
	if status != "" {
		errU = ul.vendorDao.UpdateVendorStatus(vendorID, status)
	} else {
		errU = ul.vendorDao.UpdateVendorActiveStatus(vendorID, isActive)
	}
	if errU != nil {
		ul.l.Error("ERROR: UpdateVendor Active/Status ", errU)
		return nil, errU
	}
	statusRes.Message = "vendor updated successfully : " + vendorInfo.VendorName
	return &statusRes, nil
}

func (ul *VendorObj) UpdateVendorActiveStatusV1(isA, status string, vendorID int64) (*dtos.Messge, error) {

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil && status == "" {
		return nil, errors.New("invalid status update")
	}
	if vendorID == 0 {
		ul.l.Error("unknown vendor", vendorID)
		return nil, errors.New("unknown vendor")
	}
	statusRes := dtos.Messge{}
	vendorInfo, errA := ul.vendorDao.GetVendorV1(vendorID)
	if errA != nil {
		ul.l.Error("ERROR: Vendor not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(vendorInfo)
	ul.l.Info("vendorInfo: ", status, isActive, string(jsonBytes))

	var errU error
	if status != "" {
		errU = ul.vendorDao.UpdateVendorStatusV1(vendorID, status)
	} else {
		errU = ul.vendorDao.UpdateVendorActiveStatusV1(vendorID, isActive)
	}
	if errU != nil {
		ul.l.Error("ERROR: UpdateVendor Active/Status ", errU)
		return nil, errU
	}
	statusRes.Message = "vendor updated successfully : " + vendorInfo.VendorName
	return &statusRes, nil
}

func (ul *VendorObj) UploadVendorImages(vendorId int64, imageFor string, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.Messge, error) {

	if imageFor == "" {
		return nil, errors.New("imageFor should not be empty visiting_card_image/pancard_img/aadhar_card_img/cancelled_check_book_img/bank_passbook_img/gst_document_img")
	}

	imageTypes := [...]string{
		"visiting_card_image",
		"pancard_img",
		"aadhar_card_img",
		"cancelled_check_book_img",
		"bank_passbook_img",
		"gst_document_img",
	}
	exists := false
	for _, imgType := range imageTypes {
		if imgType == imageFor {
			exists = true
			break
		}
	}
	if !exists {
		return nil, errors.New("imageFor not valid visiting_card_image/pancard_img/aadhar_card_img/cancelled_check_book_img/bank_passbook_img/gst_document_img")
	}

	vendorInfo, errV := ul.vendorDao.GetVendor(vendorId)
	if errV != nil {
		ul.l.Error("ERROR: vendor not found", errV)
		return nil, errV
	}
	//jsonBytes, _ := json.Marshal(vendorInfo)
	//ul.l.Info("GetVendor: ******* ", string(jsonBytes))

	baseDirectory := os.Getenv("BASE_DIRECTORY")
	uploadPath := os.Getenv("IMAGE_DIRECTORY")
	if uploadPath == "" || baseDirectory == "" {
		ul.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
	}

	imageDirectory := filepath.Join(uploadPath, "vendor", strconv.Itoa(int(vendorId)))
	fullPath := filepath.Join(baseDirectory, imageDirectory)
	ul.l.Infof("vendorId: %v imageDirectory: %s, fullPath: %s", vendorId, imageDirectory, fullPath)

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

	imageName := fmt.Sprintf("%v_%v.%v", imageFor, vendorId, extension[lengthExt-1])
	vendorImageFullPath := filepath.Join(fullPath, imageName)
	ul.l.Info("vendorImageFullPath: ", imageFor, vendorImageFullPath)

	out, err := os.Create(vendorImageFullPath)
	if err != nil {
		if utils.CheckFileExists(vendorImageFullPath) {
			ul.l.Error("updating  is already exitis: ", vendorId, uploadPath, err)
		} else {
			ul.l.Error("vendorImagePath create error: ", vendorId, uploadPath, err)
			defer out.Close()
			return nil, err
		}
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ul.l.Error("vendor upload Copy error: ", vendorId, uploadPath, err)
		return nil, err
	}

	imageDirectory = filepath.Join(imageDirectory, imageName)
	ul.l.Info("##### table to be stored path: ", imageFor, imageDirectory)
	updateQuery := fmt.Sprintf(`UPDATE vendor SET %v = '%v', status = '%v' WHERE vendor_id = '%v'`, imageFor, imageDirectory, constant.STATUS_SUBMITTED, vendorId)
	errU := ul.vendorDao.UpdateVendorImagePath(updateQuery)
	if errU != nil {
		ul.l.Error("ERROR: UpdateVehicleImagePath", vendorId, errU)
		return nil, errU
	}

	ul.l.Info("Image uploaded successfully: ", vendorInfo.VendorName)
	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", vendorInfo.VendorName)
	return &roleResponse, nil
}
