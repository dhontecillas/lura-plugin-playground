all:
	mkdir -p plugins
	# go build -buildmode=plugin -o plugins/casemodifier.so ./response_modifier
	go build -buildmode=plugin -o ../../opensource/playground-community/config/krakend/priv_hnt/casemodifier.so ./response_modifier
.PHONY: all


alpine:
	docker run -it -v "$$PWD:/app" \
		-w /app \
		krakend/builder:2.2.1 \
		go build -buildmode=plugin -o casemodifier.so ./response_modifier
	cp casemodifier.so ../../opensource/playground-community/config/krakend/priv_hnt/casemodifier.so
.PHONY: alpine

handler: 
	mkdir -p plugins
	# go build -buildmode=plugin -o plugins/casemodifier.so ./response_modifier
	go build -buildmode=plugin -o ../../opensource/playground-community/config/krakend/priv_hnt/etaghandler.so ./handler
.PHONY: handler

alpine_handler:
	docker run -it -v "$$PWD:/app" \
		-w /app \
		krakend/builder:2.2.1 \
		go build -buildmode=plugin -o etaghandler.so ./handler
	cp etaghandler.so ../../opensource/playground-community/config/krakend/priv_hnt/etaghandler.so
.PHONY: alpine_handler

client: 
	mkdir -p plugins
	# go build -buildmode=plugin -o plugins/casemodifier.so ./response_modifier
	go build -buildmode=plugin -o ../../opensource/playground-community/config/krakend/priv_hnt/censor.so ./client
.PHONY: client

alpine_client:
	docker run -it -v "$$PWD:/app" \
		-w /app \
		krakend/builder:2.2.1 \
		go build -buildmode=plugin -o censor.so ./client
	cp censor.so ../../opensource/playground-community/config/krakend/priv_hnt/censor.so
.PHONY: alpine_handler
