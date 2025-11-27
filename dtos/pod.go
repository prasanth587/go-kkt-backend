package dtos

type ManagePodReq struct {
	TripSheetID         int64   `json:"trip_sheet_id"`
	TripSheetNum        string  `json:"trip_sheet_num"`
	LRNumber            string  `json:"lr_number"`
	CustomerID          int64   `json:"customer_id"`
	CustomerName        string  `json:"customer_name"`
	PodUpload           string  `json:"pod_doc"`
	UnPoadingCharges    float64 `json:"unloading_charges"`
	UnPoadingDate       string  `json:"unloading_date"`
	PaidBy              string  `json:"paid_by"`
	PODSubmitedDate     string  `json:"pod_submited_date"`
	LateSubmissionDebit float64 `json:"late_submission_debit"`
	PodRemark           string  `json:"pod_remark"`
	PodStatus           string  `json:"pod_status"`
	OrgId               int64   `json:"org_id"`
	TripType            string  `json:"trip_type"`
	HaltingAmount       float64 `json:"halting_amount"`
	RunningKM           float64 `json:"running_km"`
	HaltingDays         float64 `json:"halting_days"`
}

type UpdateManagePodReq struct {
	TripSheetID         int64   `json:"trip_sheet_id"`
	TripSheetNum        string  `json:"trip_sheet_num"`
	LRNumber            string  `json:"lr_number"`
	CustomerID          int64   `json:"customer_id"`
	CustomerName        string  `json:"customer_name"`
	PodUpload           string  `json:"pod_doc"`
	UnPoadingCharges    float64 `json:"unloading_charges"`
	UnPoadingDate       string  `json:"unloading_date"`
	PaidBy              string  `json:"paid_by"`
	PODSubmitedDate     string  `json:"pod_submited_date"`
	LateSubmissionDebit float64 `json:"late_submission_debit"`
	PodRemark           string  `json:"pod_remark"`
	PodStatus           string  `json:"pod_status"`
	OrgId               int64   `json:"org_id"`
	TripType            string  `json:"trip_type"`
	HaltingAmount       float64 `json:"halting_amount"`
	RunningKM           float64 `json:"running_km"`
	HaltingDays         float64 `json:"halting_days"`
}

type ManagePod struct {
	PodId               int64   `json:"pod_id"`
	TripSheetID         int64   `json:"trip_sheet_id"`
	LoadingPointIDs     string  `json:"loading_points"`
	UnLoadingPointIDs   string  `json:"un_loading_points"`
	TripSheetNum        string  `json:"trip_sheet_num"`
	LRNumber            string  `json:"lr_number"`
	SubmitedDate        string  `json:"submited_date"`
	CustomerName        string  `json:"customer_name"`
	CustomerID          int64   `json:"customer_id"`
	SendBy              string  `json:"send_by"`
	PodUpload           string  `json:"pod_doc"`
	UnPoadingCharges    float64 `json:"unloading_charges"`
	UnPoadingDate       string  `json:"unloading_date"`
	PaidBy              string  `json:"paid_by"`
	PODSubmitedDate     string  `json:"pod_submited_date"`
	LateSubmissionDebit float64 `json:"late_submission_debit"`
	PodRemark           string  `json:"pod_remark"`
	PodStatus           string  `json:"pod_status"`
	OrgId               int64   `json:"org_id"`
	TripType            string  `json:"trip_type"`
	HaltingAmount       float64 `json:"halting_amount"`
	KilometersCovered   float64 `json:"kilometers_covered"`
	HaltingDays         float64 `json:"halting_days"`
}
type ManagePods struct {
	ManagePod *[]ManagePod `json:"pods"`
	Total     int64        `json:"total"`
	Limit     int64        `json:"limit"`
	OffSet    int64        `json:"offset"`
}

type PodLoadUnLoad struct {
	TripSheetId    int64  `json:"trip_sheet_id"`
	LoadingPointId int64  `json:"loading_point_id"`
	CityCode       string `json:"city_code"`
	CityName       string `json:"city_name"`
	Type           string `json:"type"`
}

type LoadUnLoadLoc struct {
	LoadingPointId int64  `json:"loading_point_id"`
	CityCode       string `json:"city_code"`
	CityName       string `json:"city_name"`
}
