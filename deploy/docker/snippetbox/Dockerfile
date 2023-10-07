FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /snippetbox ./cmd/web/

FROM alpine:latest

COPY --from=build /snippetbox /app/snippetbox

CMD ["/app/snippetbox"]

EXPOSE 4000
