FROM golang:1.24-alpine AS registration-builder

COPY . /app
WORKDIR /app

RUN go mod download && go mod verify
RUN go build -v -o /app/registration cmd/service/main.go
RUN go build -v -o /app/registration_create_user cmd/create_user/main.go

FROM alpine:latest AS registration

ENV GIN_MODE release

WORKDIR /app

COPY --from=registration-builder /app/registration /app/registration
COPY --from=registration-builder /app/registration_create_user /app/registration_create_user
COPY --from=registration-builder /app/web /app/web
COPY --from=registration-builder /app/LICENSE.md /app/LICENSE.md
COPY --from=registration-builder /app/README.md /app/README.md

RUN apk --no-cache add curl

ENTRYPOINT ["/app/registration"]
