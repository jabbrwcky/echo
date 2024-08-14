FROM golang:1 as build
WORKDIR /go/src/app/
COPY ./ ./
RUN go build -o echo ./...

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /go/src/app/echo .
USER 65532:65532
CMD ["./echo"]
