Sizing Guide for Application Router
===================================

<!-- toc -->

- [Idle](#idle)
- [Test Setup](#test-setup)
- [HTTP Traffic](#http-traffic)
- [Web Socket Traffic](#web-socket-traffic)
- [Memory Configuration](#memory-configuration)

<!-- tocstop -->

In this guide we provide measurements done in different application router scenarios. You can use them to approximately calculate the amount of memory that would be required by the application router. The tables contain the exact results from the measurements with Node.js v6.9.1. It is a good idea to provide higher numbers for productive usage.

All measurements are with authentication. If you have additional session content and want to count the session memory consumption please take a look at what is stored in the session - described in README's [Session Contents](../README.md#session-contents) section. You will need to add the calculated session size taking into account the number of different users and the session timeout. In our tests only the JWT token took ~4KB.

## Idle
The memory consumption for an idle application router is around 50 MB.

## Test Setup

The application router runs in a container with limited amount of memory. Swap is turned off.
The test client creates new sessions on the server with a step of 100.
No more than 100 users request the application router at a given time
(e.g. 100 sessions are initialized and become idle, then 100 more session are created and become idle ...).
The test ends when an *Out of Memory* event occurs, causing the container to be stopped.
The number of created sessions before the process ends is taken.

## HTTP Traffic

There are 2 separate test scenarios depending on what is done after a session is created:
- Scenario (1)
  - A 'Hello World' static resource is being served.
- Scenario (2)
  - A 'Hello World' static resource is being served.
  - A static resource of 84.78kb (compressed by application router to 28.36kb) is being served.
  - A backend which returns a payload of 80kb (compressed by application router to 58kb) is being called.
  - Another backend which returns a payload of 160kb (compressed by application router to 116kb) is being called.

Memory Limit | Max Sessions - Scenario (1) | Max Sessions - Scenario (2)
------------ | --------------------------- | ---------------------------
256MB        | 5 300                       | 800
512MB        | 13 300                      | 2 300
1GB          | 30 100                      | 8 400
2GB          | 65 500                      | 19 500
4GB          | 134 900                     | 46 400
8GB          | 275 500                     | 102 300

## Web Socket Traffic

There are 2 separate test scenarios depending on what is done after a session is created:
- Scenario (1)
  - A 'Hello World' static resource is being served.
  - A single 'Hello' message is sent and then received through a web socket connection.
- Scenario (2)
  - A 'Hello World' static resource is being served.
  - A backend which returns a payload of 80kb over a web socket is being called.
  - Another backend which returns a payload of 160kb over a web socket is being called.

**Note**: Web sockets require a certain amount of file handles to be available to the process - it is approximately two times the number of the sessions.
In Cloud Foundry the default value is 16384.

Memory Limit | Max Sessions - Scenario (1) | Max Sessions - Scenario (2)
------------ | --------------------------- | ---------------------------
256MB        | 600                         | 300
512MB        | 1 100                       | 500
1GB          | 3 100                       | 800
2GB          | 6 500                       | 1 400
4GB          | 13 300                      | 2 900
8GB          | 20 700                      | 6 100

**Note**: `--max-old-space-size` restricts the amount of memory used in the JavaScript heap.
Its default value is below 2GB. So in order to use the full resources that has been provided to the application,
the value of this restriction should be set to a number equal to the memory limit of the whole application.

For example, if the application memory is limited to 2GB, set the V8 heap limit like this in the `package.json`:
```
    "scripts": {
        "start": "node --max-old-space-size=2048 node_modules/@sap/approuter/approuter.js"
    }
```

## Memory Configuration

Application router process should run with at least 256MB memory. It may require more memory depending on the application.
These aspects influence memory usage:
- concurrent connections
- active sessions
- JWT token size
- backend session cookies
