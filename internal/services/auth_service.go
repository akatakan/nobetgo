package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
	"github.com/akatakan/nobetgo/util"
)

type AuthService struct {
	employeeRepo *repositories.EmployeeRepository
	tokenRepo    *repositories.PasswordResetTokenRepository
	jwtSecret    string
}

func NewAuthService(empRepo *repositories.EmployeeRepository, tokenRepo *repositories.PasswordResetTokenRepository, jwtSecret string) *AuthService {
	return &AuthService{
		employeeRepo: empRepo,
		tokenRepo:    tokenRepo,
		jwtSecret:    jwtSecret,
	}
}

func (s *AuthService) Login(username, password string) (string, string, error) {
	employee, err := s.employeeRepo.GetByUsername(username)
	if err != nil {
		return "", "", errors.New("geçersiz kullanıcı adı veya şifre")
	}

	if !util.CheckPasswordHash(password, employee.PasswordHash) {
		return "", "", errors.New("geçersiz kullanıcı adı veya şifre")
	}

	token, err := util.GenerateToken(employee.ID, employee.Role, s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	return token, employee.Role, nil
}

func (s *AuthService) GenerateResetToken(email string) (string, error) {
	employee, err := s.employeeRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("bu e-posta adresi ile kayıtlı kullanıcı bulunamadı")
	}

	// Generate a secure random token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(b)

	resetToken := &core.PasswordResetToken{
		EmployeeID: employee.ID,
		Token:      tokenStr,
		ExpiresAt:  time.Now().Add(1 * time.Hour), // 1 hour expiry
	}

	if err := s.tokenRepo.Create(resetToken); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *AuthService) ResetPassword(tokenStr, newPassword string) error {
	token, err := s.tokenRepo.GetByToken(tokenStr)
	if err != nil {
		return errors.New("geçersiz veya süresi dolmuş token")
	}

	if !token.IsValid() {
		return errors.New("geçersiz veya süresi dolmuş token")
	}

	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update employee password
	employee, err := s.employeeRepo.GetByID(token.EmployeeID)
	if err != nil {
		return err
	}

	employee.PasswordHash = hashedPassword
	if err := s.employeeRepo.Update(employee); err != nil {
		return err
	}

	// Mark token as used
	if err := s.tokenRepo.MarkAsUsed(token.ID); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ChangePassword(employeeID uint, oldPassword, newPassword string) error {
	employee, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return errors.New("kullanıcı bulunamadı")
	}

	if !util.CheckPasswordHash(oldPassword, employee.PasswordHash) {
		return errors.New("mevcut şifre hatalı")
	}

	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	employee.PasswordHash = hashedPassword
	return s.employeeRepo.Update(employee)
}
