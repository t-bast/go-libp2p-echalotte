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
protobuf: protodeps $(PROTO_GO_FILES)

protodeps:
	go get -u github.com/gogo/protobuf/protoc-gen-gogofaster

%.pb.go: %.proto
	protoc -I=$(GOPATH)/src --proto_path=$(GOPATH)/src:. --gogofaster_out=\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:. $<

# Publish the package.
publish:
	gx-go rewrite --undo