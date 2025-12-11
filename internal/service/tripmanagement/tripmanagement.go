package tripmanagement

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
	// "go-transport-hub/internal/service/notification" // Temporarily disabled
	"go-transport-hub/utils"
)

type TripSheetObj struct {
	l            *log.Logger
	dbConnMSSQL  *mssqlcon.DBConn
	tripSheetDao daos.TripSheetDao
}

var (
	ErrUnableToPingDB = errors.New("unable to ping database")
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

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *TripSheetObj {
	return &TripSheetObj{
		l:            l,
		dbConnMSSQL:  dbConnMSSQL,
		tripSheetDao: daos.NewTripSheetObj(l, dbConnMSSQL),
	}
}

func (trp *TripSheetObj) CreateTripSheetHeader(tripSheetReq dtos.CreateTripSheetHeader) (*dtos.Messge, error) {

	err := trp.validateTrip(tripSheetReq)
	if err != nil {
		trp.l.Error("ERROR: CreateTripSheetHeader", err)
		return nil, err
	}
	trp.l.Info("trip sheet type: ", tripSheetReq.TripSheetType)

	// Parse VehicleSize to VehicleSizeID if VehicleSizeID is not set but VehicleSize is
	if tripSheetReq.VehicleSizeID == 0 && tripSheetReq.VehicleSize != "" {
		if vehicleSizeID, err := strconv.ParseInt(tripSheetReq.VehicleSize, 10, 64); err == nil {
			tripSheetReq.VehicleSizeID = vehicleSizeID
		} else {
			trp.l.Error("ERROR: Failed to parse VehicleSize to VehicleSizeID: ", tripSheetReq.VehicleSize, err)
			return nil, fmt.Errorf("invalid vehicle_size value: %s", tripSheetReq.VehicleSize)
		}
	}

	createTripSheetQuery := ""
	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_SCHEDULED_TRIP) ||
		strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_ADHOC_TRIP) {
		tripSheetReq.ZonalName = "_"
	}
	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LINE_HUAL_SCHEDULED_TRIP) ||
		strings.EqualFold(tripSheetReq.TripSheetType, constant.LINE_HUAL_ADHOC_TRIP) {
		tripSheetReq.LoadHoursType = "_"
	}
	tripSheetReq.LoadStatus = constant.STATUS_CREATED
	//tripSheetReq.OpenTripDateTime = strings.Replace(tripSheetReq.OpenTripDateTime, "T", " ", 1)

	createTripSheetQuery = trp.tripSheetDao.CreateTripSheetBuildQuery(tripSheetReq)

	tripSheetId, errD := trp.tripSheetDao.CreateTripSheet(tripSheetReq.TripSheetNum, createTripSheetQuery)
	if errD != nil {
		trp.l.Error("ERROR: CreateTripSheet", errD)
		return nil, errD
	}
	trp.l.Info("CreateTripSheet: ", tripSheetId)

	// Saving loading points
	for _, loadingPointId := range tripSheetReq.LoadingPointIDs {
		errD := trp.tripSheetDao.SaveTripSheetLoadingPoint(uint64(tripSheetId), loadingPointId, constant.LOADING_POINT)
		if errD != nil {
			trp.l.Error("ERROR: SaveTripSheetLoadingPoint", errD)
			return nil, errD
		}
	}

	// Saving unloading points
	for _, unLoadingPointId := range tripSheetReq.UnLoadingPointIDs {
		errD := trp.tripSheetDao.SaveTripSheetLoadingPoint(uint64(tripSheetId), unLoadingPointId, constant.UN_LOADING_POINT)
		if errD != nil {
			trp.l.Error("ERROR: SaveTripSheetUnLoadingPoint", errD)
			return nil, errD
		}
	}

	// Send notification for trip sheet creation
	// Temporarily disabled
	// notificationSvc := notification.New(trp.l, trp.dbConnMSSQL)
	// if err := notificationSvc.NotifyTripSheetCreated(int64(tripSheetReq.OrgId), tripSheetId, tripSheetReq.TripSheetNum, tripSheetReq.UserLoginId); err != nil {
	// 	trp.l.Error("ERROR: Failed to send trip creation notification: ", err)
	// 	// Don't fail the request if notification fails
	// }

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("Trip Sheet successfully: %s", tripSheetReq.TripSheetNum)
	return &response, nil
}

func (br *TripSheetObj) GetTripSheets(orgId int64, tripSheetId string, limit string, offset, tripStatus, tripSearchTex, fromDate, toDate, podRequired, podReceived string) (*dtos.TripSheetRes, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	whereQuery := br.tripSheetDao.BuildWhereQuery(orgId, tripSheetId, tripStatus, tripSearchTex, fromDate, toDate, podRequired, podReceived)
	res, errA := br.tripSheetDao.GetTripSheets(orgId, whereQuery, limitI, offsetI)
	if errA != nil {
		br.l.Error("ERROR: GetTripSheets", errA)
		return nil, errA
	}

	loadingpointEntries := dtos.TripSheetRes{}
	loadingpointEntries.TripSheet = res
	loadingpointEntries.Total = br.tripSheetDao.GetTotalCount(whereQuery)
	loadingpointEntries.Limit = limitI
	loadingpointEntries.OffSet = offsetI
	return &loadingpointEntries, nil
}

func (br *TripSheetObj) GetTripStats(orgId int64, tripSheetId string, fromDate, toDate, tripStatus, tripSearchText, podRequired, podReceived string) (*dtos.TripStatsRes, error) {

	whereQuery := br.tripSheetDao.BuildWhereQuery(orgId, tripSheetId, tripStatus, tripSearchText, fromDate, toDate, podRequired, podReceived)
	res, errA := br.tripSheetDao.GetTripStats(whereQuery)
	if errA != nil {
		br.l.Error("ERROR: GetTripSheets", errA)
		return nil, errA
	}
	count := len(*res)

	tripStats := dtos.TripStatsRes{}
	tripStats.TripStats = res
	tripStats.Total = count

	return &tripStats, nil
}

func (trp *TripSheetObj) GetTripSheetChallanInfo(tripSheetId int64, loginId string) (*dtos.TripStatsChallanInfo, error) {

	tripSheetInfo, errV := trp.tripSheetDao.GetTripSheet(tripSheetId)
	if errV != nil {
		trp.l.Error("ERROR: TripSheet not found", tripSheetInfo, errV)
		return nil, errV
	}
	jsonBytes, _ := json.Marshal(tripSheetInfo)
	trp.l.Info("TripSheet INFO: ******* ", string(jsonBytes))

	tripStats := dtos.TripStatsChallanInfo{}
	tripStats.AddressTop = "NO.5, SHANMUGA NAGAR, ATTANTHANGAL, REDHILLS, CHENNAI-600052"
	tripStats.Title = "KK TRANSPORT"
	tripStats.Month = trp.getCurrentMonth()
	tripStats.SerialNumber = fmt.Sprintf("KKT%v", time.Now().Format("20060102150405"))
	tripStats.LoadingDate = trp.loadingDate(tripSheetInfo.OpenTripDateTime)
	tripStats.LRNumber = tripSheetInfo.LRNumber
	tripStats.VehicleNumber = tripSheetInfo.VehicleNumber
	tripStats.VehicleSize = trp.GetVehicleSizeType(tripSheetInfo)

	var fromCityCodes, toCityCodes []string
	for _, p := range *tripSheetInfo.LoadingPointIDs {
		fromCityCodes = append(fromCityCodes, p.CityCode)
	}
	for _, p := range *tripSheetInfo.UnLoadingPointIDs {
		toCityCodes = append(toCityCodes, p.CityCode)
	}
	tripStats.FromLocations = strings.Join(fromCityCodes, " - ")
	tripStats.ToLocations = strings.Join(toCityCodes, " - ")

	tripStats.ContactNumber = tripSheetInfo.MobileNumber
	tripStats.DriverName = tripSheetInfo.DriverName

	if tripSheetInfo.Vendor != nil {
		tripStats.TransportName = tripSheetInfo.Vendor.VendorName
	}

	tripStats.VehicleHire = tripSheetInfo.VendorTotalHire
	tripStats.VehicleAdvance = tripSheetInfo.VendorAdvance
	tripStats.Mamul = tripSheetInfo.VendorMonul
	tripStats.Balance1 = tripStats.VehicleHire - (tripStats.VehicleAdvance + tripStats.Mamul)

	tripStats.LoadingUNCharges = tripSheetInfo.VendorLoadUnLoadAmount
	tripStats.HaltingCharges = tripSheetInfo.VendorHaltingPaid
	tripStats.Balance2 = (tripStats.Balance1 + tripStats.LoadingUNCharges + tripStats.HaltingCharges)

	tripStats.OfficeBalance3 = tripStats.Balance2
	tripStats.OfficeCommission = tripSheetInfo.VendorCommission
	tripStats.OfficeTotal = tripStats.Balance2 - tripSheetInfo.VendorCommission

	loginDetails, errU := trp.tripSheetDao.GetUserLoginDetails(loginId)
	if errU != nil {
		trp.l.Error("ERROR: GetUserLoginDetails", tripSheetId, loginId, errU)
	}

	// those are needs to be discussed.....
	tripStats.ReceivedDate = tripSheetInfo.CustomerPaymentReceivedDate
	tripStats.CompanyName = "KK Transport"
	tripStats.IssuedBy = "KK Transport"
	tripStats.OfficeApprovedBy = fmt.Sprintf("%s %s", loginDetails.FirstName, loginDetails.LastName)
	tripStats.OfficeAuthorizedBy = "User root"
	tripStats.Notes = ""

	trp.l.Info("loginId ", loginId)

	return &tripStats, nil
}

func (trp *TripSheetObj) getCurrentMonth() string {
	return time.Now().Format("January 2006")
}

func (trp *TripSheetObj) loadingDate(loadingDate string) string {
	t, err := time.Parse("2006-01-02 15:04", loadingDate)
	if err != nil {
		trp.l.Error("loadingDate error ", loadingDate, err)
		return utils.GetCurrentDateStr()
	}
	return t.Format("02-01-2006")
}

func (trp *TripSheetObj) GetVehicleSizeType(tripSheetInfo *dtos.TripSheet) string {

	vehicleType, err := trp.tripSheetDao.GetVehicleSizeType(tripSheetInfo.VehicleSizeID)
	if err != nil {
		trp.l.Error("ERROR vehicleType error ", tripSheetInfo.VehicleSizeID, err)
		return tripSheetInfo.VehicleSize
	}
	return fmt.Sprintf("%s - %s", vehicleType.VehicleSize, vehicleType.VehicleType)
}
