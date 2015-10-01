Anchor is a script-based configuration management tool that adheres to UNIX philosophy.

Anchor CM is just a tool
========================

By itself, Anchor does not do anything. The secret sauce is in the cmdfile. Anchor CM is just a script preprocessor which enables you to process data beforehand and feed it as flags to corresponding CLI software.

Use interface, not leaky abstraction
------------------------------------

Anchor CM tool invokes the command-line interface of the remote node (physical or virtual machine).

Chef, Puppet, Ansible and Salt are leaky abstractions like ORM, forever chasing a moving target which is the CLI of programs that actually matter. 

Push-based workflow
-------------------

A push-based workflow of configuration management is easier to reason about. Once you have gathered your artifacts like source code, binaries, files, images, etc to your control machine, for example a local computer or VM in public cloud, those artifacts can be pushed to your remote node (the one being configured) through SSH.

Pushing artifacts to multiple hosts in an efficient manner is left to the user or can be delegated to another project as separate orchestration tool.


Simple error handling
---------------------

The CM tool must stop execution at the first occurrence of error.


Anchor on server-side, easyssh on client-side
---------------------------------------------

Anchor is used to configure remote servers, [easyssh](https://github.com/ibmendoza/anchor/blob/master/easyssh/main.go) is used to scp files to remote server and run ssh command on the client side


Usage
-----

Assuming you have already uploaded anchor executable to remote servers as well as the corresponding cmdfile, you may also use Anchor on the client side to batch execute commands on the remote servers.

#### uptime will be executed on remote servers
```
#cmd.txt

[code]
RUN easyssh -user root -server 192.168.56.101 -keypath id_rsa -cmd uptime
RUN easyssh -user root -server 192.168.56.102 -keypath id_rsa -cmd uptime
```

#### On client side (Windows or Linux)

```
C:\mygo\src\github.com\ibmendoza\easyssh\example>anchor cmd.txt
```

#### Output

```
C:\mygo\src\github.com\ibmendoza\easyssh\example>anchor cmd.txt
RUN
easyssh -user root -server 192.168.56.101 -keypath id_rsa -cmd uptime
==> OUTPUT:  03:39:03 up 4 min,  0 users,  load average: 0.01, 0.05, 0.04


RUN
easyssh -user root -server 192.168.56.102 -keypath id_rsa -cmd uptime
==> OUTPUT:  03:39:03 up 3 min,  0 users,  load average: 0.01, 0.03, 0.02



If any error appears, cmdfile run is not completed
```

License
-------

MIT
