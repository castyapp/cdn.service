# Casty cdn.service
This project handles storage for s3 buckets like: avatars, posters, subtitles and etc... on minio s3 bucket

## Prerequisites
* Minio S3 Bucket Server [Getting Start](https://github.com/minio/minio)

## Clone the project
```bash
$ git clone https://github.com/CastyLab/cdn.service.git
```

## Configuraition
Make a copy of `config.example.yml` for your own configuration. save it as `config.yml` file.
```bash
$ cp config.example.yml config.yml
```

### s3 configuration
Put your running s3 bucket server details
```yaml
endpoint: "s3.example.com"
region: "" # not required if using minio
use_https: true
insecure_skip_verify: true
access_key: ""
secret_key: ""
sentry_dsn: "" # for logging in sentry * Optional
```

You're ready to Go!

## Run project with go compiler
You can simply run the project with following command
```bash
$ go run server.go
```

or if you're considering building the project use
```bash
$ go build -o server .
```

## or Build/Run the docker image
```bash
$ docker build . --tag=casty.cdn

$ docker run -dp --restart=always 5555:5555 casty.cdn
```
