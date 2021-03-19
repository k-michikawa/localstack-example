FROM golang:1.15.6

WORKDIR /go/app
COPY src/ .
RUN go install
CMD ["go", "run", "."]
