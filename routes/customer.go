package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/customer"
)

func customerHub(router *httprouter.Router, recoverHandler alice.Chain) {
	// router.POST("/v1/create/customer", wrapHandler(recoverHandler.ThenFunc(createCustomer)))
	// router.GET("/v1/:orgId/customers", wrapHandler(recoverHandler.ThenFunc(getCustomers)))
	// router.PUT("/v1/customer/:customerId/update/status", wrapHandler(recoverHandler.ThenFunc(updateCustomerActiveStatus)))
	// router.POST("/v1/customer/:customerId/update", wrapHandler(recoverHandler.ThenFunc(updateCustomer)))
	// router.POST("/v1/customer/:customerId/upload", wrapHandler(recoverHandler.ThenFunc(uploadCustomerImages)))

	router.POST("/v2/create/customer", wrapHandler(recoverHandler.ThenFunc(createCustomerV1)))
	router.GET("/v2/:orgId/customers", wrapHandler(recoverHandler.ThenFunc(getCustomersV1)))
	router.POST("/v2/customer/:customerId/update", wrapHandler(recoverHandler.ThenFunc(updateCustomerV1)))
	router.PUT("/v2/customer/:customerId/update/status", wrapHandler(recoverHandler.ThenFunc(updateCustomerActiveStatusV1)))
	router.POST("/v2/customer/:customerId/upload", wrapHandler(recoverHandler.ThenFunc(uploadCustomerImagesV1)))

}

// func createCustomer(w http.ResponseWriter, r *http.Request) {
// 	rd := logAndGetContext(w, r)
// 	cus := customer.New(rd.l, rd.dbConnMSSQL)
// 	var customerReq dtos.CustomerReq
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&customerReq)
// 	if err != nil {
// 		rd.l.Error("createcustomer: request error: ", err.Error())
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	res, err := cus.CreateCustomer(customerReq)
// 	if err != nil {
// 		rd.l.Error("CreateCustomer error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	writeJSONStruct(res, http.StatusOK, rd)
// }

func createCustomerV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	cus := customer.New(rd.l, rd.dbConnMSSQL)
	var customerReq dtos.CustomersReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&customerReq)
	if err != nil {
		rd.l.Error("createcustomer: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := cus.CreateCustomerV1(customerReq)
	if err != nil {
		rd.l.Error("CreateCustomer error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

// func getCustomers(w http.ResponseWriter, r *http.Request) {
// 	rd := logAndGetContext(w, r)
// 	keys := r.URL.Query()
// 	limit := keys.Get("limit")
// 	offset := keys.Get("offset")
// 	orgID, isErr := GetIDFromParams(w, r, "orgId")
// 	rd.l.Info("getCustomers", "limit", limit, "offset", offset, "orgID", orgID)
// 	if !isErr {
// 		return
// 	}
// 	customerId := keys.Get("customer_id")
// 	searchText := strings.TrimSpace(keys.Get("search"))

// 	cus := customer.New(rd.l, rd.dbConnMSSQL)
// 	res, err := cus.GetCustomers(orgID, customerId, limit, offset, searchText)
// 	if err != nil {
// 		rd.l.Error("GetCustomers error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	writeJSONStruct(res, http.StatusOK, rd)
// }

func getCustomersV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	searchText := strings.TrimSpace(keys.Get("search"))
	rd.l.Info("getCustomers", "limit", limit, "offset", offset, "orgID", orgID, "searchText:", searchText)
	if !isErr {
		return
	}
	customerId := keys.Get("customer_id")
	cus := customer.New(rd.l, rd.dbConnMSSQL)
	res, err := cus.GetCustomersV1(orgID, customerId, limit, offset, searchText)
	if err != nil {
		rd.l.Error("GetCustomers error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

// func updateCustomer(w http.ResponseWriter, r *http.Request) {
// 	rd := logAndGetContext(w, r)
// 	cus := customer.New(rd.l, rd.dbConnMSSQL)
// 	var customerObj dtos.CustomerUpdate
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&customerObj)
// 	if err != nil {
// 		rd.l.Error("updatecustomer: request error: ", err.Error())
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	customerId, isErr := GetIDFromParams(w, r, "customerId")
// 	if !isErr {
// 		return
// 	}
// 	rd.l.Info("updatecustomer", "customerId", customerId)

// 	res, err := cus.UpdateCustomer(customerId, customerObj)
// 	if err != nil {
// 		rd.l.Error("UpdateCustomer error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	writeJSONStruct(res, http.StatusOK, rd)
// }

func updateCustomerV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	cus := customer.New(rd.l, rd.dbConnMSSQL)
	var customerObj dtos.CustomersUpdateReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&customerObj)
	if err != nil {
		rd.l.Error("updatecustomer: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	customerId, isErr := GetIDFromParams(w, r, "customerId")
	if !isErr {
		return
	}
	rd.l.Info("updatecustomer", "customerId", customerId)

	res, err := cus.UpdateCustomerV1(customerId, customerObj)
	if err != nil {
		rd.l.Error("UpdateCustomer error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

// func updateCustomerActiveStatus(w http.ResponseWriter, r *http.Request) {
// 	rd := logAndGetContext(w, r)
// 	keys := r.URL.Query()
// 	isActive := keys.Get("isActive")
// 	status := keys.Get("status")
// 	customerID, isErr := GetIDFromParams(w, r, "customerId")
// 	rd.l.Info("updateCustomerActiveStatus", "customerID", customerID, "isActive", isActive)
// 	if !isErr {
// 		return
// 	}
// 	cus := customer.New(rd.l, rd.dbConnMSSQL)
// 	res, err := cus.UpdateCustomerActiveStatus(isActive, status, customerID)
// 	if err != nil {
// 		rd.l.Error("UpdateCustomerActiveStatus error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	writeJSONStruct(res, http.StatusOK, rd)
// }

func updateCustomerActiveStatusV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	status := keys.Get("status")
	customerID, isErr := GetIDFromParams(w, r, "customerId")
	rd.l.Info("updateCustomerActiveStatusV1", "customerID", customerID, "isActive", isActive)
	if !isErr {
		return
	}
	cus := customer.New(rd.l, rd.dbConnMSSQL)
	res, err := cus.UpdateCustomerActiveStatusV1(isActive, status, customerID)
	if err != nil {
		rd.l.Error("UpdateCustomerActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

// func uploadCustomerImages(w http.ResponseWriter, r *http.Request) {
// 	rd := logAndGetContext(w, r)
// 	customer := customer.New(rd.l, rd.dbConnMSSQL)
// 	customerId, isErr := GetIDFromParams(w, r, "customerId")
// 	if !isErr {
// 		return
// 	}
// 	keys := r.URL.Query()
// 	imageFor := keys.Get("imageFor")
// 	// Get the file from the request
// 	file, header, err := r.FormFile("image")
// 	if err != nil {
// 		rd.l.Error("uploadcustomerImages error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	defer file.Close()

// 	res, err := customer.UploadCustomerImages(customerId, imageFor, file, header)
// 	if err != nil {
// 		rd.l.Error("UploadCustomerImages error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	writeJSONStruct(res, http.StatusOK, rd)
// }

func uploadCustomerImagesV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	customer := customer.New(rd.l, rd.dbConnMSSQL)
	customerId, isErr := GetIDFromParams(w, r, "customerId")
	if !isErr {
		return
	}
	keys := r.URL.Query()
	period := keys.Get("period")
	agrementNo := keys.Get("aggrement_no")
	aggrementName := keys.Get("aggrement_name")
	aggrementType := keys.Get("aggrement_type")
	remark := keys.Get("remark")
	aggrementId := keys.Get("aggrement_id")

	req := dtos.CustomersAggrementUpload{}
	req.AggrementName = aggrementName
	req.Period = period
	req.AggrementType = aggrementType
	req.Remark = remark
	req.AggrementNo = agrementNo
	req.AggrementId = aggrementId
	req.CustomerId = customerId

	jsonBytes, _ := json.Marshal(req)
	rd.l.Info("CustomersAggrementUpload: ", string(jsonBytes))

	// Get the file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		rd.l.Error("uploadcustomerImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	res, err := customer.UploadCustomerImagesV1(customerId, req, file, header)
	if err != nil {
		rd.l.Error("UploadCustomerImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
