package xsecurity

// TLSConfig is used to generate a TLSClientConfig that's useful for talking to
// Consul using TLS.
type TLSConfig struct {
	// Address is the optional address of the Consul server. The port, if any
	// will be removed from here and this will be set to the ServerName of the
	// resulting config.
	Address string

	// CAFile is the optional path to the CA certificate used for Consul
	// communication, defaults to the system bundle if not specified.
	CAFile string

	// CAPath is the optional path to a directory of CA certificates to use for
	// Consul communication, defaults to the system bundle if not specified.
	CAPath string

	// CAPem is the optional PEM-encoded CA certificate used for Consul
	// communication, defaults to the system bundle if not specified.
	CAPem []byte

	// CertFile is the optional path to the certificate for Consul
	// communication. If this is set then you need to also set KeyFile.
	CertFile string

	// CertPEM is the optional PEM-encoded certificate for Consul
	// communication. If this is set then you need to also set KeyPEM.
	CertPEM []byte

	// KeyFile is the optional path to the private key for Consul communication.
	// If this is set then you need to also set CertFile.
	KeyFile string

	// KeyPEM is the optional PEM-encoded private key for Consul communication.
	// If this is set then you need to also set CertPEM.
	KeyPEM []byte

	// InsecureSkipVerify if set to true will disable TLS host verification.
	InsecureSkipVerify bool
}
