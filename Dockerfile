FROM golang:alpine  as builder 

WORKDIR /app
COPY . .
RUN go mod download 

RUN CGO_ENABLED=0 GOOS=linux go build   -o taiga-hooker /app/cmd
# Set the time zone


FROM alpine:3.18.3 as production
RUN apk add tzdata
WORKDIR /app

ENV TZ=Asia/Tbilisi
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


COPY --from=builder /app/taiga-hooker .



EXPOSE 8080


CMD ["./taiga-hooker"]
