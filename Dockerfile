FROM golang:1.13
LABEL maintainer=Caaalabash
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY ./ /app
RUN cd /app && go build app.go 
