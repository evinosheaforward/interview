FROM golang:latest

ADD  . /go/src/interview
RUN go install interview/src/filestat
#ENV GOPATH=/golib/
#RUN /usr/local/go/bin/go get github.com/julienschmidt/httprouter
#CMD /usr/local/go/bin/go run /code/app/main.go

EXPOSE 8080
