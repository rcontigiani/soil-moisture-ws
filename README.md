##  Soil Moisture Web Service
Soil moisture web service defines a set of API to retrieve data from AWS DynamoDB Table with a defined structure. 

This is a part of a automated sprinkler project, this part is strictly related to [Soil Moisture Lambda](https://github.com/rcontigiani/soil-moisture-lambda.git) that populate DynamoDB table.

## DynamoDB Table Structure

The DynamoDB Table must have these keys:
- **Type** (*string*) - Partition Key
- **Date** (*number*) - Sort Key

It must be in **eu-west-1** region and named **Sprinkler**, which will be replaced with a configurations in next releases.

## Environment Configuration

To allow connection between the web service and AWS DynamoDB a environment file is required. The file must be structured as follows:

```
AWS_ACCESS_KEY_ID=<aws access key>
AWS_SECRET_ACCESS_KEY=<aws secret>
region=eu-west-1
output=json
```

This informations can be retrieved in your `~/.aws/config` and  `~/.aws/credentials` files.

## Api

### Health Check

Used to retrieve service status

**URL** : `/healthCheck`

**Method** : `GET`

**Auth required** : NO

**Response**

```json
{
    "status": true
}
```

### Get Last

Used to retrieve the most recent record

**URL** : `/getLast`

**Method** : `POST`

**Auth required** : NO

**Request Body**

```json
{
	"Type": "SM",
	"DateStart": null,
	"DateEnd": null
}
```

**Response**

```json
{
    "Id": "0a9f6670-188e-4e7b-b31a-7bc1b0c767c2",
    "Date": 1590357294,
    "Value": 55,
    "Type": "SM"
}
```

### Get Range

Used to retrieve a subset of rows

**URL** : `/getRange`

**Method** : `POST`

**Auth required** : NO

**Request Body**

```json
{
	"Type": "SM",
	"DateStart": 1590356289,
	"DateEnd": 1590356905
}
```

**Response**

```json
[
    {
        "Id": "21cf5919-148d-4804-859c-0cffd591487a",
        "Date": 1590356289,
        "Value": 32,
        "Type": "SM"
    },
    {
        "Id": "5a1b6c29-ba5f-4e50-a342-1f8122c7c895",
        "Date": 1590356905,
        "Value": 35,
        "Type": "SM"
    }
]
```


## Deploy

This solution can be executed in local environment or in docker environment.

### Locally

To run the solution locally is needed *Go* installed and configured in your machine and AWS Cli avalaible with configured profile. 

If you launch the solution locally, you don't have to use the environment file because the credentials are taken from the files in `~/.aws` folder. 

##### Run directly

Run directly the code
`go run main.go`

##### Build and Run

Build the project
`go build`

Run the executable file
`./soil-moisture-ws`

### Docker

##### Compose

The project can be launched with docker compose running `docker-compose up` command in the terminal if you have the source code.

##### Directly with Docker

Otherwise, if you have pulled the image from [docker hub](https://hub.docker.com/repository/docker/rcontigiani/soil-moisture-ws) with `docker pull rcontigiani/soil-moisture-ws:latest` or you have build it in your machine with this command `docker build -t soilmoisturews:<version> . -f Dockerfile` you can start the container using the following command: `docker run -it --env-file=.env -p 8080:8080 <image id>`.

The image id can be retrieved using `docker images` command in your terminal.

###### Notes

Clearly, as you can see in the previous commands, the environment file is needed such if you run the project both with compose that directly with docker. The only difference between the two methodologies is that compose implicitly takes file named .env in the root directory of the project as environment configuration file, while running with docker directly the environment file must be explained in the start command. 
