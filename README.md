# cj-app
Go 1.16 - gin 1.7.1

IMPORTANT NOTES:

    1. If the backend require any database please use any IN-MEMORY or SQLLite database Unless mentioned in Questions Otherwise 

PROJECT START STEPS:

    Pre-requisites:
    1. Install need go version 1.1 to be installed in your system.
    2. Keep all the migrations inside the db/migration-goose or db/migration-migrate folder (check the respective folder and file name format for example migrations)
    3. goose/golang-migrate to be installed to apply any migrations

    Steps:
    1. To apply migrations via goose, do the following:
        1.a. Go to the project root directory.
        1.b. Run the command `goose -dir db/migration-goose db_type "user:password@host:port/database_name?query" up` (Change the parameters as per your specification)
    2. To apply migrations via golang-migrate, do the following:
        2.a. Go to the project root directory.
        2.b. Run the command `migrate -path db/migration-migrate -database "db_type://user:password@host:port/database_name?query" -verbose up` (Change the parameters as per your specification)
    3. To run this application, do the following:
        2.a. Go to the project root directory.
        2.b. Run the following commands to build the app:
        	- go build -o main . 
        2.c. Run the following command(s) in the terminal/command line to run the app:    
          - ./main
    
    4. Go to http://localhost:8080 in your browser to view it.
    
    CLOUD-IDE SETUP STEPS(follow the below steps in case you are using the Cloud IDE instead of your Local IDE):
	1. Please run the below commands from the project root to setup MySQL and MongoDB in this workspace:
		- chmod 0755 ./database-setup.sh
		- sh ./database-setup.sh
	2. In case you want to connect to MySQL or MongoDB, kindly use the following credentials in your application:
		2.a. MySQL
			- host: localhost
			- port: 3306
			- username: root
			- password: admin
			- database: db
		2.b. MongoDB
			- host: localhost
			- port: 27017
			- username: root
			- password: admin
			- database: db





## Added By Test-Taker
Below Instructions are added by test taker i.e. in this case me(NAVEEN GARG).
This is done in case the checker is unable to run the code on his system.
Once you are in the root directory do:

RUN the below command:
go build main.go
./main

Migration has been setup already using the below commands:
cd db/migration-goose
goose sqlite3 ./memory-sqlite.db create questions sql
goose sqlite3 ./memory-sqlite.db create quiz sql
// content of created questions and quiz sql file is updated as required
goose sqlite3 ./memory-sqlite.db up
goose sqlite3 ./memory-sqlite.db status



