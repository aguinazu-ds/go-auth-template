FROM golang:1.22-alpine as builder

WORKDIR /app
RUN apk add --no-cache make nodejs npm

COPY . ./
RUN make install
RUN make build

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/bin/go-auth-template /go-auth-template
COPY --from=builder /app/.env /.env

EXPOSE 3000
CMD ["./go-auth-template"]

