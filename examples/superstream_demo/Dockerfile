FROM golang:1.16 AS build

ARG service_name

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app ./${service_name}

FROM scratch
COPY --from=build /app /app
EXPOSE ${PORT}

ENTRYPOINT ["/app"]
