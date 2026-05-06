package models

type SimulationRequest struct {
	Lokasi         string  `json:"lokasi"`
	TagihanBulanan float64 `json:"tagihan_bulanan"`
}

type TechnicalRecommendation struct {
	PanelCount        int     `json:"panel_count"`
	SystemCapacityKwp float64 `json:"system_capacity_kwp"`
	RequiredAreaM2    float64 `json:"required_area_m2"`
}

type Financials struct {
	TotalInvestment    float64 `json:"total_investment"`
	MonthlySavings     float64 `json:"monthly_savings"`
	PaybackPeriodYears float64 `json:"payback_period_years"`
	BreakEvenYear      int     `json:"break_even_year"`
	CoverageRatio      float64 `json:"coverage_ratio"`
}

type Assumptions struct {
	ElectricityTariff float64 `json:"electricity_tariff"`
	SystemEfficiency  float64 `json:"system_efficiency"`
	AnnualDegradation float64 `json:"annual_degradation"`
}

type ROIChartData struct {
	Year             int     `json:"year"`
	CumulativeProfit float64 `json:"cumulative_profit"`
}

type ConfidenceFactors struct {
	DataCompleteness       float64 `json:"data_completeness"`
	IrradiationVariability float64 `json:"irradiation_variability"`
	FinancialAssumptions   float64 `json:"financial_assumptions"`
}

type DecisionRules struct {
	PaybackPeriod map[string]string `json:"payback_period"`
}

type DecisionBasis struct {
	Metric      string  `json:"metric"`
	Value       float64 `json:"value"`
	AppliedRule string  `json:"applied_rule"`
	Result      string  `json:"result"`
}

type SimulationResponse struct {
	Status                  string                  `json:"status"`
	Reasoning               string                  `json:"reasoning"`
	Confidence              string                  `json:"confidence"`
	ConfidenceScore         float64                 `json:"confidence_score"`
	ConfidenceFactors       ConfidenceFactors       `json:"confidence_factors"`
	Financials              Financials              `json:"financials"`
	TechnicalRecommendation TechnicalRecommendation `json:"technical_recommendation"`
	Assumptions             Assumptions             `json:"assumptions"`
	ROIChartData            []ROIChartData          `json:"roi_chart_data"`
	Warnings                []string                `json:"warnings"`
	DecisionRules           DecisionRules           `json:"decision_rules"`
	DecisionBasis           DecisionBasis           `json:"decision_basis"`
	// Display provides optional UI-formatted strings separated from engine data
	Display map[string]string `json:"display,omitempty"`
}
