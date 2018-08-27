FROM golang:1.10
EXPOSE 8080

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH


RUN mkdir -p $GOPATH/src/filestat
VOLUME ./src/:/go/src/
# now copy your app to the proper build path
ADD ./src $GOPATH/src/

# should be able to build now
WORKDIR $GOPATH/src/filestat
RUN go build -o filestat .
