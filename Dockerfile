FROM --platform=$BUILDPLATFORM golang:1.27 AS build
ARG TARGETARCH
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /ahab

FROM scratch
COPY --from=build /ahab /ahab
EXPOSE 8080
ENTRYPOINT ["/ahab"]
