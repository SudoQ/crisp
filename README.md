# satchel
[![Travis](https://img.shields.io/travis/SudoQ/satchel.svg)](https://travis-ci.org/SudoQ/satchel)
[![Docker Stars](https://img.shields.io/docker/stars/sudoq/satchel.svg)](https://hub.docker.com/r/sudoq/satchel/)
[![Docker Pulls](https://img.shields.io/docker/pulls/sudoq/satchel.svg)](https://hub.docker.com/r/sudoq/satchel/)

Periodically scrapes provided URL and hosts the data as a RESTful HTTP API

#Docker usage
```
$ docker pull sudoq/satchel:master
$ docker run -p 80:8080 --rm satchel <URL>
```

#Example
```
$ docker run -p 80:8080 --rm satchel https://api.github.com/events
$ curl localhost
```
