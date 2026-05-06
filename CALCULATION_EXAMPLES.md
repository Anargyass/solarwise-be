# Calculation Example & Logic Documentation

## Visual Flow (Mermaid) Untuk Pembaca Awam

Bagian ini fokus ke visual, supaya orang awam bisa mengikuti alur tanpa harus membaca rumus dulu.

### Bahasa Isyarat (Legenda Cepat)

- AMAN: kondisi lolos, lanjut.
- CEK: hasil masih valid, tapi sensitif terhadap perubahan asumsi.
- STOP: proses/skenario berhenti karena tidak memenuhi syarat.

### Diagram 1 - Alur Utama End-to-End

~~~mermaid
flowchart TB
    %% ====== style ======
    classDef process fill:#E8F4FD,stroke:#1E88E5,color:#0D47A1,stroke-width:1px;
    classDef decision fill:#FFF3E0,stroke:#FB8C00,color:#E65100,stroke-width:1px;
    classDef stop fill:#FFEBEE,stroke:#E53935,color:#B71C1C,stroke-width:1px;
    classDef ok fill:#E8F5E9,stroke:#43A047,color:#1B5E20,stroke-width:1px;

    A([Mulai]):::ok --> B[Input: lokasi dan tagihan bulanan]:::process
    B --> C{Input valid?}:::decision
    C -- Tidak --> C1[STOP: validasi gagal]:::stop
    C -- Ya --> D[Geocoding lokasi: Google Maps]:::process
    D --> E{Koordinat ditemukan?}:::decision
    E -- Tidak --> E1[STOP: lokasi tidak ditemukan]:::stop
    E -- Ya --> F[Ambil radiasi: NASA POWER]:::process
    F --> G{Radiasi > 0?}:::decision
    G -- Tidak --> G1[STOP: data radiasi invalid]:::stop
    G -- Ya --> H[Hitung kebutuhan energi dan required kWp]:::process
    H --> I[Generate kandidat: 2 sampai min(required, 8) kWp]:::process

    I --> J{Loop tiap kandidat}:::decision
    J --> K[Konversi ke panel: ceil(kWp/0.55)]:::process
    K --> L[Hitung actual kWp dan luas atap]:::process
    L --> M{Luas atap <= 40 m2?}:::decision
    M -- Tidak --> M1[STOP skenario: gagal constraint atap]:::stop --> J
    M -- Ya --> N[Hitung produksi bulanan]:::process
    N --> O[Hitung effective offset: min(produksi, kebutuhan x 75 persen)]:::process
    O --> P[Hitung saving bulanan dan saving tahunan]:::process
    P --> Q{Saving bulanan > 0?}:::decision
    Q -- Tidak --> Q1[STOP skenario: saving invalid]:::stop --> J
    Q -- Ya --> R[Hitung ROI dan break-even]:::process
    R --> S{ROI <= 15 tahun?}:::decision
    S -- Tidak --> S1[STOP skenario: ROI terlalu lama]:::stop --> J
    S -- Ya --> T[Hitung coverage ratio]:::process
    T --> U[Simpan ke daftar skenario valid]:::ok
    U --> J

    J --> V{Ada skenario valid?}:::decision
    V -- Tidak --> V1[STOP: no feasible scenario]:::stop
    V -- Ya --> W[Pilih best: ROI terkecil, tie-break coverage terbesar]:::ok
    W --> X[Isi semua output existing dari best scenario]:::process
    X --> Y[Generate status, reasoning, warnings, confidence, decision rules]:::process
    Y --> Z([Selesai: kirim response]):::ok
~~~

### Diagram 2 - Keputusan Status Kelayakan

~~~mermaid
flowchart TB
    classDef decision fill:#FFF3E0,stroke:#FB8C00,color:#E65100,stroke-width:1px;
    classDef output fill:#E8F5E9,stroke:#43A047,color:#1B5E20,stroke-width:1px;

    A[Input: ROI dari best scenario]:::decision --> B{ROI < 7?}:::decision
    B -- Ya --> C[Status: Sangat Layak]:::output
    B -- Tidak --> D{ROI <= 10?}:::decision
    D -- Ya --> E[Status: Layak]:::output
    D -- Tidak --> F[Status: Tidak Layak]:::output
~~~

### Diagram 3 - Ringkasan Super Singkat (Untuk Slide Presentasi)

~~~mermaid
flowchart LR
    A[Input pengguna] --> B[Ambil data lokasi dan radiasi]
    B --> C[Uji banyak skenario kapasitas]
    C --> D[Buang skenario yang tidak feasible]
    D --> E[Pilih skenario terbaik]
    E --> F[Keluarkan hasil finansial dan teknis]
~~~

### Contoh Narasi Yang Mudah Diikuti

Contoh input:
- Lokasi: Surabaya
- Tagihan: Rp 1.000.000

Urutan berpikir untuk orang awam:
1. Sistem cek data dasar dulu.
2. Sistem ambil potensi matahari lokasi.
3. Sistem mencoba beberapa ukuran PLTS, bukan satu tebakan.
4. Skenario yang atapnya tidak cukup atau ROI terlalu lama langsung dihentikan.
5. Dari skenario yang lolos, dipilih yang paling cepat balik modal.
6. Jika ROI hampir sama, dipilih yang menutup kebutuhan listrik lebih besar.
7. Hasil akhir ditampilkan dalam format yang mudah dibaca: status, alasan, angka finansial, rekomendasi teknis.

## Contoh Perhitungan Lengkap

### Input Data
```
Lokasi: Surabaya
Tagihan Bulanan: Rp 1,000,000
```

### Step 1: Geocoding (Google Maps)
- Input: "Surabaya"
- Output: Latitude = -7.2504, Longitude = 112.7688

### Step 2: Solar Radiation Data (NASA POWER API)
- Koordinat: -7.2504, 112.7688
- Parameter: ALLSKY_SFC_SW_DWN (ANN - Annual Average)
- Output: Radiasi Harian = ~5.2 kWh/m²/hari (rata-rata tahunan)

### Step 3: Hitung Kebutuhan Kapasitas

```
Tarif PLN (tarifPLN) = Rp 1,444 per kWh

Kebutuhan Bulanan (kWh) = Tagihan Bulanan / Tarif
                        = Rp 1,000,000 / Rp 1,444
                        = 692.51 kWh

Kebutuhan Harian (kWh) = Kebutuhan Bulanan / 30
                       = 692.51 / 30
                       = 23.08 kWh/hari

Kapasitas Sistem (kWp) = Kebutuhan Harian / (Radiasi × Performance Ratio)
                       = 23.08 / (5.2 × 0.75)
                       = 23.08 / 3.9
                       = 5.92 kWp
```
# DEPRECATED — Calculation examples moved

This file has been deprecated. See README.md at repository root for consolidated calculation explanations and examples.
```
Harga Per kWp = Rp 15,000,000


Kapasitas Aktual (kWp) = 11 × 0.55
                       = 6.05 kWp

Luas Atap = 11 × 2.5
          = 27.5 m²

Estimasi Biaya = 6.05 × Rp 15,000,000
               = Rp 90,750,000
```

### Step 5: Perhitungan Produksi Energi

```
Performance Ratio = 0.75 (75% efficiency)
Hari Per Bulan = 30

Produksi Bulanan (kWh) = Kapasitas × Radiasi × Performance Ratio × Hari
                       = 6.05 × 5.2 × 0.75 × 30
                       = 703.35 kWh

Saving Per Bulan = Produksi × Tarif
                 = 703.35 × Rp 1,444
                 = Rp 1,015,237
```

**Catatan:** Jika saving > tagihan bulanan, maka saving = tagihan bulanan
Dalam hal ini: Rp 1,015,237 ≈ Rp 1,000,000 (dianggap maksimal tagihan)

### Step 6: Perhitungan ROI

```
Saving Tahunan = Saving Bulanan × 12
               = Rp 1,000,000 × 12
               = Rp 12,000,000

ROI (tahun) = Estimasi Biaya / Saving Tahunan
            = Rp 90,750,000 / Rp 12,000,000
            = 7.5625 tahun

Kelayakan:
- ROI < 8 tahun = SANGAT LAYAK ✅
- Dalam hal ini: 7.56 tahun < 8 tahun
- Status: SANGAT LAYAK
```

### Step 7: Generate ROI Chart Data (15 Tahun)

| Tahun | Akumulasi Profit | Penjelasan |
|-------|------------------|-----------|
| 0 | -Rp 90,750,000 | Investasi awal (negatif) |
| 1 | -Rp 78,750,000 | -Rp 90,750,000 + Rp 12,000,000 |
| 2 | -Rp 66,750,000 | -Rp 78,750,000 + Rp 12,000,000 |
| 3 | -Rp 54,750,000 | -Rp 66,750,000 + Rp 12,000,000 |
| 4 | -Rp 42,750,000 | -Rp 54,750,000 + Rp 12,000,000 |
| 5 | -Rp 30,750,000 | -Rp 42,750,000 + Rp 12,000,000 |
| 6 | -Rp 18,750,000 | -Rp 30,750,000 + Rp 12,000,000 |
| 7 | -Rp 6,750,000 | -Rp 18,750,000 + Rp 12,000,000 |
| 8 | Rp 5,250,000 | -Rp 6,750,000 + Rp 12,000,000 (BREAK EVEN) |
| 9 | Rp 17,250,000 | Rp 5,250,000 + Rp 12,000,000 |
| 10 | Rp 29,250,000 | Rp 17,250,000 + Rp 12,000,000 |
| 11 | Rp 41,250,000 | Rp 29,250,000 + Rp 12,000,000 |
| 12 | Rp 53,250,000 | Rp 41,250,000 + Rp 12,000,000 |
| 13 | Rp 65,250,000 | Rp 53,250,000 + Rp 12,000,000 |
| 14 | Rp 77,250,000 | Rp 65,250,000 + Rp 12,000,000 |
| 15 | Rp 89,250,000 | Rp 77,250,000 + Rp 12,000,000 |

**Observasi:**
- Pada tahun ke-7 masih negatif: -Rp 6,750,000
- Pada tahun ke-8 profit positif: +Rp 5,250,000
- Break-even point: ~7.56 tahun (sesuai ROI calculation)
```json
{
  "kelayakan": "SANGAT LAYAK",
  "estimasi_biaya": 90750000,
  "saving_per_bulan": 1000000,
  "roi_tahun": 7.5625,
  "rekomendasi": {
    "jumlah_panel": 11,
    "kapasitas_kwp": 6.05,
    "luas_atap_m2": 27.5
  },
  "roi_chart_data": [
    {"tahun": 0, "akumulasi_profit": -90750000},
    {"tahun": 1, "akumulasi_profit": -78750000},
    {"tahun": 2, "akumulasi_profit": -66750000},
    {"tahun": 3, "akumulasi_profit": -54750000},
    {"tahun": 4, "akumulasi_profit": -42750000},
    {"tahun": 5, "akumulasi_profit": -30750000},
    {"tahun": 6, "akumulasi_profit": -18750000},
    {"tahun": 7, "akumulasi_profit": -6750000},
    {"tahun": 8, "akumulasi_profit": 5250000},
    {"tahun": 9, "akumulasi_profit": 17250000},
    {"tahun": 10, "akumulasi_profit": 29250000},
    {"tahun": 11, "akumulasi_profit": 41250000},
    {"tahun": 15, "akumulasi_profit": 89250000}
  ],
  "warnings": []
}
```

## Constants & Konstanta

```go
const (
    tarifPLN          = 1444.0       // Rp 1.444 per kWh (Indonesia PLN standard)
    performanceRatio  = 0.75         // 75% efficiency (standard solar panel)
    kapasitasPerPanel = 0.55         // 0.55 kWp per panel
    luasPerPanel      = 2.5          // 2.5 m² per panel
    hargaPerKWp       = 15000000.0   // Rp 15 juta per kWp (estimate)
    chartYears        = 15           // 15 tahun untuk chart visualization
)
```

## Edge Cases & Warnings

### Case 1: Large Capacity
```
Kondisi: Jumlah Panel > 20
Warning: "Kapasitas besar, pastikan atap Anda kuat dan memiliki luas minimal X m²"
Alasan: Memastikan customer aware dengan kebutuhan luas atap yang besar
```

### Case 2: Low Solar Radiation
```
Kondisi: Radiasi < 4.0 kWh/m²/hari
Warning: "Intensitas matahari di lokasi ini rendah, mempengaruhi kecepatan balik modal"
Alasan: Produksi energi akan lebih rendah, ROI akan lebih lama
Lokasi Contoh: Daerah dengan banyak awan sepanjang tahun
```

### Case 3: High Bill
```
Kondisi: Tagihan sangat tinggi (> Rp 5,000,000)
Tidak ada warning khusus, tapi investasi akan lebih besar
Sistem akan merekomendasikan kapasitas lebih besar untuk maksimal savings
```

## Formulas Summary

```
1. Kebutuhan Bulanan (kWh) = Tagihan Bulanan / Tarif
2. Kebutuhan Harian (kWh) = Kebutuhan Bulanan / 30
3. Kapasitas Terukur (kWp) = Kebutuhan Harian / (Radiasi × Performance Ratio)
4. Jumlah Panel = ceil(Kapasitas Terukur / Kapasitas Per Panel)
5. Kapasitas Aktual (kWp) = Jumlah Panel × Kapasitas Per Panel
6. Luas Atap (m²) = Jumlah Panel × Luas Per Panel
7. Estimasi Biaya = Kapasitas Aktual × Harga Per kWp
8. Produksi Bulanan (kWh) = Kapasitas Aktual × Radiasi × Performance Ratio × 30
9. Saving Per Bulan = Produksi Bulanan × Tarif
10. Saving Tahunan = Saving Per Bulan × 12
11. ROI (tahun) = Estimasi Biaya / Saving Tahunan
12. Akumulasi Profit (Tahun N) = -Estimasi Biaya + (N × Saving Tahunan)
```

## API Response Status Codes

```
200 OK - Request berhasil diproses
400 Bad Request - Input tidak valid atau API eksternal tidak menemukan lokasi
500 Internal Server Error - Error tidak terduga di server
```
