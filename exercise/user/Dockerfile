FROM alpine:latest

RUN mkdir /user

COPY bin/user-api /user
COPY api/etc/user.yaml /user

RUN chmod +x /user/user-api

WORKDIR /user

ENTRYPOINT ["./user-api", "-f", "user.yaml"]

#
#echo -e "https://mirrors.aliyun.com/alpine/v3.15/main\nhttps://mirrors.aliyun.com/alpine/v3.15/community" > /etc/apk/repositories
#apk update -v
#
#CMD ["ping" ,"baidu.com"]
#
#RUN mkdir /go && cd /go \
#    && wget --no-check-certificate https://golang.google.cn/dl/go1.21.0.linux-amd64.tar.gz \
#    && tar -C /usr/local -zxf go1.21.0.linux-amd64.tar.gz \
#    && rm -rf /go/go1.21.0.linux-amd64.tar.gz \
#    && mkdir /lib64 \
#    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
#
#ENV GOPATH /go
#ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
#
#CMD ["ping" ,"baidu.com"]