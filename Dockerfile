FROM golang:1.20 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /since

FROM scratch

COPY --from=builder /since /since

WORKDIR /project

ENTRYPOINT ["/since"]
