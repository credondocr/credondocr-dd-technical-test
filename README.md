# credondocr-dd-technical-test


## Tech Stack
- Golang 
- Postgresql
- Gin
- viper
- uuid
- gorm
- Docker

## Setup
In order tu setup and run the application just use the following commands:
`docker-compose up`

The application not need any migration or initial script, everything is ready out of the box using docker

If you want to make some code changes, please recreate the orquestation with:
`docker-compose up --build --force-recreate`

## Architecture
![](https://i.imgur.com/Wdv8do7.png)
## Workflow

- The request hit the endpoint `/release`
- The application check if the requested dates are cached, if any day from params is stored in the database, those days will be skipped to reduce the cost.
- If the requested days between dates is less than 25 (excluding the cached days), the application will request the data using the `daily method`, if not, the application will request the data using the `montly method`.
- Since the data is ready, the query will request the information using the params (from, until and artist) to get the data with format (the format happends on database to reduce the for loops)
- Finally the application return the response with the requested data.


## Final Thoughts
This challenge was really fun, I enjoy a lot during the process because I did 2 different approaches. This opportunity was really good to build something using go which is good to my carear and future.
Thanks Esteban, the meetings, troubleshooting and Q&A was really interesting and I learn a lot with the challenge so is profit to me