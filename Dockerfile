FROM golang:1.23.6-alpine

WORKDIR /app 

COPY go.mod go.sum ./
RUN go mod download 

COPY . . 

# RUN go install github.com/cosmtrek/air@latest

COPY ./air /usr/local/bin/air

RUN chmod +x /usr/local/bin/air

CMD ["air"]
