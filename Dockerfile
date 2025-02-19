FROM golang:1.24-alpine AS registration-builder

COPY . /app
WORKDIR /app

RUN apk --no-cache add curl

RUN go mod download && go mod verify
RUN go build -v -o /app/registration cmd/main.go

FROM alpine AS registration

ENV GIN_MODE release

COPY --from=registration-builder /app/registration /app/registration
COPY --from=registration-builder /app/web /app/web
COPY --from=registration-builder /app/LICENSE.md /app/LICENSE.md
COPY --from=registration-builder /app/README.md /app/README.md

ENTRYPOINT ["/app/registration"]
