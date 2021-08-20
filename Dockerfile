FROM golang:1.16.3-buster as build

RUN wget https://codejudge-starter-repo-artifacts.s3.ap-south-1.amazonaws.com/backend-project/python/pre-build.sh
RUN chmod 775 ./pre-build.sh
RUN sh pre-build.sh

RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash

RUN apt-get update

RUN apt-get install -y migrate

RUN migrate -version

RUN go get -u github.com/pressly/goose/cmd/goose

RUN wget https://codejudge-starter-repo-artifacts.s3.ap-south-1.amazonaws.com/backend-project/database/db-setup.sh
RUN chmod 775 ./db-setup.sh
RUN sh db-setup.sh

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
EXPOSE 8080
RUN go build -o main . 

RUN wget https://codejudge-starter-repo-artifacts.s3.ap-south-1.amazonaws.com/backend-project/golang/run-2.sh
RUN chmod 755 ./run-2.sh
CMD sh run-2.sh