FROM iron/go:1.7

ADD ./auth-server /app/

WORKDIR /app

RUN chmod +x auth-server

ENTRYPOINT ["./auth-server"]