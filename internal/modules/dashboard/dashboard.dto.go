package dashboard

type ParkingFlowQuery struct {
	LotID *uint64 `form:"lotId"`
	Date  string  `form:"date"`
}

type DashboardSummaryResponse struct {
	TodayIn         int64   `json:"todayIn"`
	TodayOut        int64   `json:"todayOut"`
	CurrentVehicles int64   `json:"currentVehicles"`
	Capacity        int64   `json:"capacity"`
	AvailableSlots  int64   `json:"availableSlots"`
	OccupancyRate   float64 `json:"occupancyRate"`
}

type HourlyFlowResponse struct {
	Hour string `json:"hour"`
	In   int64  `json:"in"`
	Out  int64  `json:"out"`
}

type PeakTimeResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type DashboardInsightsResponse struct {
	PeakTime        PeakTimeResponse `json:"peakTime"`
	OccupancyStatus string           `json:"occupancyStatus"`
	StatusMessage   string           `json:"statusMessage"`
}

type ParkingFlowResponse struct {
	Date       string                    `json:"date"`
	LotID      *uint64                   `json:"lotId"`
	LotName    string                    `json:"lotName"`
	Summary    DashboardSummaryResponse  `json:"summary"`
	HourlyFlow []HourlyFlowResponse      `json:"hourlyFlow"`
	Insights   DashboardInsightsResponse `json:"insights"`
}

type RevenueDateQuery struct {
	LotID *uint64 `form:"lotId"`
	Date  string  `form:"date" binding:"required"` // yyyy-mm-dd
}
