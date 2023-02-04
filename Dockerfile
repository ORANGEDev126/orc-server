FROM golang:1.16-buster AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
COPY config.json /config.json
COPY /static/*.* /static/*.*
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /orc-server

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /orc-server /orc-server
COPY --from=build /static/*.* /static/*.*
COPY --from=build /config.json /config.json

EXPOSE 8000
EXPOSE 80
USER nonroot:nonroot

ENTRYPOINT ["/orc-server"]