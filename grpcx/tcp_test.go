package grpcx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTCPConnectOK(t *testing.T) {
	tlss := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer tlss.Close()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer s.Close()

	tlsAddr := tlss.Listener.Addr().String()
	plainAddr := s.Listener.Addr().String()

	tests := []struct {
		timeout    time.Duration
		addr       string
		wantErr    bool
		errContain string
	}{
		{time.Second, tlsAddr, false, ""},
		{time.Second, plainAddr, false, ""},
		{time.Nanosecond, tlsAddr, true, "timeout"},
		{time.Nanosecond, plainAddr, true, "timeout"},
		{time.Second, "127.0.0.1:10011", true, "refused"},
		{time.Nanosecond, "127.0.0.1:10011", true, "timeout"},
		{time.Second, "10.9.8.7:1234", true, "timeout"},
		{time.Nanosecond, "10.9.8.7:1234", true, "timeout"},
	}
	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			err := TCPConnectOK(tt.timeout, tt.addr)
			t.Logf("addr %q err = %v", tt.addr, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("TCPConnectOK() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Contains(t, err.Error(), tt.errContain)
			}
		})
	}
}
