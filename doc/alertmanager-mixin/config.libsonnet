{
  _config+:: {
    local c = self,
    // alertmanagerSelector is inserted as part of the label selector in
    // PromQL queries to identify metrics collected from Alertmanager
    // servers.
    alertmanagerSelector: 'job="alertmanager"',

    // alertmanagerClusterLabels is a string with comma-separated
    // labels that are common labels of instances belonging to the
    // same Alertmanager cluster. Include not only enough labels to
    // identify cluster members, but also all common labels you want
    // to keep for resulting cluster-level alerts.
    alertmanagerClusterLabels: 'job',

    // alertmanagerQuerySelector is a matcher selector built from 'alertmanagerClusterLabels'.
    // It is used to filter queries on grafana dashboards
    alertmanagerQuerySelector: std.join(',', ['%s="$%s"' % [label, label] for label in std.split(c.alertmanagerClusterLabels, ',')]),

    // alertmanagerNameLabels is a string with comma-separated
    // labels used to identify different alertmanagers within the same
    // alertmanager HA cluster.
    // If you run Alertmanager on Kubernetes with the Prometheus
    // Operator, you can make use of the configured target labels for
    // nicer naming:
    // alertmanagerName: 'namespace,pod'
    alertmanagerNameLabels: 'instance',

    // alertmanagerName is an identifier for alerts that is built from 'alertmanagerNameLabels'
    alertmanagerName: std.join('/', ['{{$labels.%s}}' % [label] for label in std.split(c.alertmanagerNameLabels, ',')]),

    // alertmanagerNameDashboards is used to identify alertmanagers on panels legends,
    // and is built from the 'alertmanagerNameLabels'
    alertmanagerNameDashboards: std.join('/', ['{{%s}}' % [label] for label in std.split(c.alertmanagerNameLabels, ',')]),

    // alertmanagerClusterName is inserted into annotations to name an
    // Alertmanager cluster. All labels used here must also be present
    // in alertmanagerClusterLabels above.
    alertmanagerClusterName: '{{$labels.job}}',

    // alertmanagerCriticalIntegrationsRegEx is matched against the
    // value of the `integration` label to determine if the
    // AlertmanagerClusterFailedToSendAlerts is critical or merely a
    // warning. This can be used to avoid paging about a failed
    // integration that is itself not used for critical alerts.
    // Example: @'pagerduty|webhook'
    alertmanagerCriticalIntegrationsRegEx: @'.*',

    dashboardNamePrefix: 'Alertmanager / ',
    dashboardTags: ['alertmanager-mixin'],
  },
}
