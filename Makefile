generate-cert:
	@mkdir -p certs
	@openssl req -subj '/CN=localhost/O=Registry Demo/C=US' -new -newkey rsa:2048 -days 365 -nodes -x509 -keyout ./certs/server.key -out ./certs/server.crt

vendorize:
	@rm -rf go.sum ; \
		rm -rf vendor ;\
		go mod tidy ; \
		go mod vendor


protocol:
	docker run --rm 													\
		-v $$(pwd):$$(pwd)												\
		-w $$(pwd) znly/protoc											\
		-I $$(pwd)/pkg/protocol											\
		$$(find $$(pwd)/pkg/protocol/ -type f -name "*.proto" | xargs)	\
		--go_out=$$(pwd)/pkg/protocol
