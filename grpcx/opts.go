package grpcx

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

func NewManualVerifyPeerCertificate(tlsSkip bool, fingerprint, dnsName string) *ManualVerifyPeerCertificate {
	return &ManualVerifyPeerCertificate{
		DNSName:       dnsName,
		errCh:         make(chan error, 1),
		SkipTLSVerify: tlsSkip,
		Fingerprint:   SHA1Fingerprint(fingerprint),
	}
}

type ManualVerifyPeerCertificate struct {
	DNSName       string
	SkipTLSVerify bool
	errCh         chan error
	// checker of fingerprint from server cert
	Fingerprint SHA1Fingerprint
}

func (v ManualVerifyPeerCertificate) releaseErr(err error) {
	select {
	case v.errCh <- err:
	default:
		panic(fmt.Sprint("ManualVerifyPeerCertificate: no one listens (or errCh is closed) to chan with errors - got error", err))
	}
}

func (v ManualVerifyPeerCertificate) VerifyPeerCertificateOption(ctx context.Context, cancel context.CancelFunc) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	go func() {
		defer close(v.errCh)
		select {
		case <-ctx.Done():
			log.Println("parent context is done")
			cancel()
			return
		case err, ok := <-v.errCh:
			if !ok {
				log.Println("premature close of a channel with errors")
				cancel()
				return
			}
			if err != nil {
				log.Println("failed valid cert with err", err)
				cancel()
				return
			}
		}
	}()
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		opts := x509.VerifyOptions{
			// TODO: add rootCAs
			CurrentTime:   time.Now(),
			DNSName:       v.DNSName,
			Intermediates: x509.NewCertPool(),
		}
		// Coped code from https://github.com/golang/go/blob/1419ca7cead4438c8c9f17d8901aeecd9c72f577/src/crypto/tls/handshake_client.go#L835
		certs := make([]*x509.Certificate, len(rawCerts))
		for i, asn1Data := range rawCerts {
			cert, err := x509.ParseCertificate(asn1Data)
			if err != nil {
				v.releaseErr(errors.Wrap(err, "failed parse cert from server"))
				return errors.New("tls: failed to parse certificate from server: " + err.Error())
			}
			certs[i] = cert
		}

		if !v.SkipTLSVerify {
			for _, cert := range certs[1:] {
				opts.Intermediates.AddCert(cert)
			}
			_, err := certs[0].Verify(opts)
			// certErr := x509.CertificateInvalidError{}
			// if errors.As(err, &certErr) {
			// 	switch certErr.Reason {
			// 	case x509.Expired:
			// 		log.Println("Expired")
			// 	default:
			// 		log.Println("!>!>!>!> Failed verify cert", err)
			// 		return err
			// 	}
			// } else {
			// 	log.Println("!>!>!>!> (not expected error) Failed verify cert", err)
			// 	return err
			// }
			if err != nil {
				v.releaseErr(errors.Wrap(err, "failed TLS verify cert"))
				return err
			}
		}

		if !v.Fingerprint.Empty() {
			if !v.Fingerprint.Match(rawCerts[0]) {
				err := errors.New("fingerprint does not match")
				v.releaseErr(err)
				return err
			}
		}

		// TODO: more checks if need
		// certs[0].DNSNames
		// certs[0].Subject,
		// certs[0].VerifyHostname(...)
		return nil
	}
}
