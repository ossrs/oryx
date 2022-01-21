# srs-terraform

A control panel for SRS, to terraform the open-source video platform

## Usage

Build UI:

```bash
npm build
```

Run server with UI:

```bash
npm start
```

Access the browser: http://localhost:2022/mgmt

## Development

Run the backend:

```
npm start
```

Run the ui:

```
cd ui && npm start
```

Access the browser: http://localhost:3000

## Ports

* [SRS ports](https://github.com/ossrs/srs/blob/develop/trunk/doc/Resources.md#ports) `tcp://1935`, `tcp://1985`, `tcp://8080`, `udp://8000`, `tcp://8088`, `tcp://1990`, `udp://8935`, `tcp://554`, `tcp://8936`, `udp://10080`, `udp://1989`.
* SRS terraform mgmt port is `tcp://2022` that mount at `/mgmt`.

