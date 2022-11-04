openssl req -newkey rsa:2048 -nodes -x509 -days 365 -out ca.crt -keyout ca.key -addext 'subjectAltName = DNS:gateway-token.io' -subj "/C=CN/ST=State/L=City/O=org/OU=dev/CN=gateway-token.io/"
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -addext 'subjectAltName = DNS:gateway-token.io' -subj "/C=CN/ST=State/L=City/O=org/OU=dev/CN=gateway-token.io/"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:gateway-token.io") -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -addext 'subjectAltName = DNS:microservice1.io' -subj "/C=CN/ST=State/L=City/O=org/OU=dev/CN=microservice1.io/"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:microservice1.io") -in client.csr -CA ca.crt -CAkey ca.key -out client.crt -days 365 -sha256 -CAcreateserial
