# STAGE: build
FROM golang:1.14.2-stretch AS builder

COPY . /src
WORKDIR /src

RUN make all

WORKDIR /dist
RUN cp /src/github-migration-cli .

# STAGE: final
FROM debian:stretch-slim

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

WORKDIR /app
COPY --from=builder /dist/github-migration-cli .

# Run your program under Tini
# Configure runtime behavior
ENV PATH="/app:${PATH}"
ENTRYPOINT ["/tini", "--"]
CMD ["github-migration-cli"]