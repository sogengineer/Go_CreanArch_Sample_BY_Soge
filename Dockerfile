FROM golang:1.22-alpine as builder

ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk upgrade --update && \
    apk --no-cache add git

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download


COPY ./src ${ROOT}
RUN CGO_ENABLED=0 GOOS=linux go build -o $ROOT/binary



FROM scratch as prod

ENV ROOT=/go/src/app
WORKDIR ${ROOT}
COPY --from=builder ${ROOT}/binary ${ROOT}

EXPOSE 8080
CMD ["/go/src/app/binary"]