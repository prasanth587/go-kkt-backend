package routes

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/vehicle"
)

func vehicleHub(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/vehicle", wrapHandler(recoverHandler.ThenFunc(createVehicle)))
	router.GET("/v1/:orgId/vehicles", wrapHandler(recoverHandler.ThenFunc(getVehicles)))
	router.PUT("/v1/vehicle/:vid/update/status", wrapHandler(recoverHandler.ThenFunc(updateVehicleActiveStatus)))
	router.POST("/v1/vehicle/:vid/update", wrapHandler(recoverHandler.ThenFunc(updateVehicle)))
	router.POST("/v1/vehicle/:vid/upload", wrapHandler(recoverHandler.ThenFunc(uploadVehicleImages)))

	// Vehicle size types
	router.POST("/v1/create/vehicle/size/type", wrapHandler(recoverHandler.ThenFunc(createVehicleSizeTypes)))
	router.GET("/v1/:orgId/vehicle/size/type", wrapHandler(recoverHandler.ThenFunc(getVehicleSizeTypes)))
	router.POST("/v1/vehicle/:vid/size/type/update", wrapHandler(recoverHandler.ThenFunc(updateVehicleSizeType)))
	router.PUT("/v1/vehicle/:vid/size/type/status", wrapHandler(recoverHandler.ThenFunc(updateVehicleSizeTypeActiveStatus)))
}

func createVehicle(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	veh := vehicle.New(rd.l, rd.dbConnMSSQL)
	var vehObj dtos.VehicleReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vehObj)
	if err != nil {
		rd.l.Error("createVehicle: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := veh.CreateVehicle(vehObj)
	if err != nil {
		rd.l.Error("CreateVehicle error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getVehicles(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getDrivers", "limit", limit, "offset", offset, "orgID", orgID)
	if !isErr {
		return
	}
	veh := vehicle.New(rd.l, rd.dbConnMSSQL)
	res, err := veh.GetVehicles(orgID, limit, offset)
	if err != nil {
		rd.l.Error("GetVehicles error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateVehicle(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	veh := vehicle.New(rd.l, rd.dbConnMSSQL)
	var vehicleObj dtos.VehicleUpdate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vehicleObj)
	if err != nil {
		rd.l.Error("updateVehicle: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	vehicleId, isErr := GetIDFromParams(w, r, "vid")
	if !isErr {
		return
	}
	rd.l.Info("updateVehicle", "vehicleId", vehicleId)

	res, err := veh.UpdateVehicle(vehicleId, vehicleObj)
	if err != nil {
		rd.l.Error("UpdateVehicle error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateVehicleActiveStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	vehicleID, isErr := GetIDFromParams(w, r, "vid")
	rd.l.Info("updateVehicleActiveStatus", "vehicleID", vehicleID, "isActive", isActive)
	if !isErr {
		return
	}

	veh := vehicle.New(rd.l, rd.dbConnMSSQL)
	res, err := veh.UpdateVehicleActiveStatus(isActive, vehicleID)
	if err != nil {
		rd.l.Error("UpdateVehicleActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func uploadVehicleImages(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	vehicle := vehicle.New(rd.l, rd.dbConnMSSQL)
	vehicleId, isErr := GetIDFromParams(w, r, "vid")
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

	res, err := vehicle.UploadVehicleImages(vehicleId, imageFor, file, header)
	if err != nil {
		rd.l.Error("UploadVehicleImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func createVehicleSizeTypes(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	veh := vehicle.New(rd.l, rd.dbConnMSSQL)
	var vehSizeObj dtos.VehicleSizeType
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vehSizeObj)
	if err != nil {
		rd.l.Error("createVehicleSizeTypes: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := veh.CreateVehicleSizeTypes(vehSizeObj)
	if err != nil {
		rd.l.Error("createVehicleSizeTypes error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getVehicleSizeTypes(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getDrivers", "limit", limit, "offset", offset, "orgID", orgID)
	if !isErr {
		return
	}
	veh := vehicle.New(rd.l, rd.dbConnMSSQL)
	res, err := veh.GetVehicleSizeTypes(limit, offset)
	if err != nil {
		rd.l.Error("GetVehicleSizeTypes error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateVehicleSizeType(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	veh := vehicle.New(rd.l, rd.dbConnMSSQL)
	var vehicleSizeObj dtos.VehicleSizeTypeUpdate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vehicleSizeObj)
	if err != nil {
		rd.l.Error("updateVehicleSizeType: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	vehicleSizeId, isErr := GetIDFromParams(w, r, "vid")
	if !isErr {
		return
	}
	rd.l.Info("updateVehicleSizeType", "vehicleSizeId", vehicleSizeId)

	res, err := veh.UpdateVehicleSizeType(vehicleSizeId, vehicleSizeObj)
	if err != nil {
		rd.l.Error("updateVehicleSizeType error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateVehicleSizeTypeActiveStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	vehicleID, isErr := GetIDFromParams(w, r, "vid")
	rd.l.Info("updateVehicleActiveStatus", "vehicleID", vehicleID, "isActive", isActive)
	if !isErr {
		return
	}

	veh := vehicle.New(rd.l, rd.dbConnMSSQL)
	res, err := veh.UpdateVehicleSizeTypeActiveStatus(isActive, vehicleID)
	if err != nil {
		rd.l.Error("UpdateVehicleSizeTypeActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
