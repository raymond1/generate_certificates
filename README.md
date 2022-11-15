This Makefile uses openssl to generate certificates meant to be used for local testing.

Requirements:

OpenSSL.
GNU Make.

Usage:

In the root directory, type in 

```
make
```

Various SSL-related files and directories will be output at this point.

To start the test server, in the root directory type

```
node test_server/server.js
```
