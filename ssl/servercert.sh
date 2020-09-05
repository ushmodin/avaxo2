openssl genrsa -out localhost.key 2048
openssl req -new -key localhost.key -out localhost.csr
openssl x509 -req -in localhost.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out localhost.crt -days 825 -sha256 -extfile servercert.ext
