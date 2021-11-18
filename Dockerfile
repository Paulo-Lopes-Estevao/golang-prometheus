FROM golang as builder

WORKDIR /app

ENV GOOS=linux
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOARCH=amd64

RUN go mod init github.com/Paulo-Lopes-Estevao/prometheus_gopher

COPY . /app/

RUN go build -o main

EXPOSE 9000

RUN chmod +x ./main

CMD ["./main"]