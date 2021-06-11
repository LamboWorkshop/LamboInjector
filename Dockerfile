FROM golang:latest

WORKDIR /lamboInjector

ADD . .

RUN make dep
RUN make

CMD ./build/lamboInjector