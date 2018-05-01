generate.static :
	statik -src=static && echo "Generated!"
build : test
	go build infrastructure/web.go
test : generate.static
	go test ./...
run : generate.static
	go run infrastructure/web.go