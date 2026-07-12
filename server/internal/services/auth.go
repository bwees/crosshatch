package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"crosshatch/internal/database/models"
	"crosshatch/internal/dtos"
	"crosshatch/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

const SessionDuration = 30 * 24 * time.Hour

// statusError is a domain error that carries the HTTP status it maps to, so the
// transport layer can return it as-is instead of translating each case.
type statusError struct {
	status int
	msg    string
}

func (e statusError) Error() string     { return e.msg }
func (e statusError) StatusCode() int   { return e.status }
func (e statusError) DetailMsg() string { return e.msg }

var (
	ErrInvalidCredentials = statusError{http.StatusUnauthorized, "invalid username or password"}
	ErrUsernameTaken      = statusError{http.StatusConflict, "username already taken"}
	ErrSetupComplete      = statusError{http.StatusConflict, "initial setup already complete"}
	ErrForbidden          = statusError{http.StatusForbidden, "admin privileges required"}
)

type AuthService struct {
	users    *repositories.UserRepository
	sessions *repositories.SessionRepository
}

func NewAuthService(users *repositories.UserRepository, sessions *repositories.SessionRepository) *AuthService {
	return &AuthService{users: users, sessions: sessions}
}

func (s *AuthService) SetupRequired() (bool, error) {
	count, err := s.users.CountUsers()
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (s *AuthService) Setup(dto dtos.CreateUserDto) (*models.User, error) {
	required, err := s.SetupRequired()
	if err != nil {
		return nil, err
	}
	if !required {
		return nil, ErrSetupComplete
	}

	return s.createUser(dto.Username, dto.Password, true)
}

func (s *AuthService) CreateUser(dto dtos.CreateUserDto) (*models.User, error) {
	return s.createUser(dto.Username, dto.Password, dto.IsAdmin)
}

func (s *AuthService) createUser(username, password string, isAdmin bool) (*models.User, error) {
	existing, err := s.users.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{Username: username, PasswordHash: string(hash), IsAdmin: isAdmin}
	if err := s.users.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) ListUsers() ([]models.User, error) {
	return s.users.GetUsers()
}

func (s *AuthService) DeleteUser(id uint) error {
	return s.users.DeleteUser(id)
}

func (s *AuthService) Login(dto dtos.LoginDto) (*models.User, string, error) {
	user, err := s.users.GetUserByUsername(dto.Username)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := generateToken()
	if err != nil {
		return nil, "", err
	}

	session := &models.Session{
		TokenHash: hashToken(token),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(SessionDuration),
	}
	if err := s.sessions.CreateSession(session); err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Logout(token string) error {
	if token == "" {
		return nil
	}
	return s.sessions.DeleteSession(hashToken(token))
}

func (s *AuthService) Authenticate(token string) (*models.User, error) {
	if token == "" {
		return nil, ErrInvalidCredentials
	}

	session, err := s.sessions.GetSession(hashToken(token))
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrInvalidCredentials
	}
	if time.Now().After(session.ExpiresAt) {
		_ = s.sessions.DeleteSession(session.TokenHash)
		return nil, ErrInvalidCredentials
	}

	user, err := s.users.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
