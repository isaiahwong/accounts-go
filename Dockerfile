FROM golang:1.13-alpine as builder

RUN apk add --update nodejs npm
RUN apk add --update npm

WORKDIR /accounts
COPY go.mod . 
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/accounts

FROM alpine
COPY --from=builder /go/bin/accounts /go/bin/accounts

ENTRYPOINT ["/go/bin/accounts"]