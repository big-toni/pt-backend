FROM golang:1.19

WORKDIR /src
COPY . /src
RUN go build -o pt-backend
EXPOSE 3000
CMD [ "./pt-backend" ]

