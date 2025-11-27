package dtos

type TripSheetXlsRes struct {
	TripSheet *[]TripSheetXlsV1 `json:"trip_sheets"`
	Total     int64             `json:"total"`
	Limit     int64             `json:"limit"`
	OffSet    int64             `json:"offset"`
}

type TripSheetXls struct {
	TripSheetID      int64  `json:"trip_sheet_id"`
	TripSheetNum     string `json:"trip_sheet_num"`
	TripType         string `json:"trip_type"`
	TripSheetType    string `json:"trip_sheet_type"`
	LoadHoursType    string `json:"load_hours_type"`
	OpenTripDateTime string `json:"open_trip_date_time"`
	// Foreign Keys
	BranchID          int64                 `json:"branch_id"`
	CustomerId        int64                 `json:"customer_id"`
	CustomerName      string                `json:"customer_name"`
	CustomerCode      string                `json:"customer_code"`
	VendorID          int64                 `json:"vendor_id"`
	LoadingPointIDs   *[]TripSheetLoading   `json:"loading_points"`
	UnLoadingPointIDs *[]TripSheetUnLoading `json:"un_loading_point"`
	// Vehicle Details
	VehicleCapacityTon string `json:"vehicle_capacity_ton"`
	VehicleNumber      string `json:"vehicle_number"`
	VehicleSize        string `json:"vehicle_size"`
	MobileNumber       string `json:"mobile_number"`
	DriverName         string `json:"driver_name"`
	DriverLicenseImage string `json:"driver_license_image"`
	LRGateImage        string `json:"lr_gate_image"`
	LRNumber           string `json:"lr_number"`

	// Customer Section
	CustomerCloseTripDateTime              string  `json:"customer_close_trip_date_time"`
	CustomerInvoiceNo                      string  `json:"customer_invoice_no"`
	CustomerPaymentReceivedDate            string  `json:"customer_payment_received_date"`
	CustomerRemark                         string  `json:"customer_remark"`
	CustomerBillingRaisedDate              string  `json:"customer_billing_raised_date"`
	CustomerBaseRate                       float64 `json:"customer_base_rate"`
	CustomerKMCost                         float64 `json:"customer_km_cost"`
	CustomerToll                           float64 `json:"customer_toll"`
	CustomerExtraHours                     float64 `json:"customer_extra_hours"`
	CustomerExtraKM                        float64 `json:"customer_extra_km"`
	CustomerTotalHire                      float64 `json:"customer_total_hire"`
	CustomerDebitAmount                    float64 `json:"customer_debit_amount"`
	PodRequired                            int64   `json:"pod_required"`
	CustomerPerLoadHire                    float64 `json:"customer_per_load_hire"`
	CustomerRunningKM                      float64 `json:"customer_running_km"`
	CustomerPerKMPrice                     float64 `json:"customer_per_km_price"`
	CustomerPlacedVehicleSize              string  `json:"customer_placed_vehicle_size"`
	CustomerLoadCancelled                  string  `json:"customer_load_cancelled"`
	CustomerReportedDateTimeForHaltingCalc string  `json:"customer_reported_date_time_for_halting_calc"`

	// Vendor Section
	VendorPaidDate         string  `json:"vendor_paid_date"`
	VendorCrossDock        string  `json:"vendor_cross_dock"`
	VendorRemark           string  `json:"vendor_remark"`
	VendorBaseRate         float64 `json:"vendor_base_rate"`
	VendorKMCost           float64 `json:"vendor_km_cost"`
	VendorToll             float64 `json:"vendor_toll"`
	VendorTotalHire        float64 `json:"vendor_total_hire"`
	VendorAdvance          float64 `json:"vendor_advance"`
	VendorDebitAmount      float64 `json:"vendor_debit_amount"`
	VendorBalanceAmount    float64 `json:"vendor_balance_amount"`
	VendorPaidBy           string  `json:"vendor_paid_by"`
	VendorLoadUnLoadAmount float64 `json:"vendor_load_unload_amount"`
	VendorHaltingDays      float64 `json:"vendor_halting_days"`
	VendorHaltingPaid      float64 `json:"vendor_halting_paid"`
	VendorExtraDelivery    float64 `json:"vendor_extra_delivery"`
	VendorBreakDown        string  `json:"vendor_break_down"`
	PodReceived            int64   `json:"pod_received"`
	LoadStatus             string  `json:"load_status"`
	ZonalName              string  `json:"zonal_name"`
}

type TripSheetXlsV1 struct {
	TripSheetID       int64  `json:"trip_sheet_id"`
	TripSheetNum      string `json:"trip_sheet_num"`
	TripType          string `json:"trip_type"`
	TripSheetType     string `json:"trip_sheet_type"`
	LoadHoursType     string `json:"load_hours_type"`
	OpenTripDateTime  string `json:"open_trip_date_time"`
	BranchID          int64  `json:"branch_id"`
	CustomerId        int64  `json:"customer_id"`
	CustomerName      string `json:"customer_name"`
	CustomerCode      string `json:"customer_code"`
	VendorID          int64  `json:"vendor_id"`
	FromLocations     string `json:"from_locations"`
	ToLocations       string `json:"to_locations"`
	CustomerInvoiceNo string `json:"customer_invoice_no"`
	LRNumber          string `json:"lr_number"`
	LoadStatus        string `json:"load_status"`
	VehicleNumber     string `json:"vehicle_number"`
	VendorName        string `json:"vendor_name"`
}

type LoadUnLoadObj struct {
	TripSheetId    int64  `json:"trip_sheet_id"`
	LoadingPointId int64  `json:"loading_point_id"`
	CityCode       string `json:"city_code"`
	CityName       string `json:"city_name"`
	Type           string `json:"type"`
}

type DownloadTripSheetXls struct {
	TripSheetID        int64  `json:"trip_sheet_id"`
	TripSheetNum       string `json:"trip_sheet_num"`
	TripType           string `json:"trip_type"`
	TripSheetType      string `json:"trip_sheet_type"`
	LoadHoursType      string `json:"load_hours_type"`
	OpenTripDateTime   string `json:"open_trip_date_time"`
	CustomerName       string `json:"customer_name"`
	CustomerCode       string `json:"customer_code"`
	VendorID           int64  `json:"vendor_id"`
	LoadingPointIDs    string `json:"loading_points"`
	UnLoadingPointIDs  string `json:"un_loading_point"`
	VehicleCapacityTon string `json:"vehicle_capacity_ton"`
	VehicleNumber      string `json:"vehicle_number"`
	VehicleSize        string `json:"vehicle_size"`
	MobileNumber       string `json:"mobile_number"`
	DriverName         string `json:"driver_name"`
	DriverLicenseImage string `json:"driver_license_image"`
	LRGateImage        string `json:"lr_gate_image"`
	LRNumber           string `json:"lr_number"`

	// Customer Section
	CustomerCloseTripDateTime              string  `json:"customer_close_trip_date_time"`
	CustomerInvoiceNo                      string  `json:"customer_invoice_no"`
	CustomerPaymentReceivedDate            string  `json:"customer_payment_received_date"`
	CustomerRemark                         string  `json:"customer_remark"`
	CustomerBillingRaisedDate              string  `json:"customer_billing_raised_date"`
	CustomerBaseRate                       float64 `json:"customer_base_rate"`
	CustomerKMCost                         float64 `json:"customer_km_cost"`
	CustomerToll                           float64 `json:"customer_toll"`
	CustomerExtraHours                     float64 `json:"customer_extra_hours"`
	CustomerExtraKM                        float64 `json:"customer_extra_km"`
	CustomerTotalHire                      float64 `json:"customer_total_hire"`
	CustomerDebitAmount                    float64 `json:"customer_debit_amount"`
	PodRequired                            int64   `json:"pod_required"`
	CustomerPerLoadHire                    float64 `json:"customer_per_load_hire"`
	CustomerRunningKM                      float64 `json:"customer_running_km"`
	CustomerPerKMPrice                     float64 `json:"customer_per_km_price"`
	CustomerPlacedVehicleSize              string  `json:"customer_placed_vehicle_size"`
	CustomerLoadCancelled                  string  `json:"customer_load_cancelled"`
	CustomerReportedDateTimeForHaltingCalc string  `json:"customer_reported_date_time_for_halting_calc"`

	// Vendor Section
	VendorPaidDate         string  `json:"vendor_paid_date"`
	VendorCrossDock        string  `json:"vendor_cross_dock"`
	VendorRemark           string  `json:"vendor_remark"`
	VendorBaseRate         float64 `json:"vendor_base_rate"`
	VendorKMCost           float64 `json:"vendor_km_cost"`
	VendorToll             float64 `json:"vendor_toll"`
	VendorTotalHire        float64 `json:"vendor_total_hire"`
	VendorAdvance          float64 `json:"vendor_advance"`
	VendorDebitAmount      float64 `json:"vendor_debit_amount"`
	VendorBalanceAmount    float64 `json:"vendor_balance_amount"`
	VendorPaidBy           string  `json:"vendor_paid_by"`
	VendorLoadUnLoadAmount float64 `json:"vendor_load_unload_amount"`
	VendorHaltingDays      float64 `json:"vendor_halting_days"`
	VendorHaltingPaid      float64 `json:"vendor_halting_paid"`
	VendorExtraDelivery    float64 `json:"vendor_extra_delivery"`
	VendorBreakDown        float64 `json:"vendor_break_down"`
	PodReceived            int64   `json:"pod_received"`
	LoadStatus             string  `json:"load_status"`
	ZonalName              string  `json:"zonal_name"`
}

type XlsUpdateMessge struct {
	Success           *[]Messge `json:"success"`
	ErrorRows         *[]Messge `json:"error_trips"`
	ErrorTripsCount   int       `json:"error_trips_conut"`
	SuccessTripsCount int       `json:"success_trips_conut"`
}

type TripSheetUpdateData struct {
	CustomerInvoiceNo                      string  `json:"customer_invoice_no"`
	CustomerCloseTripDateTime              string  `json:"customer_close_trip_date_time"`
	CustomerPaymentReceivedDate            string  `json:"customer_payment_received_date"`
	CustomerRemark                         string  `json:"customer_remark"`
	CustomerBillingRaisedDate              string  `json:"customer_billing_raised_date"`
	CustomerBaseRate                       float64 `json:"customer_base_rate"`
	CustomerKMCost                         float64 `json:"customer_km_cost"`
	CustomerToll                           float64 `json:"customer_toll"`
	CustomerExtraHours                     float64 `json:"customer_extra_hours"`
	CustomerExtraKM                        float64 `json:"customer_extra_km"`
	CustomerTotalHire                      float64 `json:"customer_total_hire"`
	CustomerDebitAmount                    float64 `json:"customer_debit_amount"`
	CustomerPerLoadHire                    float64 `json:"customer_per_load_hire"`
	CustomerRunningKM                      float64 `json:"customer_running_km"`
	CustomerPerKMPrice                     float64 `json:"customer_per_km_price"`
	CustomerPlacedVehicleSize              string  `json:"customer_placed_vehicle_size"`
	CustomerLoadCancelled                  string  `json:"customer_load_cancelled"`
	CustomerReportedDateTimeForHaltingCalc string  `json:"customer_reported_date_time_for_halting_calc"`
	LoadStatus                             string  `json:"load_status"`
}

type UpdateDraftInvoiceData struct {
	CustomerCloseTripDateTime              string  `json:"customer_close_trip_date_time"`
	CustomerPaymentReceivedDate            string  `json:"customer_payment_received_date"`
	CustomerRemark                         string  `json:"customer_remark"`
	CustomerBillingRaisedDate              string  `json:"customer_billing_raised_date"`
	CustomerBaseRate                       float64 `json:"customer_base_rate"`
	CustomerKMCost                         float64 `json:"customer_km_cost"`
	CustomerToll                           float64 `json:"customer_toll"`
	CustomerExtraHours                     float64 `json:"customer_extra_hours"`
	CustomerExtraKM                        float64 `json:"customer_extra_km"`
	CustomerTotalHire                      float64 `json:"customer_total_hire"`
	CustomerDebitAmount                    float64 `json:"customer_debit_amount"`
	CustomerPerLoadHire                    float64 `json:"customer_per_load_hire"`
	CustomerRunningKM                      float64 `json:"customer_running_km"`
	CustomerPerKMPrice                     float64 `json:"customer_per_km_price"`
	CustomerLoadCancelled                  string  `json:"customer_load_cancelled"`
	CustomerReportedDateTimeForHaltingCalc string  `json:"customer_reported_date_time_for_halting_calc"`
	LoadStatus                             string  `json:"load_status"`
}

type TripSheetDraftPull struct {
	TripSheetID         int64   `json:"trip_sheet_id"`
	TripSheetNum        string  `json:"trip_sheet_num"`
	CustomerPerLoadHire float64 `json:"customer_per_load_hire"`
	CustomerRunningKM   float64 `json:"customer_running_km"`
	CustomerPerKMPrice  float64 `json:"customer_per_km_price"`
	CustomerBaseRate    float64 `json:"customer_base_rate"`
	TripType            string  `json:"trip_type"`
	CustomerKMCost      float64 `json:"customer_km_cost"`
	CustomerToll        float64 `json:"customer_toll"`
}

type CustomerInvoice struct {
	ID            int64   `json:"id"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceRef    string  `json:"invoice_ref"`
	WorkType      string  `json:"work_type"`
	WorkStartDate string  `json:"work_start_date"`
	WorkDateDate  string  `json:"work_date_date"`
	DocumentDate  string  `json:"document_date"`
	InvoiceStatus string  `json:"invoice_status"`
	InvoiceAmount float64 `json:"invoice_amount"`
	InvoiceDate   string  `json:"invoice_date"`
	PaymentDate   string  `json:"payment_date"`
	TripRef       string  `json:"trip_ref"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type CreateDraftInvoice struct {
	InvoiceRefID  string  `json:"invoice_ref"`
	WorkType      string  `json:"work_type"`
	WorkStartDate string  `json:"work_start_date"`
	WorkEndDate   string  `json:"work_end_date"`
	DocumentDate  string  `json:"document_date"`
	InvoiceStatus string  `json:"invoice_status"`
	InvoiceAmount float64 `json:"invoice_amount"`
	TripRef       string  `json:"trip_ref"`
	CustomerName  string  `json:"customer_name"`
	CustomerId    int64   `json:"customer_id"`
	CustomerCode  string  `json:"customer_code"`
}

type XlsDraftInvoiceMessge struct {
	Success           *[]Messge `json:"success"`
	ErrorRows         *[]Messge `json:"error_trips"`
	ErrorTripsCount   int       `json:"error_trips_conut"`
	SuccessTripsCount int       `json:"success_trips_conut"`
	Msg               string    `json:"msg"`
}
