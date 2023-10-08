## Explanation
As we talked about code generation from swagger and grpc I've chosen to use grpc with gateway (so it exposes the required rest api call).
The idea is to make every dependency replaceable so the business operations do not know which kind of storage they use or that
an event is sent when an account is updated. I added this in a custom middleware, another way to do it could be as a grpc interceptor but would allow
less control (or if this operation is called in some other way in the future the event wouldn't be sent).

The event itself it's also a proto message, using nats encoded client. This encoded client is useful to avoid parsing but limiting, as I didn't see
a way to ack or nack an event, so I used the normal client on the receiver. The event uses protobuf "oneof" which is a way to have different kind
of types (kind of generics) in a proto message, so we could add extra details for different audit logs (account created, deleted, etc...)

The services can be built natively or with docker via the Makefiles or with docker compose.
It will require to provide valid configs (the `config_dev.json` in each service are valid configs) as `config.json`. If not using compose
the service names un the configs uris (postgres, nats), should be changed to `localhost` or whichever remote host.

To run everything with docker compose also create a valid `.env`, the `.env_example` is a valid configuration that corresponds to
`config_dev.json` in the services. The run with `docker compose -f docker-compose.nats.yaml -f docker-compose.services.yaml up`

Finally, I've tried to complete all the bonus points, but couldn't do the jaeger part as I'm not too familiar tracing over events, so
I couldn't find a simple way to send trace details over nats.

## Introduction

This test aims at demonstrating your capacity to develop a service from scratch in Go, using both synchronous and asynchronous communication. We expect from your code to be:
* Simple and concise
* Smart
* Properly documented

Also, your technical choices, if any, should be described. To go through this exercise, we expect you to be familiar with SQL, messaging patterns such as consumer queues and sequencing, and replication.

As a global advice, do not write complex code, keep things as simple as possible and try to stick to the basics.

## Assignment

As your system is getting bigger and bigger, you now need to integrate [audit logs](https://www.dnsstuff.com/what-is-audit-log) to keep track of who is accessing/mutating which resource at a given time. To do so, you'll build a placeholder API which should emit an event to a message broker which will then be consumed by an audit logs worker to be eventually stored into a database. Please find below an architecture diagram to illustrate the communication flows:

![Diagram](assets/diagram.png)

Note that we will not focus on tracking mutations themselves, such as storing previous and current state of mutated resources. Let's keep things simple.

## Walkthrough

### Useful commands

Multiple `docker-compose` files have already been written for you. To make things easier, there is no network definitions so you can call the services on localhost. You might or not want to use them. If you prefer to use a different message broker from NSQ or NATS, it's up to you. But we expect another `docker-compose` file to be added to your repository. Only PostgreSQL is a requirement.

* Start NATS + PostgreSQL: `docker-compose -f docker-compose.nats.yaml up -d`
* Start NSQ + PostgreSQL: `docker-compose -f docker-compose.nsq.yaml up -d`
* Start Jaeger: `docker-compose -f docker-compose.jaeger.yaml up -d`

### Structure your workspace

Start with setting up a proper structure for your workspace. We expect files to be organized in a coherent way and to be easy to locate. Also, both the placeholder API and the audit logs consumer should stand as 2 different go modules.

### Placeholder API

Write the placeholder api. You can give it any name you want and use any libraries or not. In the diagram, we assessed that the original call made by the client was a mutation intent on an accounts with id `id`. You can definitely change that, it's only an example. This API should consist of:
* A router
* A single route that matches a given pattern
* A producer to emit to the message broker

At this time, it means that you should already have in mind the data model of your audit logs to make sure that you can properly answer the following questions:
* Who ?
* What ?
* At which time ?

### Audit logs consumer

Write the code for the consumer. It should listen to a given topic/channel and store the audit logs to a database. Your consumer should be:
* Horizontally scalable
* Failure tolerant
* Self healing

We expect from you to come up with a data model for storing the audit logs and write SQL migrations accordingly. Migrations should either be run independantly from your service, or at it's startup time.

### Containerize your applications

Write the dockerfiles to containerize both your systems.

### Running the system

Once done, we expect you to update this section and describe how to start the system.

## Bonus points (2 max.)

* Integrate build/run steps inside of the `docker-compose.yaml` file, or a dedicated one.
* Integrate structured logging, using a third party library, with basic settings such as verbosity and format.
* Write integration tests for your repository layer.
* [Only if it was all too easy for you and you got plenty of time] Use the [Jaeger](https://www.jaegertracing.io/) deployment to trace your code end-to-end.

## Requirements

* We expect you to fork this repository before starting working on it.
* Your code should be, at least partially, unit tested. We do not expect every single line of code to be but would like to evaluate your knowledge in testing.

Happy coding 🍀