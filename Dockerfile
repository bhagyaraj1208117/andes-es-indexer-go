FROM golang:1.20.5 as builder

RUN apt-get update && apt-get install -y

WORKDIR /dharitri
COPY . .

WORKDIR /dharitri/cmd/elasticindexer

RUN go build -o elasticindexer

# ===== SECOND STAGE ======
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y

RUN useradd -m -u 1000 appuser
USER appuser

COPY --from=builder /dharitri/cmd/elasticindexer /dharitri

EXPOSE 22111

WORKDIR /dharitri

ENTRYPOINT ["./elasticindexer"]
CMD ["--log-level", "*:DEBUG"]
