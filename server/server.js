import * as express from "express"
import * as https_package from "node:https"
import * as fs_package from "node:fs"
console.log(Object.keys(fs_package))


let fs = fs_package.default

const options = {
  key: fs.readFileSync('output/simple.dev/server.pem'),
  cert: fs.readFileSync('output/simple.dev/server_bundle.crt')
};


let https = https_package.default

https.createServer(options, (req, res) => {
  res.writeHead(200);
  res.end('hello world\n');
}).listen(443);

// console.log(Object.keys(https))

// let app = express.default()
// let port = 3000

// app.use(express.static('public'))
// app.listen(port, (error)=>{
// 	console.log('Server running on port ', port)
// })
/*
const https = require('node:https');
const fs = require('node:fs');

const options = {
  key: fs.readFileSync('test/fixtures/keys/agent2-key.pem'),
  cert: fs.readFileSync('test/fixtures/keys/agent2-cert.pem')
};

https.createServer(options, (req, res) => {
  res.writeHead(200);
  res.end('hello world\n');
}).listen(8000);
*/
