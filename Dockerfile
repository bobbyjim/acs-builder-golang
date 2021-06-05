FROM golang:1.16.4 as builder
WORKDIR /go/src/github.com/bobbyjim/acs-builder
RUN go get -d -v github.com/gorilla/mux
COPY main.go  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o acs-builder .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/bobbyjim/acs-builder/acs-builder .
CMD ["./acs-builder"]

