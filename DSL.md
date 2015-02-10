DSL specification (v1.0.0)
-----------------

#Cmdfile

```script
[testvariable]
ipaddr = RUN ifconfig eth0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}'

[code]

# http://www.cyberciti.biz/faq/how-to-find-out-the-ip-address-assigned-to-eth0-and-display-ip-only/
# RUN ifconfig eth0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}'

# double-dash flag
RUNFLAG echo @testvariable

# single-dash flag
RUNFLAG echo %testvariable

# OUTPUT:

# RUNFLAG
# echo @testvariable
# echo --ipaddr 192.168.1.102
# ==> OUTPUT: --ipaddr 192.168.1.102

# RUNFLAG
# echo %testvariable
# echo -ipaddr 192.168.1.102
# ==> OUTPUT: -ipaddr 192.168.1.102
```

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

Makes use of % and @ to differentiate between single/double dash flags

Case 1 (% denotes single dash flag)

http://developer.dnsimple.com/certificates/#configure

Example CLI:

curl  -H 'X-DNSimple-Token: <email>:<token>' \
      -H 'Accept: application/json' \
      -X PUT \
      -H 'Content-Type: application/json' \
      https://api.dnsimple.com/v1/domains/example.com/certificates/2/configure
	  

Corresponding cmdfile:
	 
[curl]

H = 'X-DNSimple-Token: <email>:<token>'
H = 'Accept: application/json'
X = PUT
H = 'Content-Type: application/json'

[code]
	 
RUNFLAG curl https://api.dnsimple.com/v1/domains/example.com/certificates/2/configure %curl



Case 2 (@ denotes double dash flag)


Example CLI:

VBoxManage modifyvm tklinux --vrdeauthtype external

Corresponding cmdfile:

[vbox]

vrdeauthtype = external

[code]

RUNFLAG VBoxManage modifyvm tklinux @vbox
```

Sample Command

NOTE: You cannot use both % and @ in one line 

```
sudo docker run -d -m 100m -e DEVELOPMENT=1 \
	-e BRANCH=example-code \
	-v $(pwd):/app/bin:ro \
	--name app appserver
```

can be rewritten as cmdfile like so:

```
[docker]

m = 100m
e = DEVELOPMENT=1
e = BRANCH=example-code
v = $(pwd):/app/bin:ro

[code]

RUNFLAG sudo docker run -d --name app appserver %docker
```

INCLUDE

```
INCLUDE /path/to/cmdfile

# executes cmdfile (division of labor)
# note that cmdfile invoked by INCLUDE executes in its own context
# that is, it is not seen by other cmdfile, hence it is self-contained
```

TEMPLATE

```
[code]

# TEMPLATE template file, json file, destination file, shell script file

TEMPLATE test.tmpl test.json testconfig.txt testecho.bat
```
