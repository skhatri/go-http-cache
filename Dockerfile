FROM golang:1.21 as builder
RUN mkdir /build
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0
RUN go mod vendor
RUN go build -o app

FROM scratch

COPY --from=builder /build/app /app
COPY --from=builder /build/config.yaml /conf/config.yaml
ENV CONF_FILE=/conf/config.yaml
EXPOSE 8070
CMD ["/app"]

