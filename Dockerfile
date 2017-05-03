FROM golang
ADD . /go/src/jenkins-ci
WORKDIR /go/src/jenkins-ci
RUN go get && go install jenkins-ci
ENTRYPOINT /go/bin/jenkins-ci