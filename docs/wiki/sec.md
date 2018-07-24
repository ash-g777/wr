Certificates

All of wr manager's interfaces (CLI, web interface and REST API) use TLS for security, preventing any information sent over the network from being read by others. This is important because, amongst other things, wr passes your environment variables around so that the commands you want to run will work as expected. Your set of environment variables might include passwords and other sensitive information.

To make it easy to use wr, by default, the first time you wr manager start, wr will create a CA certificate and use it to sign the server certificate that your browser, CLI and REST API need to trust. Trust in the server certificate is gained by trusting the CA certificate.

In practical terms that means that the first time you view the web interface, you will need to click through the security warning and make an exception to trust wr's CA.

Likewise, if you use the REST API you will need to add wr's generated ca.pem file to your root CAs.

The CLI commands are built to automatically use wr's generated ca.pem.
Using your own certificate

Rather than have wr generate a self-signed certificate and then have to allow security exceptions, you can use your own certificate and key. You might sign these with your company's internal CA for an internal domain, with an IP pointing to the machine where you will wr manager start.

Just set the managercertfile, managerkeyfile, managercafile and managercertdomain options in your wr config file (see the example config file for details).

If you're deploying wr to OpenStack, it is recommended that you create a certificate and key just for wr, since these will be copied into your OpenStack environment.

If you use infoblox for DNS management, wr cloud deploy --set_domain_ip can be used to update the IP for the domain your certificate is valid for to point to the new OpenStack server where wr manager is now running.

Changing IPs like this isn't necessary for wr CLI interaction, but is needed for using the web interface or REST API if you want to avoid man-in-the-middle attacks.

Alternatively, if you alter your local machine's hosts to have your domain resolve to localhost, and then forward that to the machine where the manager is running, or where wr cloud deploy has created ssh forwarding, your browser will be happy since you are accessing the web interface using the correct domain, but it doesn't matter what IP address is associated with that domain.
Sanger users

A note for people at the Sanger Institute: It is possible to get an infoblox account that has permissions to control a wild-carded internal sub-domain for your team. With this in place you can set managercafile to /usr/share/ca-certificates/sanger.ac.uk/Genome_Research_Ltd_Certificate_Authority-cert.pem, managercertdomain to something like "wrMyUserName.teamname.sanger.ac.uk" and use the --set_domain_ip option to wr manager start or wr cloud deploy successfully.

If you access the web interface on a MacBook via an SSH tunnel (using the forwarding suggested at the bottom of README.md), you can alter /private/etc/hosts (you will need sudo privileges) by appending the line 127.0.0.1 wrMyUserName.teamname.sanger.ac.uk, and then you can access https://wrMyUserName.teamname.sanger.ac.uk/ in your webbrowser with no security warnings.
Authentication

wr is used to run arbitrary command lines, and these run as the same user as was used to start the manager. To take advantage of the user security model provided by your operating system, where you only have permission to run things and read and write certain files based on your user account, wr follows a single user model. It is intended that each user start their own wr manager, with potentially many people running a manager on the same machine.

Having a trusted certificate and using TLS provides a secure channel over which to communicate, but wr also needs to stop other people with access to your machine from using your manager: you need to be authenticated before using any of the interfaces.

Because wr is a simple single-user system, there is no "account" or password to set up in wr. Instead, when you wr manager start, a unique token is generated by the manager and returned to you. It is valid for the lifetime of the manager. To authenticate, you simply supply the token.

The token is stored in a file that is only readable by your user account and used automatically by the CLI commands.

It is also presented to you when you're told how to access the web interface, forming part of the URL.

The REST API guide tells you how to supply the token when using that.

Note that anyone you give the token to, or anyone you enable to read the token file, has access to your manager (if they have network access to where it is running). If there's some kind of security breach and you need to "revoke" a token that's been shared, simply stop and restart the manager.