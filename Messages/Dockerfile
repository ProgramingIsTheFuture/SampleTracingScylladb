FROM golang:1.17

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"

COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o messages
ENTRYPOINT ["./messages"]
