package services

import (
	"fmt"
	"math"

	"solar-backend/internal/clients"
	"solar-backend/internal/models"
)

const (
	tarifPLN          = 1444.0 // Rp 1.444 per kWh (standar PLN)
	performanceRatio  = 0.75
	kapasitasPerPanel = 0.55
	luasPerPanel      = 2.5
	hargaPerKWp       = 15000000.0
	chartYears        = 15
)

func CalculateFeasibility(req models.SimulationRequest) (*models.SimulationResponse, error) {
	if req.Lokasi == "" {
		return nil, fmt.Errorf("lokasi is required")
	}
	if req.TagihanBulanan <= 0 {
		return nil, fmt.Errorf("tagihan_bulanan must be greater than 0")
	}

	// Gunakan Google Maps Geocoding API untuk mendapatkan koordinat
	lat, lon, err := clients.GetCoordinatesFromGoogleMaps(req.Lokasi)
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinates: %w", err)
	}

	// Ambil data radiasi matahari dari NASA POWER API
	radiasi, _, err := clients.GetClimatologyData(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("failed to get climatology data: %w", err)
	}
	if radiasi <= 0 {
		return nil, fmt.Errorf("invalid radiation data")
	}

	// Layer 1: Energy model
	// Gunakan faktor musiman sederhana agar produksi tidak terlalu optimistis.
	seasonalFactor := 0.9
	effectiveRadiation := radiasi * seasonalFactor

	// Hitung kebutuhan kapasitas berdasarkan tagihan bulanan.
	kebutuhanBulananKwh := req.TagihanBulanan / tarifPLN
	kebutuhanHarianKwh := kebutuhanBulananKwh / 30.0
	requiredKwp := kebutuhanHarianKwh / (effectiveRadiation * performanceRatio)

	// Layer 2: System constraint
	// DESIGN LIMIT: batas offset ini bukan constraint fisik, hanya batas desain
	// agar sistem tidak over-sized terhadap profil konsumsi yang dimodelkan.
	maxSystemKwp := 8.0
	defaultMaxRoofArea := 40.0
	maxOffsetRatio := 0.75
	maxOffsetKwh := kebutuhanBulananKwh * maxOffsetRatio

	maxCandidateKwp := math.Min(requiredKwp, maxSystemKwp)
	if maxCandidateKwp < 2.0 {
		maxCandidateKwp = 2.0
	}

	candidateKwp := make([]float64, 0, int(math.Ceil(maxCandidateKwp/kapasitasPerPanel)))
	for kwp := 2.0; kwp <= maxCandidateKwp+1e-9; kwp += kapasitasPerPanel {
		candidateKwp = append(candidateKwp, kwp)
	}
	if len(candidateKwp) == 0 {
		candidateKwp = append(candidateKwp, 2.0)
	}

	type scenarioResult struct {
		jumlahPanel       int
		actualKwp         float64
		luasAtap          float64
		produksiBulanan   float64
		effectiveOffset   float64
		savingPerBulan    float64
		savingTahunan     float64
		roiTahun          float64
		coverageRatio     float64
		estimasiBiaya     float64
		breakEvenYear     int
		requestedKwp      float64
		effectiveScore    float64
	}

	validScenarios := make([]scenarioResult, 0, len(candidateKwp))
	for _, requestedKwp := range candidateKwp {
		jumlahPanel := int(math.Ceil(requestedKwp / kapasitasPerPanel))
		actualKwp := float64(jumlahPanel) * kapasitasPerPanel
		luasAtap := float64(jumlahPanel) * luasPerPanel

		// Constraint fisik: luas atap.
		if luasAtap > defaultMaxRoofArea {
			continue
		}

		// Layer 3: Economic model
		estimasiBiaya := actualKwp * hargaPerKWp
		produksiBulananKwh := actualKwp * effectiveRadiation * performanceRatio * 30.0
		effectiveOffsetKwh := math.Min(produksiBulananKwh, maxOffsetKwh)
		savingPerBulan := effectiveOffsetKwh * tarifPLN
		if savingPerBulan <= 0 {
			continue
		}

		savingTahunan := savingPerBulan * 12.0
		roiTahun := estimasiBiaya / savingTahunan

		coverageRatio := effectiveOffsetKwh / kebutuhanBulananKwh
		breakEvenYear := int(math.Ceil(roiTahun))
		effectiveScore := 0.7*roiTahun - 0.3*coverageRatio
		if roiTahun > 15.0 {
			effectiveScore += 5.0
		}

		validScenarios = append(validScenarios, scenarioResult{
			jumlahPanel:     jumlahPanel,
			actualKwp:       actualKwp,
			luasAtap:        luasAtap,
			produksiBulanan: produksiBulananKwh,
			effectiveOffset: effectiveOffsetKwh,
			savingPerBulan:  savingPerBulan,
			savingTahunan:   savingTahunan,
			roiTahun:        roiTahun,
			coverageRatio:   coverageRatio,
			estimasiBiaya:   estimasiBiaya,
			breakEvenYear:   breakEvenYear,
			requestedKwp:    requestedKwp,
			effectiveScore:  effectiveScore,
		})
	}

	if len(validScenarios) == 0 {
		return nil, fmt.Errorf("no feasible scenario found within roof area and ROI constraints")
	}

	// Step 3: Pilih skenario terbaik menggunakan scoring gabungan.
	// Score menyeimbangkan ROI dan coverage, sehingga tidak bias ke sistem kecil.
	best := validScenarios[0]
	for i := 1; i < len(validScenarios); i++ {
		s := validScenarios[i]
		if s.effectiveScore < best.effectiveScore {
			best = s
			continue
		}
		if math.Abs(s.effectiveScore-best.effectiveScore) <= 0.1 && s.coverageRatio > best.coverageRatio {
			best = s
		}
	}

	jumlahPanel := best.jumlahPanel
	actualKwp := best.actualKwp
	luasAtap := best.luasAtap
	produksiBulananKwh := best.produksiBulanan
	effectiveOffsetKwh := best.effectiveOffset
	savingPerBulan := best.savingPerBulan
	savingTahunan := best.savingTahunan
	roiTahun := best.roiTahun
	estimasiBiaya := best.estimasiBiaya
	breakEvenYear := best.breakEvenYear
	coverageRatio := best.coverageRatio
	capacityCapped := requiredKwp > best.requestedKwp

	// Apply rounding ke 2 desimal untuk nilai finansial dan ROI
	estimasiBiaya = math.Round(estimasiBiaya*100) / 100
	savingPerBulan = math.Round(savingPerBulan*100) / 100
	roiTahun = math.Round(roiTahun*100) / 100
	actualKwp = math.Round(actualKwp*100) / 100
	luasAtap = math.Round(luasAtap*100) / 100
	produksiBulananKwh = math.Round(produksiBulananKwh*100) / 100
	effectiveOffsetKwh = math.Round(effectiveOffsetKwh*100) / 100

	// Tentukan status dan alasan (Bahasa Indonesia) berdasarkan ROI dan sertakan angka
	status := "Tidak Layak"
	reasoning := ""
	if roiTahun < 7.0 {
		status = "Sangat Layak"
		reasoning = fmt.Sprintf("Investasi tergolong kategori sangat layak dengan estimasi balik modal %.2f tahun, jauh di bawah ambang 7 tahun. Pengembalian cepat ini didukung penghematan stabil sebesar %.0f Rp per bulan (%.0f%% offset konsumsi) dengan kapasitas sistem %.2f kWp. Proyeksi ini mengasumsikan kondisi radiasi normal dan tarif listrik stabil.", roiTahun, savingPerBulan, coverageRatio*100, actualKwp)
	} else if roiTahun <= 10.0 {
		status = "Layak"
		reasoning = fmt.Sprintf("Investasi tergolong kategori pengembalian menengah (%.2f tahun), yang masih berada dalam batas kelayakan investasi energi, namun belum masuk kategori optimal (<7 tahun). Penghematan stabil %.0f Rp per bulan (%.0f%% offset konsumsi) dengan kapasitas sistem %.2f kWp. Keputusan investasi bergantung pada preferensi Anda terhadap jangka panjang pengembalian modal.", roiTahun, savingPerBulan, coverageRatio*100, actualKwp)
	} else {
		status = "Tidak Layak"
		reasoning = fmt.Sprintf("Investasi tidak masuk kategori layak karena estimasi balik modal %.2f tahun melebihi ambang 10 tahun. Meskipun sistem dapat menutup %.0f%% konsumsi (%.0f Rp per bulan), periode pengembalian modal yang panjang melampaui batasan kelayakan ekonomi. Skalabilitas dan optimasi biaya sangat perlu untuk memperbaiki kelayakan investasi ini.", roiTahun, coverageRatio*100, savingPerBulan)
	}

	// Generate ROI chart data untuk 15 tahun dengan degradasi tahunan 0.5%.
	roiChartData := generateROIChartData(estimasiBiaya, savingTahunan, chartYears)

	warnings := make([]string, 0, 6)

	// Priority 1: ROI near threshold (most important - affects feasibility decision)
	distTo7 := math.Abs(roiTahun - 7.0)
	distTo10 := math.Abs(roiTahun - 10.0)
	if distTo7 < 0.5 {
		warnings = append(warnings, fmt.Sprintf("ROI %.2f tahun berada dekat ambang batas %.1f tahun untuk kategori sangat layak, sehingga sensitivitas terhadap perubahan asumsi (tarif, radiasi) cukup tinggi.", roiTahun, 7.0))
	} else if distTo10 < 0.5 {
		warnings = append(warnings, fmt.Sprintf("ROI %.2f tahun berada dekat ambang batas %.1f tahun antara layak dan tidak layak, sehingga perubahan kecil pada asumsi dapat mempengaruhi rekomendasi investasi.", roiTahun, 10.0))
	}

	// Priority 2: Assumptions uncertainty (affects confidence and projection accuracy)
	uncertaintyWarning := "Perkiraan mengasumsikan radiasi matahari dan tarif listrik konstan; produksi aktual dapat bervariasi musiman."
	warnings = append(warnings, uncertaintyWarning)

	// Priority 3: Shading and technical data missing (can impact future accuracy)
	shadingWarning := "Estimasi tidak mempertimbangkan potensi shading dari bangunan atau pohon sekitar."
	warnings = append(warnings, shadingWarning)

	// Priority 4: Coverage ratio informational
	if coverageRatio < 1.0 {
		warnings = append(warnings, fmt.Sprintf("Sistem ini dirancang untuk menutup sekitar %.0f%% dari konsumsi listrik bulanan Anda, bukan 100%%, untuk menjaga efisiensi biaya dan keterbatasan area atap.", coverageRatio*100))
	}

	// Priority 5: Capacity capped informational
	if capacityCapped {
		warnings = append(warnings, fmt.Sprintf("Kapasitas evaluasi dibatasi hingga %.1f kWp untuk menjaga kesesuaian skala residensial dan keterbatasan area atap umum (30–40 m²).", maxSystemKwp))
	}

	// Priority 6: Low radiation informational
	if radiasi < 4.0 {
		warnings = append(warnings, "Intensitas matahari di lokasi ini rendah, mempengaruhi kecepatan balik modal.")
	}

	// Struktur financials dengan pembulatan dan coverage ratio
	financials := models.Financials{
		TotalInvestment:    estimasiBiaya,
		MonthlySavings:     savingPerBulan,
		PaybackPeriodYears: roiTahun,
		BreakEvenYear:      breakEvenYear,
		CoverageRatio:      math.Round(coverageRatio*100) / 100,
	}

	// Struktur assumptions (numeric)
	assumptions := models.Assumptions{
		ElectricityTariff: tarifPLN,
		SystemEfficiency:  performanceRatio,
		AnnualDegradation: 0.005,
	}

	// Confidence: compute numeric score (0..1) and factors
	// data completeness: check presence of key inputs
	dataCompletenessScore := 1.0
	if req.TagihanBulanan <= 0 || req.Lokasi == "" {
		dataCompletenessScore = 0.5
	}

	// irradiation variability: higher radiation -> more stable
	irradiationScore := 0.6
	if radiasi >= 5.0 {
		irradiationScore = 1.0
	} else if radiasi >= 4.0 {
		irradiationScore = 0.8
	} else if radiasi >= 3.0 {
		irradiationScore = 0.5
	} else {
		irradiationScore = 0.3
	}

	// financial assumptions: penalize if defaults used (no shading, orientation, consumption profile)
	// Score 0.5 reflects lack of custom data inputs
	financialAssumptionsScore := 0.5
	// (If in future user provides shading, orientation, or hourly consumption data, increase score)

	// Weighted confidence score: reweight to emphasize financial assumptions importance
	// weights: 0.3 data completeness, 0.3 irradiation, 0.4 financial assumptions
	confidenceScore := 0.3*dataCompletenessScore + 0.3*irradiationScore + 0.4*financialAssumptionsScore
	// clamp and round
	if confidenceScore < 0.0 {
		confidenceScore = 0.0
	}
	if confidenceScore > 1.0 {
		confidenceScore = 1.0
	}
	confidenceScore = math.Round(confidenceScore*100) / 100

	// map numeric to label (conservative: HIGH only if very high confidence)
	confidenceLabel := "MEDIUM"
	if confidenceScore >= 0.85 {
		confidenceLabel = "HIGH"
	} else if confidenceScore >= 0.65 {
		confidenceLabel = "MEDIUM"
	} else {
		confidenceLabel = "LOW"
	}

	confidenceFactors := models.ConfidenceFactors{
		DataCompleteness:       dataCompletenessScore,
		IrradiationVariability: irradiationScore,
		FinancialAssumptions:   financialAssumptionsScore,
	}

	// Optional display map (UI formatted strings) separated from engine data
	displayMap := map[string]string{
		"total_investment":    fmt.Sprintf("Rp %.0f", estimasiBiaya),
		"monthly_savings":     fmt.Sprintf("Rp %.0f", savingPerBulan),
		"payback_period":      fmt.Sprintf("%.2f tahun", roiTahun),
		"system_capacity_kwp": fmt.Sprintf("%.2f kWp", actualKwp),
	}

	// Decision rules (transparent)
	decisionRules := models.DecisionRules{
		PaybackPeriod: map[string]string{
			"<7":   "Sangat Layak",
			"7-10": "Layak",
			">10":  "Tidak Layak",
		},
	}

	// Determine applied decision rule and basis
	appliedRule := ""
	decisionResult := ""
	if roiTahun < 7.0 {
		appliedRule = "<7"
		decisionResult = decisionRules.PaybackPeriod[appliedRule]
	} else if roiTahun <= 10.0 {
		appliedRule = "7-10"
		decisionResult = decisionRules.PaybackPeriod[appliedRule]
	} else {
		appliedRule = ">10"
		decisionResult = decisionRules.PaybackPeriod[appliedRule]
	}

	decisionBasis := models.DecisionBasis{
		Metric:      "payback_period_years",
		Value:       roiTahun,
		AppliedRule: appliedRule,
		Result:      decisionResult,
	}

	return &models.SimulationResponse{
		Status:            status,
		Reasoning:         reasoning,
		Confidence:        confidenceLabel,
		ConfidenceScore:   confidenceScore,
		ConfidenceFactors: confidenceFactors,
		Financials:        financials,
		TechnicalRecommendation: models.TechnicalRecommendation{
			PanelCount:        jumlahPanel,
			SystemCapacityKwp: actualKwp,
			RequiredAreaM2:    luasAtap,
		},
		Assumptions:   assumptions,
		ROIChartData:  roiChartData,
		Warnings:      warnings,
		DecisionRules: decisionRules,
		DecisionBasis: decisionBasis,
		Display:       displayMap,
	}, nil
}

// generateROIChartData membuat array data chart ROI selama 15 tahun
// dengan akumulasi profit (Saving Tahunan - Investasi Awal)
func generateROIChartData(investasi float64, savingTahunan float64, years int) []models.ROIChartData {
	chartData := make([]models.ROIChartData, years+1)
	degradationRate := 0.005

	// Year 0: only initial investment (negative)
	chartData[0] = models.ROIChartData{
		Year:             0,
		CumulativeProfit: math.Round(-investasi*100) / 100,
	}

	// Year 1..N: cumulative profit memakai saving yang terdegradasi setiap tahun.
	akumulasi := -investasi
	for i := 1; i <= years; i++ {
		savingYearN := savingTahunan * math.Pow(1-degradationRate, float64(i))
		akumulasi += savingYearN
		chartData[i] = models.ROIChartData{
			Year:             i,
			CumulativeProfit: math.Round(akumulasi*100) / 100,
		}
	}

	return chartData
}
