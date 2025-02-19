FROM golang:1.24-alpine AS registration-builder

COPY . /app
WORKDIR /app

RUN apk --no-cache add curl

RUN go mod download && go mod verify
RUN go build -v -o /app/registration cmd/main.go

FROM alpine AS registration

ENV GIN_MODE release

COPY --from=registration-builder /app/registration /app/registration

ENTRYPOINT ["/app/registration"]
