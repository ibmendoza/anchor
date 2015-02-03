This specification is not yet complete (subject to change). Note that I designed this spec with Go language in mind. Just ignore Go-specific keywords if you are going to implement this other than Go.

DSL specification (syntax and rules)
------------------------------------

Here is the specification for writing a cmdfile.

#Comment

```
# single-line comment
; so is this
# sorry no double slash // (conflicts with AWS-style file reference)

# there is no block comment like in Go /* ... */
```

#Data variables section

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

LICENSE

```
LICENSE MIT

# Optional license definition
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

RUNFLAG

```
RUNFLAG
=======

Makes use of ? and * to differentiate between single/double dash flags

Case 1 (Question mark denotes single dash flag)

http://developer.dnsimple.com/certificates/#configure

curl  -H 'X-DNSimple-Token: <email>:<token>' \
      -H 'Accept: application/json' \
      -X PUT \
      -H 'Content-Type: application/json' \
      https://api.dnsimple.com/v1/domains/example.com/certificates/2/configure
	  
Here is the cmdfile:
	 
[curl]

H = 'X-DNSimple-Token: <email>:<token>'
H = 'Accept: application/json'
X = PUT
H = 'Content-Type: application/json'

[code]
	 
RUNFLAG curl https://api.dnsimple.com/v1/domains/example.com/certificates/2/configure ?curl



Case 2 (Asterisk denotes double dash flag)


Command prompt: VBoxManage modifyvm tklinux --vrdeauthtype external

Here is the cmdfile:

[vbox]

vrdeauthtype = external

[code]

RUNFLAG VBoxManage modifyvm tklinux *vbox
```

Sample Command

NOTE: You cannot use both ? and * in one line 

```
sudo docker run -d -m 100m -e DEVELOPMENT=1 \
	-e BRANCH=example-code \
	-v $(pwd):/app/bin:ro \
	--name app appserver
```

can be rewritten as cmdfile like the following.

```
[docker]

m = 100m
e = DEVELOPMENT=1
e = BRANCH=example-code
v = $(pwd):/app/bin:ro

[code]

RUNFLAG sudo docker run -d --name app appserver ?docker
```

INCLUDE

```
INCLUDE /path/to/cmdfile

# executes cmdfile (division of labor)
# note that cmdfile invoked by INCLUDE executes in its own context
# that is, it is not seen by other cmdfile, hence it is self-contained
```
