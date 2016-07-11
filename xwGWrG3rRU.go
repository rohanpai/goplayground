package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
)

func main() {
	// To get started with SSL/TLS follow the instructions for adding SSL/TLS
	// support in RabbitMQ with a private certificate authority here:
	//
	// http://www.rabbitmq.com/ssl.html
	//
	// Then in your rabbitmq.config, disable the plain AMQP port, verify clients
	// and fail if no certificate is presented with the following:
	//
	//   [
	//   {rabbit, [
	//     {tcp_listeners, []},     % listens on 127.0.0.1:5672
	//     {ssl_listeners, [5671]}, % listens on 0.0.0.0:5671
	//     {ssl_options, [{cacertfile,"/path/to/your/testca/cacert.pem"},
	//                    {certfile,"/path/to/your/server/cert.pem"},
	//                    {keyfile,"/path/to/your/server/key.pem"},
	//                    {verify,verify_peer},
	//                    {fail_if_no_peer_cert,true}]}
	//     ]}
	//   ].

	cfg := new(tls.Config)

	// The self-signing certificate authority's certificate must be included in
	// the RootCAs to be trusted so that the server certificate can be verified.
	//
	// Alternatively to adding it to the tls.Config you can add the CA's cert to
	// your system's root CAs.  The tls package will use the system roots
	// specific to each support OS.  Under OS X, add (drag/drop) your cacert.pem
	// file to the 'Certificates' section of KeyChain.app to add and always
	// trust.
	//
	// Or with the command line add and trust the DER encoded certificate:
	//
	//   security add-certificate testca/cacert.cer
	//   security add-trusted-cert testca/cacert.cer
	//
	// If you depend on the system root CAs, then use nil for the RootCAs field
	// so the system roots will be loaded.

	cfg.RootCAs = x509.NewCertPool()

	if ca, err := ioutil.ReadFile("testca/cacert.pem"); err == nil {
		cfg.RootCAs.AppendCertsFromPEM(ca)
	}

	// Move the client cert and key to a location specific to your application
	// and load them here.

	if cert, err := tls.LoadX509KeyPair("client/cert.pem", "client/key.pem"); err == nil {
		cfg.Certificates = append(cfg.Certificates, cert)
	}

	// Server names are validated by the crypto/tls package, so the server
	// certificate must be made for the hostname in the URL.  Find the commonName
	// (CN) and make sure the hostname in the URL matches this common name.  Per
	// the RabbitMQ instructions for a self-signed cert, this defautls to the
	// current hostname.
	//
	//   openssl x509 -noout -in server/cert.pem -subject
	//
	// If your server name in your certificate is different than the host you are
	// connecting to, set the hostname used for verification in
	// ServerName field of the tls.Config struct.

	conn, err := amqp.DialTLS("amqps://server-name-from-certificate/", cfg)

	log.Print("conn: %v, err: %v", conn, err)
}
