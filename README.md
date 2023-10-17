<h1 align="center">Code Racer</h1>

<h2 align="center">
A Remote Code Execution Engine built with Go and ❤️
</h2>

<p align="center">
    <a href="https://github.com/The-Flash/code-racer/commits/master">
    <img src="https://img.shields.io/github/last-commit/The-Flash/code-racer.svg?style=for-the-badge&logo=github&logoColor=white"
         alt="GitHub last commit">
    <a href="https://github.com/The-Flash/code-racer/issues">
    <img src="https://img.shields.io/github/issues/The-Flash/code-racer.svg?style=for-the-badge&logo=github&logoColor=white"
         alt="GitHub issues">
    <a href="https://github.com/The-Flash/code-racer/pulls">
    <img src="https://img.shields.io/github/issues-pr-raw/The-Flash/code-racer.svg?style=for-the-badge&logo=github&logoColor=white"
         alt="GitHub pull requests">
</p>


## Setup

Create a .env file with the following:

```
VERSION=0.0.3
APP_NAME=thefl45h/code-racer
```

This is necessary for building the docker container and using the Makefile

## Environment variables

* ```PORT```(optional): port to serve requests. Defaults to 8000

* ```RUNNERS_PATH```(required): Path to runners script. This should be an absolute path to the ```runners``` directory in the repo. This will be mounted in each executor

* ```MNTFS```(required): Path to mounted filesystem. Code racer will create directories and files for code execution. This should be an absolute path. This will be mounted in each executor. Recommend using an absolute path to the ```mntfs``` folder in this repo.

* ```NOSOCKET```(optional): Path to ```nosocket``` binary. This should be an absolute path. If this env variable is not set, set ```DISABLE_NETWORKING``` to ```1```.

* ```PULL_IMAGES```(optional): ```true/false```. Attempt to pull images or not. May be useful of executor images are already available locally. 

* ```DISABLE_NETWORKING```(optional): ```1/0```. If ```1``` running tasks will not be able to access the internet and make network calls. ```0``` running tasks/process will have access to the internet.

## Running in development

Create an ```.env.development``` with the above environment variables. Refer to ```.env.sample```

We strongly recommend docker/docker-compose for a smooth development
process.

This project includes a (devcontainer)[https://code.visualstudio.com/docs/devcontainers/containers].

In the dev container, you can use the ```air``` command to start up the server with hot reloading.


## Running in production

Create an ```.env.production``` with the above environment variables. Refer to ```.env.sample```



Use the root ```docker-compose.yml``` file to create a production ready docker image

### Building docker image

```
docker-compose build
```

### Starting container

```
docker-compose up
```

## What is ```nosocket```?

```nosocket``` is a custom binary that uses ```libseccomp``` to block socket calls.
This is a security measure to prevent tasks submitted from accessing the internet.


## API Documentation

## Public API

```
POST https://code-racer-api.codewithflash.com/api/v1/execute
GET https://code-racer-api.codewithflash.com/api/v1/runtimes
```

**Runtimes endpoint**

```GET /api/v1/runtimes```: Endpoint to list available runtimes

Response

```json
[
    {
        "language": "python3"
    },
    {
        "language": "node"
    }
]
```

**Execute endpoint**

```POST /api/v1/execute```. Endpoint to execute code

* ```language```(required): language to use for execution. Refer to ```/api/v1/runtimes``` for available languages

* ```entrypoint```(required): the file path to start execution from. This path must be present in ```files```
* ```files```(required): An array of files containing code to be used for execution.
* ```file.name```(required): filename
* ```file.content```(required): content of file


```json
{
    "language": "node",
    "entrypoint": "main.js",
    "files": [
        {
            "name": "main.js",
            "content": "console.log('Code Racer')"
        }
    ]
}
```


```json
{
    "stderr": "",
    "stdout": "Code Racer\n",
    "executionTime": "438.757875ms",
    "preparationTime": "102.106667ms"
}
```

## Supported Languages

```python```, ```node```, ```go```

More to come