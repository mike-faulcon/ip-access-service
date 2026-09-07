grpcurl -plaintext \
-d '{
"ip": "8.8.8.8",
"allowedCountries": ["US", "CA"]
}' \
localhost:9090 \
access.v1.AccessService/CheckAccess