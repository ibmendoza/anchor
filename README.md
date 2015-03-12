Anchor is a script-based configuration management tool that adheres to UNIX philosophy.

Anchor CM is just a tool
========================

By itself, Anchor does not do anything. The secret sauce is in the cmdfile. Anchor CM is just a script preprocessor which enables you to compute data beforehand and feed it as flags to corresponding infrastructure software.

Use interface, not leaky abstraction
------------------------------------

Anchor CM tool invokes the command-line interface (CLI) of infrastructure software (e.g. bash, serf, consul, OpenStack, AWS, gcutil, VBoxManage, etc).

Chef, Puppet, Ansible and Salt are leaky abstractions like ORM, forever chasing a moving target which is the CLI of programs that actually matter. 

Division of labor
-----------------

Specific tasks can be encapsulated as roles and can be included in a script as cmdfile

Push-based workflow
-------------------

A push-based workflow of configuration management is easier to reason about. Once you have gathered your artifacts like source code, binaries, files, images, etc to your control machine, for example a local computer or VM in public cloud, those artifacts can be pushed to your remote machine/host (the one being configured) through SSH.

Pushing artifacts to multiple hosts in an efficient manner is left to the user or can be delegated to another project as separate orchestration tool.

Don't reinvent the wheel
------------------------

Your gool old SSH tools like ssh-keygen, ssh-agent, scp and others are enough to implement a push-based workflow

Infrastructure as data
----------------------

A shell script is the simplest form of infrastructure as data but it has limitations. 

Dockerfile is simple and intuitive but it is limited to Docker containers only. 

Taking a cue from Dockerfile, we can design a script file called but not to be executed by your favorite shell. Instead, it will be executed by our CM runtime such that it follows the principles described in this README file.

The cmdfile is a virtual script since it is preprocessed by the CM runtime.

See https://github.com/ibmendoza/anchor/blob/master/DSL.md for DSL specification.

Infrastructure as code
----------------------

Like developer code, a cmdfile must be tested and stored in a version control system along with config files and others. Large binary files like OS images must be stored elsewhere (file store or object storage)

Simple error handling
---------------------

The CM tool must stop execution at the first occurrence of error.

Tool-agnostic
-------------

Any software tool or runtime can be used as long as it exposes a CLI (command line interface). There is no need to prepackage the tool into a CM module since the tool itself acts as the module. Anchor CM is simply a glue that holds and orchestrates a certain process or workflow. Hence, you can use any of your favorite shell or scripting language

Linux and Windows
-----------------

Anchor is primarily designed for Linux. Using PowerShell in Windows has not been tested. Your contributions are welcome.

Usage
-----

```
go get github.com/ibmendoza/anchor

anchor /path/to/cmdfile
```

Mutable vs Immutable Infrastructure
===================================

Immutable infrastructure as defined by Codeship (http://blog.codeship.com/immutable-deployments/) works best for servers that work at the application tier (see http://en.wikipedia.org/wiki/Multitier_architecture). However, it is not a silver bullet particularly at the data tier. Configuration management is a fine-grained solution to situations where an all-or-nothing old-server-teardown/spin-new-server workflow is not acceptable.

In short, configuration management is complementary to immutable infrastructure. You just need to know the differences  where it matters.

How to Contribute
-----------------

- Speak out your use cases. I will look into it if it's meant to be included in the Anchor kernel.

- A community of cmdfile akin to playbooks and recipes. For example, Vagrant workflow can be reduced to cmdfiles


License
-------

MIT
