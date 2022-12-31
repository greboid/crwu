FROM ghcr.io/greboid/dockerfiles/golang as compose

ENV GOBIN /build

RUN set -eux; \
    VER=$(git -c 'versionsort.suffix=-' ls-remote --exit-code --refs --sort='version:refname' --tags https://github.com/docker/compose '*.*.*' | tail -n 1 | cut -d '/' -f 3); \
    git clone -c advice.detachedHead=false --depth=1 -b $VER --single-branch https://github.com/docker/compose /src;

WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go install -a -trimpath -ldflags="-extldflags \"-static\" -buildid= -s -w" ./...

FROM ghcr.io/greboid/dockerfiles/golang as build

ENV GOBIN /build

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download

COPY . /src
RUN CGO_ENABLED=0 GOOS=linux go install -a -trimpath -ldflags="-extldflags \"-static\" -buildid= -s -w" ./...

COPY --from=compose /build/cmd /rootfs/composemount/docker-compose

RUN set -eux; \
    mkdir -p /rootfs/bin; \
    mkdir -p /rootfs/lib; \
    cp /usr/sbin/chroot /rootfs/bin/chroot; \
    cp -R /lib /rootfs/composemount/lib; \
    cp /bin/mount /rootfs/bin/mount; \
    cp /bin/mkdir /rootfs/bin/mkdir; \
    cp /build/executor /rootfs/executor; \
    cp /build/listener /rootfs/listener

FROM ghcr.io/greboid/dockerfiles/baseroot

COPY --from=build /rootfs /
