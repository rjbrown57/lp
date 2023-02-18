FROM alpine:latest
COPY lp /lp
ENTRYPOINT ["./lp"]
