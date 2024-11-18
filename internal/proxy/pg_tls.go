package proxy

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/GoogleCloudPlatform/cloud-sql-proxy/v2/cloudsql"
	"github.com/jackc/pgx/v5/pgproto3"
)

// TLS <-----> Plain Text Connection

// PlainText {
//   Pipe
// 	   TLS
// }

type PgTLSListener struct {
	inner net.Listener
	cfg   *tls.Config
	log   cloudsql.Logger
}

func (pgSrv PgTLSListener) Addr() net.Addr {
	return pgSrv.inner.Addr()
}

func (pgSrv PgTLSListener) Close() error {
	return pgSrv.inner.Close()
}

func (pgSrv PgTLSListener) Accept() (net.Conn, error) {
	innerConn, err := pgSrv.inner.Accept()
	if err != nil {
		return nil, err
	}

	backend := pgproto3.NewBackend(innerConn, innerConn)

	conn := tlsConn{
		listener: &pgSrv,
		inner:    innerConn,
		backend:  backend,
	}

	pgSrv.log.Debugf("accepted new connection, negotiating TLS")
	tlsConn, err := conn.handleStartup()
	if err != nil {
		return nil, err
	}

	return tlsConn, nil
}

type tlsConn struct {
	listener *PgTLSListener
	inner    net.Conn
	backend  *pgproto3.Backend
}

func (conn *tlsConn) Read(b []byte) (int, error) {
	n, err := conn.inner.Read(b)
	return n, err
}
func (conn *tlsConn) Write(b []byte) (int, error) {
	n, err := conn.inner.Write(b)
	return n, err
}

func (conn *tlsConn) Close() error {
	return conn.inner.Close()
}

func (conn *tlsConn) handleStartup() (*tls.Conn, error) {

	startupMsg, err := conn.backend.ReceiveStartupMessage()

	if err != nil {
		return nil, fmt.Errorf("error receiving startup msg: %w", err)
	}

	switch startupMsg.(type) {
	case *pgproto3.SSLRequest:
		_, err = conn.inner.Write([]byte("S")) // confirm SSL
		if err != nil {
			return nil, fmt.Errorf("error sending deny SSL request: %w", err)
		}

		conn.listener.log.Debugf("upgrading client connection to TLS")

		tlsConn := tls.Server(conn.inner, conn.listener.cfg)

		return tlsConn, nil

	default:
		return nil, fmt.Errorf("unknown startup msg: %#v", startupMsg)
	}
}

// type TempDialer struct {
// 	inner net.Dialer
// }

// func (d TempDialer) Dial(ctx context.Context, inst string, opts ...cloudsqlconn.DialOption) (net.Conn, error) {
// 	return d.inner.DialContext(ctx, "tcp", "localhost:5432")
// }

// func (d TempDialer) EngineVersion(ctx context.Context, inst string) (string, error) {
// 	return "POSTGRES_15", nil
// }

// func (d TempDialer) Close() error {
// 	return nil
// }
