#!/bin/bash
openssl req -newkey rsa:2048 -nodes -x509 -days 365 -out ca.crt -keyout ca.key -addext 'subjectAltName = DNS:gateway-token.io' -subj "/C=CN/ST=State/L=City/O=org/OU=dev/CN=gateway-token.io/"
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -addext 'subjectAltName = DNS:gateway-token.io' -subj "/C=CN/ST=State/L=City/O=org/OU=dev/CN=gateway-token.io/"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:gateway-token.io") -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256
openssl genrsa -out client1.key 2048
openssl req -new -key client1.key -out client1.csr -addext 'subjectAltName = DNS:microservice1.io' -subj "/C=CN/ST=State/L=City/O=org/OU=dev/CN=microservice1.io/"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:microservice1.io") -in client1.csr -CA ca.crt -CAkey ca.key -out client1.crt -days 365 -sha256 -CAcreateserial
openssl genrsa -out client2.key 2048
openssl req -new -key client2.key -out client2.csr -addext 'subjectAltName = DNS:microservice2.io' -subj "/C=CN/ST=State/L=City/O=org/OU=dev/CN=microservice2.io/"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:microservice2.io") -in client2.csr -CA ca.crt -CAkey ca.key -out client2.crt -days 365 -sha256 -CAcreateserial
