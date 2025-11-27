package dtos

type DashBoardRes struct {
	LoadUnloadPoints *LoadUnloadPoints `json:"load_unload_points"`
	Employee         *Employee         `json:"employee"`
	Customer         *CustomerStats    `json:"customer"`
	Vendor           *VendorStats      `json:"vendor"`
	Vehicle          *VehicleStats     `json:"vehicle"`
	TripCreated      *TripStatusStats  `json:"trip_created"`
	TripSubmitted    *TripStatusStats  `json:"trip_submitted"`
	TripDelivered    *TripStatusStats  `json:"trip_delivered"`
	TripClosed       *TripStatusStats  `json:"trip_closed"`
	TripCompleted    *TripStatusStats  `json:"trip_completed"`
	TripGraphStats   *[]TripSummary    `json:"trip_graph_stats"`
	InvoiceStats     InvStatusResponse `json:"invoice_stats"`
}

type Employee struct {
	TotalCount    int64 `json:"total_count"`
	ActiveCount   int64 `json:"active_count"`
	InActiveCount int64 `json:"in_active_count"`
}
type LoadUnloadPoints struct {
	TotalCount    int64 `json:"total_count"`
	ActiveCount   int64 `json:"active_count"`
	InActiveCount int64 `json:"in_active_count"`
}

type CustomerStats struct {
	TotalCount    int64 `json:"total_count"`
	ActiveCount   int64 `json:"active_count"`
	InActiveCount int64 `json:"in_active_count"`
}
type VendorStats struct {
	TotalCount    int64 `json:"total_count"`
	ActiveCount   int64 `json:"active_count"`
	InActiveCount int64 `json:"in_active_count"`
}
type VehicleStats struct {
	TotalCount    int64 `json:"total_count"`
	ActiveCount   int64 `json:"active_count"`
	InActiveCount int64 `json:"in_active_count"`
}

type TripStatusStats struct {
	TotalCount            int64 `json:"total_count"`
	LocalScheduledTrip    int64 `json:"local_scheduled_trip"`
	LocalAdhocTrip        int64 `json:"local_adhoc_trip"`
	LineHaulScheduledTrip int64 `json:"line_haul_scheduled_trip"`
	LineHaulAdhocTrip     int64 `json:"line_haul_adhoc_trip"`
}

type TripCreated struct {
	TotalCount            int64 `json:"total_count"`
	LocalScheduledTrip    int64 `json:"local_scheduled_trip"`
	LocalAdhocTrip        int64 `json:"local_adhoc_trip"`
	LineHaulScheduledTrip int64 `json:"line_haul_scheduled_trip"`
	LineHaulAdhocTrip     int64 `json:"line_haul_adhoc_trip"`
}
type TripSubmitted struct {
	TotalCount            int64 `json:"total_count"`
	LocalScheduledTrip    int64 `json:"local_scheduled_trip"`
	LocalAdhocTrip        int64 `json:"local_adhoc_trip"`
	LineHaulScheduledTrip int64 `json:"line_haul_scheduled_trip"`
	LineHaulAdhocTrip     int64 `json:"line_haul_adhoc_trip"`
}
type TripDelivered struct {
	TotalCount            int64 `json:"total_count"`
	LocalScheduledTrip    int64 `json:"local_scheduled_trip"`
	LocalAdhocTrip        int64 `json:"local_adhoc_trip"`
	LineHaulScheduledTrip int64 `json:"line_haul_scheduled_trip"`
	LineHaulAdhocTrip     int64 `json:"line_haul_adhoc_trip"`
}
type TripClosed struct {
	TotalCount            int64 `json:"total_count"`
	LocalScheduledTrip    int64 `json:"local_scheduled_trip"`
	LocalAdhocTrip        int64 `json:"local_adhoc_trip"`
	LineHaulScheduledTrip int64 `json:"line_haul_scheduled_trip"`
	LineHaulAdhocTrip     int64 `json:"line_haul_adhoc_trip"`
}

type TripGraphStats struct {
	Date       string         `json:"date"`
	TotalTrips int            `json:"totalTrips"`
	ByType     map[string]int `json:"byType"`
}

type TripCountWithType struct {
	Date          string `json:"date"`
	TripSheetType string `json:"trip_sheet_type"`
	TripCount     int64  `json:"tripCount"`
}

type TripSummary struct {
	Date       string           `json:"date"`
	TotalTrips int64            `json:"totalTrips"`
	ByType     map[string]int64 `json:"byType"`
}
