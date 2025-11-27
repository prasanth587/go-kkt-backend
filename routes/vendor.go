package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/vendor"
)

func vendorHub(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/vendor", wrapHandler(recoverHandler.ThenFunc(createvendor)))
	router.GET("/v1/:orgId/vendors", wrapHandler(recoverHandler.ThenFunc(getVendors)))
	router.PUT("/v1/vendor/:vendorId/update/status", wrapHandler(recoverHandler.ThenFunc(updatevendorActiveStatus)))
	router.POST("/v1/vendor/:vendorId/update", wrapHandler(recoverHandler.ThenFunc(updateVendor)))
	router.POST("/v1/vendor/:vendorId/upload", wrapHandler(recoverHandler.ThenFunc(uploadvendorImages)))

	router.POST("/v2/create/vendor", wrapHandler(recoverHandler.ThenFunc(createvendorV1)))
	router.GET("/v2/:orgId/vendors", wrapHandler(recoverHandler.ThenFunc(getVendorsV1)))
	router.PUT("/v2/vendor/:vendorId/update/status", wrapHandler(recoverHandler.ThenFunc(updatevendorActiveStatusV1)))
	router.POST("/v2/vendor/:vendorId/update", wrapHandler(recoverHandler.ThenFunc(updateVendor1)))
	router.POST("/v2/upload/vendor/doc", wrapHandler(recoverHandler.ThenFunc(uploadVendorImagesV1)))

}

func createvendor(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ven := vendor.New(rd.l, rd.dbConnMSSQL)
	var venObj dtos.VendorReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&venObj)
	if err != nil {
		rd.l.Error("createvendor: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := ven.CreateVendor(venObj)
	if err != nil {
		rd.l.Error("CreateVendor error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func createvendorV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ven := vendor.New(rd.l, rd.dbConnMSSQL)
	var venObj dtos.VendorRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&venObj)
	if err != nil {
		rd.l.Error("createvendorV1: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := ven.CreateVendorV1(venObj)
	if err != nil {
		rd.l.Error("CreateVendor error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getVendors(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	searchText := strings.TrimSpace(keys.Get("search"))
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getVendors", "limit", limit, "offset", offset, "orgID", orgID)
	if !isErr {
		return
	}
	vendorId := keys.Get("vendor_id")
	ven := vendor.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.GetVendors(orgID, vendorId, limit, offset, searchText)
	if err != nil {
		rd.l.Error("GetVendors error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getVendorsV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	searchText := strings.TrimSpace(keys.Get("search"))
	vendorId := keys.Get("vendor_id")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getVendors", "limit", limit, "offset", offset, "orgID", orgID, "searchText", searchText, "vendorId", vendorId)
	if !isErr {
		return
	}

	ven := vendor.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.GetVendorsV1(orgID, vendorId, limit, offset, searchText)
	if err != nil {
		rd.l.Error("GetVendors error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateVendor(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ven := vendor.New(rd.l, rd.dbConnMSSQL)
	var vendorObj dtos.VendorUpdate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vendorObj)
	if err != nil {
		rd.l.Error("updatevendor: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	vendorId, isErr := GetIDFromParams(w, r, "vendorId")
	if !isErr {
		return
	}
	rd.l.Info("updatevendor", "vendorId", vendorId)

	res, err := ven.UpdateVendor(vendorId, vendorObj)
	if err != nil {
		rd.l.Error("UpdateVendor error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateVendor1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ven := vendor.New(rd.l, rd.dbConnMSSQL)
	var vendorObj dtos.VendorV1Update
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vendorObj)
	if err != nil {
		rd.l.Error("updatevendor: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	vendorId, isErr := GetIDFromParams(w, r, "vendorId")
	if !isErr {
		return
	}
	rd.l.Info("updatevendor", "vendorId", vendorId)

	res, err := ven.UpdateVendorV1(vendorId, vendorObj)
	if err != nil {
		rd.l.Error("UpdateVendor error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updatevendorActiveStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	status := keys.Get("status")
	vendorID, isErr := GetIDFromParams(w, r, "vendorId")
	rd.l.Info("updatevendorActiveStatus", "vendorID", vendorID, "isActive", isActive)
	if !isErr {
		return
	}

	ven := vendor.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.UpdateVendorActiveStatus(isActive, status, vendorID)
	if err != nil {
		rd.l.Error("updatevendorActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
func updatevendorActiveStatusV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	status := keys.Get("status")
	vendorID, isErr := GetIDFromParams(w, r, "vendorId")
	rd.l.Info("updatevendorActiveStatus", "vendorID", vendorID, "isActive", isActive)
	if !isErr {
		return
	}

	ven := vendor.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.UpdateVendorActiveStatusV1(isActive, status, vendorID)
	if err != nil {
		rd.l.Error("updatevendorActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func uploadvendorImages(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	vendor := vendor.New(rd.l, rd.dbConnMSSQL)
	vendorId, isErr := GetIDFromParams(w, r, "vendorId")
	if !isErr {
		return
	}
	keys := r.URL.Query()
	imageFor := keys.Get("imageFor")
	// Get the file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		rd.l.Error("uploadvendorImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	res, err := vendor.UploadVendorImages(vendorId, imageFor, file, header)
	if err != nil {
		rd.l.Error("UploadVendorImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func uploadVendorImagesV1(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	vendor := vendor.New(rd.l, rd.dbConnMSSQL)
	keys := r.URL.Query()
	vendorCode := keys.Get("vendor_code")
	imageFor := keys.Get("imageFor")
	vendorId := keys.Get("vendor_id")
	vehicleID := keys.Get("vehicle_id")

	req := dtos.VendorAndVehicleUpload{}
	req.RCExpiryDoc = keys.Get("rc_expiry_doc")
	req.InsuranceDoc = keys.Get("insurance_doc")
	req.PUCCExpiryDoc = keys.Get("pucc_expiry_doc")
	req.NPExpireDoc = keys.Get("np_expire_doc")
	req.FitnessExpiryDoc = keys.Get("fitness_expiry_doc")
	req.TaxExpiryDoc = keys.Get("tax_expiry_doc")
	req.MPExpireDoc = keys.Get("mp_expire_doc")
	req.PancardImg = keys.Get("pancard_img")
	req.BankPassbookORChequeImg = keys.Get("bank_passbook_or_cheque_img")
	req.TdsDeclaration = keys.Get("tds_declaration")

	req.VehicleID = keys.Get("vehicle_id")
	req.VendorID = keys.Get("vendor_id")

	jsonBytes, _ := json.Marshal(req)
	rd.l.Info("VendorAndVehicleUpload: ", string(jsonBytes))

	// Get the file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		rd.l.Error("uploadcustomerImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	res, err := vendor.UploadVendorImagesV1(vendorCode, imageFor, file, header, vendorId, vehicleID)
	if err != nil {
		rd.l.Error("UploadCustomerImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
