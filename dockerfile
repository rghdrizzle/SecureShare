FROM golang:alpine3.19 AS build

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o secureshare

FROM scratch AS runtime

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /app/secureshare ./

EXPOSE 3333

CMD [ "./secureshare" ]