package vendor

import (
	"errors"

	"go-transport-hub/dtos"
)

func (vh *VendorObj) validateVendor(vendorReg dtos.VendorRequest) error {

	if vendorReg.VendorName == "" {
		vh.l.Error("Error VendorName: vendor name should not empty")
		return errors.New("vendor name should not empty")
	}
	if vendorReg.VendorCode == "" {
		vh.l.Error("Error VendorCode: vendor code should not empty")
		return errors.New("vendor code should not empty")
	}
	if vendorReg.OwnerName == "" {
		vh.l.Error("Error OwnerName: vendor code should not empty")
		return errors.New("owner name should not empty")
	}
	if vendorReg.AddressLine1 == "" {
		vh.l.Error("Error AddressLine1: address should not empty")
		return errors.New("address should not empty")
	}
	if vendorReg.PreferredOperatingRoutes == "" {
		vh.l.Error("Error PreferredOperatingRoutes: PreferredOperatingRoutes should not empty")
		return errors.New("preferred operating routes should not empty")
	}
	if vendorReg.PANNumber == "" {
		vh.l.Error("Error PANNumber: PreferredOperatingRoutes should not empty")
		return errors.New("pan number should not empty")
	}
	if vendorReg.TDSDeclaration == "" {
		vh.l.Error("Error TDSDeclaration: PreferredOperatingRoutes should not empty")
		return errors.New("TDS declaration should not empty")
	}
	if vendorReg.Remark == "" {
		vh.l.Error("Error Remark: PreferredOperatingRoutes should not empty")
		return errors.New("remark should not empty")
	}
	if vendorReg.BankAccountHolderName == "" {
		vh.l.Error("Error: BankAccountHolderName should not empty")
		return errors.New("bank account holder name should not empty")
	}
	if vendorReg.BankAccountNumber == "" {
		vh.l.Error("Error: BankAccountNumber should not empty")
		return errors.New("bank account number should not empty")
	}
	if vendorReg.BankName == "" {
		vh.l.Error("Error: BankName should not empty")
		return errors.New("bank name should not empty")
	}
	if vendorReg.BankIFSCCode == "" {
		vh.l.Error("Error: BankIFSCCode should not empty")
		return errors.New("IFSC code should not empty")
	}
	if vendorReg.BankPassbookOrChequeImg == "" {
		vh.l.Error("Error: BankPassbookOrChequeImg should not empty")
		return errors.New("bank passbook Or cheque image should not empty")
	}

	//

	if len(vendorReg.ContactInfo) != 0 {
		for _, contactInfo := range vendorReg.ContactInfo {
			if contactInfo.ContactPersonName == "" {
				vh.l.Error("Error: ContactPersonName contact person name")
				return errors.New("contact person name should not empty")
			}
			if contactInfo.Post == "" {
				vh.l.Error("Error: post contact person name")
				return errors.New("post should not empty")
			}
		}
	}

	if len(vendorReg.Vehicles) != 0 {
		for _, vehicleObj := range vendorReg.Vehicles {
			if vehicleObj.VehicleNumber == "" {
				vh.l.Error("Error: VehicleNumber should not empty")
				return errors.New("vehicle number should not empty")
			}
			if vehicleObj.VehicleType == "" {
				vh.l.Error("Error: VehicleType should not empty")
				return errors.New("vehicle type should not empty")
			}
			if vehicleObj.VehicleModel == "" {
				vh.l.Error("Error: VehicleModel should not empty")
				return errors.New("vehicle model should not empty")
			}
			if vehicleObj.VehicleMake == "" {
				vh.l.Error("Error: VehicleMake should not empty")
				return errors.New("vehicle make should not empty")
			}
			if vehicleObj.VehicleModel == "" {
				vh.l.Error("Error: VehicleModel should not empty")
				return errors.New("vehicle model should not empty")
			}
			// Vehicle document uploads (RC, Insurance, PUCC, NP, Fitness, Tax, MP) are now optional
			if vehicleObj.VehicleSize == "" {
				vh.l.Error("Error: VehicleSize should not empty")
				return errors.New("vehicle size should not empty")
			}
		}
	}
	return nil
}
