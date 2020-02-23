FROM golang:alpine as builder

LABEL stage=intermediate

ARG PROJECT_PATH

RUN apk add --update \
        git \
        openssh-client \
        ca-certificates

# Change workdir
WORKDIR /go/src/${PROJECT_PATH}

# Copy code
COPY . .

# Build
RUN go get
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags=-w -o /go/bin/app

FROM scratch

#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/app /

ENTRYPOINT ["/app"]
