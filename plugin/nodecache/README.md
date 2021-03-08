# nodecache

## Name

*nodecache* - cache record in client-side, avoid DNS cluster down.

## Description

CoreDNS Node Cache plugin is a plugin for node cache, which can provider dns local service(like nscd, because nscd cannot affect Golang code)


## Syntax

~~~ txt
nodecache
~~~

## Examples

Start a server on the default port and load the *validation* plugin with a whoami plugin.

~~~ corefile
. {
    nodecache 
    forward
}
~~~