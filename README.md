### Opsgenie Heartbeat
A docker image to send heartbeat request to Opsgenie api periodically or single shot.

```bash
$ docker run --rm ghcr.io/icaliskanoglu/opsgenie-heartbeat:latest
```

The result will then look like this:

```bash
WARN[0003] Setting up periodic heartbeat                 Interval(minutes)=5
INFO[0003] Sending ping!                                 Heartbeat=IhsanTest
```

## Configuration

| Environment variable | Mandatory | Default value               | Description                                                                                                                              |
|----------------------|-----------|-----------------------------|------------------------------------------------------------------------------------------------------------------------------------------|
| `NAME`               | `true`    |                             | Name of the heartbeat                                                                                                                    |
| `BASE_URL`           | `false`   | `https://api.opsgenie.com`  | Opsgenie REST API base url. <br/>US: `https://api.opsgenie.com` <br/> EU: `https://api.eu.opsgenie.com`                                  |
| `API_KEY`            | `true`    |                             | Authentication key for Opsgenie Rest API                                                                                                 |
| `ALERT_PRIORITY`     | `false`   | `P3`                        | Specifies the alert priority for heartbeat expiration alert. If this is not provided, default priority is P3                             |
| `INTERVAL`           | `false`   | `5`                         | Specifies how often a heartbeat message should be expected.                                                                              |
| `INTERVAL_UNIT`      | `false`   | `minutes`                   | Interval specified as `minutes`, `hours` or `days`                                                                                       |
| `ENABLED`            | `false`   | `true`                      | Enable/disable heartbeat monitoring                                                                                                      |
| `PERIODIC`           | `false`   | `false`                     | The flag for to send heartbeat one time or periocically                                                                                  |
| `TEAM`               | `false`   |                             | Owner team of the heartbeat, consisting id and/or name of the owner team                                                                 |
| `DESCRIPTION`        | `false`   |                             | An optional description of the heartbeat                                                                                                 |
| `ALERT_MESSAGE`      | `false`   | `HeartbeatName is expired`  | Specifies the alert message for heartbeat expiration alert. If this is not provided, default alert message is `HeartbeatName is expired` |
| `ALERT_TAGS`         | `false`   |                             | Specifies the alert tags for heartbeat expiration alert                                                                                  |

##### Docker Compose
```yaml
version: "3"
services:
  grafana:
    image: ghcr.io/icaliskanoglu/opsgenie-heartbeat:latest
    restart: always
    environment:
      - NAME="Sample Heartbeat"
      - API_KEY="******"
```
