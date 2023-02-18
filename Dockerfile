FROM ubuntu:22.04 as base

FROM scratch
COPY --from=base /tmp/ /tmp/
COPY lp /lp
ENTRYPOINT ["./lp"]
