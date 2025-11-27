package lrreceipt

import (
	"errors"

	"go-transport-hub/dtos"
)

func (mp *LRReceipt) validateLR(lrReq dtos.LRReceiptReq) error {
	if lrReq.TripSheetNum == "" {
		mp.l.Error("Error trip sheet number should not empty")
		return errors.New("trip sheet number should not empty")
	}
	if lrReq.LRNumber == "" {
		mp.l.Error("Error LRNumber: lr should not empty")
		return errors.New("lr number should not empty")
	}

	if lrReq.TripDate == "" {
		mp.l.Error("Error TripDate : trip date should not empty")
		return errors.New("trip date should not empty")
	}

	if lrReq.VehicleNumber == "" {
		mp.l.Error("Error VehicleNumber: Vehicle number should not empty")
		return errors.New("vehicle number should not empty")
	}
	if lrReq.VehicleSize == "" {
		mp.l.Error("Error VehicleSize: Vehicle size should not empty")
		return errors.New("vehicle size should not empty")
	}

	if lrReq.InvoiceNumber == "" {
		mp.l.Error("Error InvoiceNumber: invoice number should not empty")
		return errors.New("invoice number should not empty")
	}
	if lrReq.InvoiceValue == 0 {
		mp.l.Error("Error InvoiceValue: invoice value should not empty")
		return errors.New("invoice value should not empty")
	}
	if lrReq.ConsigneeName == "" {
		mp.l.Error("Error ConsigneeName: consignee name should not empty")
		return errors.New("consignee name should not empty")
	}
	if lrReq.ConsigneeAddress == "" {
		mp.l.Error("Error ConsigneeAddress: consignee address should not empty")
		return errors.New("consignee address should not empty")
	}

	if lrReq.ConsignorName == "" {
		mp.l.Error("Error ConsignorName: consignor name should not empty")
		return errors.New("consignor name should not empty")
	}
	if lrReq.ConsignorAddress == "" {
		mp.l.Error("Error ConsignorAddress: consignor address should not empty")
		return errors.New("consignor address should not empty")
	}
	if lrReq.GoodsType == "" {
		mp.l.Error("Error GoodsType: goods type should not empty")
		return errors.New("goods type should not empty")
	}
	if lrReq.GoodsWeight == "" {
		mp.l.Error("Error GoodsType: goods weight should not empty")
		return errors.New("goods weight should not empty")
	}
	if lrReq.DriverName == "" {
		mp.l.Error("Error DriverName: driver name should not empty")
		return errors.New("driver name should not empty")
	}
	if lrReq.DriverMobileNumber == "" {
		mp.l.Error("Error DriverMobileNumber: driver mobile number should not empty")
		return errors.New("driver mobile number should not empty")
	}

	return nil
}

func (mp *LRReceipt) validateUpdateLR(lrReq dtos.LRReceiptUpdateReq) error {
	if lrReq.TripSheetNum == "" {
		mp.l.Error("Error trip sheet number should not empty")
		return errors.New("trip sheet number should not empty")
	}
	if lrReq.LRNumber == "" {
		mp.l.Error("Error LRNumber: lr should not empty")
		return errors.New("lr should not empty")
	}

	if lrReq.TripDate == "" {
		mp.l.Error("Error TripDate : trip date should not empty")
		return errors.New("trip date should not empty")
	}

	if lrReq.VehicleNumber == "" {
		mp.l.Error("Error VehicleNumber: Vehicle number should not empty")
		return errors.New("vehicle number should not empty")
	}
	if lrReq.VehicleSize == "" {
		mp.l.Error("Error VehicleSize: Vehicle size should not empty")
		return errors.New("vehicle size should not empty")
	}

	if lrReq.InvoiceNumber == "" {
		mp.l.Error("Error InvoiceNumber: invoice number should not empty")
		return errors.New("invoice number should not empty")
	}
	if lrReq.InvoiceValue == 0 {
		mp.l.Error("Error InvoiceValue: invoice value should not empty")
		return errors.New("invoice value should not empty")
	}
	if lrReq.ConsigneeName == "" {
		mp.l.Error("Error ConsigneeName: consignee name should not empty")
		return errors.New("consignee name should not empty")
	}
	if lrReq.ConsigneeAddress == "" {
		mp.l.Error("Error ConsigneeAddress: consignee address should not empty")
		return errors.New("consignee address should not empty")
	}

	if lrReq.ConsignorName == "" {
		mp.l.Error("Error ConsignorName: consignor name should not empty")
		return errors.New("consignor name should not empty")
	}
	if lrReq.ConsignorAddress == "" {
		mp.l.Error("Error ConsignorAddress: consignor address should not empty")
		return errors.New("consignor address should not empty")
	}
	if lrReq.GoodsType == "" {
		mp.l.Error("Error GoodsType: goods type should not empty")
		return errors.New("goods type should not empty")
	}
	if lrReq.GoodsWeight == "" {
		mp.l.Error("Error GoodsType: goods weight should not empty")
		return errors.New("goods weight should not empty")
	}
	if lrReq.DriverName == "" {
		mp.l.Error("Error DriverName: driver name should not empty")
		return errors.New("driver name should not empty")
	}
	if lrReq.DriverMobileNumber == "" {
		mp.l.Error("Error DriverMobileNumber: driver mobile number should not empty")
		return errors.New("driver mobile number should not empty")
	}

	return nil
}
