.PHONY: help

test:
	go test -v -run='^\QTest_Check_' ./shm
	go test -v -run='^\QTest_Check_' ./dataStructure/speedyArray
cover:
	go test -cover -run='^\QTest_Check_' ./shm
	go test -cover -run='^\QTest_Check_' ./dataStructure/speedyArray
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "This mechanism is a suite of tests designed to ensure that"
	@echo "the packages are functioning correctly and"
	@echo "to identify any issues that may exist."
	@echo ""
	@echo "Available targets:"
	@echo "  test     - unit test"
	@echo "  cover    - coverage test"
	@echo ""