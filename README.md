`sshare`
========

a system to upload screenshots, probably to a webserver

proper readme pending

Generating a certificate
------------------------

While testing, you should only have to do this once. For production, you should already have a
certificate from your webserver, or obtain a new one some other way. Refer to "Configuring the
server".

``` bash
# Generate the certificate. Fill in sshare.dev for Common Name.
openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out cert.pem
# Creates files:
#   cert.pem
#   key.pem


```