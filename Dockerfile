FROM golang:1.7
WORKDIR $GOPATH/src/github.com/loadimpact/demo-site
ADD . .
RUN go get ./... && go install .
EXPOSE 8000
CMD ["demo-site"]
