This piece of software uses openssl to generate TLS certificates meant to be used for local testing. The idea behind this piece of software is that someone who is developing software and needs an https connection can do so with a single command followed by several configuration steps.

Usage:
```
go run generate_certificates.go <domain.name>
```

After running the command, you will generate an output folder. Under the output folder, there will be a root_authority folder containing the root certificate in the file root.crt, amongst other files. Add this certificate to the list of certificates in Keychain Access in MacOS. Then, always trust the certificate. Then, add <domain.name> to your /etc/hosts file. For me, the line looks like: 
```
127.0.0.1       <domain.name>
```




It was tested on Mac OS only.

Requirements:

OpenSSL.
Go.
NodeJS.
Chrome or similar browser.

Usage:

In the root directory, type in 

```
go run generate_certificates.go <domain.name>
```

Various TLS-related files and directories will be output at this point in the "output" directory.

Go into the server/server.js file and find "server.pem" and "server_bundle.crt". In server.js, point the "key" and "cert" option fields of the 

```
node server/server.js
```

STEPS:
step_1: Make all the private keys
step_2: Make all the certs
step_3: Trust the root cert in the Mac keychain
step_4: Enable security.enterprise_roots.enabled, or make sure that you have an enterprise type security enabled for the browser of your choice.

...


Notes:
Only root server needs to be trusted.

Troubleshooting:
SEC_ERROR_REUSED_ISSUER_AND_SERIAL in Firefox. To fix this, try deleting the intermediate certificate, the server certificate and the server bundle certificate, regenerating and reloading the test website.

//SiteSecurityServiceState.txt in the profile folder
//security.disable_button.openCertManager
//about:config