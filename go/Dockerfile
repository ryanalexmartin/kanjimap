FROM golang:1.21

WORKDIR /app

RUN mkdir /vue
RUN mkdir /vue/dist

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8081

CMD ["sh", "-c", "sleep 10 && ./main"]
# CMD ["sh", "-c", "tail -f /dev/null"] 
