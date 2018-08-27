# FROM golang:latest
#
# ADD  . /go/src/interview
# RUN go get -v github.com/gocql/gocql/...
# RUN go get interview/src/filestat
# #ENV GOPATH=/golib/
# #RUN /usr/local/go/bin/go get github.com/julienschmidt/httprouter
# #CMD /usr/local/go/bin/go run /code/app/main.go

FROM golang:1.8

# RUN mkdir -p /go/src/filestat
# RUN cd /go/src/filestat

ADD ./src /go/src/
ADD ./data /data

RUN go get github.com/lib/pq/...
#
# RUN cd ../

EXPOSE 8080
