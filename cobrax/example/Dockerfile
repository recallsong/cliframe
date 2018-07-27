FROM alpine:3.6

RUN mkdir -p /example-app/conf
Copy bin/linux-amd64-example /example-app/example

WORKDIR /example-app
VOLUME /example-app/conf

CMD ["/example-app/example"]