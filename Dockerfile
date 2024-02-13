FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN make build

FROM debian:latest
RUN apt-get update && apt-get install -y ca-certificates
RUN useradd pugdns
WORKDIR /app
COPY --from=build /app/pugdns /app
USER pugdns:pugdns

ENTRYPOINT [ "/app/pugdns" ]
