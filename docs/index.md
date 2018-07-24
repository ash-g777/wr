# Welcome to the wr wiki!

The [README.md](http://ash-g777.viewdocs.io/wr/README) contains a basic starting guide for using `wr`, and as explained there, `wr` itself has detailed help texts for advanced usage.

A complete from-scratch walkthrough of using wr in OpenStack to carry out a software workflow can be found [here](https://ash-g777.viewdocs.io/wr/wiki/opnstk).

## Gotchas

Using `wr` with OpenStack requires that you source your openstack rc file; see `wr cloud deploy -h`.

`wr cloud deploy` has a default `--os`, but it may not be suitable for your particular installation of OpenStack. Don't forget that you can change the default by setting cloudos in your `wr` config file.

If you run in to problems, the first thing to do is check your log file. By default this will be `~/.wr_production/log` (or log.openstack for a cloud deployment, available after tearing down).

### Sanger

If you're at the Sanger Institute and want to use 'wr' with OpenStack, you'll need to use a flavor regex of:
'''
    ^m.*$
'''
You'll probably also want to use Sanger's DNS IPs, to resolve local domains.

It'll be easiest if you set these and other cloud options in your config file (~/.wr_config.yml):

cloudflavor: "^m.*$"
clouddns: "172.18.255.1,172.18.255.2,172.18.255.3"
# (the following are the defaults and don't need to be set)
cloudcidr: "192.168.0.0/18"
cloudgateway: "192.168.0.1"
cloudos: "Ubuntu Xenial"
clouduser: "ubuntu"
cloudram: 2048

If you use docker, you will have to configure it to not conflict with the Sanger's network or the network that wr will create for you. For example, the script you supply to wr cloud deploy --script might start:

sudo mkdir -p /etc/docker/
sudo bash -c "echo '{ \"bip\": \"192.168.3.3/24\", \"dns\": [\"8.8.8.8\",\"8.8.4.4\"], \"mtu\": 1380 }' > /etc/docker/daemon.json"
[further commands for installing docker]

S3 Mounts

wr add has a --mounts option that lets you mount S3 buckets prior to executing your commands. The mounts option can first be tested with wr mount. For this to work, you will need a working ~/.s3cfg or ~/.aws/credentials file: if a tool like s3cmd works for you, so will wr mount. Using --mounts on a cloud deployment will work automagically by copying over your S3 config file(s) to spawned servers.

Mounts can be done a number of different ways, and make the commands you add cleaner and simpler.

For example, instead of doing something like (on an image where s3cmd has been installed and configured):

echo 's3cmd get s3://inbucket/path/input.file && myexe -i input.file > output.file && s3cmd put output.file s3://outbucket/path/output.file' | wr add

You could (not requiring s3cmd be installed):

echo 'myexe -i inputs/input.file > outputs/output.file' | wr add --mount_json '[{"Mount":"inputs","Targets":[{"Path":"inbucket/path"}]},{"Mount":"outputs","Targets":[{"Path":"outbucket/path","Write":true}]}]'

Or even nicer:

echo 'myexe -i input.file > output.file' | wr add --mounts 'ur:inbucket/path,cw:outbucket/path'

(Note that for direct use as a working directory like this, we ought to enable caching on the writable target. Without caching we can only do serial writes and for more complicated commands things may not work as expected.)

You could have a text file with many of these short and sweet command lines, and specify the --mounts just once as an option to wr add. Performance will also be higher than using s3cmd or s3fs et al.

If an individual command will read the same data multiple times, enable per-command caching (which gets deleted once the cmd completes):

--mounts 'cr:inbucket/path'

If multiple different commands could run on the same machine and access the same data, put the cache in a fixed location (where it won't ever get deleted by wr; be careful about doing this for writable mounts!; this is also slower than than the previous scenario if you don't read whole files):

--mount_json '[{"Targets":[{"Path":"inbucket/path","CacheDir":"/tmp/mounts_cache"}]}]'

NB: Do not try and mount them at the same location: it won't work! Give them unique mount points, but the same cache location.

Unlike s3cmd, wr's mount options support "profiles", useful if you need to mount multiple buckets that have different configurations. In your ~/.s3cfg file, after the [default] section add more named sections with the necessary settings, then select the section (or "profile") to use by saying profile@bucket when specifying your bucket, where 'profile' is the name of the desired section.
Sanger

NPG have put a public bucket together containing lots of reference-related files that you might want to use. Eg. If you will run samtools to do something with cram files you might:

--mount_json '[{"Targets":[{"Path":"inbucket/path"},{"Path":"npg-repository","CacheDir":"/tmp/mounts_cache"}]}]'

And then in the JSON you supply to wr add -f say something like:

{"cmd":"samtools view ...","env":["REF_PATH=cram_cache/%2s/%2s/%s"]}

Inside the npg-repository bucket you'll also find reference indexes for use by bwa, samtools and other software. For tools like samtools that need the index file and the original fasta file in the same directory, you can take advantage of the multiplexing possible in --mounts:

--mounts 'ur:inbucket/path,cr:npg-repository/references/Homo_sapiens/GRCh38_15_noEBV/all/fasta,cr:npg-repository/references/Homo_sapiens/GRCh38_15_noEBV/all/samtools'

(Now your cmd will see Homo_sapiens.GRCh38_15_noEBV.fa and Homo_sapiens.GRCh38_15_noEBV.fa.fai in the current directory, along with your input files.)
iRODS @ Sanger

If you need to process data in OpenStack that is stored in iRODS, your best bet is probably to copy the data to S3 first, and then use S3 mounts as described above.

Because putting files in to S3 (ceph) happens at about 40MB/s from an OpenStack node but only about 20MB/s from a farm node (while reading from iRODS is a similar speed from both), you may prefer to do these copy jobs in OpenStack. That means bringing up instances with the iRODS clients installed and authentication sorted out.

The following guide assumes you have non-interactive (non-Kerberos) authentication already configured and working on the farm.

First, create a bash script with the commands neccessary to install the iRODS client, and enable resolution of Sanger's iRODS servers:

irods.script:

wget -qO - https://packages.irods.org/irods-signing-key.asc | sudo apt-key add -
echo "deb [arch=amd64] https://packages.irods.org/apt/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/renci-irods.list
sudo apt-get update
sudo apt-get install irods-icommands -y
echo "search sanger.ac.uk internal.sanger.ac.uk" | sudo tee -a /etc/resolv.conf

Now create an OpenStack-specific version of the environment file that excludes any local paths:

grep -Ev "plugins|certificate" ~/.irods/irods_environment.json > ~/.irods/irods_environment.json.openstack

One time only, we need to create an OpenStack-specific iRODS authentication file:

    wr cloud deploy --os "Ubuntu Trusty" --config_files '~/.irods/irods_environment.json.openstack:~/.irods/irods_environment.json' --script irods.script
    ssh -i ~/.wr_production/cloud_resources.openstack.key ubuntu@[ip address from step 1]
    iinit
    [enter your password and then as quickly as possible - time is important - carry out steps 5-7]
    exit
    sftp -i ~/.wr_production/cloud_resources.openstack.key ubuntu@[ip address from step 1]
    get .irods/.irodsA
    exit
    mv .irodsA ~/.irods/.irodsA.openstack
    wr cloud teardown

From now on, when we wish to do iRODS -> S3 copy jobs, we just have to be sure to copy over these irods files to the servers we create, and install the iRODS client on them (currently that means having to use Ubuntu Trusty, since that's what they have built the client for), eg.:

    wr cloud deploy --config_files '~/.irods/irods_environment.json.openstack:~/.irods/irods_environment.json,~/.irods/.irodsA.openstack:~/.irods/.irodsA,~/.s3cfg'
    echo "iget /seq/123/123.bam" | wr add --mounts 'cw:s3seq/123' --cloud_os "Ubuntu Trusty" --cloud_script irods.script

(Note that this doesn't work without caching turned on because random writes are not supported without caching.)# Testing out the viewdocs program (?)

## Contents
- [Multiple Cloud Deployments](http://ash-g777.viewdocs.io/wr/wiki/mul_cloud_deps)
- [OpenStack Walkthrough](http://ash-g777.viewdocs.io/wr/wiki/opnstk)
- [Rest API](http://ash-g777.viewdocs.io/wr/wiki/rest_api)
- [Security](http://ash-g777.viewdocs.io/wr/wiki/sec/)

