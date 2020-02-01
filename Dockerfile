FROM golang:1.13
LABEL maintainer=Caaalabash
WORKDIR /app
COPY ./ /app
RUN cd /app && go build app.go 
