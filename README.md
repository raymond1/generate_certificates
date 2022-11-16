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
node server/server.js
```

STEPS:
step_1: Make all the private keys
step_2: Make all the certs
step_3: Trust the certs in Mac.

Troubleshooting:
...
