FROM alpine

MAINTAINER Denis Tsai <denistsai@blazingorb.com>

RUN mkdir -p /app

COPY artifacts/bin/johnnyfive.linux.amd64 /usr/bin/johnnyfive

WORKDIR /app

ENTRYPOINT ["/usr/bin/johnnyfive"]
