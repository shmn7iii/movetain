FROM golang:1.17.8 as builder
WORKDIR /movetain
COPY . /movetain
RUN go build -o main && chmod +x ./main
CMD ["./main"]
