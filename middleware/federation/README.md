# federation

The *federation* middleware enables federated (https:///localhost) queries to be resolved via the
kubernetes middleware.

Enabling *federation* without also having *kubernetes* is a noop.

## Syntax

~~~
federation [ZONES...] {
    NAME DOMAIN
    fallthrough
~~~

* Each **NAME** and **DOMAIN** defines federation membership. One entry for each. Duplicate **NAME**
  will silently overwrite any previous value.
* `fallthrough` if the query is *not* a federation domain allow falling through to the next
  middleware. You probably always want fallthrough.

## Examples

Here we handle all service requests in the `prod` and `stage` federations. We need to `fallthrough`
to call into the *kubernetes* middleware.

~~~ txt
. {
    kubernetes cluster.local 
    federation cluster.local {
        fallthrough
        prod prod.feddomain.com
        stage stage.feddomain.com
    }
}
~~~
