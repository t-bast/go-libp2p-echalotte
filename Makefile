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

# Publish the package.
publish:
	gx-go rewrite --undo