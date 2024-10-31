
test:
	LOGGER=info go test -v -cover ./... | grep 'RUN\|PASS\|FAIL'