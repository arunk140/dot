```
Note: Client Server Connection will not work without Valid Certs (TLS)
```

### Generate Certs with Valid SANs (eg. with FQDN as localhost)

```
openssl ecparam -genkey -name prime256v1 -out server.key
openssl req -new -x509 -key server.key -out server.pem -days 3650 -addext "subjectAltName = DNS:localhost"
```