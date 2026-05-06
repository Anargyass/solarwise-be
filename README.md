<div align="center">

# SolarWise

_Optimizing Solar Energy ROI with Precision Data_

![Next.js](https://img.shields.io/badge/Next.js-15-black?logo=next.js)
![License](https://img.shields.io/badge/License-TBD-lightgrey)
![Status](https://img.shields.io/badge/Status-Active-brightgreen)

</div>

## Ringkasan
- SolarWise adalah platform decision-support untuk ROI PLTS atap berbasis data lokasi dan radiasi.
- Repositori ini menyediakan backend API dan engine kalkulasi SolarWise.
- Output utama: proyeksi energi, ROI, rekomendasi teknis, dan confidence.

## Fitur Utama
- Precision Mapping: integrasi Google Maps API untuk analisis berbasis lokasi.
- Real-time Solar Data: menggunakan NASA POWER API untuk data radiasi.
- Financial Forecasting: kalkulasi otomatis penghematan, ROI, dan payback period.
- Decision Engine: memberikan status, confidence, dan reasoning terstruktur.
- Interactive Dashboard: UI minimal dan bersih untuk visualisasi data.

## Tech Stack
- Frontend: Next.js 15 (Turbopack), TypeScript, Tailwind CSS.
- Backend: Golang (custom calculation engine).
- APIs: NASA POWER API, Google Maps Platform.
- Deployment: TBD / VPS / Cloud.

## Arsitektur Sistem
- Alur data: User -> Frontend -> Backend (Go) -> External APIs -> Processing -> Response -> Dashboard.
- Pemisahan tanggung jawab: UI layer, calculation engine, external data sources.

## API Reference
- Base URL: http://localhost:8080
- Endpoint: POST /v1/simulation
- Content-Type: application/json

### Request Body (JSON)
- `lokasi` (string): nama kota/alamat untuk geocoding.
- `tagihan_bulanan` (number): estimasi tagihan listrik bulanan (Rupiah).

### Response Ringkas
- `status`, `reasoning`, `confidence`, `financials`, `technical_recommendation`, `assumptions`, `warnings`.
- `decision_rules`, `decision_basis`, `roi_chart_data`, `display`.

### Error Response
```json
{"error": "message"}
```

### Status Codes
- 200: sukses.
- 400: input tidak valid atau lokasi tidak ditemukan.
- 500: error internal.

### Contoh Request
```bash
curl -X POST http://localhost:8080/v1/simulation \
  -H "Content-Type: application/json" \
  -d '{"lokasi":"Jakarta","tagihan_bulanan":1500000}'
```

## Konfigurasi
- `GOOGLE_MAPS_API_KEY` wajib di-set untuk Google Maps Geocoding.
- `.env` didukung untuk local development (format `KEY=VALUE`).
- NASA POWER API tidak memerlukan API key.

## Getting Started
### Prerequisites
- Go 1.22+
- Node.js dan pnpm/npm (jika menjalankan frontend SolarWise)

### Setup Backend
```bash
git clone https://github.com/username/solarwise-be.git
cd solarwise-be
go mod download
```

### Menjalankan Backend (Go)
```bash
go run ./cmd/api
```

### Smoke Test
```bash
curl -X POST http://localhost:8080/v1/simulation \
  -H "Content-Type: application/json" \
  -d '{"lokasi":"Jakarta","tagihan_bulanan":1500000}'
```

## Struktur Proyek
- [cmd/api/main.go](cmd/api/main.go): entrypoint HTTP server.
- [internal/handlers/simulation.go](internal/handlers/simulation.go): validasi request dan handler.
- [internal/services/calculator.go](internal/services/calculator.go): engine kalkulasi ROI.
- [internal/clients/googlemaps.go](internal/clients/googlemaps.go): integrasi Google Maps.
- [internal/clients/nasa.go](internal/clients/nasa.go): integrasi NASA POWER.
- [internal/models/simulation.go](internal/models/simulation.go): kontrak request/response.

## Status dan Lisensi
- Status: Active (prototype).
- License: TBD.

## Lampiran Teknis: Calculator ROI Engine (Transparansi)
Konten teknis berikut melampirkan seluruh inti penjelasan dari [calculator.md](calculator.md) secara jelas dan bertahap. Klik untuk membuka.

<details>
<summary>Lihat lampiran lengkap</summary>

### Calculator ROI Engine - Solarwise
- Dokumen ini menjelaskan alur perhitungan ROI pada fungsi `CalculateFeasibility`.
- Fokus: metode, rumus, dan alasan kenapa pendekatan ini efektif.

### Pengantar Untuk Pembaca Awam
#### Apa yang sebenarnya dihitung Solarwise?
- Apakah investasi PLTS atap layak?
- Berapa lama balik modal?
- Seberapa besar tagihan listrik bisa berkurang?
- Sistem mencoba beberapa ukuran PLTS lalu memilih yang paling masuk akal.

#### Kamus Istilah Sederhana
- PLTS atap: sistem panel surya di atap rumah atau bangunan.
- kWh: satuan energi listrik (pemakaian/produksi).
- kWp: ukuran kapasitas puncak sistem surya.
- Radiasi matahari: potensi sinar di lokasi, memengaruhi produksi.
- Performance ratio: faktor efisiensi realistis setelah rugi-rugi.
- Offset: porsi kebutuhan listrik yang ditutupi panel surya.
- Coverage ratio: persentase kebutuhan bulanan yang tertutup produksi.
- Estimasi investasi: biaya awal pemasangan.
- Saving per bulan: penghematan listrik per bulan.
- ROI (tahun): estimasi waktu balik modal.
- Break-even year: tahun saat total penghematan menutup investasi awal.

#### Alur Logika Dalam Bahasa Sederhana
1. Pengguna isi lokasi dan tagihan listrik bulanan.
2. Sistem cari koordinat lokasi dan data potensi matahari.
3. Sistem hitung kapasitas PLTS yang dibutuhkan.
4. Sistem coba beberapa ukuran PLTS (multi-skenario).
5. Tiap skenario diuji: luas atap, produksi listrik, saving per bulan, ROI.
6. Skenario tidak masuk akal dibuang (atap tidak cukup, ROI terlalu lama).
7. Skenario terbaik dipilih: ROI tercepat, lalu coverage lebih tinggi jika ROI mirip.
8. Hasil akhir: status kelayakan, alasan, angka finansial, rekomendasi teknis, warning, data chart.

#### Kenapa pendekatan ini ramah untuk pengguna awam?
- Tidak mengandalkan satu tebakan ukuran panel.
- Lebih realistis: batas atap dan batas ekonomi diperhitungkan.
- Keputusan transparan: alasan kelayakan berbasis angka nyata.
- Output mudah dipahami: status, reasoning, warning, ringkasan finansial.

### 1. Tujuan Engine
- Menghitung kelayakan investasi PLTS atap dari input lokasi dan tagihan bulanan.
- Output mengikuti struktur response yang sudah ada, tetapi nilai berasal dari skenario terbaik.

### 2. Parameter Tetap (Asumsi Dasar)
- Tarif listrik PLN: `tarifPLN = 1444` Rupiah per kWh.
- Performance ratio: `performanceRatio = 0.75`.
- Kapasitas per panel: `kapasitasPerPanel = 0.55` kWp.
- Luas per panel: `luasPerPanel = 2.5` m2.
- Harga instalasi: `hargaPerKWp = 15,000,000` Rupiah per kWp.
- Horizon chart ROI: `chartYears = 15` tahun.

Parameter batas simulasi:
- Kapasitas minimum kandidat: 2 kWp.
- Kapasitas maksimum kandidat: 8 kWp.
- Batas area atap: 40 m2.
- Batas offset konsumsi: 75%.
- Batas ROI layak simulasi: maksimal 15 tahun.

### 3. Validasi Input
1. `lokasi` wajib terisi.
2. `tagihan_bulanan` harus lebih besar dari 0.
- Jika gagal, fungsi mengembalikan error.

### 4. Akuisisi Data Eksternal
#### 4.1 Geocoding lokasi
- Lokasi diubah menjadi koordinat latitude dan longitude via Google Maps Geocoding API.

#### 4.2 Data radiasi
- Data klimatologi diambil dari NASA POWER API.
- Nilai radiasi harus positif.

### 5. Rumus Dasar Kebutuhan Energi
#### 5.1 Konsumsi bulanan dan harian
$$
\text{kebutuhanBulananKwh} = \frac{\text{tagihanBulanan}}{\text{tarifPLN}}
$$

$$
\text{kebutuhanHarianKwh} = \frac{\text{kebutuhanBulananKwh}}{30}
$$

#### 5.2 Kebutuhan kapasitas ideal awal
$$
\text{requiredKwp} = \frac{\text{kebutuhanHarianKwh}}{\text{radiasi} \times \text{performanceRatio}}
$$
- Nilai ini adalah estimasi teoritis sebelum constraint teknis dan finansial.

### 6. Metode Multi-Skenario Kapasitas
- Engine membangun beberapa kandidat kapasitas, bukan single-capacity.

#### 6.1 Rentang kandidat
- Kandidat dibangkitkan dari 2 kWp sampai:

$$
\min(\text{requiredKwp}, 8)
$$

- Step 1 kWp.
- Contoh:
  - `requiredKwp = 5.7` -> kandidat: 2, 3, 4, 5.
  - `requiredKwp = 9.2` -> kandidat: 2, 3, 4, 5, 6, 7, 8.
  - `requiredKwp < 2` -> tetap diuji minimal 2 kWp.

### 7. Simulasi Tiap Kandidat
- Untuk setiap `requestedKwp`, seluruh metrik dihitung ulang.

#### 7.1 Konversi kapasitas ke panel
Jumlah panel dibulatkan ke atas:
$$
\text{jumlahPanel} = \left\lceil \frac{\text{requestedKwp}}{\text{kapasitasPerPanel}} \right\rceil
$$

Kapasitas aktual setelah pembulatan panel:
$$
\text{actualKwp} = \text{jumlahPanel} \times \text{kapasitasPerPanel}
$$

Luas atap:
$$
\text{luasAtap} = \text{jumlahPanel} \times \text{luasPerPanel}
$$

#### 7.2 Filter constraint teknis
- Skenario ditolak bila:
$$
\text{luasAtap} > 40
$$

#### 7.3 Produksi energi bulanan
$$
\text{produksiBulananKwh} = \text{actualKwp} \times \text{radiasi} \times \text{performanceRatio} \times 30
$$

#### 7.4 Offset efektif (dibatasi 75%)
Batas offset:
$$
\text{maxOffsetKwh} = \text{kebutuhanBulananKwh} \times 0.75
$$

Offset efektif:
$$
\text{effectiveOffsetKwh} = \min(\text{produksiBulananKwh}, \text{maxOffsetKwh})
$$

#### 7.5 Penghematan bulanan
$$
\text{savingPerBulan} = \text{effectiveOffsetKwh} \times \text{tarifPLN}
$$
- Jika `savingPerBulan <= 0`, skenario ditolak.

#### 7.6 Saving tahunan, ROI, dan break-even
$$
\text{savingTahunan} = \text{savingPerBulan} \times 12
$$

$$
\text{roiTahun} = \frac{\text{estimasiBiaya}}{\text{savingTahunan}}
$$

Dengan:
$$
\text{estimasiBiaya} = \text{actualKwp} \times \text{hargaPerKWp}
$$

Break-even year:
$$
\text{breakEvenYear} = \lceil \text{roiTahun} \rceil
$$

#### 7.7 Filter constraint finansial
- Skenario ditolak bila:
$$
\text{roiTahun} > 15
$$

#### 7.8 Coverage ratio (versi benar)
$$
\text{coverageRatio} = \frac{\text{effectiveOffsetKwh}}{\text{kebutuhanBulananKwh}}
$$
- Coverage berbasis energi lebih representatif dibanding nominal rupiah.

### 8. Pemilihan Skenario Terbaik
- Aturan pemilihan:
  1. ROI paling kecil.
  2. Jika ROI mirip (selisih <= 0.1 tahun), pilih coverage ratio lebih tinggi.
- Tujuan: balik modal tercepat tanpa mengorbankan kontribusi energi.

### 9. Pengisian Output Existing
- Semua field response diisi dari `best scenario`, termasuk:
  - `status`
  - `reasoning`
  - `financials`
  - `technical_recommendation`
  - `assumptions`
  - `warnings`
  - `decision_rules`
  - `decision_basis`
  - `display`
  - `roi_chart_data`
  - `confidence` dan `confidence_factors`
- Struktur response tidak berubah, engine internal menjadi lebih robust.

### 10. Klasifikasi Kelayakan
- `roiTahun < 7`: Sangat Layak.
- `7 <= roiTahun <= 10`: Layak.
- `roiTahun > 10`: Tidak Layak.
- Reasoning dibuat dinamis dalam Bahasa Indonesia dengan angka hasil simulasi.

### 11. Warnings dan Confidence
#### 11.1 Warnings
- ROI dekat ambang 7 atau 10 tahun.
- Asumsi radiasi dan tarif konstan.
- Shading belum dimodelkan.
- Coverage < 100%.
- Kapasitas dibatasi range evaluasi.
- Radiasi rendah.

#### 11.2 Confidence score
Skor confidence:
$$
\text{confidenceScore} = 0.3 \times \text{dataCompletenessScore} + 0.3 \times \text{irradiationScore} + 0.4 \times \text{financialAssumptionsScore}
$$

Mapping label:
- HIGH jika >= 0.85.
- MEDIUM jika >= 0.65.
- LOW jika < 0.65.

### 12. ROI Chart Data
- Fungsi `generateROIChartData` membentuk profit kumulatif tahun ke-0 sampai tahun ke-15.
- Tahun 0: `-investasi`.
- Tahun ke-i:

$$
\text{cumulativeProfit}_i = -\text{investasi} + i \times \text{savingTahunan}
$$

- Nilai dipakai untuk visualisasi kapan kurva melewati titik impas.

### 13. Kenapa Metode Ini Efektif untuk Solarwise
1. Lebih realistis dari single-point estimate.
2. Constraint diperlakukan sebagai filter, bukan pemotongan paksa.
3. Optimasi langsung pada tujuan bisnis utama (ROI minimum).
4. Tetap menjaga manfaat energi (tie-breaker coverage ratio).
5. Transparan dan mudah dijelaskan ke user non-teknis.
6. Skalabel untuk iterasi produk.

### 14. Keterbatasan Saat Ini
- Radiasi dianggap representatif untuk rata-rata produksi.
- Tarif listrik dianggap konstan.
- Belum ada model shading, orientasi panel, dan konsumsi per jam.
- Belum menghitung OPEX tahunan detail.

### 15. Ringkasan Metode
1. Estimasi kebutuhan energi dari tagihan.
2. Bangun kandidat kapasitas.
3. Simulasikan metrik teknis-finansial tiap kandidat.
4. Filter kandidat yang tidak feasible.
5. Pilih skenario terbaik (ROI minimum, lalu coverage).
6. Isi response existing dari skenario terbaik.

- Pendekatan ini menghasilkan rekomendasi investasi yang lebih robust, explainable, dan actionable.

</details>
