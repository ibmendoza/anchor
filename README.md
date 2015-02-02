Anchor is a Go-based implementation of my specification for a script-based configuration management (CM) tool. You may also implement this tool using your favorite programming language (Python, Ruby, Perl, Node, C, etc) as long as it follows these principles:

Use interface, not leaky abstraction
------------------------------------

A CM tool simply invokes the command-line interface (CLI) of infrastructure software (e.g. Bash, Serf, Consul, OpenStack, AWS, gcutil, VBoxManage, etc). 

Chef, Puppet, Ansible and Salt are leaky abstractions like ORM, forever chasing a moving target which is the CLI of programs that actually matter. This is the CM version of the subjective debate between ORM and SQL. A script-based configuration management tool is in the same vein as the SQL argument, just as there is jQuery versus Angular/Ember debate. Bear in mind that software design is inherently subjective.

But don't get me wrong. If you are fine with your current CM tools, then by all means use it. This is not meant to start a senseless flame war. As I have said in my blog http://itjumpstart.wordpress.com,

"Objectivity is falsifiability. Subjectivity is individuality".

Rather than reading this post as a rant against those CM tools, consider this project as an alternative approach to configuration management.

GUI is optional
---------------

A CM tool itself is primarily CLI-based. If the program to be called has an HTTP API (e.g. Proxmox), a CM tool can simply use curl to orchestrate it. A GUI is best for users who have a visual preference but it must serve a large percentage of common use cases

Separation of code and data
---------------------------

Compute data beforehand so it is clear in our script which is code and which is data

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

There is no need for modules like those in Ansible, since the CM tool is just a thin wrapper for calling the CLI of your desired program. As long as any software exposes a CLI, you can use that to automate a certain process. As a result, there is no impedance mismatch between your CM tool and infrastructure software.

See https://github.com/ibmendoza/anchor/blob/master/DSL.md for DSL syntax and rules.

Infrastructure as code
----------------------

Like developer code, a CM script must be tested and stored in a version control system along with config files, just as a front-end developer stores HTML, CSS, JavaScript files and other artifacts. Large binary files like OS images must be stored elsewhere (file store or object storage)

Simple error handling
---------------------

The CM tool must stop execution at the first occurrence of error. 

Linux-only
----------

For subjective reasons, this tool is designed for first-class operating system only (Linux) although you can manage your remote Linux servers from a Windows control machine using Git Bash (included in GitHub for Windows).

Usage
-----

```
go get github.com/ibmendoza/anchor

anchor /path/to/cmdfile
```

How you can contribute
----------------------

- Anchor is a Go-based reference implementation of this script-based configuration management specification. Contribute your ideas. If you do not like Go, write your own implementation using your favorite programming language.

- A community of cmdfile akin to playbooks and recipes.


License
-------

MIT
