Anchor is a Go-based implementation of my specification for a script-based configuration management (CM) tool. You may also implement this tool using your favorite programming language (Python, Ruby, Perl, Node, C, etc) as long as it follows this specification:

- Use interface, not leaky abstraction

A CM tool simply invokes the command-line interface (CLI) of a program or runtime (e.g. Bash, Serf, Consul, OpenStack, AWS, gcutil, VBoxManage, etc)

- GUI is optional

A CM tool itself is primarily CLI-based. If the program to be called has an HTTP API (e.g. Proxmox), a CM tool can simply use curl to orchestrate it. A GUI is best for users who have a visual preference but it must serve a large percentage of common use cases

- Separation of code and data

Compute data beforehand so it is clear in our script which is code and which is data

- Division of labor

Specific tasks can be encapsulated as roles and can be included in a script

- Push-based model

A push-based model of configuration management is easier to reason about. Once you have gathered your artifacts (like source code, binaries, files, images, etc) to your local computer (the so-called control machine), those deployables can be pushed to your remote machine (the one being configured) through SSH

- Don't reinvent the wheel

Your gool old SSH tools like ssh-keygen, ssh-agent, scp and others are enough to implement a push-based workflow

- Infrastructure as data

A CM tool does not assume the role of infrastructure programs like AWS, OpenStack, Proxmox, VirtualBox, Bash, etc. Instead, it leverages the CLI of those programs and define the orchestration script as data in the form of cmdfile. Hence, a CM tool must wrap a thin layer of control logic through a DSL (domain specific language)

- Infrastructure as code

Like developer code, a CM script must be tested and stored in a version control system


See DSL.md for DSL syntax and rules.


