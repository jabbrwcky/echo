FROM golang:alpine as build
WORKDIR /go/src/echo/
COPY ./ ./
RUN go mod download
RUN CGO_ENABLE=0 go build -o echo-server ./...

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /go/src/echo/echo-server /echo-server
USER 65532:65532
ENTRYPOINT ["/echo-server"]
