# Hive-Of-Things Hub

The Hub is the reference implementation of the Hub for the *Hive of Things*. It acts as an intermediary between IoT devices 'Things' and consumers using a hub-and-spokes architecture. Consumers interact with Things through Hub services without connecting directly to the Thing device. The Hub leverages [dapr](https://dapr.io/) for the infrastructure.

## Project Status

Status: The status of the Hub is In Development. It is undergoing a rewrite to use **dapr** for infrastructure.

HiveOT is based on the [W3C WoT TD 1.1 specification](https://www.w3.org/TR/wot-thing-description11/). See [docs/README-TD] for more information.

## Audience

This project is aimed at software developers and system implementors that are working on secure IoT devices. Users choose to not run servers on Things and instead use a hub and spokes model which greatly reduces the security risk post by traditional IoT devices.

## Objective

The primary objective of HiveOT is to provide a solution to use 'internet of things' devices in a highly secure manner. The Hub supports this objective by not allowing servers to run on IoT devices and by isolating IoT devices from the wider network via a secure Hub.

> The HiveOT mandate is: 'Things Do Not Run Servers'.

The secondary objective is to simplify development of IoT devices for the web of things. HiveOT supports this by requiring only minimal features to operate on an IoT device. No server is used and the Hub handles authentication and authorization on behalf of the device. This simplifies the IoT device development and allows allocating most of the resources to the actual device operation.

The third objective is to follow the WoT and other open standard where possible.

## Summary

This document describes a technical overview of the Hub.

Security is big concern with today's IoT devices. The Internet of Things contains billions of devices that when not properly secured can be hacked. Unfortunately the reality is that the security of many of these devices leaves a lot to be desired. Many devices are vulnerable to attacks and are never upgraded with security patches. This problem is only going to get worse as more IoT devices are coming to market. Imagine a botnet of a billion devices on the Internet ready for use by unscrupulous
actors.

This 'HiveOT Hub' repository provides core services to securely interact with IoT devices and consumers. This includes certificate management, authentication, authorization, provisioning, message bus service and directory service.

HiveOT compatible IoT devices (Things) therefore do not need to implement these features. This improves security as IoT devices do not run servers and are not directly accessible. They can remain isolated from the wider network and only require an outgoing connection to the Hub message bus. This in turn reduces required device resources such as memory and CPU (and cost). An additional benefit is that consumers receive a consistent user experience independent of the IoT device provider as all
interaction takes place via the Hub interface.

Note that since HiveOT Things interact via the Hub message bus, they are still vulnerable to insecurities as a result of bugs in handling those messages. The Hub message bus can somewhat mitigate this by validating the messages against their schema. (this is not currently implemented)

HiveOT is based on the 'WoT' (Web of Things) open standard developed by the W3C organization. It aims to be compatible with this standard.

The communication infrastructure for the services is provided by 'dapr'. Dapr supports http and grpc communication methods for service invocation with middleware and components for authentication, logging, resiliency, pub/sub, and more. Dapr is distributed and includes support for security, scalability and extensibility out of the box. Using dapr also enables a wide selection of storage, messaging, and integration solutions.

### Services & Adapters

All Hub functionality is provided through plugins that are services and protocol adapters.

Services and Protocol adapters can be written in any programming language but must follow some simple guidelines to integrate with the dapr infrastructure. The [writing-plugins.md] document describes how to write new plugins. Existing services/plugins can also serve as an example.

To aid application development, all API's are defined with gRPC. 'protoc', The gRPC compiler, can generate client libraries in many different programming languages.

Plugin development is further simplified by the use of dapr extensions for state management, bindings, and messaging.

# Launching

The launcher is responsible for starting and stopping services and their dapr sidecars, and for configuration of dapr based on the included templates.

The services to launch are defined in the launcher.yaml file along with their metadata such as port or unix domain socket for communication between service and sidecar.

# Core Use-cases And Services

## IoT Device Registers With Hub

Device uses idprov protocol to obtain authentication certificate.

* device => [idprov gateway] <-> [cert svc]

## IoT Device Makes A TD Document Available

Device publishes a Thing Description document. Since there are multiple consumers a pub/sub mechanism must be used.

* device => [pubsub gateway] <-> [TD recorder svc] -> [thingstore svc]
  alt:
* device => [td gateway] -> [pubsub svc] -> [TD recorder svc] -> [thingstore svc]

## IoT Device Publish Thing Properties and Events

Device publishes events with thing value changes. Since there are multiple consumers a pub/sub mechanism must be used.

* device => [pubsub gateway] <-> [Event recorder svc] -> [history svc]
  alt:
* device => [event gateway] -> [pubsub svc] -> [event recorder svc] -> [history svc]

## IoT Device Receives Actions

* device <= [sub gateway] <- pub action
  alt:
* device <= [websocket gateway] <- [sub svc] <- pub action
  alt:
* device => [action gateway]   (poll)

## Consumer Publishes Thing Action

Consumer publishes action for Things

* consumer => [pubsub gateway] -> [action recorder svc] -> [history svc]
  alt:
* consumer => [action gateway] -> [pub action] -> [action recorder svc] -> [history svc]

## Consumer Retrieves Things

* consumer <=> [thing gateway] <-> [thingstore svc]

## Consumer Retrieves Values

* consumer <=> [value gateway] <-> [historystore svc]

## Consumer Subscribes to Value Updates

* consumer <=> [pubsub gateway]

### Launcher Service

The Hub launcher is responsible for starting and stopping Hub services and their dapr sidecars.

### certs: Certificate Management

The certs service provides a commandline interface for managing certificates.

- the Hub self-signed CA certificate. Can be added to the browser for local use.
- the Hub server certificate, signed by the CA.
- the Hub plugin client certificate, signed by the CA. Intended for machine-to-machine authentication.
- the IoT device certificates, used by the 'idprov' service during the provisioning process.

### gateway: Gateway Services

Gateway services map the external Hub REST API to internal Hub micro-services using middleware for authentication. The gateway service listens on port 443. For testing, port 8443 is used.

### idprov: Provisioning Service

IoT devices that support the [idprov protocol](https://github.com/hiveot/idprov-standard) can automatically discover the provisioning server on the local network using the DNS-SD protocol and initiate the provisioning process. When accepted, a CA signed device (client) certificate is issued.

The device certificate supports machine to machine authentication between IoT device and Hub Services such as the message bus. See [idprov service](https://github.com/hiveot/hub/tree/main/idprov) for more information.

### authn: Authentication Service

The authentication service manages users and issues access and refresh tokens.
It provides a CLI to add/remove users and a service with a REST API to handle authentication request and issue tokens. See [authn service](https://github.com/hiveot/hub/tree/main/authn) for more information.

Authentication is implemented through dapr middleware. IoT devices and Hub services are agnostic to the authentication mechanism used.

### authz: Authorization Service

The authorization service manages role based access control using groups of consumers and Things.
Consumers that are in the same group as a Thing have permission to access the Thing based on their role as viewer, operator, manager, administrator or thing. See the [authorization service](https://github.com/hiveot/hub/tree/main/authz) for more information.

Authorization is implemented through dapr middleware. IoT devices are agnostic to the authorization mechanism used. The authorization service can be used by other services for fine grained authorization control.

### mosquittomgr: Message Bus Manager and Mosquitto auth plugin (deprecated)

This mosquittomgr service is replaced by dapr pub/sub.

Interaction with Things takes place via a message bus. [Exposed Things](https://www.w3.org/TR/wot-architecture/#exposed-thing-and-consumed-thing-abstractions) publish their TD document and events onto the bus and subscribe to action messages. Consumers can subscribe to these messages and publish actions to the Thing.

The Mosquitto manager configures the Mosquitto MQTT broker (server) including authentication and authorization of things, services and consumers. See the [mosquittomgr service](https://github.com/hiveot/hub/tree/main/mosquittomgr) for more information.

IoT devices must be able to connect to the message bus through TLS and use client certificate authentication. The Hub library provides protocol bindings to accomplish this.

### thingdir: Directory Service

The directory service captures TD document publications and lets consumer list and query for known Things. It uses the Authorization service to filter the TD's that a consumer is allowed to see. See the [directory service](https://github.com/hiveot/hub/tree/main/thingdir) for more information.

The directory service is intended for use by consumers. IoT devices only need to use the pub/sub API to publish TDs and events, and subscribe to actions.

## Build

This step can be skipped if you are using the pre-built binaries.

### Build From Source

To build the core and bundled plugins from source, a Linux system with golang 1.15+ and make tools must be available on the target system. 3rd party plugins are out of scope for these instructions and can require nodejs, python and golang.

Prerequisites:

1. Golang 1.15 or newer
2. GCC Make

Build from source (tentative):

```sh
$ git clone https://github.com/hiveot/hub
$ cd hub
$ make all
```

After the build is complete, the distribution binaries can be found in the 'dist/bin' folder and configuration files in dist/config.

The makefile also support a quick install for the current user:

```sh
make install
```

This copies the binaries and config to the ~/bin/hivehub location as described in the manual install section below. Executables are always replaced but only new configuration files are installed. Existing configuration remains untouched.

Additional plugins are built similarly:

```bash
$ git clone https://github.com/hiveot/{plugin}
$ cd {plugin}
$ make all 
$ make install                    (to install as user to ~/bin/hivehub/...)
```

## Installation (draft)

The HiveOT Hub is designed to run on Linux based computers. It might be able to work on other platforms but at this stage this is not tested nor a priority.

### System Requirements

The Hub can run on most small to large Intel and Arm based systems.

The minimal requirement for the Hub is 100MB of RAM and an Intel Celeron, or ARMv7 CPU. Additional resources might be required for some plugins. See plugin documentation.

Dapr for Infrastructure:

The Hub uses the Dapr runtime for infrastructure. The default configuration is the minimal slim stand-alone configuration, meaning no docker, no kubernetes, and no Redis. For large scale deployments it is recommended to configure dapr for use with docker containers and Kubernetes. This is managed by the launcher configuration through the use of configuration templates.

### Install From Package Manager

Installation from package managers is currently not available.

### Manual Install As User

The Hub can be installed and run as a dedicated user or system user. This section describes to install the Hub in a dedicated user home directory.

0. Download or build the binaries. See the build section for more info.
1. Create a user, for example a 'hivehub' user. Login as that user.
2. Create the hub folder structure

```sh
mkdir -p ~/bin/dapr
mkdir -p ~/bin/hivehub/bin
mkdir -p ~/bin/hivehub/config
mkdir -p ~/bin/hivehub/logs 
mkdir -p ~/bin/hivehub/certs 
mkdir -p ~/bin/hivehub/certstore
```

3. Install dapr as user. Do not initialize yet.

```
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | DAPR_INSTALL_DIR="$HOME/bin/dapr" /bin/bash
```

4. Copy the application binaries into the ~/bin/hivehub/bin folder and default configuration in the ~/bin/hivehub/config folder

```sh
cp bin/* ~/bin/hivehub/bin
cp config/* ~/bin/hivehub/config
```

5. Generate the certificates using the certs CLI

```sh
cd ~/bin/hivehub
bin/certs certbundle   
```

6. Run

```sh
bin/launcher start
```

If desired, this can be started using systemd. Use the init/hivehub.service file.

### Install To System (tenative)

For systemd installation to run as user 'hivehub'. When changing the user and folders make sure to edit the init/hivehub.service file accordingly. From the dist folder run:

1. Create the folders and install the files

```sh
sudo mkdir /opt/hivehub/
sudo mkdir -P /etc/hivehub/certs/ 
sudo mkdir -P /var/lib/hivehub/certstore/ 
sudo mkdir /var/log/hivehub/   

# Install HiveOT configuration and systemd
# download and extract the binaries tarfile in a temp for and copy the files:
tar -xf hivehub.tgz
sudo cp config/* /etc/hivehub
sudo vi /etc/hivehub/hub.yaml    - and edit the config, log, plugin folders
sudo cp init/hivehub.service /etc/systemd/system
sudo cp bin/* /opt/hivehub
```

2. Setup the system user and permissions

```sh
sudo adduser --system --no-create-home --home /opt/hivehub --shell /usr/sbin/nologin --group hivehub
sudo chown -R hivehub:hivehub /etc/hivehub
sudo chown -R hivehub:hivehub /var/log/hivehub
sudo chown -R hivehub:hivehub /var/lib/hivehub

sudo systemctl daemon-reload
```

3. Install dapr

User must have sudo access.

```sh
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash
```

4. Start the hub

```sh
sudo service hivehub start
```

5. Autostart the hub after startup

```sh
sudo systemctl enable hivehub
```

## Configuration

All Hub services will run out of the box with their default configuration. To change the default network and folder locations edit the 'config/hub.yaml' configuration file.

Hub services load their common configuration from the hub.yaml file in the config folder. This file MUST exist as it contains the message bus connection information for use by plugins. If no address is configured, the host outbound IP address is determined during startup. For hosts with multiple addresses, the address to use can be configured in hub.yaml

Plugins can have their own plugin specific configuration file in the config folder. Plugins must be able to run without a configuration file.

## Launching

The Hub can be launched manually by invoking the 'launcher' app in the Hub bin folder. eg:

```shell
~/bin/hivehub/bin/launcher
```

The services to start are defined in the config/launcher.yaml configuration file. When adding services, this file needs to be updated with the new service executable name.

A systemd launcher is provided that can be configured to launch on startup for systemd compatible Linux systems. See 'init/hivehub.service'

```shell
sudo cp init/hivehub.service /etc/systemd/system
sudo vi /etc/systmd/system/hivehub.service      (edit user and working directory)
sudo systemctl enable hivehub
sudo systemctl start hivehub
```

# Contributing

Contributions to HiveOT projects are always welcome. There are many areas where help is needed, especially with documentation and building plugins for IoT and other devices. See [CONTRIBUTING](docs/CONTRIBUTING.md) for guidelines.

# Credits

This project builds on the Web of Things (WoT) standardization by the W3C.org standards organization. For more information https://www.w3.org/WoT/

This project is inspired by the Mozilla Thing draft API [published here](https://iot.mozilla.org/wot/#web-thing-description). However, the Mozilla API is intended to be implemented by Things and is not intended for Things to register themselves. The HiveOT Hub will therefore deviate where necessary.

The open source [dapr](https://docs.dapr.io/) project provides the infrastructure for the Hub. After making an inventory of [infrastructure challenges](docs/infrastructure-challenges.md) from the first iteration, it was found that dapr solved most of them and then some.

Many thanks go to JetBrains for sponsoring the HiveOT open source project with development tools.  
