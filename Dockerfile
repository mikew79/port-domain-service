FROM golang:1.18

WORKDIR /go/src/github.com/mikew79/port-domain-service
RUN mkdir -p /go/src/github.com/mikew79/port-domain-service
COPY . /go/src/github.com/mikew79/port-domain-service
RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build ./service/main.go

CMD ["./main", "-dbname=domainPortsDb","-port=7000"]