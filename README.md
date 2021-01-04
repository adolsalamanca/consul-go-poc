## Consul used for leader election using Go

<p>
This is a simple probe of concept that uses Consul KV to make leader election between several instances of a service that exposes a simple API.
</p>
<p>
All services APIs would be reachable by its hostname but sometimes we have tasks that would need to be performed by just one of the services, this is why consul is useful here.
</p>

## What is consul?

<p>
Consul is a distributed, highly-available, and multi-datacenter aware tool for service discovery, configuration, and orchestration. 
</p>

<p>
Consul enables rapid deployment, configuration, and maintenance of service-oriented architectures at massive scale. For more information, please see:
</p>


* [Consul documentation](https://duckduckgo.com)
* [Consul on Github](https://github.com/hashicorp/consul)


## Running the POC

#### How to start consul agents
```
docker run -d -p 8500:8500 --name=dev-consul -e CONSUL_BIND_INTERFACE=eth0 consul
docker run -d -e CONSUL_BIND_INTERFACE=eth0 consul agent -dev -join=172.17.0.2
docker run -d -e CONSUL_BIND_INTERFACE=eth0 consul agent -dev -join=172.17.0.2
```

## References

*  [Consul docker images](https://hub.docker.com/_/consul)
*  [Leader election inspiration](https://clivern.com/leader-election-with-consul-and-golang/)


## Modules
* Echo
* Consul




