# syntax=docker/dockerfile:1
# Build args default to the versions pinned in boilerplate/versions.env (scripts pass them in).
ARG GO_IMAGE=golang:1.27-alpine
ARG RUNTIME_IMAGE=gcr.io/distroless/static-debian13:nonroot

FROM ${GO_IMAGE} AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
# Tests run inside the build: an image that exists has passed its tests.
RUN CGO_ENABLED=0 go test ./... \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/hello-api .

FROM ${RUNTIME_IMAGE}
COPY --from=build /out/hello-api /hello-api
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/hello-api"]
