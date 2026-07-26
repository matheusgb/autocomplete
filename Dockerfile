FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o autocomplete .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=build /app/autocomplete ./autocomplete
COPY --from=build /app/index.html ./index.html
COPY --from=build /app/script.js ./script.js
COPY --from=build /app/styles.css ./styles.css

EXPOSE 8080

CMD ["./autocomplete"]
