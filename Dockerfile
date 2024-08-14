FROM golang:1 as build
WORKDIR /go/src/app/
COPY ./ ./
RUN go mod download && go build -o echo ./...

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /go/src/app/echo .
USER 65532:65532
CMD ["./echo"]
