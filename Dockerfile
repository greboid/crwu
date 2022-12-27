FROM ghcr.io/greboid/dockerfiles/golang as build

ENV GOBIN /build

COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go install -a -trimpath -ldflags="-extldflags \"-static\" -buildid= -s -w" ./...

FROM ghcr.io/greboid/dockerfiles/base

COPY --from=build --chown=65532:65532 /build /
