# Same version as our go.mod
FROM golang:1.23-alpine

WORKDIR /svc

# Copy Go code into the container.
COPY cmd cmd
COPY internal internal

# Copy go.mod file across and download dependencies.
COPY go.* .
RUN go mod download

# Build the Go code.
RUN go build -o /bin/svc cmd/*.go

EXPOSE 8081

USER root
ENTRYPOINT ["/bin/svc"]
