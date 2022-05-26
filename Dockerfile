FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src

RUN cd app/system &&  make build

FROM alpine

COPY --from=builder /src/app/system/bin /app

RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.4/main\nhttp://mirrors.aliyun.com/alpine/v3.4/community" >  /etc/apk/repositories \
&& apk update && apk add tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Shanghai/Asia" > /etc/timezone \
&& apk del tzdata

WORKDIR /app

EXPOSE 8000
EXPOSE 9000

CMD ["./cmd","-conf","/app/configs"]
