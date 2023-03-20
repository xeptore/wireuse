FROM docker.io/library/golang:1 AS build

RUN apt-get update && apt-get upgrade -y && apt-get install -y libssl-dev pkg-config

RUN update-ca-certificates

RUN useradd --create-home --uid 1000 nonroot

USER nonroot

WORKDIR /home/nonroot

COPY --chown=nonroot:nonroot go.sum go.mod ./

ENV GOPROXY=https://goproxy.io,direct

RUN go mod download -x

COPY --chown=nonroot:nonroot . .

RUN make build

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=build --chown=nonroot:nonroot /home/nonroot/bin/ingest ingest

ENV TZ=UTC

ENTRYPOINT [ "./ingest" ]