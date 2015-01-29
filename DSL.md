DSL specification (syntax and rules)
------------------------------------

Here is the specification for writing a cmdfile.

#Comment

```
# single-line comment

# there is no block comment like in Go /* ... */
```

#Data variables

```
# data section must be put above [code] section
# sample VirtualBox key/value flags
# uses Windows ini-style config

[system]

description = "Turnkey Linux 13.0 x64"

# vboxmanage list ostypes
ostype = Debian_64

# memory is in MB
memory = 256
ioapic = on
cpus = 4
pae = on
hwvirtex = on 
nestedpaging = on 
vtxvpid = on 
largepages = on

[network]

nic1 = bridged
nic2 = hostonly
```

#Pre section 

```
[pre]

# variables to be computed beforehand must be put here
# Examples include reading environment variables

# save result to variable prefixed with $
# e.g. output of consul keygen etc

# functions are either built-in functions or
# any CLI-based programs

# read file contents if path to a file

# example from AWS
# aws ec2 authorize-security-group-ingress --group-name MySecurityGroup \
# --ip-permissions file://ip_perms.json

$ip_perms = /path/to/file	
$keygen = consul keygen
```

#Code section

```

# Must be put as last section of cmdfile
[code]

FROM <>
MAINTAINER <>
RUN ...

```

#Keywords

FROM

```
FROM <OS name>

# Optional
# Specify the OS of physical/virtual machine
```

MAINTAINER

```
MAINTAINER <full name> or <organization>

# Optional
# Name of cmdfile writer
```

ENV

```
ENV variablename value

# Set environment variable with value
```

GO

```
# Go-specific commands

GO chdir directoryname

GO getenv environmentvariable

GO hostname
```

ANKO

```
# Go-specific since it has no built-in scripting

ANKO /path/to/ankoscriptfile
```

RUN

```script
RUN VBoxManage modifyvm tklinux -memory 512MB
RUN echo "Hello world" > /home/go/echoFile.txt
RUN echo "Hello world2" >> /home/go/echoFile.txt
RUN echo "Hello world ..."
RUN sed -i.bak s/^\(VAR5=\).*/\1VALUE10/ file.cfg
RUN chmod +x cmd.sh
```

This specification is not yet complete so it is subject to change...
