FROM nimmis/alpine-golang:1.9.1

RUN mkdir -p /goroot/src/github.com/zuuby/zuuby-ipfs && \

    apk update && apk upgrade && \

    apk add git && \

    go get -u github.com/ipfs/ipfs-update && \

    ipfs-update install v0.4.13 && \

    which ipfs && \

    ipfs init && \

    apk del git && \

    rm -rf /var/cache/apk/*


WORKDIR /goroot/src/github.com/zuuby/zuuby-ipfs

COPY . .

EXPOSE 5000

ENTRYPOINT go run cmd/zuupfs/main.go
