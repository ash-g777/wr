# REST-API

Normally you use the `wr` sub-commands on the command line or status web interface to interact with wr's job queue.
Alternatively you could write your own clients in go, using the jobqueue API directly. This would give you the maximum performance and the full feature set.

For clients written in other languages, however, you can use wr's REST API, sending and receiving JSON over HTTPS. The REST API is quite limited at the moment. If you find you need more functionality, please submit an issue.

## Quick Start

Add jobs to the queue with some defaults for the properties:

```shell
$ wr manager start
info: wr manager started on localhost:11301, pid 9046
info: wr's web interface can be reached at https://localhost:11302/?token=Mbg3vfN0xo-BrabfkdX-3u_s4X4BoKwXNQOGyzmiCzM

$ https_proxy='' curl --cacert ~/.wr_production/ca.pem -H "Content-Type: application/json" -H "Authorization: Bearer Mbg3vfN0xo-BrabfkdX-3u_s4X4BoKwXNQOGyzmiCzM" -X POST -d '[{"cmd":"sleep 5 && echo mymsg && false","memory":"5M"},{"cmd":"sleep 5","cpus":1}]' 'https://localhost:11302/rest/v1/jobs/?rep_grp=myid&cpus=2&memory=3G&time=5s'

[{"Key":"58cef10e7a340c3b7fa09ea304a3cb98","RepGroup":"myid","DepGroups":null,"Dependencies":null,"Cmd":"sleep 5 && echo mymsg && false","State":"ready","Cwd":"","CwdBase":"/tmp","HomeChanged":false,"Behaviours":"","Mounts":"","ExpectedRAM":5,"ExpectedTime":5,"RequestedDisk":0,"OtherRequests":null,"Cores":2,"PeakRAM":0,"Exited":false,"Exitcode":0,"FailReason":"","Pid":0,"Host":"","HostID":"","HostIP":"","Walltime":0,"CPUtime":0,"Started":-62135596800,"Ended":-62135596800,"StdErr":"","StdOut":"","Env":null,"Attempts":0,"Similar":0},{"Key":"ea4266bc0fc7dd493caeb2fff1963e4a","RepGroup":"myid","DepGroups":null,"Dependencies":null,"Cmd":"sleep 5","State":"ready","Cwd":"","CwdBase":"/tmp","HomeChanged":false,"Behaviours":"","Mounts":"","ExpectedRAM":3072,"ExpectedTime":5,"RequestedDisk":0,"OtherRequests":null,"Cores":1,"PeakRAM":0,"Exited":false,"Exitcode":0,"FailReason":"","Pid":0,"Host":"","HostID":"","HostIP":"","Walltime":0,"CPUtime":0,"Started":-62135596800,"Ended":-62135596800,"StdErr":"","StdOut":"","Env":null,"Attempts":0,"Similar":0}]
```

[ 5 seconds later... ]

Get the status of all incomplete jobs:

```shell
$ https_proxy='' curl --cacert ~/.wr_production/ca.pem -H "Authorization: Bearer Mbg3vfN0xo-BrabfkdX-3u_s4X4BoKwXNQOGyzmiCzM" https://localhost:11302/rest/v1/jobs/

[{"Key":"58cef10e7a340c3b7fa09ea304a3cb98","RepGroup":"myid","DepGroups":null,"Dependencies":null,"Cmd":"sleep 5 && echo mymsg && false","State":"buried","Cwd":"/wr_cwd/5/8/c/ef10e7a340c3b7fa09ea304a3cb98591997052/cwd","CwdBase":"/tmp","HomeChanged":false,"Behaviours":"","Mounts":"","ExpectedRAM":5,"ExpectedTime":5,"RequestedDisk":0,"OtherRequests":null,"Cores":2,"PeakRAM":13,"Exited":true,"Exitcode":1,"FailReason":"command exited non-zero","Pid":25927,"Host":"vr-2-2-02","HostID":"","HostIP":"172.17.27.150","Walltime":5.002266249,"CPUtime":0,"Started":1524225096,"Ended":1524225101,"StdErr":"","StdOut":"","Env":null,"Attempts":1,"Similar":0}]
```

Get the status of all jobs within a given reporting group:

```shell
$ https_proxy='' curl --cacert ~/.wr_production/ca.pem -H "Authorization: Bearer Mbg3vfN0xo-BrabfkdX-3u_s4X4BoKwXNQOGyzmiCzM" 'https://localhost:11302/rest/v1/jobs/myid?std=true'

[{"Key":"58cef10e7a340c3b7fa09ea304a3cb98","RepGroup":"myid","DepGroups":null,"Dependencies":null,"Cmd":"sleep 5 && echo mymsg && false","State":"buried","Cwd":"/wr_cwd/5/8/c/ef10e7a340c3b7fa09ea304a3cb98591997052/cwd","CwdBase":"/tmp","HomeChanged":false,"Behaviours":"","Mounts":"","ExpectedRAM":5,"ExpectedTime":5,"RequestedDisk":0,"OtherRequests":null,"Cores":2,"PeakRAM":13,"Exited":true,"Exitcode":1,"FailReason":"command exited non-zero","Pid":25927,"Host":"vr-2-2-02","HostID":"","HostIP":"172.17.27.150","Walltime":5.002266249,"CPUtime":0,"Started":1524225096,"Ended":1524225101,"StdErr":"","StdOut":"mymsg","Env":null,"Attempts":1,"Similar":0},{"Key":"ea4266bc0fc7dd493caeb2fff1963e4a","RepGroup":"myid","DepGroups":null,"Dependencies":null,"Cmd":"sleep 5","State":"complete","Cwd":"/wr_cwd/e/a/4/266bc0fc7dd493caeb2fff1963e4a234007716/cwd","CwdBase":"/tmp","HomeChanged":false,"Behaviours":"","Mounts":"","ExpectedRAM":3072,"ExpectedTime":5,"RequestedDisk":0,"OtherRequests":null,"Cores":1,"PeakRAM":13,"Exited":true,"Exitcode":0,"FailReason":"","Pid":25940,"Host":"vr-2-2-02","HostID":"","HostIP":"172.17.27.150","Walltime":5.001559525,"CPUtime":0,"Started":1524225096,"Ended":1524225101,"StdErr":"","StdOut":"","Env":null,"Attempts":1,"Similar":0}]
```

Get the status of particular jobs by their Keys:

```shell
$ https_proxy='' curl --cacert ~/.wr_production/ca.pem -H "Authorization: Bearer Mbg3vfN0xo-BrabfkdX-3u_s4X4BoKwXNQOGyzmiCzM" 'https://localhost:11302/rest/v1/jobs/58cef10e7a340c3b7fa09ea304a3cb98,ea4266bc0fc7dd493caeb2fff1963e4a'

[{"Key":"58cef10e7a340c3b7fa09ea304a3cb98","RepGroup":"myid","DepGroups":null,"Dependencies":null,"Cmd":"sleep 5 && echo mymsg && false","State":"buried","Cwd":"/wr_cwd/5/8/c/ef10e7a340c3b7fa09ea304a3cb98591997052/cwd","CwdBase":"/tmp","HomeChanged":false,"Behaviours":"","Mounts":"","ExpectedRAM":5,"ExpectedTime":5,"RequestedDisk":0,"OtherRequests":null,"Cores":2,"PeakRAM":13,"Exited":true,"Exitcode":1,"FailReason":"command exited non-zero","Pid":25927,"Host":"vr-2-2-02","HostID":"","HostIP":"172.17.27.150","Walltime":5.002266249,"CPUtime":0,"Started":1524225096,"Ended":1524225101,"StdErr":"","StdOut":"","Env":null,"Attempts":1,"Similar":0},{"Key":"ea4266bc0fc7dd493caeb2fff1963e4a","RepGroup":"myid","DepGroups":null,"Dependencies":null,"Cmd":"sleep 5","State":"complete","Cwd":"/wr_cwd/e/a/4/266bc0fc7dd493caeb2fff1963e4a234007716/cwd","CwdBase":"/tmp","HomeChanged":false,"Behaviours":"","Mounts":"","ExpectedRAM":3072,"ExpectedTime":5,"RequestedDisk":0,"OtherRequests":null,"Cores":1,"PeakRAM":13,"Exited":true,"Exitcode":0,"FailReason":"","Pid":25940,"Host":"vr-2-2-02","HostID":"","HostIP":"172.17.27.150","Walltime":5.001559525,"CPUtime":0,"Started":1524225096,"Ended":1524225101,"StdErr":"","StdOut":"","Env":null,"Attempts":1,"Similar":0}]
```

## Security

All requests must be via https. http queries return nothing. By default, wr generates its own self-signed certificate, as well as its own CA, and your client can gain trust in wr's certificate by adding `~/.wr_production/ca.pem` to the list of root CAs to check. To you use your own certificates, see the [security guide](http://ash-g777.viewdocs.io/wr/wiki/sec/).

To authenticate you must provide an `Authorization: Bearer` header with the token that `wr manager start` tells you about. This token is also stored, by default, in the file `~/.wr_production/client.token`.

See the Quick Start above for examples of providing both of these.

## GET /rest/v1/jobs/

Gets the status of jobs in the queue. With no identifier supplied, gets the status of all incomplete jobs in the queue. Alternatively, suffix a comma-separated list of identifiers (job rep groups or keys) to get the status of those jobs only (including complete ones).

Possible url query parameters are:

* **std=true** : retrieve any stdout/err of failed jobs
* **env=true** : retrieve any environment variables the jobs were set to run with
* **limit=[int]** : group similar jobs together to reduce the number of objects returned, showing this number of jobs per group
* **state=[delayed|ready|reserved|running|lost|buried|dependent|complete]** : only get the status of jobs in this state (ignored if using job keys for identifiers)

The server will return a JSON string describing an array of job status objects with these properties:

* Key           string
* RepGroup      string
* DepGroups     []string
* Dependencies  []string
* Cmd           string
* State         JobState
* Cwd           string
* CwdBase       string
* HomeChanged   bool
* Behaviours    string
* Mounts        string
* MonitorDocker string
* ExpectedRAM   int (Megabytes)
* ExpectedTime  float64 (seconds)
* RequestedDisk int (Gigabytes)
* Cores         int
* PeakRAM       int
* OtherRequests []string
* Exited        bool
* Exitcode      int
* FailReason    string
* Pid           int
* Host          string
* Walltime      float64
* CPUtime       float64
* Started       int64
* Ended         int64
* StdErr        string
* StdOut        string
* Attempts      uint32
* Similar       int

## POST /rest/v1/jobs/

Add jobs to the queue.

Post a JSON string that describes an array of "job" objects with these properties, one for each job you wish to add:

* **cmd=[string]** (required)
* **cwd=[string]** (defaults to /tmp)
* **cwd_matters=[boolean]**
* **change_home=[boolean]**
* **mounts=[slice of mount configs]** (as per the JSON accepted by `wr mount --mount_json`)
* **req_grp=[string]**
* **monitor_docker=[string]** (--name or --cidfile of docker container cmd will run, or ? (%3F) to monitor the first docker container to start running after cmd starts to run)
* **memory=[string]** (int followed by a unit, such as 1G for 1 Gigabyte; defaults to 1G)
* **time=[string]** (int followed by a unit, such as 1h for 1 hour; defaults to 1h)
* **cpus=[int]** (defaults to 1)
* **disk=[int]** (in Gigabytes)
* **override=[int]** (in the range 0..2)
* **priority=[int]** (in the range 0..255)
* **retries=[int]** (in the range 0..255)
* **rep_grp=[string]**
* **dep_grps=[string array]**
* **deps=[string array]**
* **on_failure=[slice of behaviour objects]**
* **on_success=[slice of behaviour objects]**
* **on_exit=[slice of behaviour objects]**
* **env=[string array]**
* **cloud_flavor=[string]**
* **cloud_os=[string]**
* **cloud_username=[string]**
* **cloud_script=[string]** (path to file on the machine where the manager is running; see the upload endpoint)
* **cloud_config_files=[string]** (comma seperated list of source:dest config file paths; you can use ~/ prefixes for dest paths that should go to the home directory; source paths must exist on the machine where the manager is running; see the upload endpoint)
* **cloud_ram=[int]** (in Megabytes; defaults to 1000)

A behaviour object has one of the these key:value pairs:

* {"cleanup":true}
* {"cleanup_all":true}
* {"run":"unix command line"}

URL query parameters define default properties of jobs, and have the same names as the properties of the "job" object above, except that "cmd" can't be supplied as a default. For dep_grps, deps and env, which normally take string array, provide a comma-separated string. mounts, on_failure, on_success and on_exit values should be supplied as url query escaped JSON strings.

## PUT /rest/v1/upload/

Upload files to the machine where the manager is running. It is not intended that you use this for many files or large files!

When you add jobs you can specify the path to a "cloud_script" or to "cloud_config_files", but the files must exist at that path on the machine where the manager is running. Before adding such a job, you can upload your files to the manager's machine first using this endpoint.

To specify a path relative to the home directory, you can use tilda. Eg:

```shell
$ https_proxy='' curl --cacert ~/.wr_production/ca.pem https://localhost:11302/rest/v1/upload/?path=~/my_cloud_script -H "Authorization: Bearer Mbg3vfN0xo-BrabfkdX-3u_s4X4BoKwXNQOGyzmiCz" --upload-file my_cloud_script

{"path":"/home/ubuntu/my_cloud_script"}

$ https_proxy='' curl --cacert ~/.wr_production/ca.pem https://localhost:11302/rest/v1/jobs/ -H "Content-Type: application/json" -H "Authorization: Bearer Mbg3vfN0xo-BrabfkdX-3u_s4X4BoKwXNQOGyzmiCz" -X POST -d '[{"cmd":"cat /tmp/file_created_by_my_cloud_script","cloud_script":"~/my_cloud_script"}]'
```

Possible url query parameters are:

* **path=[string]** : path to save the uploaded file to. Can be prefixed with tilda to specify a path relative to the home directory; otherwise should be an absolute path. If not supplied, a unique path based on a MD5 checksum of the file's contents, rooted in the configured manageruploaddir is chosen for you. 

The server will return a JSON string describing a map with the key "path" and a value of the absolute path of the uploaded file.

## GET /rest/v1/warnings/

Get any warnings produced when trying to use the scheduler. These are not explicitly tied to particular jobs, though it is typically in attempting to schedule jobs to run that warnings occur. There is currently no standard format for the warning messages: they are free-form text.

The server will return a JSON string describing an array of warning objects with these properties:

* Msg           string (unique amongst the objects in the array)
* FirstDate     int64 (date that the Msg was first sent in seconds since unix epoch)
* LastDate      int64 (date that the Msg was last sent in seconds since unix epoch)
* Count         int (number of times the message has been sent since your last GET)

In wr's web interface, messages have to be manually "dismissed" by the user or they won't go away. By contrast, the act of GETting this url will "dismiss" the messages, ie. delete them.

## GET /rest/v1/servers/

For cloud deployments that are spawning servers, get the details of any servers that can no longer be ssh'd to and thus seem dead.

The server will return a JSON string describing an array of "bad server" objects with these properties:

* ID      string
* Name    string
* IP      string
* Date    int64 (the date that the server went bad in seconds since Unix epoch)
* IsBad   bool (always true, except via the websocket)
* Problem string

If Problem is a non-blank string, then wr will never try to re-use the server and you should DELETE it (after investigating it if desired).

Otherwise, there's a possibility that the server will later come back to life, in which case a subsequent GET will no longer include a "bad server" object with that ID.

## DELETE /rest/v1/servers/?id=[badserver.ID]

Confirm that one of the servers reported by GET is dead. If it still exists, wr will try to terminate it. The id parameter is required.

(This will only act on servers that wr already thinks are bad.)