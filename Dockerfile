FROM golang:1.20

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o main .

EXPOSE 4040

CMD ["/app/main"]
