package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/internal/service/dashboard"
)

func dashboardAPI(router *httprouter.Router, recoverHandler alice.Chain) {
	router.GET("/v1/:orgId/dashboard/stats", wrapHandler(recoverHandler.ThenFunc(dashboardStats)))
}

func dashboardStats(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()

	startDate := keys.Get("start_date")
	endDate := keys.Get("end_date")
	employee := keys.Get("employee")
	vendor := keys.Get("vendor")
	vehicle := keys.Get("vehicle")
	customer := keys.Get("customer")
	loadUnloadPoints := keys.Get("load_unload_points")
	tripCreated := keys.Get("trip_created")
	tripSubmitted := keys.Get("trip_submitted")
	tripDelivered := keys.Get("trip_delivered")
	tripClosed := keys.Get("trip_closed")
	tripCompleted := keys.Get("trip_completed")
	tripGraph := keys.Get("trip_graph")
	inventory := keys.Get("inventory")

	rd.l.Info("startDate", startDate, "endDate", endDate, "employee", employee, "vendor", vendor,
		"vehicle", vehicle, "customer", customer, "loadUnloadPoints", loadUnloadPoints,
		"tripCreated", tripCreated, "tripSubmitted", tripSubmitted, "tripDelivered", tripDelivered, "tripClosed", tripClosed, "inventory", inventory)

	ven := dashboard.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.GetStats(startDate, endDate, employee, vendor, vehicle, customer, loadUnloadPoints,
		tripCreated, tripSubmitted, tripDelivered, tripClosed, tripGraph, tripCompleted, inventory)
	if err != nil {
		rd.l.Error("updatevendorActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
