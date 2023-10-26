build:
	@go build -o client client.go dht.go mdns.go; \
	go build -o bootstrap bootstrap.go