FROM golang:1.13
RUN mkdir -p /build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rev4 .

FROM alpine:latest
RUN mkdir -p /launch
WORKDIR /launch
COPY --from=0 /build/rev4 .
CMD ["./rev4"]
