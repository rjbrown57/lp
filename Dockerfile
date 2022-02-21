FROM ubuntu:20.04 as base

FROM scratch
COPY --from=base /tmp/ /tmp/
COPY lp /lp
ENTRYPOINT ["./lp"]
