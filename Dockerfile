FROM golang:1.20

RUN mkdir /app
ADD ./app
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY main_test.go .

RUN go mod download
RUN go build -o main .

EXPOSE 4040

CMD ["./main"]
