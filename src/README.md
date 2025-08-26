# SAMM Grafana AWS Datasource

This plugin communicates with AWS API to collect information about "Workspaces" and "Appstream" services.

![Add Data Source](https://raw.githubusercontent.com/samanamonitor/samm-grafanaaws-datasource/master/src/img/config.png)

## Configuration

### Region
To configure the plugin, the only mandatory paramenter is the Region. Each region of AWS requires a different datasource.

### Authentication
For authentication, the administrator has two options:
* Have an AWS Access Key, a Access Secret and an optional Access Token to be able to communicate with AWS.
* If Grafana is installed in an AWS ec2 instance with a role that allows management of Workspaces and Appstream, the authentication values can be left empty.

### Retries
The additional settings configure the behavior of the AWS SDK libraries when communicating with AWS API. There are situations that require the plugin to retry the requests to AWS in case of errors like "Throttling".

### Cache
Some environments are very large and the queries may take a long time to load or they can even fail while downloading. For this reason, the plugin has a Cache implemented internally that will keep the objects for a configurable amount of time. Also, if there are communication issues that cause an interruption of the load of the objects, the plugin will try to resume the load from where it failed instead of starting from the beggining, as long as the cache has not expired.
