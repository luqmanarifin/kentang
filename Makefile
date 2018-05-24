pretty:
	# gofmt -d -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -d -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	go tool vet .

