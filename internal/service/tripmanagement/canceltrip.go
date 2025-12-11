package tripmanagement

import (
	"go-transport-hub/constant"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/notification"
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

	// Get trip sheet info for notification
	tripSheetInfo, errV := trp.tripSheetDao.GetTripSheet(tripSheetId)
	if errV == nil {
		// Send notification for trip cancellation
		notificationSvc := notification.New(trp.l, trp.dbConnMSSQL)
		if err := notificationSvc.NotifyTripSheetStatusChanged(tripSheetInfo.OrgID, tripSheetId, tripSheetInfo.TripSheetNum, tripSheetInfo.LoadStatus, constant.STATUS_CANCELLED); err != nil {
			trp.l.Error("ERROR: Failed to send cancellation notification: ", err)
			// Don't fail the request if notification fails
		}
	}

	trp.l.Info("Trip cancelled successfully... ", tripSheetId)
	roleResponse := dtos.Messge{}
	roleResponse.Message = "Trip cancelled successfully!"
	return &roleResponse, nil
}
