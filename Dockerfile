# golang image where workspace (GOPATH) configured at /go.
FROM golang:1.8.3

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/fileuploaddemo/FileUploadGo

# Build the FileUploadGo command inside the container.
RUN go install github.com/fileuploaddemo/FileUploadGo

# Run the FileUploadGo command when the container starts.
ENTRYPOINT /go/bin/FileUploadGo

# http server listens on port 8080.
EXPOSE 8080