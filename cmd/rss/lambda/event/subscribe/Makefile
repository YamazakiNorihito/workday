build:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap .
	zip function.zip bootstrap