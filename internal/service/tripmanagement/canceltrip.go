package tripmanagement

import (
	"go-transport-hub/constant"
	"go-transport-hub/dtos"
)

func (trp *TripSheetObj) CancelTripSheet(tripSheetId int64) (*dtos.Messge, error) {

	trp.l.Info("CancelTripSheet", "InIt")
	errU := trp.tripSheetDao.CancelTripSheet(tripSheetId, constant.STATUS_CANCELLED)
	if errU != nil {
		errS := errU.Error()
		trp.l.Error("ERROR: CancelTripSheet", tripSheetId, errU, errS)
		return nil, errU
	}

	trp.l.Info("CancelTripSheetUpdateToPOD", "InIt")
	// Update manage_pod table too.
	errM := trp.tripSheetDao.CancelTripSheetUpdateToPOD(tripSheetId, constant.STATUS_CANCELLED)
	if errM != nil {
		errS := errM.Error()
		trp.l.Error("ERROR: CancelTripSheetUpdateToPOD", tripSheetId, errM, errS)
		return nil, errM
	}

	// tripSheetInfo, errV := trp.tripSheetDao.GetTripSheet(tripSheetId)
	// if errV == nil {
	// 	}
	// }

	trp.l.Info("Trip cancelled successfully... ", tripSheetId)
	roleResponse := dtos.Messge{}
	roleResponse.Message = "Trip cancelled successfully!"
	return &roleResponse, nil
}
