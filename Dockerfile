# builder image
FROM golang:1.18 as builder

WORKDIR /go/src/github.com/mikew79/port-domain-service
RUN mkdir -p /go/src/github.com/mikew79/port-domain-service
COPY . /go/src/github.com/mikew79/port-domain-service
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o portdomainserver ./cmd/portdomainserver/main.go

# Create a final clean image for production use
FROM alpine:3.16.2
COPY --from=builder /go/src/github.com/mikew79/port-domain-service/portdomainserver .
CMD ["./portdomainserver", "-dbname=domainPortsDb","-port=7000"]