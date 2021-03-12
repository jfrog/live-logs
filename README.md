# live-logs

## About the Live Logs Plugin
The Live Logs plugin allows you to get all JFrog products - JFrog Artifactory, JFrog Xray, JFrog Mission Control, JFrog Distribution, and JFrog Pipelines - logs with a single click.<br>
You also have the ability to `cat` and `tail -f` any log on any product node.<br>

**Important:** The Live Logs plugin is available to On-Prem customers and to Cloud customers who are Enterprise+ subscriptions. Admin permissions are required to run this plugin.

## Requirements

The plugin requires the following prerequisites,

- JFrog CLI 1.45.1+
- Artifactory 7.17.0+
- Mission Control 4.7.0+
- Xray 3.18.0+
- Distribution 2.7.0+
- Pipelines 1.13.0+

## Installation with JFrog CLI
Installing the latest version:
```
jfrog plugin install live-logs
```

Installing a specific version:
```
jfrog plugin install live-logs@version
```

Uninstalling a plugin
```
jfrog plugin uninstall live-logs
```

## Building from sources
Follow these steps to install and use this plugin with JFrog CLI.
1. Verify that JFrog CLI is installed on your machine by running ```jfrog```; if it is not installed, install it (see [Installing JFrog CLI](https://jfrog.com/getcli/).
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. `cd` into the root directory of the cloned project.
5. Run ```make build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.

## Note:
- Xray, Mission Control, Pipelines, and Distribution **only support admin access token authentication**, while, Artifactory supports all types of authentication. 
- The scope of the generated access token is limited to the corresponding product.
- For every product, a new dedicated entry will need to be added. For example, if you want to stream logs from 3 products, a separate entry will need to be configured for each product in the JFrog CLI (so is 3 entries).

## CLI Configuration by Product

### Configuring Artifactory
1. Add a new server using,

    ```
    jfrog config add
    ```
2. Add the Artifactory URL and authentication details.

### Configuring Mission Control
1. Add a new server using,

    ```
    jfrog config add
    ```
2. Add the Mission Control URL.
3. Generate the access token for Mission Control (see [Generating Admin Tokens](https://www.jfrog.com/confluence/display/JFROG/Access+Tokens#AccessTokens-GeneratingAdminTokens)).
4. Add the access token you generated for Mission Control.

### Configuring Xray
1. Add a new server using,

    ```
    jfrog config add
    ```
2. Add the Xray URL.
3. Generate the access token for Xray (see [Generating Admin Tokens](https://www.jfrog.com/confluence/display/JFROG/Access+Tokens#AccessTokens-GeneratingAdminTokens)).
4. Add the access token you generated for Xray.

### Configuring Distribution
1. Add a new server using,

    ```
    jfrog config add
    ```
2. Add the Distribution URL.
3. Generate the access token for Distribution (see [Generating Admin Tokens](https://www.jfrog.com/confluence/display/JFROG/Access+Tokens#AccessTokens-GeneratingAdminTokens)).
4. Add the access token you generated for Distribution as follows.

### Configuring Pipelines
1. Add a new server using,

    ```
    jfrog config add
    ```
2. Add the Pipelines URL.
3. Generate the access token for Pipelines (see [Generating Admin Tokens](https://www.jfrog.com/confluence/display/JFROG/Access+Tokens#AccessTokens-GeneratingAdminTokens)).
4. Add the access token you generated for Pipelines.

## Usage
### Commands
* help
    ```
    jfrog live-logs --help
    jfrog live-logs logs --help  
    jfrog live-logs config --help
    ```
* logs
  
    ```
    jfrog live-logs logs <product-id> <server-id> <node-id> <log-name> [Flags]
    ```
    - Arguments:
        - product-id - This is the ID of product, which can be one of the following,
            - rt - Artifactory
            - mc - Mission Control
            - xr - Xray
            - ds - Distribution
            - pl - Pipelines
        - server-id - This is the JFrog CLI platform server ID.
        - node-id - This is the selected product node ID.
        - log-name - This is the selected product log name.
    - Flags:
        - i: Open interactive menu **[Default: false]**
        - f: Show the log and keep following for changes **[Default: false]**
    - Example:
    ```
  $ jfrog live-logs logs local-arti 2368364e2c78 console.log -f | grep INFO
  2020-12-06T19:21:52.549Z [jfac ] [INFO ] [6469d8c8e2ece130] [a.s.b.AccessServerRegistrar:73] [pool-26-thread-1    ] - [ACCESS BOOTSTRAP] JFrog Access registrar finished.
  2020-12-06T19:21:52.612Z [jfac ] [INFO ] [7ccdb881f0258729] [s.r.NodeRegistryServiceImpl:68] [27.0.0.1-8040-exec-8] - Cluster join: Successfully joined jfevt@01eqtrgsxaztsq1yq0a9s60289 with node id a15e67cc9bed
  2020-12-06T19:21:52.622Z [jfevt] [INFO ] [152a442b8f87bacc] [access_join.go:58             ] [main                ] - Cluster join: Successfully joined the cluster [application]
  2020-12-06T19:21:52.624Z [jfevt] [INFO ] [152a442b8f87bacc] [access_join.go:58             ] [main                ] - Executing Router register at: localhost:8046 [application]
  ```
    ```
  $ jfrog live-logs logs -i
  Select JFrog CLI server id
  ✔ local-arti
  Select node id
  ✔ 2368364e2c78
  Select log name
  ✔ console.log
  2020-12-06T19:21:52.549Z [jfac ] [INFO ] [6469d8c8e2ece130] [a.s.b.AccessServerRegistrar:73] [pool-26-thread-1    ] - [ACCESS BOOTSTRAP] JFrog Access registrar finished.
  2020-12-06T19:21:52.612Z [jfac ] [INFO ] [7ccdb881f0258729] [s.r.NodeRegistryServiceImpl:68] [27.0.0.1-8040-exec-8] - Cluster join: Successfully joined jfevt@01eqtrgsxaztsq1yq0a9s60289 with node id a15e67cc9bed
  2020-12-06T19:21:52.622Z [jfevt] [INFO ] [152a442b8f87bacc] [access_join.go:58             ] [main                ] - Cluster join: Successfully joined the cluster [application]
  2020-12-06T19:21:52.624Z [jfevt] [INFO ] [152a442b8f87bacc] [access_join.go:58             ] [main                ] - Executing Router register at: localhost:8046 [application]
  ```

* config

  ```jfrog live-logs config <product-id> <server-id> [Flags]```
    - Arguments:
        - product-id - The ID of product, which can be one of the following,
            - rt - Artifactory
            - mc - Mission Control
            - xr - Xray
            - ds - Distribution
            - pl - Pipelines
        - server-id - The JFrog CLI platform server ID.
    - Flags:
        - i: Open interactive menu **[Default: false]**
    - Example:

```
$ jfrog live-logs config rt local-rt
{
  "logs": [
    "access-request.log",
    "artifactory-request.log",
    "event-request.log",
    "metadata-request.log"
  ],
  "nodes": [
    "node1",
    "node2"
  ]
}
```
```
$ jfrog live-logs config -i    
Select JFrog CLI product id
✔ rt
Select JFrog CLI server id
✔ local-rt
{
  "logs": [
    "access-request.log",
    "artifactory-request.log",
    "event-request.log",
    "metadata-request.log"
  ],
  "nodes": [
    "node1",
    "node2"
  ]
}
```
## Using JFrog CLI
If you use an argument incorrectly, the CLI will suggest the correct value.
<br>For example:
```
$ jfrog live-logs logs local-artii 2368364e2c78 console.log
[Error] server id not found [local-artii], consider using one of the following server id values [remote-arti,local-arti]
```
## Release Notes
The release notes are available [here](RELEASE.md).

## Additional Information
Live Logs are viewable through the JFrog Platform, and are configured through the system.yaml file. For more information, see [Enabling the Live Log Feature](https://www.jfrog.com/confluence/display/JFROG/Logging+-+Temp#Logging-EnablingtheLiveLogFeature).

In addition, there are two APIs used with this plugin; for more information, see [Live Logs Plugin](https://www.jfrog.com/confluence/display/JFROG/Live+Logs#LiveLogs-EnablingtheLiveLogsPlugin).

## Contributions
A big THANK YOU to the developers for coming up with this idea and building the core!

[omerkay](https://github.com/omerkay)

[hanoch-jfrog](https://github.com/hanoch-jfrog)