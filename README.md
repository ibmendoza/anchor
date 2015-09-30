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

A push-based workflow of configuration management is easier to reason about. Once you have gathered your artifacts like source code, binaries, files, images, etc to your control machine, for example a local computer or VM in public cloud, those artifacts can be pushed to your remote nodet (the one being configured) through SSH.

Pushing artifacts to multiple hosts in an efficient manner is left to the user or can be delegated to another project as separate orchestration tool.


Simple error handling
---------------------

The CM tool must stop execution at the first occurrence of error.


Usage
-----

```
go get github.com/ibmendoza/anchor

anchor /path/to/cmdfile
```

License
-------

MIT
