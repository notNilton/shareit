FROM golang:1.25-alpine AS build

WORKDIR /workspace

RUN apk add --no-cache ca-certificates tzdata git

COPY apps/backend/go.mod apps/backend/go.sum ./apps/backend/

WORKDIR /workspace/apps/backend
RUN go mod download

WORKDIR /workspace
COPY apps/backend ./apps/backend

WORKDIR /workspace/apps/backend
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/api /usr/local/bin/api

EXPOSE 8080

ENTRYPOINT ["api"]
