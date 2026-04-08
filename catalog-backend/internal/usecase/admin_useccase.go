package usecase

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAdminNotFound      = errors.New("admin not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrCannotDeleteSelf   = errors.New("cannot delete your own account")
	ErrLastAdmin          = errors.New("cannot delete the last admin")
)

type AdminUseCase struct {
	adminRepo repository.AdminRepository
	jwtSecret []byte
}

func NewAdminUseCase(adminRepo repository.AdminRepository, jwtSecret string) *AdminUseCase {
	return &AdminUseCase{
		adminRepo: adminRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// IsFirstAdmin — проверяет, нет ли ещё ни одного администратора в системе
func (uc *AdminUseCase) IsFirstAdmin(ctx context.Context) (bool, error) {
	count, err := uc.adminRepo.Count(ctx)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// Register — создаёт первого администратора (только если он единственный)
func (uc *AdminUseCase) Register(ctx context.Context, input entity.AdminCreate) (*entity.AdminResponse, string, error) {
	isFirst, err := uc.IsFirstAdmin(ctx)
	if err != nil {
		return nil, "", err
	}
	if !isFirst {
		return nil, "", errors.New("registration is closed: admin already exists")
	}
	return uc.createAdmin(ctx, input)
}

// Login — авторизация, возвращает JWT токен и данные администратора
func (uc *AdminUseCase) Login(ctx context.Context, username, password string) (*entity.AdminResponse, string, error) {
	admin, err := uc.adminRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}
	if admin == nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := uc.generateToken(admin)
	if err != nil {
		return nil, "", err
	}

	resp := admin.ToResponse()
	return &resp, token, nil
}

// CreateAdmin — создаёт нового администратора (вызывается из админ-панели)
func (uc *AdminUseCase) CreateAdmin(ctx context.Context, input entity.AdminCreate) (*entity.AdminResponse, string, error) {
	return uc.createAdmin(ctx, input)
}

func (uc *AdminUseCase) createAdmin(ctx context.Context, input entity.AdminCreate) (*entity.AdminResponse, string, error) {
	// Проверяем уникальность username
	existing, err := uc.adminRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", ErrUsernameTaken
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	admin, err := uc.adminRepo.Create(ctx, &input, string(hash))
	if err != nil {
		return nil, "", err
	}

	token, err := uc.generateToken(admin)
	if err != nil {
		return nil, "", err
	}

	resp := admin.ToResponse()
	return &resp, token, nil
}

// GetAll — список всех администраторов
func (uc *AdminUseCase) GetAll(ctx context.Context) ([]entity.AdminResponse, error) {
	admins, err := uc.adminRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]entity.AdminResponse, len(admins))
	for i, a := range admins {
		result[i] = a.ToResponse()
	}
	return result, nil
}

// UpdateAdmin — обновляет имя, логин или пароль администратора
func (uc *AdminUseCase) UpdateAdmin(ctx context.Context, id int, input entity.AdminUpdate) (*entity.AdminResponse, error) {
	// Проверяем уникальность нового username если меняется
	if input.Username != nil {
		existing, err := uc.adminRepo.GetByUsername(ctx, *input.Username)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrUsernameTaken
		}
	}

	var passwordHash *string
	if input.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		h := string(hash)
		passwordHash = &h
	}

	admin, err := uc.adminRepo.Update(ctx, id, &input, passwordHash)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}

	resp := admin.ToResponse()
	return &resp, nil
}

// DeleteAdmin — удаляет администратора (нельзя удалить себя и последнего)
func (uc *AdminUseCase) DeleteAdmin(ctx context.Context, id, callerID int) error {
	if id == callerID {
		return ErrCannotDeleteSelf
	}

	count, err := uc.adminRepo.Count(ctx)
	if err != nil {
		return err
	}
	if count <= 1 {
		return ErrLastAdmin
	}

	return uc.adminRepo.Delete(ctx, id)
}

// ValidateToken — проверяет JWT и возвращает ID администратора
func (uc *AdminUseCase) ValidateToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return uc.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	idFloat, ok := claims["admin_id"].(float64)
	if !ok {
		return 0, errors.New("invalid token payload")
	}

	return int(idFloat), nil
}

func (uc *AdminUseCase) generateToken(admin *entity.Admin) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": admin.ID,
		"username": admin.Username,
		"name":     admin.Name,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}
