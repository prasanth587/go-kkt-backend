package customer

import (
	"errors"
	"strings"

	"go-transport-hub/dtos"
)

func (trp *CustomerObj) validateCustomer(customerReq dtos.CustomerReq) error {

	if customerReq.BranchId == 0 {
		trp.l.Error("Error BranchID: select brach")
		return errors.New("select brach")
	}
	if customerReq.CustomerName == "" {
		trp.l.Error("Error CustomerName: customer name should not empty")
		return errors.New("customer name should not empty")
	}
	if customerReq.CustomerCode == "" {
		trp.l.Error("Error CustomerCode: customer code should not empty")
		return errors.New("customer code should not empty")
	}
	customerReq.CustomerCode = strings.ToUpper(customerReq.CustomerCode)

	if customerReq.Address == "" {
		trp.l.Error("Error Address: Address should not empty")
		return errors.New("address should not empty")
	}
	if customerReq.City == "" {
		trp.l.Error("Error City: city should not empty")
		return errors.New("city should not empty")
	}
	if customerReq.State == "" {
		trp.l.Error("Error State: state should not empty")
		return errors.New("state should not empty")
	}
	if customerReq.PinCode == 0 {
		trp.l.Error("Error PinCode: pincode should not empty")
		return errors.New("pincode should not empty")
	}
	if customerReq.Country == "" {
		trp.l.Error("Error Country: country should not empty")
		return errors.New("country should not empty")
	}
	if customerReq.GSTNType == "" {
		trp.l.Error("Error GSTNType: gstn type should not empty")
		return errors.New("gstn type should not empty")
	}
	if customerReq.GSTNNo == "" {
		trp.l.Error("Error GSTNNo: gstn no should not empty")
		return errors.New("gstn no should not empty")
	}
	if customerReq.PanNumber == "" {
		trp.l.Error("Error PanNumber: pan number should not empty")
		return errors.New("pan number should not empty")
	}
	if customerReq.ContactPerson == "" {
		trp.l.Error("Error ContactPerson: contact person should not empty")
		return errors.New("contact person should not empty")
	}
	if customerReq.MobileNumber == "" {
		trp.l.Error("Error ContactNumber: contact number should not empty")
		return errors.New("contact number should not empty")
	}
	if customerReq.EmailID == "" {
		trp.l.Error("Error EmailID: email should not empty")
		return errors.New("email should not empty")
	}

	if customerReq.BranchSalesPersonName == "" {
		trp.l.Error("Error BranchSalesPersonName: branch sales person not empty")
		return errors.New("branch responsible person should not empty")
	}
	if customerReq.BranchResponsiblePerson == "" {
		trp.l.Error("Error BranchResponsiblePerson: branch sales person not empty")
		return errors.New("branch sales person should not empty")
	}
	if customerReq.EmployeeCode == "" {
		trp.l.Error("Error EmployeeCode: employee code person not empty")
		return errors.New("employee code should not empty")
	}
	return nil
}

func (trp *CustomerObj) validateCustomerReq(customerReq dtos.CustomersReq) error {

	if customerReq.CustomerName == "" {
		trp.l.Error("Error CustomerName: customer name should not empty")
		return errors.New("customer name should not empty")
	}
	if customerReq.CustomerCode == "" {
		trp.l.Error("Error CustomerCode: customer code should not empty")
		return errors.New("customer code should not empty")
	}
	customerReq.CustomerCode = strings.ToUpper(customerReq.CustomerCode)

	if customerReq.Address == "" {
		trp.l.Error("Error Address: Address should not empty")
		return errors.New("address should not empty")
	}
	if customerReq.City == "" {
		trp.l.Error("Error City: city should not empty")
		return errors.New("city should not empty")
	}
	if customerReq.State == "" {
		trp.l.Error("Error State: state should not empty")
		return errors.New("state should not empty")
	}
	if customerReq.PinCode == 0 {
		trp.l.Error("Error PinCode: pincode should not empty")
		return errors.New("pincode should not empty")
	}
	if customerReq.Country == "" {
		trp.l.Error("Error Country: country should not empty")
		return errors.New("country should not empty")
	}
	if customerReq.KKTResponsibleEmpID == 0 {
		trp.l.Error("Error KKTResponsibleEmp: kkt responsible person not empty")
		return errors.New("kkt responsible person not empty")
	}
	return nil
}
