# gowiki-tutorial

This is a practice for [Writing Web Applications \- The Go Programming Language](https://golang.org/doc/articles/wiki/).

Running on Heroku https://tranquil-earth-43383.herokuapp.com

## Run locally

```
$ go get -u github.com/takatoshiono/gowiki-tutorial
$ cd $GOPATH/src/github.com/takatoshiono/gowiki-tutorial
$ go build wiki.go
$ ./wiki -p 8080
```

Visiting http://localhost:8080

## Deploy to Heroku

```
$ heroku create
$ git push heroku master
$ heroku open
```

