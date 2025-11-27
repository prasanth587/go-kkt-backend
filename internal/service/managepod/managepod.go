package managepod

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
	"go-transport-hub/utils"
)

type ManagePod struct {
	l            *log.Logger
	dbConnMSSQL  *mssqlcon.DBConn
	managePodDao daos.ManagePodDao
}

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *ManagePod {
	return &ManagePod{
		l:            l,
		dbConnMSSQL:  dbConnMSSQL,
		managePodDao: daos.NewManagePodObj(l, dbConnMSSQL),
	}
}

var lodingPointCityCacheMap = cache.New(15*time.Minute, 30*time.Minute)

func (mp *ManagePod) CreateManagePod(podReq dtos.ManagePodReq) (*dtos.Messge, error) {

	err := mp.validatePOD(podReq)
	if err != nil {
		mp.l.Error("ERROR: CreateManagePod", err)
		return nil, err
	}

	podReq.PodStatus = constant.STATUS_DELIVERED
	//podReq.TripType = constant.TRIP_TYPE_POD
	podReq.LRNumber = strings.ToUpper(podReq.LRNumber)

	err1 := mp.managePodDao.CreateManagePod(podReq)
	if err1 != nil {
		mp.l.Error("CreateManagePod not saved ", podReq.TripSheetNum, err1)
		return nil, err1
	}

	errT := mp.managePodDao.UpdateTripSheetStatus(podReq.TripSheetID, podReq.PodStatus, getStatusDateUpdateQuery(podReq.PodStatus))
	if errT != nil {
		mp.l.Error("ERROR: UpdateTripSheetStatus ", errT)
		return nil, errT
	}

	mp.l.Info("POD created successfully! : ", podReq.TripSheetNum)

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("POD created successfully: %s", podReq.TripSheetNum)
	return &response, nil
}

func (mp *ManagePod) GetPods(orgId int64, limit, offset, podId, podStatus, tripType, searchText string) (*dtos.ManagePods, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}
	query := mp.managePodDao.BuildWhereQuery(orgId, podId, podStatus, tripType, searchText)
	res, tripSheetIds, errA := mp.managePodDao.GetManagePods(orgId, query, limitI, offsetI)
	if errA != nil {
		mp.l.Error("ERROR: GetBranchs", errA)
		return nil, errA
	}
	tripSheetIdst := mp.removeDuplicates(tripSheetIds)
	mp.l.Info("tripSheetIdst: ", tripSheetIdst)

	if len(tripSheetIds) != 0 {
		idsCSV := strings.Join(tripSheetIds, ",")
		mp.l.Info("idsCSV comma: ", idsCSV)
		loadingPoints, err := mp.managePodDao.GetLoadingUnloadingPoints(idsCSV)
		if err != nil {
			mp.l.Error("ERROR: loadingPoints", err)
			return nil, err
		}
		mp.l.Info("loadingPoints: ", utils.MustMarshal(loadingPoints, "loadingPoints"))

		loadingMap, unloadingMap := mp.BuildCityMaps(loadingPoints, lodingPointCityCacheMap)

		mp.l.Info("loadingMap: ******* ", utils.MustMarshal(loadingMap, "loadingMap"))
		mp.l.Info("unloadingMap: ******* ", utils.MustMarshal(unloadingMap, "unloadingMap"))

		//

		for i := range *res {
			pod := &(*res)[i]
			tripSheetID := pod.TripSheetID
			// Set LoadingPointIDs if present in loadingMap
			if val, ok := loadingMap[tripSheetID]; ok {
				pod.LoadingPointIDs = val
			}

			// Set UnLoadingPointIDs if present in unloadingMap
			if val, ok := unloadingMap[tripSheetID]; ok {
				pod.UnLoadingPointIDs = val
			}
		}
	}

	managePods := dtos.ManagePods{}
	managePods.ManagePod = res
	managePods.Total = mp.managePodDao.GetTotalCount(query)
	managePods.Limit = limitI
	managePods.OffSet = offsetI

	return &managePods, nil
}

func (mp *ManagePod) removeDuplicates(ids []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}

func (mp *ManagePod) UpdatePod(podId int64, updateReq dtos.UpdateManagePodReq) (*dtos.Messge, error) {

	err := mp.validateUpdatePOD(updateReq)
	if err != nil {
		mp.l.Error("ERROR: UpdatePod", err)
		return nil, err
	}

	podInfo, errV := mp.managePodDao.GetManagePod(podId)
	if errV != nil {
		mp.l.Error("ERROR: pod not found", podId, errV)
		return nil, errV
	}
	mp.l.Info("podInfo: ******* ", utils.MustMarshal(podInfo, "podInfo"))
	updateReq.PodStatus = podInfo.PodStatus

	updateReq.LRNumber = strings.ToUpper(updateReq.LRNumber)

	err1 := mp.managePodDao.UpdatePod(podId, updateReq)
	if err1 != nil {
		mp.l.Error("pod not updated ", podId, err1)
		return nil, err1
	}

	errT := mp.managePodDao.UpdateTripSheetStatus(podInfo.TripSheetID, updateReq.PodStatus, getStatusDateUpdateQuery(updateReq.PodStatus))
	if errT != nil {
		mp.l.Error("ERROR: UpdateTripSheetStatus ", errT)
		return nil, errT
	}

	mp.l.Info("pod updated successfully! : ", updateReq.LRNumber, updateReq.TripSheetNum)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("pod updated successfully: %s", updateReq.TripSheetNum)
	return &roleResponse, nil
}

func getStatusDateUpdateQuery(status string) string {
	if column, ok := constant.STATUS_TO_COLUMN[status]; ok {
		return fmt.Sprintf(", %s = '%v'", column, utils.GetCurrentDatetimeStr())
	}
	return ""
}

func (mp *ManagePod) UpdatePODStatus(status string, podId int64) (*dtos.Messge, error) {

	if podId == 0 {
		mp.l.Error("unknown pod Id ", podId)
		return nil, errors.New("unknown pod Id")
	}

	tripTypes := []string{constant.STATUS_DELIVERED, constant.STATUS_CLOSED, constant.STATUS_DELETED}

	exists := false
	for _, typeT := range tripTypes {
		if typeT == status {
			exists = true
			break
		}
	}
	if !exists {
		return nil, fmt.Errorf("trip status should be %s,%s,%s ", constant.STATUS_DELIVERED, constant.STATUS_CLOSED, constant.STATUS_DELETED)
	}

	podInfo, errV := mp.managePodDao.GetManagePod(podId)
	if errV != nil {
		mp.l.Error("ERROR: pod not found", podId, errV)
		return nil, errV
	}
	mp.l.Info("PodStatus existing: ", podInfo.PodStatus, " req: ", status, " podInfo: *******: ", utils.MustMarshal(podInfo, "podInfo"))

	//trip_sheet_num,lr_number
	statusRes := dtos.Messge{}
	statusRes.Message = "pod updated successfully : " + status

	if status == constant.STATUS_DELETED {
		lrNumber := fmt.Sprintf("deleted_%s", podInfo.LRNumber)
		if strings.Contains(podInfo.LRNumber, "deleted") {
			lrNumber = podInfo.LRNumber
		}
		errU := mp.managePodDao.UpdatePODStatusDelete(podId, lrNumber, status)
		if errU != nil {
			mp.l.Error("ERROR: UpdatePODStatusDelete ", errU)
			return nil, errU
		}
		return &statusRes, nil
	}

	errU := mp.managePodDao.UpdatePODStatus(podId, status)
	if errU != nil {
		mp.l.Error("ERROR: UpdatePODStatus ", errU)
		return nil, errU
	}

	errT := mp.managePodDao.UpdateTripSheetStatus(podInfo.TripSheetID, status, getStatusDateUpdateQuery(status))
	if errT != nil {
		mp.l.Error("ERROR: UpdateTripSheetStatus ", errT)
		return nil, errT
	}

	return &statusRes, nil
}

func (ul *ManagePod) UploadPodDoc(tripSheetId int64, imageFor string, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.UploadManagePODResponse, error) {

	if imageFor == "" {
		return nil, errors.New("imageFor should not be empty.\n pod_doc")
	}

	imageTypes := [...]string{
		"pod_doc",
	}

	exists := false
	for _, imgType := range imageTypes {
		if imgType == imageFor {
			exists = true
			break
		}
	}
	if !exists {
		return nil, errors.New("imageFor not valid pod_doc")
	}

	bool := utils.CheckSizeImage(fileHeader, 10000, ul.l)
	if !bool {
		ul.l.Error("image size issue ")
		return nil, errors.New("image size issue ")
	}

	baseDirectory := os.Getenv("BASE_DIRECTORY")
	uploadPath := os.Getenv("IMAGE_DIRECTORY")
	if uploadPath == "" || baseDirectory == "" {
		ul.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
	}

	imageDirectory := filepath.Join(uploadPath, constant.POD_DIRECTORY, strconv.Itoa(int(tripSheetId)))
	fullPath := filepath.Join(baseDirectory, imageDirectory)
	ul.l.Infof("pod tripsheet id: %v imageDirectory: %s, fullPath: %s", tripSheetId, imageDirectory, fullPath)

	err := os.MkdirAll(fullPath, os.ModePerm) // os.ModePerm sets permissions to 0777
	if err != nil {
		ul.l.Error("ERROR: MkdirAll ", fullPath, err)
		return nil, err
	}

	extension := strings.Split(fileHeader.Filename, ".")
	lengthExt := len(extension)

	imageName := fmt.Sprintf("%v_%v.%v", imageFor, tripSheetId, extension[lengthExt-1])
	podImageFullPath := filepath.Join(fullPath, imageName)
	ul.l.Info("pod doc FullPath: ", imageFor, podImageFullPath)

	out, err := os.Create(podImageFullPath)
	if err != nil {
		if utils.CheckFileExists(podImageFullPath) {
			ul.l.Error("updating  is already exitis: ", tripSheetId, uploadPath, err)
		} else {
			ul.l.Error("podImageFullPath create error: ", tripSheetId, uploadPath, err)
			defer out.Close()
			return nil, err
		}
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ul.l.Error("pod upload Copy error: ", tripSheetId, uploadPath, err)
		return nil, err
	}

	imageDirectory = filepath.Join(imageDirectory, imageName)
	ul.l.Info("##### table to be stored path: ", imageFor, imageDirectory)
	updateQuery := fmt.Sprintf(`UPDATE manage_pod SET %v = '%v' WHERE trip_sheet_id = '%v'`, imageFor, imageDirectory, tripSheetId)
	errU := ul.managePodDao.UpdatePodImagePath(updateQuery)
	if errU != nil {
		ul.l.Error("ERROR: UploadPodDoc", tripSheetId, errU)
		return nil, errU
	}

	ul.l.Info("Image uploaded successfully: ", tripSheetId)
	roleResponse := dtos.UploadManagePODResponse{}
	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", tripSheetId)
	roleResponse.ImagePath = imageDirectory
	roleResponse.TripSheetNumber = tripSheetId
	return &roleResponse, nil
}
