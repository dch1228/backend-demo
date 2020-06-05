# Golang build container
FROM golang:1.13.10-alpine

WORKDIR $GOPATH/src/github.com/duchenhao/backend-demo

COPY . .

RUN CGO_ENABLED=0 go build -mod vendor -o ./bin/backend-demo

# Final container
FROM alpine:3.11.5

LABEL maintainer="chenhao.du <chenhao.du@qq.com>"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache

ENV PATH="/usr/share/backend-demo/bin:$PATH"

WORKDIR /usr/share/backend-demo

COPY --from=0 /go/src/github.com/duchenhao/backend-demo/bin/backend-demo ./bin/
COPY /conf/dev.yml /etc/backend-demo/config.yml

EXPOSE 8080

ENTRYPOINT [ "backend-demo", "run", "-c", "/etc/backend-demo/config.yml" ]