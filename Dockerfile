FROM golang
LABEL name="playground"
LABEL maintainer="flw@cpan.org"

RUN apt update && apt install -y locales locales-all
RUN go get github.com/flw-cn/playground/cmd/play

ENV LANG en_US.UTF8

ENTRYPOINT ["play"]
