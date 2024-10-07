# Same version as our go.mod
FROM golang:1.23-alpine AS build-stage

WORKDIR /svc

# Cache dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy Go code into the container.
COPY cmd cmd
COPY internal internal

# Build the Go code.
RUN go build -o /bin/svc cmd/*.go

# Final stage (minimal image) - using scratch to minimise resource footprint and increase security
FROM scratch AS final-stage

# Copy the Go binary from the build stage
COPY --from=build-stage /bin/svc /bin/svc

EXPOSE 8081 

# use first non-root user
USER 1000
ENTRYPOINT ["/bin/svc"]
