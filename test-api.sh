#!/bin/bash

# ============================================================================
# SolarWise Backend - API Simulation Endpoint Testing Script
# ============================================================================
# Script ini digunakan untuk testing endpoint /api/v1/simulation
# ============================================================================

API_URL="http://localhost:8080/api/v1/simulation"

echo "=========================================="
echo "SolarWise API Simulation Testing"
echo "=========================================="
echo ""

# Test Case 1: Surabaya dengan tagihan Rp 1,000,000
echo "Test 1: Lokasi Surabaya, Tagihan Rp 1,000,000"
echo "-------------------------------------------"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "lokasi": "Surabaya",
    "tagihan_bulanan": 1000000
  }' \
  -s | jq '.'
echo ""
echo ""

# Test Case 2: Jakarta dengan tagihan Rp 2,000,000
echo "Test 2: Lokasi Jakarta, Tagihan Rp 2,000,000"
echo "-------------------------------------------"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "lokasi": "Jakarta",
    "tagihan_bulanan": 2000000
  }' \
  -s | jq '.'
echo ""
echo ""

# Test Case 3: Bandung dengan tagihan Rp 1,500,000
echo "Test 3: Lokasi Bandung, Tagihan Rp 1,500,000"
echo "-------------------------------------------"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "lokasi": "Bandung",
    "tagihan_bulanan": 1500000
  }' \
  -s | jq '.'
echo ""
echo ""

# Test Case 4: Error - Lokasi Tidak Ditemukan
echo "Test 4: Error Case - Lokasi tidak valid"
echo "-------------------------------------------"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "lokasi": "XyzAbcNotACity",
    "tagihan_bulanan": 1000000
  }' \
  -s | jq '.'
echo ""
echo ""

# Test Case 5: Error - Tagihan Negatif
echo "Test 5: Error Case - Tagihan negatif"
echo "-------------------------------------------"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "lokasi": "Jakarta",
    "tagihan_bulanan": -1000000
  }' \
  -s | jq '.'
echo ""
echo ""

# Test Case 6: Error - Missing Required Field
echo "Test 6: Error Case - Missing lokasi"
echo "-------------------------------------------"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "tagihan_bulanan": 1000000
  }' \
  -s | jq '.'
echo ""
echo ""

echo "=========================================="
echo "Testing Complete!"
echo "=========================================="
