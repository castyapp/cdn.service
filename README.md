# Casty cdn.service
This project handles storage for s3 buckets like: avatars, posters, subtitles and etc... on minio s3 bucket

## Prerequisites
* Minio S3 Bucket Server [Getting Start](https://github.com/minio/minio)

## Clone the project
```bash
$ git clone https://github.com/castyapp/cdn.service.git
```

## Configuraition
Make a copy of `example.config.hcl` for your own configuration. save it as `config.hcl` file.
```bash
$ cp example.config.hcl config.hcl
```

### s3 configuration
Put your running s3 bucket server details
```hcl
endpoint              = "localhost:9000"
access_key            = "access_key"
use_https             = false
insecure_skip_verify  = true
secret_key            = "super-secure-secret-key"
sentry_dsn            = ""
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
