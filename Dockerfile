FROM golang:1.17-alpine as builder

WORKDIR /src
COPY . ./
RUN go build -v -o raisepr .

FROM alpine as finale

EXPOSE 9999
USER 1000

COPY --from=builder /src/raisepr /raisepr

CMD ["/raisepr"]