#!/bin/bash
# check_data.sh

cd ~/projects/bdspro-golang-microservices/tqd-service/internal/scripts/temp/qhchung2

echo "=========================================="
echo "1. TỔNG SỐ DÒNG (Features)"
echo "=========================================="
cat *.ndjson | wc -l

echo ""
echo "=========================================="
echo "2. CÁC LOẠI ĐẤT (Land Use) - QUAN TRỌNG"
echo "=========================================="
cat *.ndjson | jq -r '.properties.loai_dat_quy_hoach' | sort | uniq -c | sort -rn

echo ""
echo "=========================================="
echo "3. CÁC VĂN BẢN PHÁP LÝ (Quyet Dinh)"
echo "=========================================="
cat *.ndjson | jq -r '.properties.quyet_dinh' | sort | uniq -c | sort -rn | head -15

echo ""
echo "=========================================="
echo "4. CÁC ĐỒ ÁN (Do An)"
echo "=========================================="
cat *.ndjson | jq -r '.properties.do_an' | sort | uniq -c | sort -rn | head -15

echo ""
echo "=========================================="
echo "5. MÀU SẮC (RGB)"
echo "=========================================="
cat *.ndjson | jq -r '.properties | "\(.red),\(.green),\(.blue)"' | sort | uniq -c | sort -rn

echo ""
echo "=========================================="
echo "6. LOẠI GEOMETRY"
echo "=========================================="
cat *.ndjson | jq -r '.geometry.type' | sort | uniq -c

echo ""
echo "=========================================="
echo "7. KIỂM TRA LỖI center_lat (giá trị bất thường)"
echo "=========================================="
cat *.ndjson | jq -r '.properties.center_lat' | sort | uniq -c | head -10

echo ""
echo "=========================================="
echo "8. MẪU DỮ LIỆU (1 feature hoàn chỉnh)"
echo "=========================================="
head -1 7_37_39.ndjson | jq '.'