# live-logs

## About the Live Logs Plugin
The JFrog Platform includes an integrated Live Logs plugin, which allows customers to get the JFrog product logs (Artifactory, Xray, Distribution, and Pipelines) using the JFrog CLI Plugin. The plugin also provides the ability to `cat` and `tail -f` any log on any product node.<br>

**Note:** 
The Live Logs plugin is available to:
* Self-Hosted customers who are Enterprise | Enterprise+ subscriptions. Admin permissions are required to run this plugin. 
* Cloud customers who are Enterprise+ subscriptions. Available logs: only 'artifactory-request.log' and 'distribution-request.log'.
<br>

The log types that are supported as part of the CLI plugin for Live Logs are as follows:

Artifactory | Distribution:

    *-request.log

## Requirements

The plugin requires the following prerequisites,

- JFrog CLI 1.45.1+
- Artifactory 7.17.0+
- Xray 3.18.0+
- Distribution 2.7.0+
- Pipelines 1.13.0+

## Installation with JFrog CLI
Installing the latest version:
```
jf plugin install live-logs
```

Installing a specific version:
```
jf plugin install live-logs@version
```

Uninstalling a plugin
```
jf plugin uninstall live-logs
```

## Building from sources
Follow these steps to install and use this plugin with JFrog CLI.
1. Verify that JFrog CLI is installed on your machine by running ```jf```; if it is not installed, install it (see [Installing JFrog CLI](https://jfrog.com/getcli/).
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. `cd` into the root directory of the cloned project.
5. Run ```make build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.

## Note:
- Xray, Pipelines, and Distribution **only support admin access token authentication**, while, Artifactory supports all types of authentication. 
- The scope of the generated access token is limited to the corresponding product.
- For every product, a new dedicated entry will need to be added. For example, if you want to stream logs from 3 products, a separate entry will need to be configured for each product in the JFrog CLI (so is 3 entries).

## CLI Configuration by Product

### Configuring Artifactory
1. Add a new server using,

    ```
    jf c add
    ```
2. Add the Artifactory URL and authentication details.

### Configuring Xray
1. Add a new server using,

    ```
    jf c add
    ```
2. Add the Xray URL.
3. Generate the access token for Xray (see [Generating Admin Tokens](https://www.jfrog.com/confluence/display/JFROG/Access+Tokens#AccessTokens-GeneratingAdminTokens)).
4. Add the access token you generated for Xray.

### Configuring Distribution
1. Add a new server using,

    ```
    jf c add
    ```
2. Add the Distribution URL.
3. Generate the access token for Distribution (see [Generating Admin Tokens](https://www.jfrog.com/confluence/display/JFROG/Access+Tokens#AccessTokens-GeneratingAdminTokens)).
4. Add the access token you generated for Distribution as follows.

### Configuring Pipelines
1. Add a new server using,

    ```
    jf c add
    ```
2. Add the Pipelines URL.
3. Generate the access token for Pipelines (see [Generating Admin Tokens](https://www.jfrog.com/confluence/display/JFROG/Access+Tokens#AccessTokens-GeneratingAdminTokens)).
4. Add the access token you generated for Pipelines.

## Usage
### Commands
* help
    ```
    jf live-logs --help
    jf live-logs config --help
    jf live-logs logs --help  
    ```

* config

  ```
  jf live-logs config <product-id> <server-id> [Flags]
  ```
    - Arguments:
        - product-id - The ID of product, which can be one of the following,
            - rt - Artifactory
            - xr - Xray
            - ds - Distribution
            - pl - Pipelines
        - server-id - The JFrog CLI platform server ID.
    - Flags:
        - i: Open interactive menu **[Default: false]**
    - Example:

```
$ jf live-logs config rt local-rt
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
$ jf live-logs config -i    
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
* logs

    ```
    jf live-logs logs <product-id> <server-id> <node-id> <log-name> [Flags]
    ```
    - Arguments:
        - product-id - This is the ID of product, which can be one of the following,
            - rt - Artifactory
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
  $ jf live-logs logs rt local-arti 2368364e2c78 artifactory-service.log -f | grep INFO
    2021-03-25T04:00:00.006Z [jfrt ] [INFO ] [d76675e362ffbd6a] [.s.d.b.s.g.GarbageCollector:66] [art-exec-11         ] - Starting GC strategy 'TRASH_AND_BINARIES'
    2021-03-25T04:00:00.012Z [jfrt ] [INFO ] [d76675e362ffbd6a] [.s.d.b.s.g.GarbageCollector:68] [art-exec-11         ] - Finished GC Strategy 'TRASH_AND_BINARIES'
    2021-03-25T04:30:34.196Z [jfrt ] [INFO ] [94109ae150da76e ] [aseBundleCleanupServiceImpl:84] [art-exec-16         ] - Starting to cleanup incomplete Release Bundles
    2021-03-25T04:30:34.199Z [jfrt ] [INFO ] [94109ae150da76e ] [aseBundleCleanupServiceImpl:90] [art-exec-16         ] - Finished incomplete Release Bundles cleanup  
    ```
    ```
  $ jf live-logs logs -i
  Select JFrog CLI product id
  ✔ rt
  Select JFrog CLI server id
  ✔ local-arti
  Select node id
  ✔ 2368364e2c78
  Select log name
  ✔ artifactory-service.log
  2021-03-25T04:00:00.006Z [jfrt ] [INFO ] [d76675e362ffbd6a] [.s.d.b.s.g.GarbageCollector:66] [art-exec-11         ] - Starting GC strategy 'TRASH_AND_BINARIES'
  2021-03-25T04:00:00.012Z [jfrt ] [INFO ] [d76675e362ffbd6a] [.s.d.b.s.g.GarbageCollector:68] [art-exec-11         ] - Finished GC Strategy 'TRASH_AND_BINARIES'
  2021-03-25T04:30:34.196Z [jfrt ] [INFO ] [94109ae150da76e ] [aseBundleCleanupServiceImpl:84] [art-exec-16         ] - Starting to cleanup incomplete Release Bundles
  2021-03-25T04:30:34.199Z [jfrt ] [INFO ] [94109ae150da76e ] [aseBundleCleanupServiceImpl:90] [art-exec-16         ] - Finished incomplete Release Bundles cleanup  
  ```
  
## Using JFrog CLI
If you use an argument incorrectly, the CLI will suggest the correct value.
<br>For example:
```
$ jf live-logs logs local-artii 2368364e2c78 artifactory-service.log
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
