FROM golang

RUN go get github.com/tools/godep

ADD . /go/src/github.com/kpashka/linda
RUN cd /go/src/github.com/kpashka/linda && godep restore

RUN go install github.com/kpashka/linda

ENTRYPOINT /go/bin/linda -c $LINDA_CONFIG

EXPOSE 8080