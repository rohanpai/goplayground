package main

import (
	&#34;crypto/tls&#34;
	&#34;crypto/x509&#34;
	&#34;github.com/streadway/amqp&#34;
	&#34;io/ioutil&#34;
	&#34;log&#34;
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
	//     {ssl_options, [{cacertfile,&#34;/path/to/your/testca/cacert.pem&#34;},
	//                    {certfile,&#34;/path/to/your/server/cert.pem&#34;},
	//                    {keyfile,&#34;/path/to/your/server/key.pem&#34;},
	//                    {verify,verify_peer},
	//                    {fail_if_no_peer_cert,true}]}
	//     ]}
	//   ].

	cfg := new(tls.Config)

	// The self-signing certificate authority&#39;s certificate must be included in
	// the RootCAs to be trusted so that the server certificate can be verified.
	//
	// Alternatively to adding it to the tls.Config you can add the CA&#39;s cert to
	// your system&#39;s root CAs.  The tls package will use the system roots
	// specific to each support OS.  Under OS X, add (drag/drop) your cacert.pem
	// file to the &#39;Certificates&#39; section of KeyChain.app to add and always
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

	if ca, err := ioutil.ReadFile(&#34;testca/cacert.pem&#34;); err == nil {
		cfg.RootCAs.AppendCertsFromPEM(ca)
	}

	// Move the client cert and key to a location specific to your application
	// and load them here.

	if cert, err := tls.LoadX509KeyPair(&#34;client/cert.pem&#34;, &#34;client/key.pem&#34;); err == nil {
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

	conn, err := amqp.DialTLS(&#34;amqps://server-name-from-certificate/&#34;, cfg)

	log.Print(&#34;conn: %v, err: %v&#34;, conn, err)
}
