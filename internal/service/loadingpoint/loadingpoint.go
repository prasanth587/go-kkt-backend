package loadingpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
	"go-transport-hub/internal/service/notification"
)

type LoadingPointObj struct {
	l               *log.Logger
	dbConnMSSQL     *mssqlcon.DBConn
	loadingPointDao daos.LoadingPointDao
}

var (
	ErrUnableToPingDB = errors.New("Unable to ping database")
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

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *LoadingPointObj {
	return &LoadingPointObj{
		l:               l,
		dbConnMSSQL:     dbConnMSSQL,
		loadingPointDao: daos.NewLoadingPointObj(l, dbConnMSSQL),
	}
}

func (lp *LoadingPointObj) CreateLoadingPoint(loadingPointReq dtos.LoadingPointReq) (*dtos.Messge, error) {

	if loadingPointReq.BranchID == 0 {
		lp.l.Error("Error loadingpoint name should not empty")
		return nil, errors.New("loadingpoint name should not empty")
	}
	if loadingPointReq.CityCode == "" {
		lp.l.Error("Error CityCode: city code should not empty")
		return nil, errors.New("city code should not empty")
	}
	if loadingPointReq.CityName == "" {
		lp.l.Error("Error CityName: city name should not empty")
		return nil, errors.New("city name should not empty")
	}
	if loadingPointReq.State == "" {
		lp.l.Error("Error State: state should not empty")
		return nil, errors.New("state should not empty")
	}
	if loadingPointReq.Country == "" {
		lp.l.Error("Error Country: country should not empty")
		return nil, errors.New("country should not empty")
	}
	if loadingPointReq.MapLink == "" {
		lp.l.Error("Error MapLink: mapLink should not empty")
		return nil, errors.New("mapLink should not empty")
	}
	loadingPointReq.CityCode = strings.ToUpper(loadingPointReq.CityCode)

	loadingPointReq.IsActive = 1

	err1 := lp.loadingPointDao.CreateLoadingPoint(loadingPointReq)
	if err1 != nil {
		lp.l.Error("loadingpoint not saved ", loadingPointReq.CityCode, err1)
		return nil, err1
	}

	// Get loading point ID after creation for notification
	loadingPoints, errG := lp.loadingPointDao.GetLoadingPoints(loadingPointReq.OrgID, 1, 0, "")
	if errG == nil && loadingPoints != nil && len(*loadingPoints) > 0 {
		// Find the newly created loading point by city code
		for _, lpItem := range *loadingPoints {
			if lpItem.CityCode == loadingPointReq.CityCode {
				// Send notification for loading point creation
				// Temporarily disabled
				// notificationSvc := notification.New(lp.l, lp.dbConnMSSQL)
				// if err := notificationSvc.NotifyLoadingPointCreated(int64(loadingPointReq.OrgID), lpItem.LoadingPointID, loadingPointReq.CityName); err != nil {
				// 	lp.l.Error("ERROR: Failed to send loading point creation notification: ", err)
				// 	// Don't fail the request if notification fails
				// }
				break
			}
		}
	}

	lp.l.Info("Loading/unloading point created successfully! : ", loadingPointReq.CityCode)

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("LoadingPoint created successfully!: %s", loadingPointReq.CityCode)
	return &response, nil
}

func (br *LoadingPointObj) GetLoadingPoint(orgId int64, limit, offset, searchText, loadingPointId string) (*dtos.LoadingPointes, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	whereQuery := br.loadingPointDao.BuildWhereQuery(orgId, searchText, loadingPointId)

	res, errA := br.loadingPointDao.GetLoadingPoints(orgId, limitI, offsetI, whereQuery)
	if errA != nil {
		br.l.Error("ERROR: GetLoadingPoints", errA)
		return nil, errA
	}

	loadingpointEntries := dtos.LoadingPointes{}
	loadingpointEntries.LoadingPointEntries = res
	loadingpointEntries.Total = br.loadingPointDao.GetTotalCount(whereQuery)
	loadingpointEntries.Limit = limitI
	loadingpointEntries.OffSet = offsetI
	return &loadingpointEntries, nil
}

func (lp *LoadingPointObj) UpdateLoadingPoint(loadingPointId int64, loadingPointReq dtos.LoadingPointUpdate) (*dtos.Messge, error) {

	if loadingPointReq.BranchID == 0 {
		lp.l.Error("Error loadingpoint name should not empty")
		return nil, errors.New("loadingpoint name should not empty")
	}
	if loadingPointReq.CityCode == "" {
		lp.l.Error("Error CityCode: city code should not empty")
		return nil, errors.New("city code should not empty")
	}
	if loadingPointReq.CityName == "" {
		lp.l.Error("Error CityName: city name should not empty")
		return nil, errors.New("city name should not empty")
	}
	if loadingPointReq.State == "" {
		lp.l.Error("Error State: state should not empty")
		return nil, errors.New("state should not empty")
	}
	if loadingPointReq.Country == "" {
		lp.l.Error("Error Country: country should not empty")
		return nil, errors.New("country should not empty")
	}
	if loadingPointReq.MapLink == "" {
		lp.l.Error("Error MapLink: mapLink should not empty")
		return nil, errors.New("mapLink should not empty")
	}
	loadingPointReq.CityCode = strings.ToUpper(loadingPointReq.CityCode)

	loadingpointInfo, errV := lp.loadingPointDao.GetLoadingPoint(loadingPointId)
	if errV != nil {
		lp.l.Error("ERROR: loadingpoint not found", loadingPointId, errV)
		return nil, errV
	}
	jsonBytes, _ := json.Marshal(loadingpointInfo)
	lp.l.Info("GetLoadingPoint: ******* ", string(jsonBytes))
	loadingPointReq.IsActive = loadingpointInfo.IsActive

	err1 := lp.loadingPointDao.UpdateLoadingPoint(loadingPointId, loadingPointReq)
	if err1 != nil {
		lp.l.Error("loadingpoint not updated ", loadingPointReq.CityName, err1)
		return nil, err1
	}

	// Send notification for loading point update
	// Temporarily disabled
	// notificationSvc := notification.New(lp.l, lp.dbConnMSSQL)
	// if err := notificationSvc.NotifyLoadingPointUpdated(int64(loadingpointInfo.OrgID), loadingPointId, loadingPointReq.CityName); err != nil {
	// 	lp.l.Error("ERROR: Failed to send loading point update notification: ", err)
	// 	// Don't fail the request if notification fails
	// }

	lp.l.Info("loadingpoint updated successfully! : ", loadingPointReq.CityName)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("Loadingpoint updated successfully!: %s", loadingPointReq.CityName)
	return &roleResponse, nil
}

func (ul *LoadingPointObj) UpdateloadingpointActiveStatus(isA string, loadingpointId int64) (*dtos.Messge, error) {

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid status update")
	}
	if loadingpointId == 0 {
		ul.l.Error("unknown loadingpoint", loadingpointId)
		return nil, errors.New("unknown loadingpoint")
	}
	statusRes := dtos.Messge{}
	loadingpointInfo, errA := ul.loadingPointDao.GetLoadingPoint(loadingpointId)
	if errA != nil {
		ul.l.Error("ERROR: loadingpoint not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(loadingpointInfo)
	ul.l.Info("loadingpointInfo: ", string(jsonBytes))

	errU := ul.loadingPointDao.UpdateBranchActiveStatus(loadingpointId, isActive)
	if errU != nil {
		ul.l.Error("ERROR: UpdateloadingpointActiveStatus ", errU)
		return nil, errU
	}
	statusRes.Message = "loadingpoint updated successfully : " + loadingpointInfo.CityName
	return &statusRes, nil
}
