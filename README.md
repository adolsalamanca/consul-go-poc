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
Some of its functionality may sound familiar if you are already using or have heard about etcd. 
</p>

<p>
Consul enables rapid deployment, configuration, and maintenance of service-oriented architectures at massive scale. For more information, please see:
</p>

## Requirements






## Running the POC

First of all, you need to set up your consul cluster, simply run:
```
docker-compose -f docker-compose.yml up -d
```

After that, you can already try to use the Application, it is only needed to set up two different instances to see the magic.
Open first terminal and type:
```
go run main.go -port=8080
```

Now open another one:
```
go run main.go -port=3001
```

You will see that only of the instances is the leader, it is using consul KV store to perform this assignment and also to rotate it.



## References
* [Consul documentation](https://duckduckgo.com)
* [Consul on Github](https://github.com/hashicorp/consul)
* [Consul docker images](https://hub.docker.com/_/consul)
* [Leader election inspiration](https://clivern.com/leader-election-with-consul-and-golang/)
* [Echo](https://echo.labstack.com/)




