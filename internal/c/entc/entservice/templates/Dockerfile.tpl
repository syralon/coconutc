FROM golang:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add --no-cache tzdata
RUN apk add git

WORKDIR /opt/app

ADD . .

ENV GOPROXY=https://goproxy.cn,direct

RUN go build -ldflags "-X {{.Module}}/version.BuildTime=$(date '+%Y-%m-%dT%H:%M:%S') -X {{.Module}}/version.Version=$(git rev-parse HEAD)" ./cmd/{{.Module|basepath}}

FROM alpine:latest AS runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai

ENV TZ Asia/Shanghai

WORKDIR /opt/app

COPY --from=builder /opt/app/{{.Module|basepath}} .

ENTRYPOINT exec ./{{.Module|basepath}}