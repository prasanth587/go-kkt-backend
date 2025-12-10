package managepod

import (
	"errors"
	"fmt"

	"go-transport-hub/constant"
	"go-transport-hub/dtos"
)

func (mp *ManagePod) validatePOD(podReq dtos.ManagePodReq) error {
	if podReq.TripSheetNum == "" {
		mp.l.Error("Error trip sheet number should not empty")
		return errors.New("trip sheet number should not empty")
	}
	if podReq.LRNumber == "" {
		mp.l.Error("Error LRNumber: lr number should not empty")
		return errors.New("lr number should not empty")
	}

	if podReq.PaidBy == "" && podReq.TripType == constant.TRIP_TYPE_POD {
		mp.l.Error("Error PaidBy: paid by should not empty")
		return errors.New("paid by should not empty")
	}
	if podReq.PODSubmitedDate == "" && podReq.TripType == constant.TRIP_TYPE_POD {
		mp.l.Error("Error PODSubmitedDate: pod submite should not empty")
		return errors.New("pod submite should not empty")
	}
	// LateSubmissionDebit, UnloadingCharges, and HaltingAmount are now optional (can be 0 or empty)

	if podReq.PodRemark == "" {
		mp.l.Error("Error PodRemark: pod remark should not empty")
		return errors.New("pod remark should not empty")
	}

	if podReq.OrgId == 0 {
		mp.l.Error("Error OrgId: organization should not empty")
		return errors.New("organization should not empty")
	}

	ttype := [...]string{
		constant.TRIP_TYPE_POD, constant.TRIP_TYPE_REGULAR,
	}

	exists := false
	for _, nype := range ttype {
		if nype == podReq.TripType {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("trip type not valid. should be: %s,%s", constant.TRIP_TYPE_POD, constant.TRIP_TYPE_REGULAR)
	}

	return nil
}

func (mp *ManagePod) validateUpdatePOD(podReq dtos.UpdateManagePodReq) error {
	if podReq.TripSheetNum == "" {
		mp.l.Error("Error trip sheet number should not empty")
		return errors.New("trip sheet number should not empty")
	}
	if podReq.LRNumber == "" {
		mp.l.Error("Error LRNumber: lr should not empty")
		return errors.New("lr should not empty")
	}

	if podReq.PaidBy == "" && podReq.TripType == constant.TRIP_TYPE_POD {
		mp.l.Error("Error PaidBy: paid by should not empty")
		return errors.New("paid by should not empty")
	}
	if podReq.PODSubmitedDate == "" && podReq.TripType == constant.TRIP_TYPE_POD {
		mp.l.Error("Error PODSubmitedDate: pod submite should not empty")
		return errors.New("pod submite should not empty")
	}
	// LateSubmissionDebit, UnloadingCharges, and HaltingAmount are now optional (can be 0 or empty)

	if podReq.PodRemark == "" {
		mp.l.Error("Error PodRemark: pod remark should not empty")
		return errors.New("pod remark should not empty")
	}

	ttype := [...]string{
		constant.TRIP_TYPE_POD, constant.TRIP_TYPE_REGULAR,
	}

	exists := false
	for _, nype := range ttype {
		if nype == podReq.TripType {
			exists = true
			break
		}
	}

	if podReq.TripType == constant.TRIP_TYPE_POD {
		if podReq.PodUpload == "" {
			mp.l.Error("Error PodUpload: upload pod document")
			return errors.New("upload pod document")
		}
	}

	if !exists {
		return fmt.Errorf("trip type not valid. should be: %s,%s", constant.TRIP_TYPE_POD, constant.TRIP_TYPE_REGULAR)
	}

	return nil
}
