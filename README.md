# Memento backup system

Link: https://github.com/mementobackup/

## Description:

Memento is a backup system for remote machines. It is similar of other systems
(e.g. rsnapshot) but it has some differences

## Licence:

Memento is released under GPLv2 (see GPL.txt for details)

## Features:

It use an agent for check and download data.
It manage four distinct datasets (hour, day, week, month).
It save space with hard link creation.
It save data attributes (owner, group, permission, ACL) into database.

## Dependencies:

Client doesn't have external edependencies

## Building:
```
mkdir mclient && cd mclient
export GOPATH=`pwd`
git clone git@github.com:mementobackup/client.git .
go get github.com/gaal/go-options/options
go get github.com/go-ini/ini
go get github.com/mementobackup/common/src/common
go build mclient.go
```

## Installation:

 - Execute the client.

## Usage:

Usage is simpliest:
```
mclient --port <port> --cfg=<cfgfile> -H # hour backup
```

Where `<port>` is a port which to listen the agent.
For other options, use -h switch. If you want to use SSL connection,
you need to create a configuration file similar at this example:
```
[ssl]
certificate = ssl.crt
key = ssl.key
```

## SSL:

If you want use the SSL connection, you need:

 - Create SSL certificate with these commands:
    ```
    openssl genrsa -des3 -out memento.key 2048
    openssl rsa -in memento.key -out memento.key
    openssl req -new -key memento.key -out memento.csr
    openssl x509 -req -days 365 -in memento.csr -signkey memento.key -out memento.crt
    ```

   In particular, be sure to add the hostname of the client machine in the CN field
 - Configure server for use SSL with the same certificate.

