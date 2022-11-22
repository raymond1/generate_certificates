This piece of software uses openssl to generate TLS certificates meant to be used for local testing. The idea behind this piece of software is that someone who is developing software and needs an https connection can do so with a single command followed by several configuration steps.

Installation and set up requirements:

OpenSSL.
Go.
NodeJS.
Chrome or similar browser.

Usage:
```
go run generate_certificates.go <domain.name>
```

After running the command, you will generate an output folder. Under the output folder, there will be a root_authority folder containing the root certificate in the file root.crt, amongst other files. Add this certificate to the list of certificates in Keychain Access in MacOS. Then, always trust the certificate. Then, add <domain.name> to your /etc/hosts file. For me, the line looks like: 
```
127.0.0.1       <domain.name>
```

In Firefox, you may need to go into about:config and set "security.enterprise_roots.enabled" to true.

To test if this command has succeeded, edit the nodejs test server located in server/server.js.
The lines:
```
const options = {
  key: fs.readFileSync('output/simple.dev/server.pem'),
  cert: fs.readFileSync('output/simple.dev/server_bundle.crt')
};
```

should be modified so that simple.dev is replaced with your domain name. In other words, the key field of the options object should point to the private key that was generated for your domain name and the cert key should point to the certificate bundle server_bundle.crt.

Then, you should be able to go into your browser and type https://<domain.name> and see the message hello if everything is working.

Please note that this software, while useful, has only been tested a few times. It is very complicated and requires a lot of knowledge to truly understand. There are many situations where these commands may not work, or may not be appropriate. It is up to you as a developer to understand the security implications of using this software.

The main thing you should know is that by trusting the root certificate or the intermediate certificate that you generated on your system is that if you or anyone else who has access to the root or intermediate keys uses it to trust other certificates, then it is possible that other website urls which are signed with your root private key or the intermediate authority private key that is generated in the software can appear to your browser as legitimate websites when they are not. While this can only happen through deliberate action, a nefarious person could potentially harm you when you trust a the new root certificate that you generate with this software. Or perhaps you might accidentally mix up live and development websites and lose some data. Other scenarios are possible.

This software is free and not well tested. I don't expect it to damage people's computers, but neither can I guarantee that damage will not happen. Use at your own risk. I take no responsibility for problems may be caused as a result of you using this software.