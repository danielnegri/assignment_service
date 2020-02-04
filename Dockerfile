# Stage 1
FROM golang:1.13 AS builder
# Creating working directory
RUN mkdir -p  /go/src/github.com/surajjain36
# Copying source code to repository
COPY   .  /go/src/github.com/surajjain36/assignment_service
WORKDIR /go/src/github.com/surajjain36/assignment_service
# Installing ca certificates
RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates && rm -rf /var/lib/apt/lists/*
ENV GO111MODULE=on
# Creating go binary
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assignment-service
# Stage 2
FROM alpine
RUN apk add --no-cache openssh
# Copy ca certificates from builder
COPY --from=builder  /etc/ssl/certs /etc/ssl/certs
# Copy our static executable and dependencies from builder
COPY --from=builder /go/src/github.com/surajjain36/assignment_service/assignment-service  /
COPY --from=builder /go/src/github.com/surajjain36/assignment_service/config.yml  /

# Exposing port
EXPOSE 3000
# Run the assignment_service  binary.
ENTRYPOINT ["/assignment-service"]