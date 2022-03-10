.PHONY: test

test:
	go test -bench -v ./...

.PHONY: clean

clean:
	$(RM) -rf $(BUILD_DIR)
