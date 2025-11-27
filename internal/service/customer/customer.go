package customer

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

type CustomerObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
	customerDao daos.CustomerDao
}

var (
	USER_SUCCESS = "User logged in successfully"
)

const (
	EmployeeMobilePattern = `^((\+)?(\d{2}[-])?(\d{10}){1})?(\d{11}){0,1}?$`
)

// DatePattern
const (
	DatePattern = `^\d{1,2}\/\d{1,2}\/\d{4}$`
)

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *CustomerObj {
	return &CustomerObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		customerDao: daos.NewCustomerObj(l, dbConnMSSQL),
	}
}

// func (cus *CustomerObj) CreateCustomer(customerReq dtos.CustomerReq) (*dtos.Messge, error) {

// 	err := cus.validateCustomer(customerReq)
// 	if err != nil {
// 		cus.l.Error("ERROR: CreateCustomer", err)
// 		return nil, err
// 	}

// 	customerReq.Status = constant.STATUS_CREATED

// 	err1 := cus.customerDao.CreateCustomer(customerReq)
// 	if err1 != nil {
// 		cus.l.Error("customer not saved ", customerReq.CustomerCode, err1)
// 		return nil, err1
// 	}

// 	cus.l.Info("Customer created successfully! : ", customerReq.CustomerName)

// 	roleResponse := dtos.Messge{}
// 	roleResponse.Message = fmt.Sprintf("Customer saved successfully: %s", customerReq.CustomerName)
// 	return &roleResponse, nil
// }

func (cus *CustomerObj) CreateCustomerV1(customerReq dtos.CustomersReq) (*dtos.Messge, error) {

	err := cus.validateCustomerReq(customerReq)
	if err != nil {
		cus.l.Error("ERROR: CreateCustomer", err)
		return nil, err
	}

	customerReq.Status = constant.STATUS_CREATED

	customerId, err1 := cus.customerDao.CreateCustomerV1(customerReq)
	if err1 != nil {
		cus.l.Error("customer not saved ", customerReq.CustomerCode, err1)
		return nil, err1
	}

	cus.l.Info("Customer created successfully! : ", customerId, customerReq.CustomerName)

	if len(customerReq.ContactInfo) != 0 {
		for _, contactInfo := range customerReq.ContactInfo {
			err1 := cus.customerDao.CreateCustomerContactInfo(customerId, contactInfo)
			if err1 != nil {
				cus.l.Error("CreateCustomerContactInfo not saved ", contactInfo.ContactPerson, err1)
				return nil, err1
			}
		}
	}

	resMsg := dtos.Messge{}
	resMsg.Message = fmt.Sprintf("Customer saved successfully: %s", customerReq.CustomerName)
	return &resMsg, nil
}

// func (cus *CustomerObj) GetCustomers(orgId int64, customerId, limit string, offset, searchText string) (*dtos.CustomerEntries, error) {

// 	limitI, errInt := strconv.ParseInt(limit, 10, 64)
// 	if errInt != nil {
// 		return nil, errors.New("invalid limit")
// 	}
// 	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
// 	if errInt != nil {
// 		return nil, errors.New("invalid offset")
// 	}

// 	whereQuery := cus.customerDao.BuildWhereQuery(orgId, customerId, searchText)
// 	res, errA := cus.customerDao.GetCustomers(orgId, whereQuery, limitI, offsetI)
// 	if errA != nil {
// 		cus.l.Error("ERROR: GetCustomers", errA)
// 		return nil, errA
// 	}

// 	customerEntries := dtos.CustomerEntries{}
// 	customerEntries.CustomerEntiry = res
// 	customerEntries.Total = cus.customerDao.GetTotalCount(whereQuery)
// 	customerEntries.Limit = limitI
// 	customerEntries.OffSet = offsetI

// 	return &customerEntries, nil
// }

func (cus *CustomerObj) GetCustomersV1(orgId int64, customerId, limit string, offset, searchText string) (*dtos.CustomersResponse, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	whereQuery := cus.customerDao.BuildWhereQuery(orgId, customerId, searchText)
	res, errA := cus.customerDao.GetCustomersV1(orgId, whereQuery, limitI, offsetI)
	if errA != nil {
		cus.l.Error("ERROR: GetCustomers", errA)
		return nil, errA
	}

	for i := range *res {
		contactInfo, errA := cus.customerDao.GetCustomerContactInfo((*res)[i].CustomerId)
		if errA != nil {
			cus.l.Error("ERROR: GetCustomerContactInfo", errA)
			return nil, errA
		}
		(*res)[i].ContactInfo = contactInfo

		//
		//CustomersAggrement
		aggrements, errA := cus.customerDao.GetCustomerAggrements((*res)[i].CustomerId)
		if errA != nil {
			cus.l.Error("ERROR: GetCustomerAggrements", errA)
			return nil, errA
		}
		(*res)[i].CustomersAggrement = aggrements
	}

	customerEntries := dtos.CustomersResponse{}
	customerEntries.CustomerEntiry = res
	customerEntries.Total = cus.customerDao.GetTotalCount(whereQuery)
	customerEntries.Limit = limitI
	customerEntries.OffSet = offsetI

	return &customerEntries, nil
}

// func (cus *CustomerObj) UpdateCustomer(customerId int64, customerReq dtos.CustomerUpdate) (*dtos.Messge, error) {

// 	if customerReq.CustomerName == "" {
// 		cus.l.Error("Error CustomerName: customer name should not empty")
// 		return nil, errors.New("customer name should not empty")
// 	}
// 	if customerReq.CustomerCode == "" {
// 		cus.l.Error("Error CustomerCode: customer code should not empty")
// 		return nil, errors.New("customer code should not empty")
// 	}
// 	customerReq.CustomerCode = strings.ToUpper(customerReq.CustomerCode)

// 	if customerReq.MobileNumber == "" {
// 		cus.l.Error("Error MobileNumber: Mobile number should not empty")
// 		return nil, errors.New("mobile number should not empty")
// 	}
// 	if customerReq.ContactPerson == "" {
// 		cus.l.Error("Error ContactPerson: contact person should not empty")
// 		return nil, errors.New("contact person should not empty")
// 	}
// 	if customerReq.City == "" {
// 		cus.l.Error("Error City: city should not empty")
// 		return nil, errors.New("city should not empty")
// 	}
// 	customerInfo, errV := cus.customerDao.GetCustomer(customerId)
// 	if errV != nil {
// 		cus.l.Error("ERROR: GetEmployee not found", errV)
// 		return nil, errV
// 	}
// 	jsonBytes, _ := json.Marshal(customerInfo)
// 	cus.l.Info("GetCustomer: ******* ", string(jsonBytes))

// 	customerReq.Status = customerInfo.Status

// 	err1 := cus.customerDao.UpdateCustomer(customerId, customerReq)
// 	if err1 != nil {
// 		cus.l.Error("customer not updated ", customerReq.CustomerName, err1)
// 		return nil, err1
// 	}

// 	cus.l.Info("Customer updated successfully! : ", customerReq.CustomerName)

// 	roleResponse := dtos.Messge{}
// 	roleResponse.Message = fmt.Sprintf("Customer saved successfully: %s", customerReq.CustomerName)
// 	return &roleResponse, nil
// }

func (cus *CustomerObj) UpdateCustomerV1(customerId int64, customerReq dtos.CustomersUpdateReq) (*dtos.Messge, error) {

	if customerReq.CustomerName == "" {
		cus.l.Error("Error CustomerName: customer name should not empty")
		return nil, errors.New("customer name should not empty")
	}
	if customerReq.CustomerCode == "" {
		cus.l.Error("Error CustomerCode: customer code should not empty")
		return nil, errors.New("customer code should not empty")
	}
	customerReq.CustomerCode = strings.ToUpper(customerReq.CustomerCode)

	// if customerReq.City == "" {
	// 	cus.l.Error("Error City: city should not empty")
	// 	return nil, errors.New("city should not empty")
	// }
	customerInfo, errV := cus.customerDao.GetCustomerV1(customerId)
	if errV != nil {
		cus.l.Error("ERROR: GetEmployee not found", errV)
		return nil, errV
	}
	jsonBytes, _ := json.Marshal(customerInfo)
	cus.l.Info("GetCustomer: ******* ", string(jsonBytes))

	customerReq.Status = customerInfo.Status
	// if customerInfo.Status == constant.STATUS_CREATED {
	// 	//customerReq.Status = constant.STATUS_SUBMITTED
	// }

	err1 := cus.customerDao.UpdateCustomerV1(customerId, customerReq)
	if err1 != nil {
		cus.l.Error("customer not updated ", customerReq.CustomerName, err1)
		return nil, err1
	}

	// Deleting existing contact info for the customer
	errD := cus.customerDao.DeleteCustomerContactInfo(customerId)
	if errD != nil {
		cus.l.Error("DeleteCustomerContactInfo not deleted ", customerReq.CustomerName, customerId, errD)
		return nil, err1
	}

	// creating again contact info for the customer
	if len(customerReq.ContactInfo) != 0 {
		for _, contactInfo := range customerReq.ContactInfo {
			err1 := cus.customerDao.CreateCustomerContactInfo(customerId, contactInfo)
			if err1 != nil {
				cus.l.Error("CreateCustomerContactInfo not saved ", contactInfo.ContactPerson, err1)
				return nil, err1
			}
		}
	}

	cus.l.Info("Customer updated successfully! : ", customerReq.CustomerName)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Customer saved successfully: %s", customerReq.CustomerName)
	return &roleResponse, nil
}

// func (cus *CustomerObj) UpdateCustomerActiveStatus(isA, status string, customerID int64) (*dtos.Messge, error) {

// 	isActive, errInt := strconv.ParseInt(isA, 10, 64)
// 	if errInt != nil && status == "" {
// 		return nil, errors.New("invalid status update")
// 	}
// 	if customerID == 0 {
// 		cus.l.Error("unknown customer", customerID)
// 		return nil, errors.New("unknown customer")
// 	}
// 	statusRes := dtos.Messge{}
// 	customerInfo, errA := cus.customerDao.GetCustomer(customerID)
// 	if errA != nil {
// 		cus.l.Error("ERROR: Customer not found", errA)
// 		return nil, errA
// 	}
// 	jsonBytes, _ := json.Marshal(customerInfo)
// 	cus.l.Info("customerInfo: ", status, isActive, string(jsonBytes))

// 	var errU error
// 	if status != "" {
// 		errU = cus.customerDao.UpdateCustomerStatus(customerID, status)
// 	} else {
// 		errU = cus.customerDao.UpdateCustomerActiveStatus(customerID, isActive)
// 	}
// 	if errU != nil {
// 		cus.l.Error("ERROR: UpdateCustomer Active/Status ", errU)
// 		return nil, errU
// 	}

// 	statusRes.Message = "customer updated successfully : " + customerInfo.CustomerName
// 	return &statusRes, nil
// }

func (cus *CustomerObj) UpdateCustomerActiveStatusV1(isA, status string, customerID int64) (*dtos.Messge, error) {

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil && status == "" {
		return nil, errors.New("invalid status update")
	}
	if customerID == 0 {
		cus.l.Error("unknown customer", customerID)
		return nil, errors.New("unknown customer")
	}

	statusRes := dtos.Messge{}
	customerInfo, errA := cus.customerDao.GetCustomerV1(customerID)
	if errA != nil {
		cus.l.Error("ERROR: Customer not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(customerInfo)
	cus.l.Info("customerInfo: ", status, isActive, string(jsonBytes))

	var errU error
	if status != "" {
		stype := [...]string{
			constant.STATUS_CREATED, constant.STATUS_SUBMITTED, constant.STATUS_CANCELLED, constant.STATUS_APPROVED,
		}

		exists := false
		for _, nype := range stype {
			if nype == status {
				exists = true
				break
			}
		}
		if !exists {
			return nil, fmt.Errorf("unknown status. the status should be anyone of them %v", stype)
		}
		errU = cus.customerDao.UpdateCustomerStatusV1(customerID, status)
	} else {
		errU = cus.customerDao.UpdateCustomerActiveStatusV1(customerID, isActive)
	}
	if errU != nil {
		cus.l.Error("ERROR: UpdateCustomer Active/Status ", errU)
		return nil, errU
	}

	statusRes.Message = "customer updated successfully : " + customerInfo.CustomerName
	return &statusRes, nil
}

// func (cus *CustomerObj) UploadCustomerImages(customerId int64, imageFor string, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.Messge, error) {

// 	if imageFor == "" {
// 		return nil, errors.New("imageFor should not be empty. imageFor:agreement_doc_image")
// 	}

// 	imageTypes := [...]string{
// 		"agreement_doc_image",
// 	}
// 	exists := false
// 	for _, imgType := range imageTypes {
// 		if imgType == imageFor {
// 			exists = true
// 			break
// 		}
// 	}
// 	if !exists {
// 		return nil, errors.New("imageFor should not be empty. imageFor:agreement_doc_image")
// 	}

// 	customerInfo, errV := cus.customerDao.GetCustomer(customerId)
// 	if errV != nil {
// 		cus.l.Error("ERROR: customer not found", errV)
// 		return nil, errV
// 	}
// 	//jsonBytes, _ := json.Marshal(customerInfo)
// 	//cus.l.Info("GetCustomer: ******* ", string(jsonBytes))

// 	baseDirectory := os.Getenv("BASE_DIRECTORY")
// 	uploadPath := os.Getenv("IMAGE_DIRECTORY")
// 	if uploadPath == "" || baseDirectory == "" {
// 		cus.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
// 		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
// 	}

// 	imageDirectory := filepath.Join(uploadPath, "customer", strconv.Itoa(int(customerId)))
// 	fullPath := filepath.Join(baseDirectory, imageDirectory)
// 	cus.l.Infof("customerId: %v imageDirectory: %s, fullPath: %s", customerId, imageDirectory, fullPath)

// 	err := os.MkdirAll(fullPath, os.ModePerm) // os.ModePerm sets permissions to 0777
// 	if err != nil {
// 		cus.l.Error("ERROR: MkdirAll ", fullPath, err)
// 		return nil, err
// 	}

// 	extension := strings.Split(fileHeader.Filename, ".")
// 	lengthExt := len(extension)

// 	imageName := fmt.Sprintf("%v_%v.%v", imageFor, customerId, extension[lengthExt-1])
// 	customerImageFullPath := filepath.Join(fullPath, imageName)
// 	cus.l.Info("customerImageFullPath: ", imageFor, customerImageFullPath)

// 	out, err := os.Create(customerImageFullPath)
// 	if err != nil {
// 		if utils.CheckFileExists(customerImageFullPath) {
// 			cus.l.Error("updating  is already exitis: ", customerId, uploadPath, err)
// 		} else {
// 			cus.l.Error("customerImagePath create error: ", customerId, uploadPath, err)
// 			defer out.Close()
// 			return nil, err
// 		}
// 	}
// 	defer out.Close()

// 	_, err = io.Copy(out, file)
// 	if err != nil {
// 		cus.l.Error("customer upload Copy error: ", customerId, uploadPath, err)
// 		return nil, err
// 	}

// 	imageDirectory = filepath.Join(imageDirectory, imageName)
// 	cus.l.Info("##### table to be stored path: ", imageFor, imageDirectory)
// 	updateQuery := fmt.Sprintf(`UPDATE customer SET %v = '%v', status = '%v' WHERE customer_id = '%v'`, imageFor, imageDirectory, constant.STATUS_SUBMITTED, customerId)
// 	errU := cus.customerDao.UpdateCustomerImagePath(updateQuery)
// 	if errU != nil {
// 		cus.l.Error("ERROR: UpdateVehicleImagePath", customerId, errU)
// 		return nil, errU
// 	}

// 	cus.l.Info("Image uploaded successfully: ", customerInfo.CustomerName)
// 	roleResponse := dtos.Messge{}
// 	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", customerInfo.CustomerName)
// 	return &roleResponse, nil
// }

func (cus *CustomerObj) UploadCustomerImagesV1(customerId int64, aggReq dtos.CustomersAggrementUpload, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.Messge, error) {

	if aggReq.AggrementName == "" {
		return nil, errors.New("aggrement name should not empty")
	}
	if aggReq.AggrementType == "" {
		return nil, errors.New("aggrement type should not empty")
	}
	if aggReq.Period == "" {
		return nil, errors.New("aggrement periode should not empty")
	}

	customerInfo, errV := cus.customerDao.GetCustomerV1(customerId)
	if errV != nil {
		cus.l.Error("ERROR: customer not found", errV)
		return nil, errV
	}

	baseDirectory := os.Getenv("BASE_DIRECTORY")
	uploadPath := os.Getenv("IMAGE_DIRECTORY")
	if uploadPath == "" || baseDirectory == "" {
		cus.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
	}

	imageDirectory := filepath.Join(uploadPath, "customer", strconv.Itoa(int(customerId)))
	fullPath := filepath.Join(baseDirectory, imageDirectory)
	cus.l.Infof("customerId: %v imageDirectory: %s, fullPath: %s", customerId, imageDirectory, fullPath)

	err := os.MkdirAll(fullPath, os.ModePerm) // os.ModePerm sets permissions to 0777
	if err != nil {
		cus.l.Error("ERROR: MkdirAll ", fullPath, err)
		return nil, err
	}

	extension := strings.Split(fileHeader.Filename, ".")
	lengthExt := len(extension)

	imageFor := aggReq.AggrementNo

	imageName := fmt.Sprintf("%v_%v.%v", imageFor, customerId, extension[lengthExt-1])
	customerImageFullPath := filepath.Join(fullPath, imageName)
	cus.l.Info("customerImageFullPath: ", imageFor, customerImageFullPath)

	out, err := os.Create(customerImageFullPath)
	if err != nil {
		if utils.CheckFileExists(customerImageFullPath) {
			cus.l.Error("updating  is already exitis: ", customerId, uploadPath, err)
		} else {
			cus.l.Error("customerImagePath create error: ", customerId, uploadPath, err)
			defer out.Close()
			return nil, err
		}
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		cus.l.Error("customer upload Copy error: ", customerId, uploadPath, err)
		return nil, err
	}

	imageDirectory = filepath.Join(imageDirectory, imageName)
	cus.l.Info("##### table to be stored path: ", imageFor, customerImageFullPath, imageDirectory)

	// Save/Update Aggrement table
	if aggReq.AggrementId != "" && aggReq.AggrementId != "0" {
		errU := cus.customerDao.UpdateCustomerAggrement(aggReq, imageDirectory)
		if errU != nil {
			cus.l.Error("ERROR: UpdateCustomerAggrement", customerId, aggReq.AggrementName, errU)
			return nil, errU
		}
	} else if aggReq.AggrementId == "" || aggReq.AggrementId == "0" {
		errU := cus.customerDao.SaveCustomerAggrement(aggReq, imageDirectory)
		if errU != nil {
			cus.l.Error("ERROR: SaveCustomerAggrement", customerId, errU)
			return nil, errU
		}
	}

	//customerStatus update
	customerStatus := customerInfo.Status
	if customerInfo.Status == constant.STATUS_CREATED {
		customerStatus = constant.STATUS_SUBMITTED
	}
	errU := cus.customerDao.UpdateCustomerStatusV1(customerId, customerStatus)
	if errU != nil {
		cus.l.Error("ERROR: UpdateCustomerStatusV1.....", customerId, aggReq.AggrementName, errU)
		return nil, errU
	}

	cus.l.Info("Image uploaded successfully: ", customerInfo.CustomerName, aggReq.AggrementName)
	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", aggReq.AggrementName)
	return &roleResponse, nil
}
