generate-cert:
	@mkdir -p certs
	@openssl req -subj '/CN=localhost/O=Registry Demo/C=US' -new -newkey rsa:2048 -days 365 -nodes -x509 -keyout ./certs/server.key -out ./certs/server.crt