# STAGE: build
FROM golang:1.14.2-stretch AS builder

COPY . /src
WORKDIR /src

RUN make all

WORKDIR /dist
RUN cp /src/gitsink .

# STAGE: final
FROM alpine

# Add Tini
RUN apk add --no-cache tini
# Tini is now available at /sbin/tini

WORKDIR /app
COPY --from=builder /dist/gitsink .

# Run your program under Tini
# Configure runtime behavior
ENV PATH="/app:${PATH}"
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["gitsink"]
