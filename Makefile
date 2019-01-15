# Protobuf files.
PROTO_FILES=$(shell find pb -name '*.proto')
PROTO_GO_FILES=$(PROTO_FILES:.proto=.pb.go)

# Install the gx package manager.
gx:
	go get -u github.com/whyrusleeping/gx
	go get -u github.com/whyrusleeping/gx-go

# Install dependencies.
deps: gx 
	gx --verbose install --global
	gx-go rewrite

# Install mocking tools.
mockdeps:
	go get -u github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen

# Generate mocks.
mockgen: mockdeps
	go generate

# Build protobuf definitions.
protobuf: $(PROTO_GO_FILES)

%.pb.go: %.proto
	protoc --proto_path=$(GOPATH)/src:. --go_out=. $<

# Publish the package.
publish:
	gx-go rewrite --undo