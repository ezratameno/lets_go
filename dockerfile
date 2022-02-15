FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN go get github.com/ezratameno/lets_go
RUN cd /build && git clone https://github.com/ezratameno/lets_go.git

RUN cd /build/lets_go/cmd/web && go build

EXPOSE 4000
RUN ls /build/lets_go/ -la

ENTRYPOINT ["/build/lets_go/cmd/web/main"]
