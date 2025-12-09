package commonsvc

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
)

type PreRequisiteObj struct {
	l               *log.Logger
	dbConnMSSQL     *mssqlcon.DBConn
	preRequisiteDao daos.PreRequisiteDao
}

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *PreRequisiteObj {
	return &PreRequisiteObj{
		l:               l,
		dbConnMSSQL:     dbConnMSSQL,
		preRequisiteDao: daos.NewPreRequisiteObj(l, dbConnMSSQL),
	}
}

func (pr *PreRequisiteObj) CreateTripPrerequisite(orgID int64, tripNum, customer, vendor, lodingpoints,
	branch, tripSheetType, tripStatus, customerCode, vendorCode, podTrips, regularTrips, lrTrips, roles,
	vehicleSizes, declarationYear, websiteScreen, permisstionLabel, employees, invoiceStatus string) (*dtos.TripPrerequisite, error) {

	var loadingPoints *[]dtos.LoadingPoints
	var customers *[]dtos.Customers
	var branches *[]dtos.BranchT
	var vendors *[]dtos.VendorT
	var tripSheetTypeList []string
	var tripType []string
	var tripStatusList []string
	var tripSheeetNumber string
	var podTripSheets *[]dtos.TripSheetInfo
	var regularTripSheets *[]dtos.TripSheetInfo
	var lrTripSheets *[]dtos.TripSheetInfo
	var rolesList *[]dtos.RoleEmpPre
	var vehileSizes *[]dtos.VehicleSizeTypePre
	var financialYears []string
	var websiteScreens []dtos.WebsiteScreen
	var permisstionLabels []dtos.PermissionLabel
	var employeeList []dtos.EmployeesPre
	var invoiceStatusList []string

	var errA error

	if lodingpoints == "true" {
		loadingPoints, errA = pr.preRequisiteDao.GetLoadingPoints(orgID)
		if errA != nil {
			pr.l.Error("ERROR: GetLoadingPoints", errA)
			return nil, errA
		}
	}

	if customer == "true" {
		customers, errA = pr.preRequisiteDao.GetCustomers(orgID)
		if errA != nil {
			pr.l.Error("ERROR: GetCustomers", errA)
			return nil, errA
		}
	}

	if branch == "true" {
		branches, errA = pr.preRequisiteDao.GetBranches(orgID)
		if errA != nil {
			pr.l.Error("ERROR: GetBranches", errA)
			return nil, errA
		}
	}

	if vendor == "true" {
		vendors, errA = pr.preRequisiteDao.GetVendors(orgID)
		if errA != nil {
			pr.l.Error("ERROR: GetVendors", errA)
			return nil, errA
		}
	}
	if tripSheetType == "true" {
		tripSheetTypeList = []string{constant.LOCAL_SCHEDULED_TRIP, constant.LOCAL_ADHOC_TRIP, constant.LINE_HUAL_SCHEDULED_TRIP, constant.LINE_HUAL_ADHOC_TRIP}
		tripType = []string{constant.ONE_WAY_TRIP, constant.ROUND_TRIP}
	}

	if tripNum == "true" {
		tripSheeetNumber = pr.GetNextTripSheetNumber()
		pr.l.Info("tripSheeetNumber:", tripNum, tripSheeetNumber)
	}

	if tripStatus == "true" {
		tripStatusList = []string{
			constant.STATUS_CREATED, constant.STATUS_SUBMITTED,
			constant.STATUS_CANCELLED, constant.STATUS_CLOSED,
			constant.STATUS_DELIVERED, constant.STATUS_COMPLETED, constant.STATUS_INVOICE_RAISED, constant.STATUS_PAID,
		}
	}

	//
	if invoiceStatus == "true" {
		invoiceStatusList = []string{
			constant.STATUS_INVOICE_DRAFT, constant.STATUS_INVOICE_RAISED, constant.STATUS_PAID, constant.STATUS_CANCELLED,
		}
	}

	var customerCodeTH, vendorCodeTH string
	if customerCode == "true" {
		customerCodeTH = strings.ToUpper(pr.GetNextCustomerCode())
	}
	if vendorCode == "true" {
		vendorCodeTH = strings.ToUpper(pr.GetNextVendorCode())
	}
	if podTrips == "true" {
		res, errP := pr.preRequisiteDao.GetPodTripInfo(1)
		if errP != nil {
			pr.l.Error("ERROR: GetPodTripInfo", errP)
			return nil, errP
		}
		podTripSheets = res
	}

	if regularTrips == "true" {
		res, errP := pr.preRequisiteDao.GetPodTripInfo(0)
		if errP != nil {
			pr.l.Error("ERROR: GetPodTripInfo 0", errP)
			return nil, errP
		}
		regularTripSheets = res
	}

	if lrTrips == "true" {
		res, errP := pr.preRequisiteDao.GetLRTripInfo()
		if errP != nil {
			pr.l.Error("ERROR: GetLRTripInfo 0", errP)
			return nil, errP
		}
		lrTripSheets = res
	}

	if roles == "true" {
		res, errP := pr.preRequisiteDao.GetEmpRoles(orgID)
		if errP != nil {
			pr.l.Error("ERROR: GetLRTripInfo 0", errP)
			return nil, errP
		}
		rolesList = res
	}

	if vehicleSizes == "true" {
		res, errP := pr.preRequisiteDao.GetVehicleSizeTypes()
		if errP != nil {
			pr.l.Error("ERROR: GetLRTripInfo 0", errP)
			return nil, errP
		}
		vehileSizes = res
	}
	if declarationYear == "true" {
		financialYears = pr.financialYearsList(5, 5)
	}

	if websiteScreen == "true" {
		res, errP := pr.preRequisiteDao.GetWebsiteScreen()
		if errP != nil {
			pr.l.Error("ERROR: GetWebsiteScreen 0", errP)
			return nil, errP
		}
		websiteScreens = *res
	}

	if permisstionLabel == "true" {
		res, errP := pr.preRequisiteDao.GetPermissionLabel()
		if errP != nil {
			pr.l.Error("ERROR: GetWebsiteScreen 0", errP)
			return nil, errP
		}
		permisstionLabels = *res
	}
	if employees == "true" {
		res, errP := pr.preRequisiteDao.GetEmployeePre()
		if errP != nil {
			pr.l.Error("ERROR: GetWebsiteScreen 0", errP)
			return nil, errP
		}
		employeeList = *res
	}

	response := dtos.TripPrerequisite{}
	response.TripSheetNumber = tripSheeetNumber
	response.CustomerCode = customerCodeTH
	response.VendorCode = vendorCodeTH
	response.TripType = tripType
	response.Branch = branches
	response.TripSheetType = tripSheetTypeList
	response.LoadingPoints = loadingPoints
	response.Customers = customers
	response.Vendor = vendors
	response.TripStatus = tripStatusList
	response.PodTrips = podTripSheets
	response.RegularTrips = regularTripSheets
	response.LRTrips = lrTripSheets
	response.Roles = rolesList
	response.VehicleSizeType = vehileSizes
	response.DeclarationYear = financialYears
	response.WebsiteScreen = websiteScreens
	response.PermissionLabel = permisstionLabels
	response.EmployeesPre = employeeList
	response.InvoiceStatus = invoiceStatusList
	return &response, nil
}

var tripSheetLayout = "200601021504"
var layout = "2006-01-02 15:04"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (br *PreRequisiteObj) GetTripSheetNumber() string {

	financialYear := getFinancialYear(time.Now())

	tripSheetNumber, err := br.preRequisiteDao.InsertAndGetTripNumber(financialYear)
	if err == nil {
		return fmt.Sprintf("%v/%v", fmt.Sprintf("%06d", tripSheetNumber), financialYear)
	}

	// Fallback
	br.l.Error("ERROR: GetTripSheetNumber ERROR: ", err)
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		br.l.Error("ERROR: time.LoadLocation", err)
	}

	dateStr := time.Now().In(loc).Format(layout)
	br.l.Info("TripSheetNumber date : ", dateStr)
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		br.l.Error("ERROR: CreateTripPrerequisite", err)
	}
	formatted := t.Format(tripSheetLayout)
	randomCode, err := GenerateRandomCode(5)
	if err != nil {
		br.l.Error("Error: CreateTripPrerequisite generating random TripSheetNumber code:", err)
	}
	return fmt.Sprintf("KKT%v%v", formatted, randomCode)
}
func (br *PreRequisiteObj) GetNextTripSheetNumber() string {
	fallBackTripSheetNumber := br.GetTripSheetNumber()
	tripSheetId, tripSheetNumber, err := br.preRequisiteDao.GetLastTripSheetRow()
	if err != nil {
		br.l.Error("fallBackTripSheetNumber generating", tripSheetId, tripSheetNumber, err)
		return fallBackTripSheetNumber
	}
	br.l.Info("last trip_sheet id", tripSheetId, tripSheetNumber)

	if tripSheetNumber != "" {
		parts := strings.Split(tripSheetNumber, "/")
		if len(parts) == 2 {
			num, err := strconv.Atoi(parts[0])
			if err != nil {
				panic(err)
			}
			num = num + 1
			financialYear := getFinancialYear(time.Now())
			newTripSheetNumber := fmt.Sprintf("%v/%v", fmt.Sprintf("%06d", num), financialYear)

			return newTripSheetNumber
		}
	}
	br.l.Warn("fallBackTripSheetNumber generating", tripSheetId, tripSheetNumber)
	return fallBackTripSheetNumber
}

func getFinancialYear(t time.Time) string {
	year := t.Year()
	if t.Month() >= time.April {
		return fmt.Sprintf("%d-%d", year, year+1)
	}
	return fmt.Sprintf("%d-%d", year-1, year)
}

func GenerateRandomCode(length int) (string, error) {
	byteArray := make([]byte, length)
	_, err := rand.Read(byteArray)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		byteArray[i] = charset[int(byteArray[i])%len(charset)]
	}

	return string(byteArray), nil
}

func (br *PreRequisiteObj) GetCustomerCode() string {

	lastInsertId, err := br.preRequisiteDao.InsertAndGetTripNumber(constant.CUSTOMER)
	if err == nil {
		return fmt.Sprintf("%v-%v", constant.CUSTOMER, fmt.Sprintf("%06d", lastInsertId))
	}

	// Fallback
	br.l.Error("ERROR: GetCustomerCode ERROR: ", err)
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		br.l.Error("ERROR: time.LoadLocation", err)
	}

	dateStr := time.Now().In(loc).Format(layout)
	br.l.Info("TripSheetNumber date : ", dateStr)
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		br.l.Error("ERROR: CreateTripPrerequisite", err)
	}
	formatted := t.Format(tripSheetLayout)
	randomCode, err := GenerateRandomCode(5)
	if err != nil {
		br.l.Error("Error: CreateTripPrerequisite generating random TripSheetNumber code:", err)
	}
	return fmt.Sprintf("KKT%v%v", formatted, randomCode)
}

func (br *PreRequisiteObj) GetNextCustomerCode() string {
	customerId, customerCode, err := br.preRequisiteDao.GetLastCustomerRow()
	if err != nil {
		// If no records exist, return the first code "000001"
		if err == sql.ErrNoRows {
			br.l.Info("No existing customers found, returning first customer code")
			return fmt.Sprintf("%v-%v", constant.CUSTOMER, fmt.Sprintf("%06d", 1))
		}
		// For other errors, also return 000001 instead of using fallback that inserts into trip_sheet_num
		br.l.Error("Error getting last customer row, returning first customer code", customerId, customerCode, err)
		return fmt.Sprintf("%v-%v", constant.CUSTOMER, fmt.Sprintf("%06d", 1))
	}
	br.l.Info("last customer:", customerId, customerCode)

	if customerCode != "" {
		parts := strings.Split(customerCode, "-")
		if len(parts) == 2 {
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				br.l.Error("Error parsing customer code number, returning first customer code:", err)
				// Return 000001 instead of using fallback
				return fmt.Sprintf("%v-%v", constant.CUSTOMER, fmt.Sprintf("%06d", 1))
			}
			num = num + 1
			newCustomerCode := fmt.Sprintf("%v-%v", constant.CUSTOMER, fmt.Sprintf("%06d", num))
			return newCustomerCode
		}
	}
	br.l.Warn("Invalid customer code format, returning first customer code", customerId, customerCode)
	// Return 000001 instead of using fallback
	return fmt.Sprintf("%v-%v", constant.CUSTOMER, fmt.Sprintf("%06d", 1))
}

func (br *PreRequisiteObj) GetVendorCode() string {

	lastInsertId, err := br.preRequisiteDao.InsertAndGetTripNumber(constant.VENDOR)
	if err == nil {
		return fmt.Sprintf("%v-%v", constant.VENDOR, fmt.Sprintf("%06d", lastInsertId))
	}

	// Fallback
	br.l.Error("ERROR: GetVendorCode ERROR: ", err)
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		br.l.Error("ERROR: time.LoadLocation", err)
	}

	dateStr := time.Now().In(loc).Format(layout)
	br.l.Info("GetVendorCode date : ", dateStr)
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		br.l.Error("ERROR: CreateTripPrerequisite", err)
	}
	formatted := t.Format(tripSheetLayout)
	randomCode, err := GenerateRandomCode(5)
	if err != nil {
		br.l.Error("Error: CreateTripPrerequisite generating random TripSheetNumber code:", err)
	}
	return fmt.Sprintf("KKT%v%v", formatted, randomCode)
}

func (br *PreRequisiteObj) GetNextVendorCode() string {
	vendorId, vendorCode, err := br.preRequisiteDao.GetLastVendorRow()
	if err != nil {
		// If no records exist, return the first code "000001"
		if err == sql.ErrNoRows {
			br.l.Info("No existing vendors found, returning first vendor code")
			return fmt.Sprintf("%v-%v", constant.VENDOR, fmt.Sprintf("%06d", 1))
		}
		// For other errors, also return 000001 instead of using fallback that inserts into trip_sheet_num
		br.l.Error("Error getting last vendor row, returning first vendor code", vendorId, vendorCode, err)
		return fmt.Sprintf("%v-%v", constant.VENDOR, fmt.Sprintf("%06d", 1))
	}
	br.l.Info("last vendor:", vendorId, vendorCode)

	if vendorCode != "" {
		parts := strings.Split(vendorCode, "-")
		if len(parts) == 2 {
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				br.l.Error("Error parsing vendor code number, returning first vendor code:", err)
				// Return 000001 instead of using fallback
				return fmt.Sprintf("%v-%v", constant.VENDOR, fmt.Sprintf("%06d", 1))
			}
			num = num + 1
			newVendorCode := fmt.Sprintf("%v-%v", constant.VENDOR, fmt.Sprintf("%06d", num))
			return newVendorCode
		}
	}
	br.l.Warn("Invalid vendor code format, returning first vendor code", vendorId, vendorCode)
	// Return 000001 instead of using fallback
	return fmt.Sprintf("%v-%v", constant.VENDOR, fmt.Sprintf("%06d", 1))
}

func (br *PreRequisiteObj) financialYearsList(yearsBefore, yearsAfter int) []string {
	currentYear := time.Now().Year()
	startYear := currentYear - yearsBefore
	endYear := currentYear + yearsAfter

	var financialYears []string
	for year := startYear; year < endYear; year++ {
		fy := fmt.Sprintf("%d-%d", year, year+1)
		financialYears = append(financialYears, fy)
	}
	return financialYears
}
