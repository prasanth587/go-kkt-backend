package tripmanagement

import (
	"encoding/json"
	"fmt"
	"slices"

	"go-transport-hub/constant"
	"go-transport-hub/dtos"

	"go-transport-hub/utils"
)

func (trp *TripSheetObj) UpdateTripSheetHeader(tripSheetId int64, tripSheetUpdateReq dtos.UpdateTripSheetHeader) (*dtos.Messge, error) {

	err := trp.validateUpdateTrip(tripSheetUpdateReq)
	if err != nil {
		trp.l.Error("ERROR: UpdateTripSheetHeader", err)
		return nil, err
	}

	tripSheetInfo, errV := trp.tripSheetDao.GetTripSheet(tripSheetId)
	if errV != nil {
		trp.l.Error("ERROR: TripSheet not found", tripSheetInfo, errV)
		return nil, errV
	}
	jsonBytes, _ := json.Marshal(tripSheetInfo)
	trp.l.Info("TripSheet: ******* ", string(jsonBytes))

	// if tripSheetUpdateReq.PodRequired == 0 {

	// }

	tripSheetUpdateReq.TripSubmittedDate = tripSheetInfo.TripSubmittedDate
	tripSheetUpdateReq.LoadStatus = tripSheetInfo.LoadStatus
	if tripSheetInfo.LoadStatus == constant.STATUS_CREATED {
		tripSheetUpdateReq.LoadStatus = constant.STATUS_SUBMITTED
		tripSheetUpdateReq.TripSubmittedDate = utils.GetCurrentDatetimeStr()
	}

	if tripSheetInfo.LoadStatus == constant.STATUS_SUBMITTED {
		tripSheetUpdateReq.TripSubmittedDate = utils.GetCurrentDatetimeStr()
	}

	err1 := trp.tripSheetDao.UpdateTripSheetHeader(tripSheetId, tripSheetUpdateReq)
	if err1 != nil {
		trp.l.Error("ERROR: UpdateTripSheetHeader ", tripSheetUpdateReq.TripSheetNum, err1)
		return nil, err1
	}

	loadUnloads, errL := trp.tripSheetDao.GetTripSheetLoadUnLoadPoints(tripSheetId)
	if errL != nil {
		trp.l.Error("ERROR: GetTripSheetLoadUnLoadPoints ", tripSheetUpdateReq.TripSheetNum, errL)
		return nil, errL
	}
	loadingPoints := make([]int64, 0)
	unloadingPoints := make([]int64, 0)

	for _, point := range *loadUnloads {
		switch point.Type {
		case constant.LOADING_POINT:
			loadingPoints = append(loadingPoints, point.LoadingPointID)
		case constant.UN_LOADING_POINT:
			unloadingPoints = append(unloadingPoints, point.LoadingPointID)
		}
	}

	updateLoading := IsNeedsToUpdateLoadUnLoadPoints(loadingPoints, tripSheetUpdateReq.LoadingPointIDs)
	updateUnLoading := IsNeedsToUpdateLoadUnLoadPoints(unloadingPoints, tripSheetUpdateReq.UnLoadingPointIDs)

	if updateLoading {
		errD := trp.tripSheetDao.DeleteTripSheetLoadUnLoadPoints(tripSheetId, constant.LOADING_POINT)
		if errD != nil {
			trp.l.Error("ERROR: DeleteTripSheetLoadUnLoadPoints ", tripSheetUpdateReq.TripSheetNum, errD)
			return nil, errD
		}

		for _, loadingPointId := range tripSheetUpdateReq.LoadingPointIDs {
			errD := trp.tripSheetDao.SaveTripSheetLoadingPoint(uint64(tripSheetId), uint64(loadingPointId), constant.LOADING_POINT)
			if errD != nil {
				trp.l.Error("ERROR: SaveTripSheetLoadingPoint", errD)
				return nil, errD
			}
		}
		trp.l.Info("upadated loading points : ", updateLoading, loadingPoints, tripSheetUpdateReq.LoadingPointIDs)
	}

	if updateUnLoading {
		errD := trp.tripSheetDao.DeleteTripSheetLoadUnLoadPoints(tripSheetId, constant.UN_LOADING_POINT)
		if errD != nil {
			trp.l.Error("ERROR: DeleteTripSheetLoadUnLoadPoints ", tripSheetUpdateReq.TripSheetNum, errD)
			return nil, errD
		}

		for _, loadingPointId := range tripSheetUpdateReq.UnLoadingPointIDs {
			errD := trp.tripSheetDao.SaveTripSheetLoadingPoint(uint64(tripSheetId), uint64(loadingPointId), constant.UN_LOADING_POINT)
			if errD != nil {
				trp.l.Error("ERROR: SaveTripSheetLoadingPoint", errD)
				return nil, errD
			}
		}
		trp.l.Info("upadated unloading points : ", updateLoading, loadingPoints, tripSheetUpdateReq.LoadingPointIDs)
	}

	// if oldStatus != tripSheetUpdateReq.LoadStatus {
	// 	var err error

	// 	switch tripSheetUpdateReq.LoadStatus {
	// 	case constant.STATUS_SUBMITTED:
	// 	case constant.STATUS_DELIVERED:
	// 	case constant.STATUS_CLOSED:
	// 	case constant.STATUS_COMPLETED:
	// 	default:
	// 	}

	// 	if err != nil {
	// 	}
	// }

	trp.l.Info("tripsheet updated successfully! : ", tripSheetUpdateReq.TripSheetNum)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("tripsheet updated successfully: %v", tripSheetUpdateReq.TripSheetNum)
	return &roleResponse, nil
}

func IsNeedsToUpdateLoadUnLoadPoints(existingArray, reqArray []int64) bool {
	if len(existingArray) == len(reqArray) {

		slices.Sort(existingArray)
		slices.Sort(reqArray)

		for i := range existingArray {
			if existingArray[i] != reqArray[i] {
				return true
			}
		}

		return false
	}

	return true
}
