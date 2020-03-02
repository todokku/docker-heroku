FROM google/cloud-sdk:alpine
RUN apk add bash go
SHELL ["bash", "-c"]
WORKDIR /root
COPY ./go/src /root/go/src
RUN go get -v github.com/gorilla/mux
RUN go build -v -o /bin/entrypoint github.com/livecodecreator/docker-heroku
ENTRYPOINT ["entrypoint"]
