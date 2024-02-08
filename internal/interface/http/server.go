package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/elinyaa/test-signer/internal"
	"github.com/elinyaa/test-signer/internal/domain/entity"
	"github.com/golang-jwt/jwt"
)

type Server struct {
	app    internal.App
	logger *log.Logger
	server *http.Server
	secret []byte
}

func NewServer(app internal.App, logger *log.Logger, addr string, secret string) *Server {
	server := &Server{
		app:    app,
		logger: logger,
		server: &http.Server{
			Addr: addr,
		},
		secret: []byte(secret),
	}
	handler := server.getHandler()
	server.server.Handler = handler
	return server
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Println("Shutting down...")
		err := s.server.Shutdown(context.Background())
		if err != nil {
			s.logger.Fatalf("Error shutting down: %v", err)
		}
	}()

	s.logger.Printf("Starting server on port %v\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) getHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/sign", s.handleSign)
	mux.HandleFunc("/verify", s.handleVerify)
	return mux
}

func (s *Server) handleSign(w http.ResponseWriter, r *http.Request) {
	userId, name, err := s.extractUserIdAndNameFromClaims(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	type request struct {
		Questions string `json:"questions"`
		Answers   string `json:"answers"`
	}
	var req request
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := s.app.SignAnswer(r.Context(), userId, name, req.Questions, req.Answers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	w.Write([]byte(resp))
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserId    int    `json:"userid"`
		Username  string `json:"username"`
		Signature string `json:"signature"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)

	ok, answers, t, err := s.app.VerifySignature(r.Context(), entity.User{ID: req.UserId, Username: req.Username}, req.Signature)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Signature not verified")
		return
	}
	w.WriteHeader(http.StatusOK)

	type response struct {
		Answers string    `json:"answers"`
		Time    time.Time `json:"time"`
	}
	resp := response{
		Answers: answers,
		Time:    t,
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) extractUserIdAndNameFromClaims(_ http.ResponseWriter, request *http.Request) (int, string, error) {
	if request.Header["Token"] != nil {
		tokenString := request.Header["Token"][0]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return s.secret, nil
		})
		if err != nil {
			return 0, "", fmt.Errorf("Error parsing token: %v", err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			userIdString := claims["userid"].(string)
			name := claims["name"].(string)

			userId, err := strconv.Atoi(userIdString)
			if err != nil {
				return 0, "", fmt.Errorf("Error parsing userid: %v", err)
			}
			return userId, name, nil
		}

		return 0, "", fmt.Errorf("Token not valid")
	}

	return 0, "", fmt.Errorf("Token not found")
}
