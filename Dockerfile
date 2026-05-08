ARG GO_VERSION=1.25.0

FROM golang:${GO_VERSION} AS build
RUN go env -w GOARCH=amd64
WORKDIR /src

## download dependencies
COPY go.sum go.mod ./
RUN go mod download

COPY . /src
RUN CGO_ENABLED=0 go build -o main ./cmd/*.go


FROM gcr.io/distroless/static-debian12
COPY --from=build /src/main /app/main
EXPOSE 8888
ENV PORT 8888
CMD ["/app/main"]
