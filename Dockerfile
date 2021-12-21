FROM golang:1.17-alpine3.14 as debug

RUN apk update && apk upgrade && \
    apk add --no-cache git dpkg gcc musl-dev

ENV GOPATH /go
ENV GO111MODULE on
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /usr/src
ADD ./src ./

RUN ls -la .

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go mod download
RUN go build -o /usr/app ./cmd/main.go

COPY ./dlv.sh /
RUN chmod +x /dlv.sh
ENTRYPOINT [ "/dlv.sh" ]

FROM alpine:3.14 as prod
COPY --from=build-env /usr/app /
CMD ./usr/app