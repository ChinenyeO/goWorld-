# Production image for the 12_protobuf UserService.
#
# The Go module lives at the repository root, so build from there:
#
#   docker build -f 12_protobuf/Dockerfile -t goworld-userservice .
#
# Configuration (DATABASE_URL, DB_*) is supplied at run time; see config.go.

# ---- Build stage -----------------------------------------------------------
# Must satisfy the "go" directive in go.mod.
FROM golang:1.26.2-bookworm AS build

WORKDIR /src

# Dependencies change less often than source, so download them in their own
# layer; it is reused until go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO disabled produces a static binary that runs on a libc-free base image.
# -trimpath and -s -w drop local paths and debug symbols from the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/userservice ./12_protobuf

# ---- Runtime stage ---------------------------------------------------------
# distroless/static contains only CA certificates, tzdata and /etc/passwd:
# no shell, no package manager, no Go toolchain. The :nonroot tag runs as
# uid/gid 65532.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/userservice /userservice

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/userservice"]
