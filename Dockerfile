FROM golang
LABEL name="playground"
LABEL maintainer="flw@cpan.org"

RUN apt update && apt install -y locales locales-all
RUN go get -v golang.org/x/tools/cmd/goimports
RUN go get -v github.com/flw-cn/playground/cmd/play

RUN echo "Asia/Shanghai" > /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata

ENV LANG en_US.UTF8

ENTRYPOINT ["play"]
