FROM golang:1.21-alpine

# RUN apk update && apk add --no-cache git tzdata g++ make ca-certificates

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf cyrus-sasl-dev mold musl-dev

# Source https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
# See https://stackoverflow.com/a/55757473/12429735RUN 
ENV USER=appuser
ENV UID=10001 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /energyUsage/shared

COPY go.* .

RUN go mod download

ADD . .

