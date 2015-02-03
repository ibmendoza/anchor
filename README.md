Anchor is under development...not ready for production.

Anchor is a script-based configuration management tool.

Use interface, not leaky abstraction
------------------------------------

A CM tool simply invokes the command-line interface (CLI) of infrastructure software (e.g. Bash, Serf, Consul, OpenStack, AWS, gcutil, VBoxManage, etc). 

Chef, Puppet, Ansible and Salt are leaky abstractions like ORM, forever chasing a moving target which is the CLI of programs that actually matter. 

Division of labor
-----------------

Specific tasks can be encapsulated as roles and can be included in a script

Push-based workflow
-------------------

A push-based workflow of configuration management is easier to reason about. Once you have gathered your artifacts (like source code, binaries, files, images, etc) to your local computer (the so-called control machine), those deployables can be pushed to your remote machine/host (the one being configured) through SSH.

Pushing artifacts to multiple hosts in an efficient manner is left to the user or can be relegated to another project as separate orchestration tool.

Don't reinvent the wheel
------------------------

Your gool old SSH tools like ssh-keygen, ssh-agent, scp and others are enough to implement a push-based workflow

Infrastructure as data
----------------------

A shell script is the simplest form of infrastructure as data but it has limitations. 

Dockerfile is simple and intuitive but it is limited to Docker containers only. 

Taking a cue from Dockerfile, we can design a script file but not to be executed by your favorite shell. Instead, it will be executed by our CM runtime such that it follows the principles described in this README file.

The CM runtime serves as control plane while the script file (hereinafter called as cmdfile) serves as data plane. The cmdfile is written in DSL (domain-specific language). In effect, the cmdfile serves as input data for the CM runtime. The actual instruction is then fed to the infrastructure program itself such as Bash, VirtualBox, gcutil, OpenStack etc.

As long as any software exposes a CLI, you can use that to automate a certain process. As a result, the tool remains simple. The secret is in the cmdfile and it is up to you however you like it to organize.

See https://github.com/ibmendoza/anchor/blob/master/DSL.md for DSL syntax and rules.

Infrastructure as code
----------------------

Like developer code, a CM script must be tested and stored in a version control system along with config files, just as a front-end developer stores HTML, CSS, JavaScript files and other artifacts. Large binary files like OS images must be stored elsewhere (file store or object storage)

Simple error handling
---------------------

The CM tool must stop execution at the first occurrence of error. 

Linux-only
----------

This tool is designed for Linux only although you can manage your remote Linux servers from a Windows control machine using Git Bash (included in GitHub for Windows).

Usage
-----

```
go get github.com/ibmendoza/anchor

anchor /path/to/cmdfile
```

How you can contribute
----------------------

- Write your own implementation using your favorite programming language.

- A community of cmdfile akin to playbooks and recipes.


License
-------

MIT
