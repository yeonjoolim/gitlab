#/bin/bash

openssl genrsa -out dev1.key 2048

openssl rsa -in dev1.key -out dev1.pub -outform PEM -pubout

echo -e 'ca\nca\nca\nca\nca\nca\n\n\n'| openssl req -new -key dev1.key -out dev1.csr

openssl x509 -req -days 365 -in dev1.csr -CA ca.crt -CAcreateserial -CAkey ca.key -out dev1.crt

