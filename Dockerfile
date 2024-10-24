FROM golang:1.23.1-bookworm as builder
LABEL authors="dylan.dinh"

WORKDIR app/

# Copy dependences tree
COPY go.mod go.sum ./

# Tidy the dependences
RUN go mod tidy

# Download go modules
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o tsa cmd/main.go


FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates sqlite3 libsqlite3-0

COPY --from=builder go/app/tsa /tsa
COPY --from=builder go/app/.env /.env
COPY --from=builder go/app/test.db /test.db

RUN chmod +x /tsa

CMD ["./tsa"]