FROM golang:1.18
LABEL authon="Dai Xin"
USER root

ENV APP_KEY=""
ENV APP_SECRET=""
ENV WX_APP_ID=""
ENV WX_APP_SECRET=""
ENV WX_TOKEN=""

WORKDIR /usr/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./.

VOLUME ["/usr/data"]
EXPOSE 8000

ENTRYPOINT app -dockerInd Y -appKey $APP_KEY -appSecret $APP_SECRET -wxAppId $WX_APP_ID -wxAppSecret $WX_APP_SECRET -wxToken $WX_TOKEN