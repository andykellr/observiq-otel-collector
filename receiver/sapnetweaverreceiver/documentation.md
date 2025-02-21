[comment]: <> (Code generated by mdatagen. DO NOT EDIT.)

# sapnetweaverreceiver

## Metrics

These are the metrics available for this scraper.

| Name | Description | Unit | Type | Attributes |
| ---- | ----------- | ---- | ---- | ---------- |
| **sapnetweaver.host.cpu_utilization** | CPU Utilization Percentage. | % | Gauge(Int) | <ul> </ul> |
| **sapnetweaver.host.memory.virtual.overhead** | Virtualization System Memory Overhead. | bytes | Gauge(Int) | <ul> </ul> |
| **sapnetweaver.host.memory.virtual.swap** | Virtualization System Swap Memory. | bytes | Gauge(Int) | <ul> </ul> |
| **sapnetweaver.host.spool_list.used** | Host Spool List Used. |  | Sum(Int) | <ul> </ul> |
| **sapnetweaver.icm_availability** | ICM Availability (color value from alert tree). |  | Sum(Int) | <ul> <li>control_state</li> </ul> |
| **sapnetweaver.locks.enqueue.count** | Count of Enqueued Locks. | {locks} | Sum(Int) | <ul> </ul> |
| **sapnetweaver.sessions.browser.count** | The number of Browser Sessions. | {sessions} | Sum(Int) | <ul> </ul> |
| **sapnetweaver.sessions.ejb.count** | The number of EJB Sessions. | {sessions} | Sum(Int) | <ul> </ul> |
| **sapnetweaver.sessions.http.count** | The number of HTTP Sessions. | {sessions} | Sum(Int) | <ul> </ul> |
| **sapnetweaver.sessions.security.count** | The number of Security Sessions. | {sessions} | Sum(Int) | <ul> </ul> |
| **sapnetweaver.sessions.web.count** | The number of Web Sessions. | {sessions} | Sum(Int) | <ul> </ul> |
| **sapnetweaver.short_dumps.rate** | The rate of Short Dumps. | {dumps/min} | Sum(Int) | <ul> </ul> |
| **sapnetweaver.work_processes.active.count** | The number of active work processes. | {work processes} | Sum(Int) | <ul> </ul> |

**Highlighted metrics** are emitted by default. Other metrics are optional and not emitted by default.
Any metric can be enabled or disabled with the following scraper configuration:

```yaml
metrics:
  <metric_name>:
    enabled: <true|false>
```

## Resource attributes

| Name | Description | Type |
| ---- | ----------- | ---- |
| sapnetweaver.instance | The SAP Netweaver instance. | Str |
| sapnetweaver.node | The SAP Netweaver node. | Str |

## Metric attributes

| Name | Description | Values |
| ---- | ----------- | ------ |
| control_state (state) | The control state color | grey, green, yellow, red |
