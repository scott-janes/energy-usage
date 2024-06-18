FROM energy-usage-base AS builder

WORKDIR /energyUsage/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLE=0 go build -tags musl -o app .

FROM alpine

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /energyUsage

COPY --from=builder /energyUsage/app/app /energyUsage/app

CMD ["/energyUsage/app"]
