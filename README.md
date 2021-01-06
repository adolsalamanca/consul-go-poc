## Leader Election using Hashicorp Consul KV

<p>
This is a simple probe of concept that uses Consul KV to make leader election between several instances of a service that exposes a simple API.
</p>
<p>
All services APIs would be reachable by its hostname but we would like some of its tasks to be performed by just one of the instances, here is where Consul shines. 
</p>

## What is Consul?

<p>
Consul is a distributed, highly-available, and multi-datacenter aware tool for service discovery, configuration, and orchestration.
Some of its functionality may sound familiar if you are already using or have heard about etcd. 
</p>

<p>
Consul enables rapid deployment, configuration, and maintenance of service-oriented architectures at massive scale. For more information, please see references section.
</p>

## Requirements
<p>
This POC assumes that the user has some tools already installed in its computer:
</p>

* Docker version 19.03.13
* Go version go1.15.5 darwin/amd64

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

You will see that only of the instances is the leader, it is using consul KV store to perform the assignment and to rotate it.
You can try to stop the Leader and one of secondary services will take the leadership now.



## Author

* **Adolfo Rodriguez** - *consul-go-poc* - [adolsalamanca](https://github.com/adolsalamanca)


## References

* [Consul documentation](https://duckduckgo.com)
* [Consul on Github](https://github.com/hashicorp/consul)
* [Consul docker images](https://hub.docker.com/_/consul)
* [Leader election inspiration](https://clivern.com/leader-election-with-consul-and-golang/)
* [Echo](https://echo.labstack.com/)


## License

This project is licensed under MIT License.




