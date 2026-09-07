curl -X POST http://localhost:8080/v1/check \
-H "Content-Type: application/json" \
-d '{"ip":"68.184.123.121", "allowedCountries":["US", "CA"]}'   