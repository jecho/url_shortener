FROM golang:1.7
MAINTAINER jecho@foxley.co
RUN mkdir -p /app
WORKDIR /app
ADD . /app
RUN go get github.com/go-sql-driver/mysql && go get github.com/gorilla/mux
RUN go build -o main .
EXPOSE 22222
CMD ["/app/main"]