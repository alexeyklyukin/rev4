FROM golang:1.13

WORKDIR /build

COPY . .

RUN mkdir -p /launch && go build -o /launch/rev4 && rm -rf /build/*

CMD ["/launch/rev4"]