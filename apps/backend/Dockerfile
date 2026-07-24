FROM golang:1.22-alpine AS build
WORKDIR /src
COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download
COPY apps/backend/ .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM alpine:3.20
COPY --from=build /out/api /usr/local/bin/api
EXPOSE 8080
ENTRYPOINT ["api"]
