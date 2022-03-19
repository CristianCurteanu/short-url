### URL shortening API

Current application provides functionality to launch and deploy a URL shortening web service. It provides access via REST and GRPC APIs, and also to access REST API via CLI

### Installation

There are several ways to launch the servers, which are described bellow

##### Using Docker

In order to run the servers using Docker, please make sure that you have Docker along with Docker Compose installed.

If you would like to both gRPC and HTTP server, please run following command:

```shell
$ make run-all
```

If you want to run separately the gRPC, then run following command:

```shell
$ make run-grpc
```

If you want to run separately the HTTP, then run following command:

```shell
$ make run-http
```

##### Standalone runtime

As prerequisited, please make sure to install MongoDB and Redis. Afterward there will be required to have environment variables set up:

```shell
MONGO_URI=<MongoDB database URI>
REDIS_HOST=<REDIS host>
REDIS_PORT=<REDIS port>
REDIS_PASSWORD=<REDIS server password>
```

These environment variables could be stored in a file and exported with `source` command (for instance `.env` file):

```shell
$ source .env
```

Then you can use the makefile commands:

```shell
# for gRPC
$ PORT=<port> make run-grpc-standalone
# or for HTTP
$ PORT=<port> make run-http-standalone
```

### The APIs

By default APIs will run on `:3000` (REST API) and `:3001` (GRPC API) ports. In order to change the ports please access `docker-compose.yml` file

#### REST

Below will be presented endpoints of REST API.

All endpoints will have same structure for error response:

```json
{
  "key": "<short-error-identified>",
  "message": "<detailed-error-message>",
}
```

##### Create Mapping Endpoint

```shell
POST /api/mappings/
{
  "url": "<url-to-be-stored>"
}
```

It will respond with:

```json
{
  "key": "<key-of-the-url>"
}
```

##### Fetch specific Mapping Endpoint

```shell
GET /api/mappings/{:key}

```

It will respond with:

```json
{
  "key": "<key-of-the-mapping>"
  "url": "<url-to-be-stored>"
}
```

##### Delete specific Mapping Endpoint

```shell
DELETE /api/mappings/{:key}

```

It will respond with:

```json
{
  "Deleted": "<deletion-status>"
}
```

#### GRPC

The structure of protobuf for gRPC is as follows

```protobuf
service MappingsService {
  rpc GetMapping (GetMappingRequest) returns (GetMappingResponse);
  rpc CreateMapping (CreateMappingRequest) returns (CreateMappingResponse);
  rpc DeleteMapping (DeleteMappingRequest) returns (DeleteMappingResponse);
}

message DeleteMappingResponse {
  string deleted = 1;
}

message DeleteMappingRequest {
  string key = 1;
}

message CreateMappingRequest {
  string url = 1;
}

message CreateMappingResponse {
  string key = 1;
}

message GetMappingRequest {
  string key = 1;
}

message GetMappingResponse {
  string key = 1;
  string url = 2;
}
```

### Known issues:
  - Authentication and Authorization
  - Multitenant URLs storage, ie. validate URL mapping for specific tenant

### The CLI

Documentation for all commands:

```shell
NAME:
   main - A new cli application

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   create-mapping  Creates a new url key mapping
   get-mapping     Fetches url key mapping
   delete-mapping  Deletes url key mapping
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

Documentation for `create-mapping` command:

```shell
NAME:
   main create-mapping - Creates a new url key mapping

USAGE:
   main create-mapping [command options] [arguments...]

OPTIONS:
   --url value  URL that will be mapped to a key
   --help, -h   show help (default: false)
```

Documentation for `get-mapping` command:

```shell
NAME:
   main get-mapping - Fetches url key mapping

USAGE:
   main get-mapping [command options] [arguments...]

OPTIONS:
   --key value  Key of the stored URL
   --help, -h   show help (default: false)
```

Documentation for `delete-mapping` command:

```shell
NAME:
   main delete-mapping - Deletes url key mapping

USAGE:
   main delete-mapping [command options] [arguments...]

OPTIONS:
   --key value  Key of the stored URL
   --help, -h   show help (default: false)
```
