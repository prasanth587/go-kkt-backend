package routes

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/attendance"
)

func attendanceEmp(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/attendance/in/:empId", wrapHandler(recoverHandler.ThenFunc(createInAttendance)))
	router.POST("/v1/create/attendance/out/:empId", wrapHandler(recoverHandler.ThenFunc(createOutAttendance)))
	router.GET("/employees/v1/attendance", wrapHandler(recoverHandler.ThenFunc(listEmployeeAttendances)))
	// router.PUT("/v1/branch/:branchId/update/status", wrapHandler(recoverHandler.ThenFunc(updateBranchActiveStatus)))
	// router.POST("/v1/branch/:branchId/update", wrapHandler(recoverHandler.ThenFunc(updateBranch)))
}

func createInAttendance(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	employeeId, isErr := GetIDFromParams(w, r, "empId")
	rd.l.Info("updateBranchActiveStatus", "employeeId", employeeId)
	if !isErr {
		return
	}

	var attendaceObj dtos.CreateInAttendance
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&attendaceObj)
	if err != nil {
		rd.l.Error("updateBranch: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	att := attendance.New(rd.l, rd.dbConnMSSQL)
	res, err := att.CreateInAttendance(attendaceObj, employeeId)
	if err != nil {
		rd.l.Error("updateBranchActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func createOutAttendance(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	employeeId, isErr := GetIDFromParams(w, r, "empId")
	rd.l.Info("updateBranchActiveStatus", "employeeId", employeeId)
	if !isErr {
		return
	}

	var attendaceObj dtos.CreateOutAttendance
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&attendaceObj)
	if err != nil {
		rd.l.Error("updateBranch: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	att := attendance.New(rd.l, rd.dbConnMSSQL)
	res, err := att.CreateOutAttendance(attendaceObj, employeeId)
	if err != nil {
		rd.l.Error("updateBranchActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func listEmployeeAttendances(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	date := keys.Get("date")
	att := attendance.New(rd.l, rd.dbConnMSSQL)
	res, err := att.GetEmployeeAttendances(limit, offset, date)
	if err != nil {
		rd.l.Error("GetBranch error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
