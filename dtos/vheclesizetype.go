package dtos

type VehicleSizeType struct {
	VehicleSize string `json:"vehicle_size"`
	VehicleType string `json:"vehicle_type"`
	IsActive    int    `json:"is_active"`
	Status      string `json:"status"`
}

type VehicleSizeTypeObj struct {
	VehicleSizeId int64  `json:"vehicle_size_id"`
	VehicleSize   string `json:"vehicle_size"`
	VehicleType   string `json:"vehicle_type"`
	IsActive      int64  `json:"is_active"`
	Status        string `json:"status"`
}
type VehicleSizeTypePre struct {
	VehicleSizeId int64  `json:"vehicle_size_id"`
	VehicleSize   string `json:"vehicle_size"`
	VehicleType   string `json:"vehicle_type"`
}

type VehicleSizeTypeEntries struct {
	VehicleSizeType *[]VehicleSizeTypeObj `json:"vehicles_size_type"`
	Total           int64                 `json:"total"`
	Limit           int64                 `json:"limit"`
	OffSet          int64                 `json:"offset"`
}

type VehicleSizeTypeUpdate struct {
	VehicleSize string `json:"vehicle_size"`
	VehicleType string `json:"vehicle_type"`
	IsActive    int    `json:"is_active"`
	Status      string `json:"status"`
}

type WebsiteScreen struct {
	WebsiteScreenID int64  `json:"website_screen_id"`
	ScreenName      string `json:"screen_name"`
	Description     string `json:"description"`
}

type PermissionLabel struct {
	PermisstionLabelID int64  `json:"permisstion_label_id"`
	PermisstionLabel   string `json:"permisstion_label"`
	Description        string `json:"description"`
}

type EmployeesPre struct {
	EmployeeId   int64  `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
}
