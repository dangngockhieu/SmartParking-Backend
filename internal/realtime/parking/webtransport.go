package parking

import (
	"backend/configs"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

type wtSession struct {
	session *webtransport.Session
}

func (s *wtSession) Send(data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	str, err := s.session.OpenUniStreamSync(ctx)
	if err != nil {
		return err
	}
	defer str.Close()

	_, err = str.Write(data)
	return err
}

func (s *wtSession) Close() error {
	return s.session.CloseWithError(0, "")
}

type Server struct {
	hub      *Hub
	server   *webtransport.Server
	certFile string
	keyFile  string
	cfg      *configs.Config
}

func NewServer(hub *Hub, certFile, keyFile string, cfg *configs.Config) *Server {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	h3Server := &http3.Server{
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
			NextProtos:   []string{"h3"},
		},
	}

	// Bắt buộc phải gọi ConfigureHTTP3Server để bật WebTransport support
	// (enable datagrams, thêm SETTINGS enable_webtransport=1, setup ConnContext)
	webtransport.ConfigureHTTP3Server(h3Server)

	return &Server{
		hub:      hub,
		certFile: certFile,
		keyFile:  keyFile,
		cfg:      cfg,
		server: &webtransport.Server{
			H3: h3Server,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				log.Println("[WT] Origin:", origin)

				return origin == cfg.FrontendURL
			},
		},
	}
}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/wt", s.handleUpgrade)

	s.server.H3.Addr = addr
	s.server.H3.Handler = mux

	log.Printf("WebTransport HTTP/3 server listening on %s", addr)

	return s.server.ListenAndServe()
}

func (s *Server) handleUpgrade(w http.ResponseWriter, r *http.Request) {
	lotIDStr := r.URL.Query().Get("lotId")
	if lotIDStr == "" {
		log.Println("WebTransport reject: missing lotId")
		http.Error(w, "missing lotId", http.StatusBadRequest)
		return
	}

	lotID64, err := strconv.ParseUint(lotIDStr, 10, 64)
	if err != nil {
		log.Println("WebTransport reject: invalid lotId:", lotIDStr)
		http.Error(w, "invalid lotId", http.StatusBadRequest)
		return
	}

	lotID := uint64(lotID64)

	sess, err := s.server.Upgrade(w, r)
	if err != nil {
		log.Println("WebTransport upgrade failed:", err)
		return
	}

	go func(sess *webtransport.Session) {
		ws := &wtSession{session: sess}

		client := &Client{
			LotID:   lotID,
			Session: ws,
		}

		s.hub.Add(client)

		defer func() {
			s.hub.Remove(ws)
			_ = sess.CloseWithError(0, "")
		}()

		for {
			stream, err := sess.AcceptStream(context.Background())
			if err != nil {
				log.Println("WebTransport accept stream stopped:", err)
				return
			}

			_, _ = io.Copy(io.Discard, stream)
			_ = stream.Close()
		}
	}(sess)
}
