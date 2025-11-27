package dtos

type LRReceiptReq struct {
	TripSheetID        int64   `json:"trip_sheet_id"`
	TripSheetNum       string  `json:"trip_sheet_num"`
	LRNumber           string  `json:"lr_number"`
	TripDate           string  `json:"trip_date"`
	DriverName         string  `json:"driver_name"`
	DriverMobileNumber string  `json:"driver_mobile_number"`
	VehicleNumber      string  `json:"vehicle_number"`
	VehicleSize        string  `json:"vehicle_size"`
	InvoiceNumber      string  `json:"invoice_number"`
	InvoiceValue       float64 `json:"invoice_value"`
	ConsignorName      string  `json:"consignor_name"`
	ConsignorAddress   string  `json:"consignor_address"`
	ConsignorGST       string  `json:"consignor_gst"`
	ConsigneeName      string  `json:"consignee_name"`
	ConsigneeAddress   string  `json:"consignee_address"`
	ConsigneeGST       string  `json:"consignee_gst"`
	GoodsType          string  `json:"goods_type"`
	GoodsWeight        string  `json:"goods_weight"`
	QuantityInPieces   string  `json:"quantity_in_pieces"`
	Remark             string  `json:"remark"`
	OrgID              int64   `json:"org_id"`
}

type LRReceiptUpdateReq struct {
	TripSheetID  int64  `json:"trip_sheet_id"`
	TripSheetNum string `json:"trip_sheet_num"`
	LRNumber     string `json:"lr_number"`
	TripDate     string `json:"trip_date"`
	// FromAddress      string  `json:"from_address"`
	// ToAddress        string  `json:"to_address"`
	VehicleNumber      string  `json:"vehicle_number"`
	VehicleSize        string  `json:"vehicle_size"`
	InvoiceNumber      string  `json:"invoice_number"`
	InvoiceValue       float64 `json:"invoice_value"`
	ConsignorName      string  `json:"consignor_name"`
	ConsignorAddress   string  `json:"consignor_address"`
	ConsignorGST       string  `json:"consignor_gst"`
	ConsigneeName      string  `json:"consignee_name"`
	ConsigneeAddress   string  `json:"consignee_address"`
	ConsigneeGST       string  `json:"consignee_gst"`
	GoodsType          string  `json:"goods_type"`
	GoodsWeight        string  `json:"goods_weight"`
	QuantityInPieces   string  `json:"quantity_in_pieces"`
	Remark             string  `json:"remark"`
	DriverName         string  `json:"driver_name"`
	DriverMobileNumber string  `json:"driver_mobile_number"`
	OrgID              int64   `json:"org_id"`
}

type LRReceipt struct {
	LRId               int64   `json:"lr_id"`
	TripSheetID        int64   `json:"trip_sheet_id"`
	TripSheetNum       string  `json:"trip_sheet_num"`
	LRNumber           string  `json:"lr_number"`
	TripDate           string  `json:"trip_date"`
	FromAddress        string  `json:"from_address"`
	ToAddress          string  `json:"to_address"`
	VehicleNumber      string  `json:"vehicle_number"`
	VehicleSize        string  `json:"vehicle_size"`
	InvoiceNumber      string  `json:"invoice_number"`
	InvoiceValue       float64 `json:"invoice_value"`
	ConsignorName      string  `json:"consignor_name"`
	ConsignorAddress   string  `json:"consignor_address"`
	ConsignorGST       string  `json:"consignor_gst"`
	ConsigneeName      string  `json:"consignee_name"`
	ConsigneeAddress   string  `json:"consignee_address"`
	ConsigneeGST       string  `json:"consignee_gst"`
	GoodsType          string  `json:"goods_type"`
	GoodsWeight        string  `json:"goods_weight"`
	QuantityInPieces   string  `json:"quantity_in_pieces"`
	Remark             string  `json:"remark"`
	DriverName         string  `json:"driver_name"`
	DriverMobileNumber string  `json:"driver_mobile_number"`
	OrgID              int64   `json:"org_id"`
}

type LRRecords struct {
	LRRecords *[]LRReceipt `json:"lr_records"`
	Total     int64        `json:"total"`
	Limit     int64        `json:"limit"`
	OffSet    int64        `json:"offset"`
}
