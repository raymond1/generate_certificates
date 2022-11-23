This piece of software uses openssl to generate TLS certificates meant to be used for local testing. The idea behind this piece of software is that someone who is developing software and needs an https connection can do so with a single command followed by several configuration steps.

# Installation and set up requirements:

* OpenSSL
* Go
* NodeJS
* Chrome or similar browser that supports enterprise security, where the browser trusts certificates that are trusted by your computer

# Usage:

There are two steps to using this software.

## Step 1
```
go run generate_certificates.go <domain.name>
```

## Step 2: Configuration

After running the command from step 1, a folder named "output" will be generated, along with files and subfolders. Under the output folder, there will be a root_authority folder containing the root certificate in the file root.crt, amongst other files. Add this certificate to the list of certificates in Keychain Access in MacOS. Then, always trust the certificate. Then, add <domain.name> to your /etc/hosts file. For me, the line looks like: 
```
127.0.0.1       <domain.name>
```

In Firefox, you will need to go into about:config and set "security.enterprise_roots.enabled" to true.

# Testing
To test if this command has succeeded, edit the NodeJS test server located in server/server.js.
The lines:
```
const options = {
  key: fs.readFileSync('output/simple.dev/server.pem'),
  cert: fs.readFileSync('output/simple.dev/server_bundle.crt')
};
```

should be modified so that simple.dev is replaced with your domain name. In other words, the key field of the options object should point to the private key that was generated for your domain name and the cert key should point to the certificate bundle server_bundle.crt. The private key for your domain name should be located in the file output/<domain.name>/server.pem and the certificate for your domain name should be located in the file output/<domain.name>/server.crt.

Then, you should be able to go into your browser and type https://<domain.name> and see the message "hello" coming from ther NodeJS server if everything is working.

# Inner Workings Overview
This software works by generating the following things:
1. The root private key, located in output/root_authority/root.pem
2. The intermediate certificate authority private key, located in output/intermediate_authority/intermediate.pem
3. The server private key, located in output/<domain.name>/server.pem
4. The self-signed root certificate, located in output/root_authority/root.crt
5. The intermediate certificate, signed by the root, located in output/intermediate_authority/intermediate.crt
6. The server certificate, signed by the intermediate authority, located in output/<domain.name>/server.crt

In addition, several intermediate steps generate other files such as certificate signing requests, OpenSSL certificate authority database files and copies of old certificates issued by the root authority and the intermediate authority.

# Other Usage Details
If you delete the entire output directory and run the script again, a new set of root, intermediate and server keys and certificates will be generated.

If a file already exists, it will not be created. So, for example, if you ran the script ```go run generate_certificates.go <domain.name> ```
once and generated a root key, an intermediate key, a server key, a root certificate, an intermediate certificate and server certificate, if you run the script again, nothing will be generated. To regenerate the server certificate(output/<domain.name>/server.crt), you will have to delete it and the run the ``go run generate_certificates.go <domain.name>``` command again. To regenerate the root certificate, you will have to delete that file and again run the ```go run generate_certificates.go <domain.name>``` command again.
