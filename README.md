# Kong Dynamic Upstream
## Summary
This is a Kong golang plugin using kong-pdk which allows setting up dynamic upstreams.
Usecase stems from the fact that you may have hundreds of endpoints which need to be addresable.

## Concept
The idea of dynamic upstream is to have pattern based routing

eg: 

here expression /bar/$upstream/
can lead to routings like 
1. https://foo.com/bar/baz/status  can route to http://baz:1234/status
2. https://foo.com/bar/qux/status  can route to http://qux:1234/status
3. https://foo.com/bar/quux/  can route to http://quxx:1234/

Here `/bar` represents a prefix which routes to `$upstream:1234`


## Config Structure

```
"config": { "upstreams": "{\"rmq\":{\"port\":\"1234\",\"expression\":\"/prefix/$upstream/\"}

```

1. Note that the `upstreams` is a single config value which needs to be set. This is however inturn an encoded json so needs to be escaped in config.
2. each dynamic upstream has a key by name.
3. port is the upstream port.
4. Expression is the pattern that needs to be extracted fron the incoming URL


## Example Config
Here is an example of having a all HAProxy status pages http://hostname:1936/stats and all RabbitMQ admin pages can be routed through a Kong using 2 ACLs

1. Setup a few Services & Routes

Setup a service `hap` pointing to a non-existing upstream or a default upstream such as a 404 upstream etc.

```
 curl -i -X POST \
   --url http://localhost:8001/services/ \
   --data 'name=hap' \
   --data 'url=http://hap-default'

```


Setup an associated route. Note the path is a prefix for our routing.

```

curl -i -X POST \
   --url http://localhost:8001/services/hap/routes \
   --data 'name=hap' \
   --data 'paths[]=/ha/*'
```


Setup a service `rmq` pointing to a non-existing upstream or a default upstream such as a 404 upstream etc.

```
curl -i -X POST \
    --url http://localhost:8001/services/ \
    --data 'name=rmq' \
    --data 'url=http://rmq-default'
```

2. Setup Dynamic Route Plugin


```
 curl -i -X POST \
      --url http://localhost:8001/plugins/ \
      --header 'content-type: application/json' \
      --data '{"name": "dynamicupstream", "config": { "upstreams": "{\"rmq\":{\"port\":\"15672\",\"expression\":\"/rmq/$upstream/\"},\"hap\":{\"port\":\"1936\",\"expression\":\"/ha/$upstream/stats/\"}}"}}'

```
