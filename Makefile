GOENV:=GO111MODULE=on

.PHONY: build
build:
	$(GOENV) CGO_ENABLED=0 go build -v -o ./bin/datamapper ./main.go

.PHONY: test
test:
	$(GOENV) CGO_ENABLED=0 go test ./... ./_test_data/generator/...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: local-test
local-test: build
	mkdir -p _test_data/local_test

	./bin/datamapper --from User --from-tag map --from-source _test_data/mapper/domain/user.go \
	 	--to User --to-tag map --to-source _test_data/mapper/transport/models.go \
	 	-d _test_data/local_test/domain_to_dto_user_converter.go \
		--cf github.com/underbek/datamapper/_test_data/mapper/convertors --cf github.com/underbek/datamapper/_test_data/mapper/other_convertors

	./bin/datamapper --from User --from-tag map --from-source _test_data/mapper/domain/user.go \
		--to User --to-tag map --to-source _test_data/mapper/transport/models.go \
		-d _test_data/local_test/domain_to_dto_user_converter.go \
		--cf github.com/underbek/datamapper/_test_data/mapper/convertors --cf github.com/underbek/datamapper/_test_data/mapper/other_convertors

	./bin/datamapper --from User --from-source github.com/underbek/datamapper/_test_data/mapper/domain \
		--to User --to-source github.com/underbek/datamapper/_test_data/mapper/transport \
		-d _test_data/local_test/domain_to_dto_user_converter.go \
		--cf github.com/underbek/datamapper/_test_data/mapper/convertors --cf github.com/underbek/datamapper/_test_data/mapper/other_convertors

	$(GOENV) go generate ./_test_data/mapper/domain
