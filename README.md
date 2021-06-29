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

## Running the application
This application run at port 8080 of localhost and and query the API via  the url:

`http://localhost:8080/releases`

**Please make sure that port 8080 is available on your machine.**

## Parameters
-  The from and until param are required
-  The from and until param are string date with format yyyy-mm-dd
-  The artist param is not required
-  The order does not matter


## Example Request

**Request**
`GET http://localhost:8080/releases?until=2021-01-01&from=2021-01-01`

```
[
    {
        "release_at": "2021-01-01",
        "songs": [
            {
                "name": "Happy New Year 2021 - Instrumental",
                "artist": "Titobeats"
            },
            {
                "name": "Suave Moods for 2021",
                "artist": "Relaxing Piano Music Girl"
            },
            {
                "name": "Little Donkey",
                "artist": "Christmas 2019"
            },
            {
                "name": "Uplifting Ambiance for 2021",
                "artist": "Coffee House Smooth Jazz Playlist"
            },
            {
                "name": "Friendly New Years Resolutions",
                "artist": "Peaceful Autumn Instrumental Jazz"
            },
            {
                "name": "You Are Not Alone",
                "artist": "Krezip"
            },
            {
                "name": "I Saw Three Ships",
                "artist": "Christmas 2019"
            },
            {
                "name": "New Years 2021",
                "artist": "Cowboy Nemo"
            },
            {
                "name": "Sine moj (Live)",
                "artist": "Roksana"
            },
            {
                "name": "We Three Kings",
                "artist": "Christmas 2019"
            }
        ]
    }
]
``` 

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