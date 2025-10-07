FROM golang:1.25 AS base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o main .

FROM gcr.io/distroless/base

COPY --from=base /app/main .

COPY --from=base /app/templates ./templates
COPY --from=base /app/static ./static

EXPOSE 8888

CMD ["./main"]