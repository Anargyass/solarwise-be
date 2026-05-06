# SolarWise Backend - API Simulation Endpoint

## Endpoint: POST `/api/v1/simulation`

Endpoint ini melakukan analisis kelayakan investasi panel surya berdasarkan lokasi dan tagihan listrik bulanan.

### Request Body

```json
{
  "lokasi": "Surabaya",
  "tagihan_bulanan": 1000000
}
```

| Parameter | Tipe | Deskripsi | Contoh |
|-----------|------|-----------|--------|
| `lokasi` | string | Nama lokasi/kota untuk geocoding | "Surabaya", "Jakarta", "Bandung" |
| `tagihan_bulanan` | number | Tagihan listrik bulanan dalam Rupiah | 1000000 |

### Response Format

```json
{
  "kelayakan": "SANGAT LAYAK",
  "estimasi_biaya": 82500000,
  "saving_per_bulan": 1000000,
  "roi_tahun": 6.875,
  "rekomendasi": {
    "jumlah_panel": 10,
    "kapasitas_kwp": 5.5,
    "luas_atap_m2": 25
  },
  "roi_chart_data": [
    { "tahun": 0, "akumulasi_profit": -82500000 },
    { "tahun": 1, "akumulasi_profit": -70500000 },
    { "tahun": 2, "akumulasi_profit": -58500000 },
    { "tahun": 3, "akumulasi_profit": -46500000 },
    { "tahun": 4, "akumulasi_profit": -34500000 },
    { "tahun": 5, "akumulasi_profit": -22500000 },
    { "tahun": 6, "akumulasi_profit": -10500000 },
    { "tahun": 7, "akumulasi_profit": 1500000 },
    { "tahun": 8, "akumulasi_profit": 13500000 },
    { "tahun": 9, "akumulasi_profit": 25500000 },
    { "tahun": 10, "akumulasi_profit": 37500000 },
    { "tahun": 11, "akumulasi_profit": 49500000 },
    { "tahun": 12, "akumulasi_profit": 61500000 },
    { "tahun": 13, "akumulasi_profit": 73500000 },
    { "tahun": 14, "akumulasi_profit": 85500000 },
    { "tahun": 15, "akumulasi_profit": 97500000 }
  ],
  "warnings": []
}
```

| Field | Deskripsi |
|-------|-----------|
| `kelayakan` | Status kelayakan: "SANGAT LAYAK" (ROI < 8 tahun), "LAYAK" (8-10 tahun), atau "TIDAK LAYAK" (> 10 tahun) |
| `estimasi_biaya` | Estimasi biaya investasi panel surya (Rp) |
| `saving_per_bulan` | Estimasi penghematan listrik per bulan (Rp) |
| `roi_tahun` | Return on Investment dalam tahun |
| `rekomendasi` | Rekomendasi sistem: jumlah panel, kapasitas (kWp), luas atap (m²) |
| `roi_chart_data` | Array data chart ROI untuk 15 tahun dengan akumulasi profit |
| `warnings` | Array warning jika ada (kapasitas besar, intensitas matahari rendah, dll) |

## Workflow Implementation

### 1. **Geocoding** (Google Maps API)
- Menerima input lokasi string (misal: "Surabaya")
- Menggunakan **Google Maps Geocoding API** untuk mendapatkan koordinat Latitude & Longitude
- Error handling jika lokasi tidak ditemukan

### 2. **Solar Data** (NASA POWER API)
- Menggunakan koordinat untuk memanggil **NASA POWER API**
- Endpoint: `https://power.larc.nasa.gov/api/temporal/climatology/point`
- Mengambil parameter `ALLSKY_SFC_SW_DWN` (rata-rata radiasi harian tahunan dalam kWh/m²/hari)

### 3. **Perhitungan Energi**
```
Produksi Energi Harian (kWh) = Radiasi × Kapasitas (kWp) × Performance Ratio (0.75)
Produksi Energi Bulanan = Produksi Harian × 30 hari
Saving per Bulan = Produksi Bulanan × Tarif Listrik (Rp 1.444/kWh)
```

### 4. **Estimasi Kapasitas**
```
Kebutuhan Harian (kWh) = Tagihan Bulanan / (Tarif × 30)
Kapasitas Terukur = Kebutuhan Harian / (Radiasi × Performance Ratio)
Jumlah Panel = ceil(Kapasitas / Kapasitas per Panel)
```

### 5. **Perhitungan ROI**
```
ROI (tahun) = Estimasi Biaya / (Saving Tahunan)
- ROI < 8 tahun = SANGAT LAYAK ✅
- 8 ≤ ROI ≤ 10 tahun = LAYAK ⚠️
- ROI > 10 tahun = TIDAK LAYAK ❌
```

### 6. **Chart Data Generation**
- Membuat perulangan 15 tahun
- Setiap tahun: Akumulasi Profit = (Tahun × Saving Tahunan) - Investasi Awal
- Tahun 0: Hanya investasi awal (negatif)

## Konstanta dan Parameter

```go
const (
	tarifPLN          = 1.444            // Rp 1.444 per kWh
	performanceRatio  = 0.75             // 75% efficiency
	kapasitasPerPanel = 0.55             // 0.55 kWp per panel
	luasPerPanel      = 2.5              // 2.5 m² per panel
	hargaPerKWp       = 15000000.0       // Rp 15 juta per kWp
	chartYears        = 15               // 15 tahun untuk chart
)
```

## Setup & Prerequisites

### 1. Environment Variables

Buat file `.env` atau set environment variable:
```bash
GOOGLE_MAPS_API_KEY=your_api_key_here
```

### 2. Mendapatkan Google Maps API Key

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Buat project baru atau pilih existing project
3. Aktifkan **Geocoding API**:
   - Navigasi ke "APIs & Services" → "Enable APIs and Services"
   - Search "Geocoding API"
   - Click "Enable"
4. Buat API Key:
   - "APIs & Services" → "Credentials"
   - "Create Credentials" → "API Key"
5. Copy API key ke `.env`

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run Server

```bash
go run ./cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## Testing

### Contoh cURL Request

```bash
curl -X POST http://localhost:8080/api/v1/simulation \
  -H "Content-Type: application/json" \
  -d '{
    "lokasi": "Surabaya",
    "tagihan_bulanan": 1000000
  }'
```

## Error Handling

Endpoint menangani berbagai error scenarios:

- **400 Bad Request**: Lokasi tidak ditemukan, input invalid
- **400 Bad Request**: API Google Maps atau NASA gagal
- **500 Internal Server Error**: Error server (uncommon dengan error handling yang baik)

### Error Response Format

```json
{
  "error": "location not found: Xyz"
}
```

## Performance & Warnings

Sistem otomatis menambahkan warnings jika:

1. **Kapasitas Besar**: Jumlah panel > 20
   - Warning: "Kapasitas besar, pastikan atap Anda kuat dan memiliki luas minimal X m2."

2. **Intensitas Matahari Rendah**: Radiasi < 4.0 kWh/m²/hari
   - Warning: "Intensitas matahari di lokasi ini rendah, mempengaruhi kecepatan balik modal."

## Flow Diagram

```
Client Request
    ↓
Validasi Input
    ↓
Google Maps Geocoding → Lat/Lon
    ↓
NASA POWER API → Radiasi Matahari
    ↓
Perhitungan Energi & ROI
    ↓
Generate Chart Data (15 tahun)
    ↓
Validasi Kelayakan
    ↓

## File Structure

│   ├── nasa.go            # NASA POWER API Client
# DEPRECATED — API Documentation moved

This file has been deprecated. See README.md at repository root for consolidated API and calculation documentation.
│   └── simulation.go       # HTTP Request Handler
├── models/
│   └── simulation.go       # Data Models & Structs
└── services/
    └── calculator.go       # Business Logic & Calculations
```

## Notes

- **Tarif Listrik**: Default Rp 1.444/kWh (sesuaikan jika diperlukan)
- **Performance Ratio**: 0.75 (75% efficiency - standard untuk panel surya)
- **Harga per kWp**: Rp 15.000.000 (estimasi, bisa disesuaikan)
- **Chart Years**: 15 tahun untuk visualisasi ROI jangka panjang

## Security Considerations

- API Key Google Maps harus di-protect (jangan commit ke repository)
- Gunakan environment variables atau secrets manager
- Implementasi rate limiting jika diperlukan (future enhancement)
- Validate input location untuk mencegah injection attacks
