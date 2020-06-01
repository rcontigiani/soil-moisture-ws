

## Build docker image
`docker build -t soilmoisturews:<version> . -f Dockerfile`

## Run docker image

`docker images`

`docker run -it --env-file=.env -p 8080:8080 soilmoisturews:<version>`
