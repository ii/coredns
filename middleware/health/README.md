# health

This module enables a simple health check endpoint.
By default it will listen on port 8080.

## Syntax

~~~
health [ADDRESS]
~~~

Optionally takes an address; the default is `:8080`. The health path is fixed to `/health`. It
will just return "OK" when CoreDNS is healthy, which currently mean: it is up and running.

This middleware only needs to be enabled once for the entire server???

~~~
health [ADDRESS] {
    middleware MIDDLEWARE..
}
~~~

* `middleware` enables healthchecking for **MIDDLEWARE**. For each middleware that implements
  a healthcheck this must be n

## Examples

~~~
health localhost:8091
~~~
