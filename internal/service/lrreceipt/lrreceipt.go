package lrreceipt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"

	"go-transport-hub/utils"
)

var lodingPointCityCacheMap = cache.New(15*time.Minute, 30*time.Minute)

type LRReceipt struct {
	l            *log.Logger
	dbConnMSSQL  *mssqlcon.DBConn
	lrReceiptDao daos.LRReceiptDao
}

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *LRReceipt {
	return &LRReceipt{
		l:            l,
		dbConnMSSQL:  dbConnMSSQL,
		lrReceiptDao: daos.NewLRReceiptObj(l, dbConnMSSQL),
	}
}

func (lr *LRReceipt) CreateLRReceipt(lrReq dtos.LRReceiptReq) (*dtos.Messge, error) {

	err := lr.validateLR(lrReq)
	if err != nil {
		lr.l.Error("ERROR: CreateLRReceipt", err)
		return nil, err
	}

	lrReq.LRNumber = strings.ToUpper(lrReq.LRNumber)

	err1 := lr.lrReceiptDao.CreateLRReceipt(lrReq)
	if err1 != nil {
		lr.l.Error("CreateLRReceipt not saved ", lrReq.TripSheetNum, err1)
		return nil, err1
	}

	// Update Trip_sheet_header for is LR is
	errT := lr.lrReceiptDao.UpdateTripSheetHeader(lrReq.TripSheetID, 1)
	if errT != nil {
		lr.l.Error("UpdateTripSheetHeader not updated", lrReq.TripSheetNum, errT)
		return nil, errT
	}


	lr.l.Info("LR Receipt created successfully! : ", lrReq.TripSheetNum)

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("LR Receipt created successfully: %s", lrReq.TripSheetNum)
	return &response, nil
}

func (lr *LRReceipt) UpdateLR(lrId int64, updateLRReq dtos.LRReceiptUpdateReq) (*dtos.Messge, error) {

	err := lr.validateUpdateLR(updateLRReq)
	if err != nil {
		lr.l.Error("ERROR: UpdateLR", err)
		return nil, err
	}

	lrInfo, errV := lr.lrReceiptDao.GetLRRecord(lrId)
	if errV != nil {
		lr.l.Error("ERROR: pod not found", lrId, errV)
		return nil, errV
	}
	lr.l.Info("lrInfo: ******* ", utils.MustMarshal(lrInfo, "lrInfo"))

	updateLRReq.LRNumber = strings.ToUpper(updateLRReq.LRNumber)

	err1 := lr.lrReceiptDao.UpdateLR(lrId, updateLRReq)
	if err1 != nil {
		lr.l.Error("LR not updated", lrId, err1)
		return nil, err1
	}

	// }

	lr.l.Info("LR updated successfully! : ", lrId, updateLRReq.LRNumber, updateLRReq.TripSheetNum)

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("LR updated successfully: %s", updateLRReq.LRNumber)
	return &response, nil
}

func (mp *LRReceipt) removeDuplicates(ids []string) []string {
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

func (lr *LRReceipt) GetLRRecords(orgId int64, limit, offset, lrId, tripSheetNum, lrNumber, tripDate, searchText string) (*dtos.LRRecords, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}
	query := lr.lrReceiptDao.BuildWhereQuery(orgId, lrId, tripSheetNum, lrNumber, tripDate, searchText)
	res, tripSheetIds, errA := lr.lrReceiptDao.GetLRRecords(orgId, query, limitI, offsetI)
	if errA != nil {
		lr.l.Error("ERROR: GetBranchs", errA)
		return nil, errA
	}

	//
	tripSheetIdst := lr.removeDuplicates(tripSheetIds)
	lr.l.Info("tripSheetIdst: ", tripSheetIdst)

	if len(tripSheetIds) != 0 {
		idsCSV := strings.Join(tripSheetIds, ",")
		lr.l.Info("idsCSV commas: ", idsCSV)
		loadingPoints, err := lr.lrReceiptDao.GetLoadingUnloadingPoints(idsCSV)
		if err != nil {
			lr.l.Error("ERROR: loadingPoints", err)
			return nil, err
		}
		//lr.l.Info("loadingPoints: ", utils.MustMarshal(loadingPoints, "loadingPoints"))

		loadingMap, unloadingMap := lr.BuildCityMaps(loadingPoints, lodingPointCityCacheMap)

		//lr.l.Info("loadingMap: ******* ", utils.MustMarshal(loadingMap, "loadingMap"))
		//lr.l.Info("unloadingMap: ******* ", utils.MustMarshal(unloadingMap, "unloadingMap"))

		for i := range *res {
			lrRes := &(*res)[i]
			tripSheetID := lrRes.TripSheetID
			// Set LoadingPointIDs if present in loadingMap
			if val, ok := loadingMap[tripSheetID]; ok {
				lrRes.FromAddress = val
			}

			// Set UnLoadingPointIDs if present in unloadingMap
			if val, ok := unloadingMap[tripSheetID]; ok {
				lrRes.ToAddress = val
			}
		}
	}

	lrResponse := dtos.LRRecords{}
	lrResponse.LRRecords = res
	lrResponse.Total = lr.lrReceiptDao.GetTotalCount(query)
	lrResponse.Limit = limitI
	lrResponse.OffSet = offsetI

	return &lrResponse, nil
}

// func (mp *ManagePod) removeDuplicates(ids []string) []string {
// 	seen := make(map[string]bool)
// 	var result []string
// 	for _, id := range ids {
// 		if !seen[id] {
// 			seen[id] = true
// 			result = append(result, id)
// 		}
// 	}
// 	return result
// }

// func (mp *ManagePod) UpdatePODStatus(status string, podId int64) (*dtos.Messge, error) {

// 	if podId == 0 {
// 		mp.l.Error("unknown pod Id ", podId)
// 		return nil, errors.New("unknown pod Id")
// 	}

// 	tripTypes := []string{constant.STATUS_DELIVERED, constant.STATUS_CLOSED, constant.STATUS_DELETED}

// 	exists := false
// 	for _, typeT := range tripTypes {
// 		if typeT == status {
// 			exists = true
// 			break
// 		}
// 	}
// 	if !exists {
// 		return nil, fmt.Errorf("trip status should be %s,%s,%s ", constant.STATUS_DELIVERED, constant.STATUS_CLOSED, constant.STATUS_DELETED)
// 	}

// 	podInfo, errV := mp.managePodDao.GetManagePod(podId)
// 	if errV != nil {
// 		mp.l.Error("ERROR: pod not found", podId, errV)
// 		return nil, errV
// 	}
// 	mp.l.Info("PodStatus existing: ", podInfo.PodStatus, " req: ", status, " podInfo: *******: ", utils.MustMarshal(podInfo, "podInfo"))

// 	//trip_sheet_num,lr_number
// 	statusRes := dtos.Messge{}
// 	statusRes.Message = "pod updated successfully : " + status

// 	if status == constant.STATUS_DELETED {
// 		lrNumber := fmt.Sprintf("deleted_%s", podInfo.LRNumber)
// 		if strings.Contains(podInfo.LRNumber, "deleted") {
// 			lrNumber = podInfo.LRNumber
// 		}
// 		errU := mp.managePodDao.UpdatePODStatusDelete(podId, lrNumber, status)
// 		if errU != nil {
// 			mp.l.Error("ERROR: UpdatePODStatusDelete ", errU)
// 			return nil, errU
// 		}
// 		return &statusRes, nil
// 	}

// 	errU := mp.managePodDao.UpdatePODStatus(podId, status)
// 	if errU != nil {
// 		mp.l.Error("ERROR: UpdatePODStatus ", errU)
// 		return nil, errU
// 	}

// 	errT := mp.managePodDao.UpdateTripSheetStatus(podInfo.TripSheetID, status)
// 	if errT != nil {
// 		mp.l.Error("ERROR: UpdateTripSheetStatus ", errT)
// 		return nil, errT
// 	}

// 	return &statusRes, nil
// }

// func (ul *ManagePod) UploadPodDoc(tripSheetId int64, imageFor string, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.Messge, error) {

// 	if imageFor == "" {
// 		return nil, errors.New("imageFor should not be empty.\n pod_doc")
// 	}

// 	imageTypes := [...]string{
// 		"pod_doc",
// 	}

// 	exists := false
// 	for _, imgType := range imageTypes {
// 		if imgType == imageFor {
// 			exists = true
// 			break
// 		}
// 	}
// 	if !exists {
// 		return nil, errors.New("imageFor not valid pod_doc")
// 	}

// 	bool := utils.CheckSizeImage(fileHeader, 10000, ul.l)
// 	if !bool {
// 		ul.l.Error("image size issue ")
// 		return nil, errors.New("image size issue ")
// 	}

// 	baseDirectory := os.Getenv("BASE_DIRECTORY")
// 	uploadPath := os.Getenv("IMAGE_DIRECTORY")
// 	if uploadPath == "" || baseDirectory == "" {
// 		ul.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
// 		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
// 	}

// 	imageDirectory := filepath.Join(uploadPath, constant.POD_DIRECTORY, strconv.Itoa(int(tripSheetId)))
// 	fullPath := filepath.Join(baseDirectory, imageDirectory)
// 	ul.l.Infof("pod tripsheet id: %v imageDirectory: %s, fullPath: %s", tripSheetId, imageDirectory, fullPath)

// 	err := os.MkdirAll(fullPath, os.ModePerm) // os.ModePerm sets permissions to 0777
// 	if err != nil {
// 		ul.l.Error("ERROR: MkdirAll ", fullPath, err)
// 		return nil, err
// 	}

// 	extension := strings.Split(fileHeader.Filename, ".")
// 	lengthExt := len(extension)

// 	imageName := fmt.Sprintf("%v_%v.%v", imageFor, tripSheetId, extension[lengthExt-1])
// 	podImageFullPath := filepath.Join(fullPath, imageName)
// 	ul.l.Info("pod doc FullPath: ", imageFor, podImageFullPath)

// 	out, err := os.Create(podImageFullPath)
// 	if err != nil {
// 		if utils.CheckFileExists(podImageFullPath) {
// 			ul.l.Error("updating  is already exitis: ", tripSheetId, uploadPath, err)
// 		} else {
// 			ul.l.Error("podImageFullPath create error: ", tripSheetId, uploadPath, err)
// 			defer out.Close()
// 			return nil, err
// 		}
// 	}
// 	defer out.Close()

// 	_, err = io.Copy(out, file)
// 	if err != nil {
// 		ul.l.Error("pod upload Copy error: ", tripSheetId, uploadPath, err)
// 		return nil, err
// 	}

// 	imageDirectory = filepath.Join(imageDirectory, imageName)
// 	ul.l.Info("##### table to be stored path: ", imageFor, imageDirectory)
// 	updateQuery := fmt.Sprintf(`UPDATE manage_pod SET %v = '%v' WHERE trip_sheet_id = '%v'`, imageFor, imageDirectory, tripSheetId)
// 	errU := ul.managePodDao.UpdatePodImagePath(updateQuery)
// 	if errU != nil {
// 		ul.l.Error("ERROR: UploadPodDoc", tripSheetId, errU)
// 		return nil, errU
// 	}

// 	ul.l.Info("Image uploaded successfully: ", tripSheetId)
// 	roleResponse := dtos.Messge{}
// 	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", tripSheetId)
// 	return &roleResponse, nil
// }
