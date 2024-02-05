FROM docker.io/library/golang:latest AS build
WORKDIR /work
COPY go.mod *.go ./
COPY public ./public
RUN CGO_ENABLED=0 go build -o web

FROM gcr.io/distroless/static-debian12:nonroot AS image
COPY --from=build /work/web /usr/local/bin/web
EXPOSE 8080
CMD ["/usr/local/bin/web"]
