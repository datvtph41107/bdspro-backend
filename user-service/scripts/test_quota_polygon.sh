#!/bin/bash

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpdGllcyI6bnVsbCwiZXhwIjoxNzg0MDE3MDkyLCJpYXQiOjE3NzYyNDEwOTIsIm9yZ2FuaXphdGlvbiI6bnVsbCwib3JpZ2luIjowLCJwbGFuIjpudWxsLCJwbGFuQXQiOm51bGwsInByb2ZpbGUiOjI4Niwicm9sZSI6IlJPTEVfVVNFUiIsInNlc3Npb24iOjQ2LCJzdWIiOjM0MiwidHlwZSI6IlJFRlJFU0gifQ.fqsstdRTbAhetq0CKsg4Y8NREqFr9D3Sd-w9oPztzto"

echo "=== TEST QUOTA WITH POLYGON ==="

# Reset quota
echo "Resetting quota..."
docker exec -it redis-quota redis-cli DEL "quota:usage:342:data_pool"
docker exec -it redis-quota redis-cli SET "quota:limit:342:data_pool" 5

# Test 5 lần
for i in {1..6}; do
  echo -e "\n--- Request $i ---"
  
  # Gọi API
  grpcurl -plaintext -d '{
    "coordinates": [
      {"lat": 21.004119047606608, "lng": 105.8255508231324},
      {"lat": 21.004166116083027, "lng": 105.82555162473699},
      {"lat": 21.004204831737468, "lng": 105.82555573885378},
      {"lat": 21.004237934824417, "lng": 105.8255611850476},
      {"lat": 21.00438212245288, "lng": 105.82558559861803},
      {"lat": 21.00426585454342, "lng": 105.82587058079025},
      {"lat": 21.004241555473996, "lng": 105.82592794266196},
      {"lat": 21.004159986286908, "lng": 105.826115435585},
      {"lat": 21.004118178209737, "lng": 105.8262071603533},
      {"lat": 21.004069613035185, "lng": 105.82630648857463},
      {"lat": 21.00401696497707, "lng": 105.82640940044035},
      {"lat": 21.003966451878057, "lng": 105.8265008069256},
      {"lat": 21.003874009561017, "lng": 105.8266675878038},
      {"lat": 21.003844978259522, "lng": 105.82671752304668},
      {"lat": 21.003820756156983, "lng": 105.82675881746496},
      {"lat": 21.003796523046006, "lng": 105.8267983436375},
      {"lat": 21.00376078965428, "lng": 105.82685215723023},
      {"lat": 21.00375881975002, "lng": 105.8268550689391},
      {"lat": 21.003756167533155, "lng": 105.82685854648969},
      {"lat": 21.003755672305928, "lng": 105.82685913590285},
      {"lat": 21.00375487994239, "lng": 105.82685931272903},
      {"lat": 21.003753196169882, "lng": 105.82685934809223},
      {"lat": 21.003751600437283, "lng": 105.82685935988201},
      {"lat": 21.003749839629293, "lng": 105.82685937166974},
      {"lat": 21.003723361475693, "lng": 105.82685937166974},
      {"lat": 21.002899565007255, "lng": 105.82679586809326},
      {"lat": 21.00336672032023, "lng": 105.8255508231324}
    ],
    "intersect": false,
    "page": 1,
    "pageSize": 20
  }' \
  -H "Authorization: Bearer $TOKEN" \
  localhost:8230 tqdpb.ParcelService/FindParcelsByPolygon 2>&1 | head -20
  
  # Kiểm tra quota
  QUOTA=$(docker exec -it redis-quota redis-cli GET "quota:usage:342:data_pool" | tr -d '\r')
  echo "Quota usage: $QUOTA / 5"
  
  # Nếu quota đã đạt 5, break
  if [ "$QUOTA" = "5" ]; then
    echo -e "\n!!! QUOTA EXHAUSTED !!!"
    break
  fi
done

echo -e "\n=== TEST COMPLETED ==="