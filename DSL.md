DSL specification (syntax and rules)
------------------------------------

Here is the specification for writing a cmdfile.

#Comment

```
# single-line comment

# there is no block comment like in Go /* ... */
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

RUN

```script
RUN VBoxManage modifyvm tklinux -memory 512MB
RUN echo "Hello world" > /home/go/echoFile.txt
RUN echo "Hello world2" >> /home/go/echoFile.txt
RUN echo "Hello world ..."
RUN sed -i.bak s/^\(VAR5=\).*/\1VALUE10/ file.cfg
RUN chmod +x cmd.sh
```


