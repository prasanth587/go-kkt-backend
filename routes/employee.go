package routes

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/employee"
)

func employeeV1(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/employee", wrapHandler(recoverHandler.ThenFunc(createEmpolyee)))
	router.GET("/v1/:orgId/employees", wrapHandler(recoverHandler.ThenFunc(getEmpolyee)))
	router.PUT("/v1/employee/:empId/update/status", wrapHandler(recoverHandler.ThenFunc(updateEmployeeActiveStatus)))
	router.POST("/v1/employee/:empId/update", wrapHandler(recoverHandler.ThenFunc(updateEmployee)))
	router.POST("/v1/employee/:empId/upload/profile", wrapHandler(recoverHandler.ThenFunc(uploadEmployeeProfile)))
	router.GET("/view/employee/:empId/profile", wrapHandler(recoverHandler.ThenFunc(viweEmployeeProfile)))
}

func employeeRole(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/role", wrapHandler(recoverHandler.ThenFunc(createRole)))
	router.GET("/v1/:orgId/roles", wrapHandler(recoverHandler.ThenFunc(getRoles)))
	router.PUT("/v1/update/:roleId/role", wrapHandler(recoverHandler.ThenFunc(updateRole)))
}

func getEmpolyee(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("GetEmpolyee", "limit", limit, "offset", offset, "orgID", orgID)
	if !isErr {
		return
	}

	emp := employee.New(rd.l, rd.dbConnMSSQL)
	res, err := emp.GetEmployee(orgID, limit, offset)
	if err != nil {
		rd.l.Error("GetEmployee error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateEmployeeActiveStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	empId, isErr := GetIDFromParams(w, r, "empId")
	rd.l.Info("UpdateStatus", "empId", empId, "isActive", isActive)
	if !isErr {
		return
	}

	emp := employee.New(rd.l, rd.dbConnMSSQL)
	res, err := emp.UpdateEmployeeActiveStatus(isActive, empId)
	if err != nil {
		rd.l.Error("UpdateEmployeeActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	emp := employee.New(rd.l, rd.dbConnMSSQL)
	var empObj dtos.UpdateEmployeeReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&empObj)
	if err != nil {
		rd.l.Error("CreateEmpolyee: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	empId, isErr := GetIDFromParams(w, r, "empId")
	if !isErr {
		return
	}
	rd.l.Info("UpdateEmployee", "empId", empId)

	res, err := emp.UpdateEmployee(empId, empObj)
	if err != nil {
		rd.l.Error("UpdateEmployee error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func createEmpolyee(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	emp := employee.New(rd.l, rd.dbConnMSSQL)
	var empObj dtos.CreateEmployeeReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&empObj)
	if err != nil {
		rd.l.Error("CreateEmpolyee: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := emp.CreateEmployee(empObj)
	if err != nil {
		rd.l.Error("CreateEmployee error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getRoles(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("GetRoles", "limit", limit, "offset", offset, "orgID", orgID)
	if !isErr {
		return
	}

	emp := employee.New(rd.l, rd.dbConnMSSQL)
	res, err := emp.GetRoles(orgID, limit, offset)
	if err != nil {
		rd.l.Error("GetRoles error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func createRole(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	emp := employee.New(rd.l, rd.dbConnMSSQL)
	var empRole dtos.EmpRole
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&empRole)
	if err != nil {
		rd.l.Error("request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := emp.CreateEmployeeRole(empRole)
	if err != nil {
		rd.l.Error("CreateEmployeeRole error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateRole(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	emp := employee.New(rd.l, rd.dbConnMSSQL)
	roleId, isErr := GetIDFromParams(w, r, "roleId")
	if !isErr {
		return
	}
	var empRole dtos.EmpRoleUpdate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&empRole)
	if err != nil {
		rd.l.Error("request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := emp.UpdateEmployeeRole(empRole, roleId)
	if err != nil {
		rd.l.Error("UpdateEmployeeRole error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func uploadEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	emp := employee.New(rd.l, rd.dbConnMSSQL)
	empId, isErr := GetIDFromParams(w, r, "empId")
	if !isErr {
		return
	}
	// Get the file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	res, err := emp.UploadEmployeeProfile(empId, file, header)
	if err != nil {
		rd.l.Error("UploadEmployeeProfile error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func viweEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	emp := employee.New(rd.l, rd.dbConnMSSQL)
	empId, isErr := GetIDFromParams(w, r, "empId")
	if !isErr {
		return
	}
	image := r.URL.Query().Get("image")
	imageByte, err := emp.ViweEmployeeProfile(empId, image)
	if err != nil {
		rd.l.Error("ViweEmployeeProfile error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(imageByte)
}
