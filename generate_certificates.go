package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var stringFragments = make(map[string]string)

//Takes in a string and runs the command in a shell
func runCommand(command string) error {
	executableCommand := convertStringIntoExecCommand(command)
	fmt.Println("runCommand:")
	fmt.Println(command)
	return executableCommand.Run()
}

//Takes in a string and produces Cmd object that can be run
func convertStringIntoExecCommand(command string) *exec.Cmd {
	arguments := strings.Split(command, " ")
	return exec.Command(arguments[0], arguments[1:]...)
}

//Incorrect
func signCertificate(signerPrivateKey string, certificateSigningRequest string, outputCertificate string) {
	//openssl ca -selfsign -keyfile root.pem -config root_ca.conf -out root.crt -in root.csr -outdir root_certificates -verbose -batch
	command := ""
	executableCommand := convertStringIntoExecCommand(command)
	executableCommand.Run()
}

//Generates a certificate based off of the root private key, root authority openssl confiration file, output filename and output directory
func generateSelfSignedCertificate(privateKey, configuration, outputCertificateFilename, certificateSigningRequest, outputDirectory string) {

	command := fmt.Sprintf("openssl ca -selfsign -keyfile %s -config %s -out %s -in %s -outdir %s -verbose -batch",
		privateKey, configuration, outputCertificateFilename, certificateSigningRequest, outputDirectory)

	fmt.Println("Inside generateSelfSignedCertificate", command)
	err := runCommand(command)
	if err != nil {
		fmt.Println("Error during generation of self-signed certificate. Command was: " + command)
		fmt.Println(err)
		os.Exit(0)
	}
}

//privateKey is a string specifying the filepath of the private key for the entity performing the sign
//outputCertificate is a string specifying the filepath of the certificate that will be generated
//configuration is a string specifying the filepath of a file containing data to be signed
func generateCertificateSigningRequest(privateKey, outputCertificate, configuration string) {
	command := fmt.Sprintf("openssl req -key %s -out %s -days 398 -new -config %s", privateKey, outputCertificate, configuration)
	fmt.Println("Inside generateCertificateSigningRequest. command is:", command)
	err := runCommand(command)
	if err != nil {
		fmt.Println("An error occurred when trying to generate the certificate signing request using the key " + privateKey + " with the configuration " + configuration)
		fmt.Println("The command was:", command)
		fmt.Println(err)
		os.Exit(0)
	}
}

//Generates a signed certificate using the openssl ca command
func generateSignedCertificate(certificateSigningRequest, outputCertificateFilepath, certificateAuthorityConfiguration, certificateAuthoritySigningKey, certificateAuthorityCertificate, outputCertificateDirectory string) {
	command := fmt.Sprintf("openssl ca -in %s -out %s -config %s -keyfile %s -cert %s -outdir %s -batch", certificateSigningRequest, outputCertificateFilepath, certificateAuthorityConfiguration, certificateAuthoritySigningKey, certificateAuthorityCertificate, outputCertificateDirectory)
	fmt.Println("Inside generateSignedCertificate the command is:", command)
	err := runCommand(command)
	if err != nil {
		fmt.Println("An error occurred when trying to generate the signed certificate. The command was: ", command)
		fmt.Println(err)
		os.Exit(0)
	}
}

//Uses OpenSSL to generate a private key
func generatePrivateKey(filename string) error {
	command := fmt.Sprintf("openssl genpkey -outform pem -out %s -algorithm rsa", filename)
	err := runCommand(command)

	if err != nil {
		fmt.Println("An error occurred when trying to generate private key " + filename + " using OpenSSL.")
		fmt.Println(err)
	}
	return err
}

//Returns true if file or directory passed in exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

//Given the domain name directory, generates the directory structure needed for this program
func makeDirectories() {
	//1)Ensure a directory named after the domain name passed in always exists in the output directory
	if !fileExists(stringFragments["outputDirectory"]) {
		fmt.Println("Generating directory: " + stringFragments["outputDirectory"])
		os.Mkdir(stringFragments["outputDirectory"], 0700)
	}

	if !fileExists(stringFragments["domainNameDirectory"]) {
		fmt.Println("Generating directory: " + stringFragments["domainNameDirectory"])
		os.Mkdir(stringFragments["domainNameDirectory"], 0700)
	}

	if !fileExists(stringFragments["rootAuthorityDirectory"]) {
		fmt.Println("Generating directory: " + stringFragments["rootAuthorityDirectory"])
		os.Mkdir(stringFragments["rootAuthorityDirectory"], 0700)
	}

	if !fileExists(stringFragments["intermediateAuthorityDirectory"]) {
		fmt.Println("Generating directory: " + stringFragments["intermediateAuthorityDirectory"])
		os.Mkdir(stringFragments["intermediateAuthorityDirectory"], 0700)
	}
}

//Takes in an output directory and generates 3 private keys, one for the root authority, one for the intermediate authority, and one for the server hosting the domain name.
func makePrivateKeys() {
	//2)Create a root authority private key if it doesn't already exist. Do not replace an existing one
	//openssl genpkey -outform pem -out root.pem -algorithm rsa
	if !fileExists(stringFragments["rootAuthorityPrivateKey"]) {
		fmt.Println("Generating root private key: " + stringFragments["rootAuthorityPrivateKey"])
		err := generatePrivateKey(stringFragments["rootAuthorityPrivateKey"])
		if err != nil {
			os.Exit(0)
		}
	}

	//3)Create an intermediate authority private key
	if !fileExists(stringFragments["intermediateAuthorityPrivateKey"]) {
		fmt.Println("Generating intermediate private key: " + stringFragments["intermediateAuthorityPrivateKey"])
		generatePrivateKey(stringFragments["intermediateAuthorityPrivateKey"])
	}

	//4)Generate a server private key
	stringFragments["serverPrivateKey"] = stringFragments["domainNameDirectory"] + "/" + stringFragments["serverPrivateKeyFilename"]
	if !fileExists(stringFragments["serverPrivateKey"]) {
		fmt.Println("Generating server private key: " + stringFragments["serverPrivateKey"])
		generatePrivateKey(stringFragments["serverPrivateKey"])
	}
}

//Copies the file in source to destination
func fileCopy(src, dst string) {
	fmt.Println("Inside fileCopy src:", src, "| dst:", dst, "|")
	bytesRead, err := ioutil.ReadFile(src)

	if err != nil {
		fmt.Println("Error reading from ", src)
		fmt.Println(err)
	}

	err = ioutil.WriteFile(dst, bytesRead, 0644)

	if err != nil {
		fmt.Println("Error writing to: ", dst)
		fmt.Println(err)
	}
}

//template: a template openssl configuration file
//output: the output file that will be guaranteed to exist. One will be generated by copying
//args: a string array of
func hydrateTemplate(template, output string, args ...any) {
	fmt.Println("Hydrating ", template, " into ", output)
	fileCopy(template, output)

	contentAsBytes, err := ioutil.ReadFile(output)
	if err != nil {
		fmt.Println("Error while reading " + output)
		fmt.Println(err)
	}

	//After copying the contents of the file needs to be altered to match input domain name
	contentsAsString := string(contentAsBytes[:])
	newFileContents := fmt.Sprintf(contentsAsString, args...)

	ioutil.WriteFile(output, []byte(newFileContents), 0644)
}

func makeServerCertificate() {
	fmt.Println("serverCSR:", stringFragments["serverCSR"])
	fmt.Println("serverCSRConfig:", stringFragments["serverCSRConfig"])
	fmt.Println("serverCSRConfigTemplate:", stringFragments["serverCSRConfigTemplate"])

	if !fileExists(stringFragments["serverCSRConfig"]) {
		hydrateTemplate(stringFragments["serverCSRConfigTemplate"], stringFragments["serverCSRConfig"], stringFragments["domainName"])
	}

	if !fileExists(stringFragments["serverConfig"]) {
		fmt.Println(stringFragments["serverConfigTemplate"])
		fmt.Println(stringFragments["serverConfig"])
		hydrateTemplate(stringFragments["serverConfigTemplate"], stringFragments["serverConfig"], stringFragments["intermediateAuthorityDatabase"], stringFragments["intermediateAuthoritySerialNumber"], stringFragments["domainName"])
	}

	if !fileExists(stringFragments["serverCSR"]) {
		fmt.Println("Generating server CSR")
		generateCertificateSigningRequest(stringFragments["serverPrivateKey"], stringFragments["serverCSR"], stringFragments["serverCSRConfig"])
	}

	if !fileExists(stringFragments["serverCertificate"]) {
		fmt.Println("Generating server certificate")
		fmt.Println(stringFragments["serverCSR"])
		fmt.Println(stringFragments["serverCertificateFilename"])
		fmt.Println(stringFragments["serverConfig"])
		fmt.Println(stringFragments["intermediateAuthorityPrivateKey"])
		fmt.Println(stringFragments["intermediateAuthorityCertificate"])
		fmt.Println(stringFragments["domainNameDirectory"])

		generateSignedCertificate(stringFragments["serverCSR"], stringFragments["serverCertificate"], stringFragments["serverConfig"], stringFragments["intermediateAuthorityPrivateKey"], stringFragments["intermediateAuthorityCertificate"], stringFragments["domainNameDirectory"])
	}
}

//This is not correct
func makeIntermediateAuthorityCertificate() {
	if !fileExists(stringFragments["intermediateAuthorityMakeInformationCSRConfig"]) {
		fileCopy(stringFragments["intermediateAuthorityMakeInformationCSRConfigTemplate"], stringFragments["intermediateAuthorityMakeInformationCSRConfig"])
	}

	//Generate the intermediate authority CSR if it doesn't already exist
	if !fileExists(stringFragments["intermediateAuthorityCSR"]) {
		fmt.Println("Generating intermediate CSR.") //This is the request from the intermediate authority to the root authority to sign its certificate
		generateCertificateSigningRequest(stringFragments["intermediateAuthorityPrivateKey"], stringFragments["intermediateAuthorityCSR"], stringFragments["intermediateAuthorityMakeInformationCSRConfig"])
	}

	//ensure the openssl configuration file for making the intermediate authority certificate is present
	if !fileExists(stringFragments["intermediateAuthorityMakeCertificateConfiguration"]) {
		hydrateTemplate(
			stringFragments["intermediateAuthorityConfigTemplate"],
			stringFragments["intermediateAuthorityMakeCertificateConfiguration"],
			stringFragments["intermediateAuthorityDatabase"],
			stringFragments["intermediateAuthoritySerialNumber"])
	}

	fmt.Println("Generating intermediate certificate")
	fmt.Println("Inside makeIntermediateAuthorityCertificate")
	fmt.Println(stringFragments["intermediateAuthorityCSR"], stringFragments["intermediateAuthorityCertificate"], stringFragments["intermediateAuthorityMakeCertificateConfiguration"], stringFragments["rootAuthorityPrivateKey"], stringFragments["rootAuthorityCertificate"], stringFragments["intermediateAuthorityDirectory"])
	generateSignedCertificate(stringFragments["intermediateAuthorityCSR"], stringFragments["intermediateAuthorityCertificate"], stringFragments["intermediateAuthorityMakeCertificateConfiguration"], stringFragments["rootAuthorityPrivateKey"], stringFragments["rootAuthorityCertificate"], stringFragments["intermediateAuthorityDirectory"])
}

func makeRootAuthorityCertificate() {
	//Stage 4
	//1)Generate root certificate
	//2)Generate intermediate certificate
	//3)Generate server certificate
	stringFragments["rootCSR"] = stringFragments["rootAuthorityDirectory"] + "/root.csr"
	if !fileExists(stringFragments["rootCSR"]) {
		fmt.Println("Generating root CSR: ", stringFragments["rootCSR"])

		stringFragments["rootAuthorityCSRConfig"] = stringFragments["rootAuthorityDirectory"] + "/" + stringFragments["rootAuthorityMakeInformationCSRConfigFilename"]
		if !fileExists(stringFragments["rootAuthorityCSRConfig"]) {
			//Copy the config file over from the templates directory to the directoryNameDirectory if it doesn't exist
			templatesDirectoryRootAuthorityCSRConfig := stringFragments["templatesDirectory"] + "/" + stringFragments["rootAuthorityMakeInformationCSRConfigFilename"]
			fileCopy(templatesDirectoryRootAuthorityCSRConfig, stringFragments["rootAuthorityCSRConfig"])
		}
		generateCertificateSigningRequest(stringFragments["rootAuthorityPrivateKey"], stringFragments["rootCSR"], stringFragments["rootAuthorityCSRConfig"])
	}

	if !fileExists(stringFragments["rootAuthorityCertificate"]) {
		if !fileExists(stringFragments["rootAuthorityMakeCertificateConfiguration"]) {
			fmt.Println("Copying root authority certificate generation configuration from ", stringFragments["rootAuthorityConfigTemplate"], " to ", stringFragments["rootAuthorityMakeCertificateConfiguration"])

			hydrateTemplate(stringFragments["rootAuthorityConfigTemplate"], stringFragments["rootAuthorityMakeCertificateConfiguration"], stringFragments["rootAuthorityDatabase"], stringFragments["rootAuthoritySerialNumber"])
		}

		fmt.Println("Generating root certificate: ", stringFragments["rootAuthorityCertificate"])
		generateSelfSignedCertificate(stringFragments["rootAuthorityPrivateKey"], stringFragments["rootAuthorityMakeCertificateConfiguration"], stringFragments["rootAuthorityCertificate"], stringFragments["rootCSR"], stringFragments["rootAuthorityDirectory"])
	}
}

func initializeStringFragments() {
	stringFragments["rootAuthorityPrivateKeyFilename"] = "root.pem"
	stringFragments["intermediateAuthorityPrivateKeyFilename"] = "intermediate.pem"
	stringFragments["serverPrivateKeyFilename"] = "server.pem"

	stringFragments["outputDirectory"] = "output"
	stringFragments["templatesDirectory"] = "templates"
	stringFragments["domainNameDirectory"] = stringFragments["outputDirectory"] + "/" + stringFragments["domainName"]

	stringFragments["rootAuthorityMakeInformationCSRConfigFilename"] = "make_root_information_csr.conf"
	stringFragments["rootAuthorityMakeCertificateFilename"] = "make_root_certificate.conf"
	stringFragments["rootAuthorityDirectory"] = stringFragments["outputDirectory"] + "/root_authority"
	stringFragments["rootAuthorityPrivateKey"] = stringFragments["rootAuthorityDirectory"] + "/" + stringFragments["rootAuthorityPrivateKeyFilename"]
	stringFragments["rootAuthorityDatabaseFilename"] = "root_database.txt"
	stringFragments["rootAuthoritySerialNumberFilename"] = "root_serial_number.txt"
	stringFragments["rootAuthorityCertificateFilename"] = "root.crt"
	stringFragments["rootAuthorityMakeCertificateConfiguration"] = stringFragments["rootAuthorityDirectory"] + "/" + stringFragments["rootAuthorityMakeCertificateFilename"]
	stringFragments["rootAuthorityDatabase"] = stringFragments["rootAuthorityDirectory"] + "/" + stringFragments["rootAuthorityDatabaseFilename"]
	stringFragments["rootAuthoritySerialNumber"] = stringFragments["rootAuthorityDirectory"] + "/" + stringFragments["rootAuthoritySerialNumberFilename"]
	stringFragments["rootAuthorityConfigTemplate"] = stringFragments["templatesDirectory"] + "/" + stringFragments["rootAuthorityMakeCertificateFilename"]
	stringFragments["rootAuthorityCertificate"] = stringFragments["rootAuthorityDirectory"] + "/" + stringFragments["rootAuthorityCertificateFilename"]

	stringFragments["intermediateAuthorityMakeInformationCSRConfigFilename"] = "make_intermediate_information_csr.conf"
	stringFragments["intermediateAuthorityDirectory"] = stringFragments["outputDirectory"] + "/intermediate_authority"
	stringFragments["intermediateAuthorityPrivateKey"] = stringFragments["intermediateAuthorityDirectory"] + "/" + stringFragments["intermediateAuthorityPrivateKeyFilename"]
	stringFragments["intermediateAuthorityMakeCertificateConfigurationFilename"] = "make_intermediate_certificate.conf"
	stringFragments["intermediateAuthorityMakeCertificateConfiguration"] = stringFragments["intermediateAuthorityDirectory"] + "/" + stringFragments["intermediateAuthorityMakeCertificateConfigurationFilename"]
	stringFragments["intermediateAuthorityMakeInformationCSRConfig"] = stringFragments["intermediateAuthorityDirectory"] + "/" + stringFragments["intermediateAuthorityMakeInformationCSRConfigFilename"]
	stringFragments["intermediateAuthorityDatabase"] = stringFragments["intermediateAuthorityDirectory"] + "/intermediate_database.txt"
	stringFragments["intermediateAuthoritySerialNumber"] = stringFragments["intermediateAuthorityDirectory"] + "/intermediate_serial_number.txt"
	stringFragments["intermediateAuthorityCSR"] = stringFragments["intermediateAuthorityDirectory"] + "/intermediate.csr"
	stringFragments["intermediateAuthorityMakeInformationCSRConfigTemplate"] = stringFragments["templatesDirectory"] + "/" + stringFragments["intermediateAuthorityMakeInformationCSRConfigFilename"]
	stringFragments["intermediateAuthorityConfigTemplate"] = stringFragments["templatesDirectory"] + "/" + stringFragments["intermediateAuthorityMakeCertificateConfigurationFilename"]
	stringFragments["intermediateAuthorityCertificate"] = stringFragments["intermediateAuthorityDirectory"] + "/intermediate.crt"

	stringFragments["serverCSR"] = stringFragments["domainNameDirectory"] + "/server.csr"
	stringFragments["serverCSRConfigFilename"] = "make_server_information_csr.conf"
	stringFragments["serverCSRConfig"] = stringFragments["domainNameDirectory"] + "/" + stringFragments["serverCSRConfigFilename"]
	stringFragments["serverCSRConfigTemplate"] = stringFragments["templatesDirectory"] + "/" + stringFragments["serverCSRConfigFilename"]

	stringFragments["serverConfigFilename"] = "make_server_certificate.conf"
	stringFragments["serverConfig"] = stringFragments["domainNameDirectory"] + "/" + stringFragments["serverConfigFilename"]
	stringFragments["serverConfigTemplate"] = stringFragments["templatesDirectory"] + "/" + stringFragments["serverConfigFilename"]
	stringFragments["serverCertificateFilename"] = "server.crt"
	stringFragments["serverCertificate"] = stringFragments["domainNameDirectory"] + "/" + stringFragments["serverCertificateFilename"]
	stringFragments["serverBundleCertificate"] = stringFragments["domainNameDirectory"] + "/server_bundle.crt"

}

func makeServerCertificateBundle() {
	fmt.Println("Generating server certificate bundle")

	serverCertificateData, err := ioutil.ReadFile(stringFragments["serverCertificate"])
	if err != nil {
		fmt.Println("Error reading server certificate during bundle generation:")
		fmt.Println(err)
		os.Exit(0)
	}

	intermediateCertificateData, err := ioutil.ReadFile(stringFragments["intermediateAuthorityCertificate"])
	if err != nil {
		fmt.Println("Error reading intermediate certificate during bundle generation:")
		fmt.Println(err)
		os.Exit(0)
	}

	rootCertificateData, err := ioutil.ReadFile(stringFragments["rootAuthorityCertificate"])
	if err != nil {
		fmt.Println("Error reading root certificate during bundle generation:")
		fmt.Println(err)
		os.Exit(0)
	}

	certificateBundleData := append(serverCertificateData, intermediateCertificateData...)
	certificateBundleData = append(certificateBundleData, rootCertificateData...)

	err = ioutil.WriteFile(stringFragments["serverBundleCertificate"], certificateBundleData, 0644)

	if err != nil {
		fmt.Println("Error writing to ", stringFragments["serverBundleCertificate"])
		fmt.Println(err)
	}
}

//Makes the database file and serial number needed for the OpenSSL ca command for
//the intermediate and root certificates
func makeDatabaseFiles() {
	//Must also ensure the files referenced in the root authority configuration file exists
	if !fileExists(stringFragments["rootAuthorityDatabase"]) {
		fmt.Println("Generating root database file:" + stringFragments["rootAuthorityDatabase"])
		rootAuthorityDatabaseFile, error := os.Create(stringFragments["rootAuthorityDatabase"])
		if error != nil {
			fmt.Println("Error while creating root authority database file:" + stringFragments["rootAuthorityDatabase"])
			fmt.Println(error)
		}

		rootAuthorityDatabaseFile.Close()
	}

	if !fileExists(stringFragments["rootAuthoritySerialNumber"]) {
		fmt.Println("Generating root serial number file:" + stringFragments["rootAuthoritySerialNumber"])
		rootAuthoritySerialNumberFile, error := os.Create(stringFragments["rootAuthoritySerialNumber"])
		if error != nil {
			fmt.Println(error)
		}

		//Serial numbers file needs to have the hexadecimal digit 01 in it when initially created.
		rootAuthoritySerialNumberFile.WriteString("01")
		rootAuthoritySerialNumberFile.Close()
	}

	if !fileExists(stringFragments["intermediateAuthorityDatabase"]) {
		fmt.Println("Generating intermediate database file:" + stringFragments["intermediateAuthorityDatabase"])
		intermediateAuthorityDatabaseFile, error := os.Create(stringFragments["intermediateAuthorityDatabase"])
		if error != nil {
			fmt.Println("Error while creating intermediate authority database file:" + stringFragments["intermediateAuthorityDatabase"])
			fmt.Println(error)
		}

		intermediateAuthorityDatabaseFile.Close()
	}

	if !fileExists(stringFragments["intermediateAuthoritySerialNumber"]) {
		fmt.Println("Generating intermediate serial number file:" + stringFragments["intermediateAuthoritySerialNumber"])
		intermediateAuthoritySerialNumberFile, error := os.Create(stringFragments["intermediateAuthoritySerialNumber"])
		if error != nil {
			fmt.Println(error)
		}

		//Serial numbers file needs to have the hexadecimal digit 01 in it when initially created.
		intermediateAuthoritySerialNumberFile.WriteString("05")
		intermediateAuthoritySerialNumberFile.Close()
	}
}

//Usage: go run generate_certificates.go <domain.name>
//domain.name will be created as a directory and files generated by generate_certificates.go will go into the directory with name "domain.name".

//In the code, the term "server" refers to the computer hosting the name domain.name
func main() {
	//Force there to be exactly two arguments, the name of the file and the domain name
	if len(os.Args) != 2 {
		fmt.Println("Error: no domain name specified.")
		fmt.Println("usage: go run generate_certificates.go <domain.name>")
		os.Exit(0)
	}

	stringFragments["domainName"] = os.Args[1]

	//Stage 1
	initializeStringFragments()

	stringFragments["domainNameDirectory"] = stringFragments["outputDirectory"] + "/" + stringFragments["domainName"]

	//Stage 2
	makeDirectories()
	makeDatabaseFiles()

	//Stage 3
	makePrivateKeys()

	//Stage 4.
	makeRootAuthorityCertificate()

	makeIntermediateAuthorityCertificate()

	makeServerCertificate()

	makeServerCertificateBundle()
}
