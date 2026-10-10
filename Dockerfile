# Plain Docker build for this fork; upstream builds its image with Nix.
FROM golang:1.27.2-alpine AS build
ARG VERSION=dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /plundrio ./cmd/plundrio

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /plundrio /bin/plundrio
ENTRYPOINT ["/bin/plundrio"]
