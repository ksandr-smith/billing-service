FROM golang:latest

RUN go version
ENV GOPATH=/
WORKDIR ./app

COPY . .

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN sed -i -e 's/\r$//' *.sh
RUN chmod +x *.sh

RUN go mod download -x
RUN go build -v ./cmd/apiserver

CMD "./apiserver"