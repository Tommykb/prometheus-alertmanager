# Alertmanager

Prometheus Alertmanager

**WARNING: The Alertmanager is still very experimental and early in its
development and design. More than any other Prometheus component, it will still
undergo frequent breaking changes, including ones that will affect its
architecture as a whole. While we do plan on making it mature and stable
eventually, use it at your own risk for now.**

Alertmanager receives alerts generated by Prometheus and takes care of the
following aspects:

* manual silencing of specific alerts
* inhibiting alerts based on alert dependencies
* aggregating alerts by labelset
* handling notification repeats
* sending alert notifications via external services:
  * email
  * generic web hook
  * [Amazon SNS](http://aws.amazon.com/sns/)
  * [Flowdock](https://www.flowdock.com/)
  * [HipChat](http://www.hipchat.com/)
  * [OpsGenie](https://www.opsgenie.com/)
  * [PagerDuty](http://www.pagerduty.com/)
  * [Pushover](https://www.pushover.net/)
  * [Slack](http://www.slack.com/)

Note: Amazon SNS notifications depend on your Amazon credentials being configured as described [here](https://github.com/aws/aws-sdk-go).

See [config/fixtures/sample.conf.input](config/fixtures/sample.conf.input) for
an example config. The full configuration schema including a documentation for
all possible options can be found in
[config/config.proto](config/config.proto). Alertmanager automatically reloads
the configuration when it changes, so restarts are not required for
configuration updates.

## Building and running

    make
    ./alertmanager -config.file=/path/to/alertmanager.conf

## Configuring Prometheus to send alerts

To make Prometheus send alerts to your Alertmanager, set the `alertmanager.url`
command-line flag on the server:

    ./prometheus -alertmanager.url=http://<alertmanager-host>:<port> <...other flags...>

Prometheus only pushes firing alerts to Alertmanager. Alertmanager expects to
receive regular pushes of firing alerts from Prometheus. Alerts which are not
refreshed for a period of `-alerts.min-refresh-period` (5 minutes by
default) are expired.

Alertmanager only shows alerts which are currently firing and pushed to
Alertmanager.

## Running tests

[![Build Status](https://travis-ci.org/prometheus/alertmanager.svg?branch=master)](https://travis-ci.org/prometheus/alertmanager)

    make test

## Caveats and roadmap

Alertmanager is still in an experimental state. Some of the known caveats which
are going to be addressed in the future:

* Alertmanager is run as a single instance and does not provide high
  availability yet. We plan on clustering multiple replicated Alertmanager
  instances to ensure reliability in the future.
* Relatedly, silence information is currently only persisted locally in a file
  and lost if you lose the machine your Alertmanager is running on.
* Alert aggregation needs to become more flexible. Currently alerts are
  aggregated based on their full labelsets. In the future, we want to allow
  grouping alerts based on a subset thereof (for example, grouping all alerts
  with one alert name and from the same job together).
* For alert dependencies, we want to support time delays: if alert A inhibits
  alert B due to a dependency and B begins firing before A, wait for a
  configurable amount of time for A to start firing as well before sending
  notifications for B. This is not yet supported.
* Alertmanager has not been tested or optimized for high alert loads yet.

## Using Docker

You can deploy the Alertmanager using the [prom/alertmanager](https://registry.hub.docker.com/u/prom/alertmanager/) Docker image.
Do contribute on this project.
For example:

```bash
docker pull prom/alertmanager

docker run -d -p 9093:9093 -v $PWD/alertmanager.conf:/alertmanager.conf \
        prom/alertmanager -config.file=/alertmanager.conf
```
