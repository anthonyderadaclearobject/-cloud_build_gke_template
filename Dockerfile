FROM golang:1.14-alpine AS build

RUN apk add --update --no-cache gcc git build-base

WORKDIR /src/
COPY main*.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/go-todo

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /bin/go-todo /bin/go-todo
ENTRYPOINT ["/bin/go-todo"]
