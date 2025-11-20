#!/bin/bash
mkdir -p certificates
openssl req -x509 -newkey rsa:2048 -keyout certificates/key.pem -out certificates/cert.pem -days 365 -nodes -subj "/C=CN/ST=Beijing/L=Beijing/O=Dev/CN=localhost"
