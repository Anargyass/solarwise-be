# Calculator ROI Engine - Solarwise

Dokumen ini menjelaskan alur perhitungan ROI pada fungsi `CalculateFeasibility` untuk prototipe Solarwise, termasuk metode, rumus, dan alasan kenapa pendekatan ini efektif.

## Pengantar Untuk Pembaca Awam

Bagian ini dibuat supaya pembaca non-teknis bisa memahami konteks dulu sebelum masuk ke rumus.

### Apa yang sebenarnya dihitung Solarwise?

Solarwise menjawab pertanyaan sederhana:
- Jika saya pasang PLTS atap, apakah investasi saya layak?
- Berapa lama balik modalnya?
- Seberapa besar tagihan listrik saya bisa berkurang?

Sistem tidak langsung memilih satu ukuran PLTS. Sistem mencoba beberapa ukuran dulu, lalu memilih yang paling masuk akal berdasarkan biaya, penghematan, luas atap, dan kecepatan balik modal.

### Kamus Istilah Sederhana

- PLTS atap: sistem panel surya yang dipasang di atap rumah/bangunan.
- kWh: satuan energi listrik (berapa banyak listrik dipakai/dihasilkan).
- kWp: ukuran kapasitas puncak sistem surya ("besar" sistem panel).
- Radiasi matahari: potensi sinar matahari di lokasi, memengaruhi produksi listrik panel.
- Performance ratio: faktor efisiensi realistis sistem setelah rugi-rugi (panas, kabel, inverter, dsb).
- Offset: porsi kebutuhan listrik yang bisa ditutupi oleh listrik dari panel surya.
- Coverage ratio: persentase kebutuhan listrik bulanan yang tertutup oleh produksi PLTS.
- Estimasi investasi: perkiraan biaya awal pemasangan sistem.
- Saving per bulan: penghematan biaya listrik per bulan setelah PLTS berjalan.
- ROI (tahun): perkiraan waktu balik modal (semakin kecil, semakin cepat balik modal).
- Break-even year: tahun ketika total penghematan menutup biaya investasi awal.

### Alur Logika Dalam Bahasa Sederhana

1. Pengguna isi lokasi dan tagihan listrik bulanan.
2. Sistem cari koordinat lokasi dan data potensi matahari di lokasi itu.
3. Sistem hitung kira-kira kapasitas PLTS yang dibutuhkan.
4. Sistem coba beberapa ukuran PLTS (multi-skenario), bukan hanya satu ukuran.
5. Tiap skenario diuji:
   - cukup/tidak luas atapnya
   - berapa listrik yang dihasilkan
   - berapa penghematan per bulan
   - berapa lama balik modal
6. Skenario yang tidak masuk akal langsung dibuang (misal atap tidak cukup, ROI terlalu lama).
7. Dari skenario yang lolos, dipilih yang terbaik:
   - utama: ROI paling cepat
   - jika ROI hampir sama: pilih yang menutup kebutuhan listrik lebih besar
8. Hasil akhir ditampilkan sebagai status kelayakan, alasan, angka finansial, rekomendasi teknis, warning, dan data chart.

### Kenapa pendekatan ini ramah untuk pengguna awam?

- Tidak mengandalkan satu tebakan ukuran panel.
- Lebih realistis karena mempertimbangkan batas fisik atap dan batas ekonomi.
- Keputusan akhir transparan: pengguna bisa melihat alasan kelayakan dari angka nyata.
- Output mudah dipahami: ada status, reasoning, warning, dan ringkasan finansial.

## 1. Tujuan Engine

Engine menghitung kelayakan investasi PLTS atap dari input pengguna:
- lokasi
- tagihan bulanan listrik

Output utama yang dihasilkan tetap mengikuti struktur response yang sudah ada (status, reasoning, financials, technical_recommendation, assumptions, warnings, decision rules, dll), tetapi nilainya berasal dari skenario terbaik hasil evaluasi multi-kapasitas.

## 2. Parameter Tetap (Asumsi Dasar)

Konstanta yang dipakai:
- Tarif listrik PLN: `tarifPLN = 1444` Rupiah per kWh
- Performance ratio sistem: `performanceRatio = 0.75`
- Kapasitas tiap panel: `kapasitasPerPanel = 0.55` kWp
- Luas tiap panel: `luasPerPanel = 2.5` m2
- Harga instalasi: `hargaPerKWp = 15,000,000` Rupiah per kWp
- Horizon chart ROI: `chartYears = 15` tahun

Parameter batas simulasi:
- Kapasitas minimum kandidat: 2 kWp
- Kapasitas maksimum kandidat: 8 kWp
- Batas area atap: 40 m2
- Batas offset konsumsi: 75%
- Batas ROI layak simulasi: maksimal 15 tahun

## 3. Validasi Input

Sebelum hitung:
1. `lokasi` wajib terisi
2. `tagihan_bulanan` harus lebih besar dari 0

Jika gagal, fungsi mengembalikan error.

## 4. Akuisisi Data Eksternal

### 4.1 Geocoding lokasi
Lokasi diubah menjadi koordinat latitude dan longitude via Google Maps Geocoding API.

### 4.2 Data radiasi
Dengan koordinat tersebut, diambil data klimatologi dari NASA POWER API. Nilai radiasi harus positif.

## 5. Rumus Dasar Kebutuhan Energi

### 5.1 Konsumsi bulanan dan harian

$$
\text{kebutuhanBulananKwh} = \frac{\text{tagihanBulanan}}{\text{tarifPLN}}
$$

$$
\text{kebutuhanHarianKwh} = \frac{\text{kebutuhanBulananKwh}}{30}
$$

### 5.2 Kebutuhan kapasitas ideal awal

$$
\text{requiredKwp} = \frac{\text{kebutuhanHarianKwh}}{\text{radiasi} \times \text{performanceRatio}}
$$

Nilai ini adalah estimasi teoritis sebelum constraint teknis dan finansial diterapkan.

## 6. Metode Multi-Skenario Kapasitas

Berbeda dari pendekatan single-capacity, engine sekarang membangun beberapa kandidat kapasitas.

### 6.1 Rentang kandidat

Kandidat dibangkitkan dari 2 kWp sampai:

$$
\min(\text{requiredKwp}, 8)
$$

Dengan step 1 kWp.

Contoh:
- jika `requiredKwp = 5.7`, kandidat: 2, 3, 4, 5
- jika `requiredKwp = 9.2`, kandidat: 2, 3, 4, 5, 6, 7, 8
- jika `requiredKwp < 2`, tetap diuji minimal 2 kWp

## 7. Simulasi Tiap Kandidat

Untuk setiap kandidat `requestedKwp`, dihitung ulang seluruh metrik.

### 7.1 Konversi kapasitas ke panel

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

### 7.2 Filter constraint teknis

Skenario langsung ditolak bila:

$$
\text{luasAtap} > 40
$$

### 7.3 Produksi energi bulanan

$$
\text{produksiBulananKwh} = \text{actualKwp} \times \text{radiasi} \times \text{performanceRatio} \times 30
$$

### 7.4 Offset efektif (dibatasi 75%)

Batas offset:

$$
\text{maxOffsetKwh} = \text{kebutuhanBulananKwh} \times 0.75
$$

Offset efektif:

$$
\text{effectiveOffsetKwh} = \min(\text{produksiBulananKwh}, \text{maxOffsetKwh})
$$

### 7.5 Penghematan bulanan

$$
\text{savingPerBulan} = \text{effectiveOffsetKwh} \times \text{tarifPLN}
$$

Jika `savingPerBulan <= 0`, skenario ditolak.

### 7.6 Saving tahunan, ROI, dan break-even

$$
\text{savingTahunan} = \text{savingPerBulan} \times 12
$$

$$
\text{roiTahun} = \frac{\text{estimasiBiaya}}{\text{savingTahunan}}
$$

dengan:

$$
\text{estimasiBiaya} = \text{actualKwp} \times \text{hargaPerKWp}
$$

Break-even year:

$$
\text{breakEvenYear} = \lceil \text{roiTahun} \rceil
$$

### 7.7 Filter constraint finansial

Skenario ditolak bila:

$$
\text{roiTahun} > 15
$$

### 7.8 Coverage ratio (versi benar)

Coverage sekarang dihitung berbasis energi, bukan nominal rupiah:

$$
\text{coverageRatio} = \frac{\text{effectiveOffsetKwh}}{\text{kebutuhanBulananKwh}}
$$

Ini lebih representatif secara teknis karena menunjukkan proporsi kebutuhan energi yang ditutupi sistem.

## 8. Pemilihan Skenario Terbaik

Setelah semua kandidat difilter, dipilih 1 skenario terbaik dengan aturan:
1. ROI paling kecil
2. Jika ROI mirip (selisih <= 0.1 tahun), pilih coverage ratio lebih tinggi

Secara logika, ini mencari kombinasi paling cepat balik modal tanpa mengorbankan kontribusi energi ketika ROI hampir setara.

## 9. Pengisian Output Existing

Semua field response yang sudah ada diisi dari `best scenario`, antara lain:
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

Artinya, struktur response tidak berubah, hanya engine internal yang menjadi lebih robust.

## 10. Klasifikasi Kelayakan

Berdasarkan `roiTahun` terbaik:
- `roiTahun < 7`: Sangat Layak
- `7 <= roiTahun <= 10`: Layak
- `roiTahun > 10`: Tidak Layak

Kemudian dibuat reasoning dinamis dalam Bahasa Indonesia dengan menyertakan angka hasil simulasi.

## 11. Warnings dan Confidence

### 11.1 Warnings
Warning utama mencakup:
- ROI dekat ambang 7 atau 10 tahun
- asumsi radiasi dan tarif konstan
- shading belum dimodelkan
- coverage < 100%
- kapasitas dibatasi range evaluasi
- radiasi rendah

### 11.2 Confidence score
Skor confidence dibangun dari 3 komponen berbobot:

$$
\text{confidenceScore} = 0.3 \times \text{dataCompletenessScore} + 0.3 \times \text{irradiationScore} + 0.4 \times \text{financialAssumptionsScore}
$$

Lalu dipetakan ke label:
- HIGH jika >= 0.85
- MEDIUM jika >= 0.65
- LOW jika < 0.65

## 12. ROI Chart Data

Fungsi `generateROIChartData` membentuk profit kumulatif tahun ke-0 sampai tahun ke-15.

- Tahun 0: `-investasi`
- Tahun ke-i:

$$
\text{cumulativeProfit}_i = -\text{investasi} + i \times \text{savingTahunan}
$$

Nilai dipakai untuk visualisasi kapan kurva melewati titik impas.

## 13. Kenapa Metode Ini Efektif untuk Solarwise

Metode ini efektif karena:

1. Lebih realistis dari single-point estimate
   - Sistem PLTS tidak ideal dihitung dari satu kapasitas saja, karena panel diskrit, keterbatasan atap, dan batas offset membuat solusi optimal tidak selalu sama dengan kapasitas teoritis.

2. Constraint diperlakukan sebagai filter, bukan pemotongan paksa
   - Dengan filter, setiap kandidat dinilai utuh; ini mencegah distorsi hasil karena hard cap yang bisa merusak rasio ekonomi.

3. Optimasi langsung pada tujuan bisnis utama
   - Tujuan pengguna biasanya balik modal secepat mungkin. Rule pemilihan berbasis ROI minimum sangat sejalan dengan objective tersebut.

4. Tetap menjaga manfaat energi
   - Tie-breaker coverage ratio memastikan ketika ROI hampir sama, sistem memilih opsi yang memberi kontribusi energi lebih tinggi.

5. Transparan dan mudah dijelaskan ke user non-teknis
   - Alur step-by-step, warning, decision rules, dan reasoning dinamis memudahkan komunikasi hasil di UI.

6. Skalabel untuk iterasi produk
   - Kerangka multi-skenario mudah diperluas ke faktor lanjutan: degradasi panel dinamis, eskalasi tarif, shading profile, variasi seasonal, atau optimasi berbasis NPV/IRR.

## 14. Keterbatasan Saat Ini

Agar ekspektasi tetap tepat, engine ini masih memakai beberapa asumsi sederhana:
- radiasi dianggap representatif untuk rata-rata produksi
- tarif listrik dianggap konstan
- belum ada model shading, orientasi panel, dan konsumsi per jam
- belum menghitung OPEX tahunan detail

Meski begitu, untuk tahap prototipe, kombinasi multi-skenario + constraint filter + objective ROI sudah sangat kuat sebagai baseline decision engine.

## 15. Ringkasan Metode

Secara singkat, metode ROI Solarwise adalah:
1. Estimasi kebutuhan energi dari tagihan
2. Bangun kandidat kapasitas
3. Simulasikan metrik teknis-finansial tiap kandidat
4. Filter kandidat yang tidak feasible
5. Pilih skenario terbaik (ROI minimum, lalu coverage)
6. Isi response existing dari skenario terbaik

Pendekatan ini memberikan hasil yang lebih robust, explainable, dan actionable untuk rekomendasi investasi PLTS rumahan.
