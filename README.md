# polyverse-security/redirect

Listen to a http or https port and redirect ("HTTP/1.1 301 Moved Permanently") to '&lt;scheme&gt;://&lt;host&gt;' while preserving the path and query string.

## Usage
```
Usage: ./redirect [options]

options:
  -bind     bind-to port. (default: :80)
  -scheme   redirect-to http or https. (default: https)
  -host     redirect-to host.
  -port     redirect-to port.
  --debug   verbose output. (default: false)
  --help    print usage. (default: false)
```

## Hello World (Scheme redirection)
The recommended deployment is to run this as a container. For example, to redirect all http://&lt;foobar&gt;/hello/world traffic to https://polyverse.io, you would launch the container with:
```
$ docker run -d -p 80:80 polyverse/redirect -bind=:80 -scheme=https
```

Any pings to http://hostname/url?query=string will get redirected to https://hostname/url?query=string (regardless of hostname.)
This mode is particularly useful when running just a dumb-redirect to the HTTPS version of the SAME URL for multi-host endpoints
(say you're running a server that can be accessed as localhost AND a different name .)

## Host based redirection

```
$ docker run -d -p 80:80 polyverse/redirect -bind=:80 -host=polyverse.io -scheme=https
```

This will redirect ALL traffic hitting the endpoint to the polyverse.io hostname:

* http://localhost/foo -> https://polyverse.io/foo
* http://localhost:5601/foo -> https://polyverse.io:5601/foo


## Port based redirection
```
$ docker run -d -p 80:80 polyverse/redirect -bind=:80 -host=polyverse.io -scheme=https
```

This will redirect ALL traffic hitting the endpoint to the polyverse.io hostname:

* http://localhost/foo -> https://polyverse.io/foo
* http://localhost:5601/foo -> https://polyverse.io:5601/foo




## Canonical Redirection (Putting it all together.)

This is where it all comes together. Say you have a bunch of domain names (polyverse-security.org, polyverse.io, www.polyverse.io, polyverse.com, polyverse.org, etc.)
that you want to support through your website. And you want ALL these stuff to be redirected to one canonical base.


```
$ docker run -d -p 80:80 polyverse/redirect -bind=:80 -host=polyverse.io -port=80 -scheme=https
```

This will redirect ALL traffic hitting the endpoint to the polyverse.io hostname:

* http://localhost/foo -> https://polyverse.io/foo
* http://localhost:5601/foo -> https://polyverse.io/foo
* http://polyverse.io:5601/foo -> https://polyverse.io/foo


