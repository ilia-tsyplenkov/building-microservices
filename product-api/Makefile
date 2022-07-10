check_install:
	which swagger || go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0
swagger: check_install
	swagger generate spec -o ./swagger.yaml --scan-models
