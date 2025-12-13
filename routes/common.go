package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	lg "log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/errs"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/commonsvc"
	"go-transport-hub/utils"
)

const (
	ERR_MSG = "ERROR_MESSAGE"
	MSG     = "MESSAGE"
)

type ResStruct struct {
	Status   string `json:"status" example:"SUCCESS" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"200" example:"500"`
	Message  string `json:"message" example:"pong" example:"could not connect to db"`
}

type Res500Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"500"`
	Message  string `json:"message" example:"could not connect to db"`
}

type Res400Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"400"`
	Message  string `json:"message" example:"Invalid param"`
}

type RequestData struct {
	l           *log.Logger
	Start       time.Time
	w           http.ResponseWriter
	r           *http.Request
	dbConnMSSQL *mssqlcon.DBConn
}

type RenderData struct {
	Data  interface{}
	Paths []string
}

type TemplateData struct {
	Data interface{}
}

func (t *TemplateData) SetConstants() {

}

func logAndGetContext(w http.ResponseWriter, r *http.Request) *RequestData {
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("X-Frame-Options", "DENY")

	//url := strings.TrimSpace("http://0.0.0.0:9008/v1/logs")

	//Set config according to the use case
	cfg := log.NewConfig("thub")
	//cfg.SetRemoteConfig(url, "", "admin")
	cfg.SetLevelStr(os.Getenv("LOG_LEVEL"))
	cfg.SetFilePathSizeStr("Full")
	cfg.SetReference(r.Header.Get("ReferenceID"))

	l := log.New(cfg)
	dbConn := new(mssqlcon.DBConn)
	dbConn.Init(l)

	//pgdbConn := new(pgsqldb.Conn)
	//pgdbConn.Init(l)

	start := time.Now().In(utils.TimeLoc())
	l.LogAPIInfo(r, 0, 0)

	return &RequestData{
		l:           l,
		Start:       start,
		r:           r,
		w:           w,
		dbConnMSSQL: dbConn,
	}
}

func redirectTo(path string, rd *RequestData) {
	rd.l.Info("Status Code:", http.StatusFound, ", Response time:", time.Since(rd.Start), ", Response: url redirect - ", path)
	rd.l.LogAPIInfo(rd.r, time.Since(rd.Start).Seconds(), http.StatusFound)
	http.Redirect(rd.w, rd.r, path, http.StatusFound)
}

func jsonifyMessage(msg string, msgType string, httpCode int) ([]byte, int) {
	var data []byte
	var Obj struct {
		Status   string `json:"status"`
		HTTPCode int    `json:"code"`
		Message  string `json:"message"`
		Err      error  `json:"error"`
	}
	Obj.Message = msg
	Obj.HTTPCode = httpCode
	switch msgType {
	case ERR_MSG:
		Obj.Status = "FAILED"

	case MSG:
		Obj.Status = "SUCCESS"
	}
	data, _ = json.Marshal(Obj)
	return data, httpCode
}

func writeJSONMessage(msg string, msgType string, httpCode int, rd *RequestData) {
	d, code := jsonifyMessage(msg, msgType, httpCode)
	writeJSONResponse(d, code, rd)
}

func writeJSONStruct(v interface{}, code int, rd *RequestData) {
	d, err := json.Marshal(v)
	if err != nil {
		writeJSONMessage("Unable to marshal data. Err: "+err.Error(), ERR_MSG, http.StatusInternalServerError, rd)
		return
	}
	writeJSONResponse(d, code, rd)
}

func writeJSONResponse(d []byte, code int, rd *RequestData) {
	rd.l.LogAPIInfo(rd.r, time.Since(rd.Start).Seconds(), code)
	if code == http.StatusInternalServerError {
		rd.l.Info(rd.r.URL, "Status Code:", code, ", Response time:", time.Since(rd.Start), rd.r.URL, " Response:", string(d))
	} else {
		rd.l.Info(rd.r.URL, "Status Code:", code, ", Response time:", time.Since(rd.Start))
	}
	rd.w.Header().Set("Access-Control-Allow-Origin", "*")
	rd.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rd.w.WriteHeader(code)
	rd.w.Write(d)
}

func writeImageResponse(imageByte []byte, code int, rd *RequestData) {
	rd.l.LogAPIInfo(rd.r, time.Since(rd.Start).Seconds(), code)
	if code == http.StatusInternalServerError {
		rd.l.Info(rd.r.URL, "Status Code:", code, ", Response time:", time.Since(rd.Start), rd.r.URL, " Response:")
	} else {
		rd.l.Info(rd.r.URL, "Status Code:", code, ", Response time:", time.Since(rd.Start))
	}
	//rd.w.Header().Set("Access-Control-Allow-Origin", "*")
	//rd.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rd.w.WriteHeader(code)
	rd.w.Header().Set("Content-Type", "image/png")
	rd.w.Write(imageByte)

}

func writeJSONMessageWithData(msg string, msgType string, httpCode int, rd *RequestData, functionName string, requestData string) {
	d, code := jsonifyMessage(msg, msgType, httpCode)
	writeJSONResponseWithData(d, code, rd, functionName, requestData)
}

func writeJSONResponseWithData(d []byte, code int, rd *RequestData, functionName string, requestData string) {
	rd.l.LogAPIInfo(rd.r, time.Since(rd.Start).Seconds(), code)
	//if code == http.StatusInternalServerError {
	//	rd.l.Info("Status Code:", code, ", Response time:", time.Since(rd.Start), " Response:", string(d))
	//} else if code == http.StatusBadRequest {
	//	rd.l.Info("Status Code:", code, ", Response time:", time.Since(rd.Start), " Response:", string(d))
	//} else {
	///	rd.l.Info("Service name : ", functionName, "Request Data : ", requestData,
	//		"Response data : ", string(d), "Status Code : ", code, "Response time : ", time.Since(rd.Start))
	//}

	rd.l.Info("Service name : ", functionName, ", Request Data : ", requestData,
		", Response data : ", string(d), ", Status Code : ", code, ", Response time : ", time.Since(rd.Start))

	rd.w.Header().Set("Access-Control-Allow-Origin", "*")
	rd.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rd.w.WriteHeader(code)
	rd.w.Write(d)
}

func writeJSONStructWithData(v interface{}, code int, rd *RequestData, functionName string, requestData string) {
	d, err := json.Marshal(v)
	if err != nil {
		writeJSONMessageWithData("Unable to marshal data. Err: "+err.Error(), ERR_MSG, http.StatusInternalServerError, rd, functionName, requestData)
		return
	}
	writeJSONResponseWithData(d, code, rd, functionName, requestData)
}

func renderJSON(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		lg.Printf("ERROR: renderJson - %q\n", err)
	}
}

func getandConvertToInt64(query url.Values, str string) int64 {
	intData, _ := strconv.ParseInt(query.Get(str), 10, 64)
	return intData
}

func getandConvertToInt(query url.Values, str string) int {
	intData, _ := strconv.Atoi(query.Get(str))
	return intData
}

// parseJson parse data to model
func parseJSON(w http.ResponseWriter, body io.ReadCloser, model interface{}) bool {
	defer body.Close()

	b, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(b, model)

	if err != nil {
		e := &errs.Error{}
		e.Message = "Error in parsing json"
		e.Err = err
		renderERROR(w, e)
		return false
	}

	return true
}

func renderERROR(w http.ResponseWriter, err *errs.Error) {
	err.Set()
	renderJSON(w, err.Code, err)
}

func renderERRORV1(w http.ResponseWriter, err *dtos.ErrorData) {
	err.Set()
	renderJSON(w, err.Code, err)
}

func parseJSONWithError(w http.ResponseWriter, body io.ReadCloser, model interface{}) (bool, error) {
	defer body.Close()
	b, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(b, model)
	if err != nil {
		e := &errs.Error{}
		e.Message = "Error in parsing json"
		e.Err = err
		//renderERROR(w, e)
		return false, err
	}
	return true, nil
}

// GetIDFromParams function is use to get ID from request
func GetIDFromParams(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	params, _ := r.Context().Value("params").(httprouter.Params)
	idStr := params.ByName(key)
	id, err := strconv.ParseInt(idStr, 10, 64)
	isErr := true

	if err != nil {
		isErr = false
		e := &dtos.ErrorData{}
		e.Message = "Invalid ID"
		e.Code = http.StatusBadRequest
		e.Err = err
		renderERRORV1(w, e)
	}

	return id, isErr
}

var baseDirectory = os.Getenv("BASE_DIRECTORY")

func ViweImage(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	image := r.URL.Query().Get("image")

	if image == "" {
		rd.l.Error("ViweImage error: image parameter is empty")
		writeJSONMessage("image parameter is required", ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	// Get BASE_DIRECTORY fresh each time in case it changes
	currentBaseDirectory := os.Getenv("BASE_DIRECTORY")
	if currentBaseDirectory == "" {
		currentBaseDirectory = baseDirectory
	}
	if currentBaseDirectory == "" {
		rd.l.Error("ViweImage error: BASE_DIRECTORY environment variable is not set")
		writeJSONMessage("BASE_DIRECTORY environment variable is not set", ERR_MSG, http.StatusInternalServerError, rd)
		return
	}

	// Construct the full image path
	// image parameter from URL is already a relative path like "t_hub_document/vendor/VEN-000001/filename"
	// BASE_DIRECTORY should be something like "/app" or "/app/uploads"
	// So we join them to get the full path
	imagePath := filepath.Join(currentBaseDirectory, image)
	rd.l.Info("ViweImage - BASE_DIRECTORY: ", currentBaseDirectory, ", image param: ", image, ", full path: ", imagePath)

	// Check if file exists before trying to read
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		rd.l.Error("ViweImage error: file does not exist at path: ", imagePath)
		// Provide a more user-friendly error message
		writeJSONMessage("Image file not found. This may happen if the server was redeployed. Please re-upload the image.", ERR_MSG, http.StatusNotFound, rd)
		return
	}

	imageByte, err := os.ReadFile(imagePath)
	if err != nil {
		rd.l.Error("ViweImage error reading file: ", err, " at path: ", imagePath)
		writeJSONMessage(fmt.Sprintf("failed to read file: %s", err.Error()), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	// Determine content type based on file extension
	ext := filepath.Ext(imagePath)
	contentType := "image/png" // default
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".pdf":
		contentType = "application/pdf"
	}

	rd.w.Header().Set("Content-Type", contentType)
	writeImageResponse(imageByte, http.StatusOK, rd)
}

func createPrerequisite(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	tpmgt := commonsvc.New(rd.l, rd.dbConnMSSQL)
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("createPrerequisite", "orgID", orgID)
	if !isErr {
		return
	}
	keys := r.URL.Query()
	tripNum := keys.Get("trip_num")
	customer := keys.Get("customer")
	vendor := keys.Get("vendor")
	lodingpoints := keys.Get("loding_points")
	branch := keys.Get("branch")
	tripSheetType := keys.Get("trip_sheet_type")
	tripStatus := keys.Get("trip_status")
	customerCode := keys.Get("customer_code")
	vendorCode := keys.Get("vendor_code")
	podTrips := keys.Get("pod_trips")
	regularTrips := keys.Get("regular_trips")
	lrTrips := keys.Get("lr_trips")
	roles := keys.Get("roles")
	vehicleSizes := keys.Get("vehicle_sizes")
	declarationYear := keys.Get("declaration_year")

	websiteScreens := keys.Get("website_screens")
	permisstionLabel := keys.Get("permisstion_label")
	employees := keys.Get("employees")
	invoiceStatus := keys.Get("invoice_status")

	res, err := tpmgt.CreateTripPrerequisite(orgID, tripNum, customer, vendor, lodingpoints, branch,
		tripSheetType, tripStatus, customerCode, vendorCode, podTrips, regularTrips, lrTrips, roles,
		vehicleSizes, declarationYear, websiteScreens, permisstionLabel, employees, invoiceStatus)
	if err != nil {
		rd.l.Error("createTripPrerequisite error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
