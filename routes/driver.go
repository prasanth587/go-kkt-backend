package routes

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/driver"
)

func driverHub(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/org/:org/driver/create", wrapHandler(recoverHandler.ThenFunc(createDriver)))
	router.GET("/v1/:orgId/driver", wrapHandler(recoverHandler.ThenFunc(getDrivers)))
	router.PUT("/v1/driver/:did/update/status", wrapHandler(recoverHandler.ThenFunc(updateDriverActiveStatus)))
	router.POST("/v1/driver/:did/update", wrapHandler(recoverHandler.ThenFunc(updateDriver)))

	router.POST("/v1/driver/:did/upload/profile", wrapHandler(recoverHandler.ThenFunc(uploadDriverImages)))

	// router.GET("/v1/:orgId/employees", wrapHandler(recoverHandler.ThenFunc(getEmpolyee)))
	// router.PUT("/v1/employee/:empId/update/status", wrapHandler(recoverHandler.ThenFunc(updateEmployeeActiveStatus)))
	// router.POST("/v1/employee/:empId/update", wrapHandler(recoverHandler.ThenFunc(updateEmployee)))
	// router.POST("/v1/employee/:empId/upload/profile", wrapHandler(recoverHandler.ThenFunc(uploadEmployeeProfile)))
	// router.GET("/view/employee/:empId/profile", wrapHandler(recoverHandler.ThenFunc(viweEmployeeProfile)))
}

func createDriver(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	dr := driver.New(rd.l, rd.dbConnMSSQL)
	var driverObj dtos.CreateDriverReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&driverObj)
	if err != nil {
		rd.l.Error("CreateDriver: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := dr.CreateDriver(driverObj)
	if err != nil {
		rd.l.Error("CreateDriver error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getDrivers(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getDrivers", "limit", limit, "offset", offset, "orgID", orgID)
	if !isErr {
		return
	}
	dr := driver.New(rd.l, rd.dbConnMSSQL)
	res, err := dr.GetDrivers(orgID, limit, offset)
	if err != nil {
		rd.l.Error("GetDrivers error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateDriverActiveStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	driverID, isErr := GetIDFromParams(w, r, "did")
	rd.l.Info("UpdateStatus", "driverId", driverID, "isActive", isActive)
	if !isErr {
		return
	}

	dr := driver.New(rd.l, rd.dbConnMSSQL)
	res, err := dr.UpdateDriverctiveStatus(isActive, driverID)
	if err != nil {
		rd.l.Error("UpdateDriverctiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateDriver(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	driver := driver.New(rd.l, rd.dbConnMSSQL)
	var driverObj dtos.UpdateDriverReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&driverObj)
	if err != nil {
		rd.l.Error("UpdateDriver: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	driverID, isErr := GetIDFromParams(w, r, "did")
	if !isErr {
		return
	}
	rd.l.Info("UpdateDriver", "driverID", driverID)

	res, err := driver.UpdateDriver(driverID, driverObj)
	if err != nil {
		rd.l.Error("UpdateDriver error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func uploadDriverImages(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	driver := driver.New(rd.l, rd.dbConnMSSQL)
	empId, isErr := GetIDFromParams(w, r, "did")
	if !isErr {
		return
	}
	keys := r.URL.Query()
	imageFor := keys.Get("imageFor")
	// Get the file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	res, err := driver.UploadDriverImages(empId, imageFor, file, header)
	if err != nil {
		rd.l.Error("UploadDriverImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
