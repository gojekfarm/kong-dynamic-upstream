FROM golang:1.13.4-alpine

RUN apk add build-base git
RUN go get github.com/Kong/go-pdk
RUN mkdir /src
ADD . /src/
WORKDIR /src
RUN echo "replace github.com/Kong/go-pdk =>  /go/src/github.com/Kong/go-pdk" >> go.mod
RUN go build -o dynamicupstream.so -buildmode=plugin handler.go
WORKDIR /
RUN wget https://github.com/Kong/go-pluginserver/archive/v0.4.0.zip
RUN unzip v0.4.0.zip
WORKDIR /go-pluginserver-0.4.0
RUN echo "replace github.com/Kong/go-pdk =>  /go/src/github.com/Kong/go-pdk" >> go.mod
RUN go install

FROM kong:2.0.4-alpine

USER root

RUN rm /usr/local/bin/go-pluginserver
COPY --from=0 /go/bin/go-pluginserver /usr/local/bin/go-pluginserver
RUN mkdir /go_plugins_dir
COPY --from=0 /src/dynamicupstream.so /go_plugins_dir/dynamicupstream.so
RUN luarocks install lua-resty-openidc
RUN luarocks install kong-enhanced-oidc

USER kong
