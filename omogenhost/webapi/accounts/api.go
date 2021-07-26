package accounts

import (
	"context"
	"encoding/base64"
	"github.com/google/logger"
	"github.com/jsannemo/omogenhost/storage"
	apipb "github.com/jsannemo/omogenhost/webapi/proto"
	"github.com/jsannemo/omogenhost/webapi/requests"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"net/mail"
	"strings"
)

type accountService struct {
}

func InitAccountService() *accountService {
	return &accountService{}
}

func validUsername(username string) bool {
	length := len(username)
	if length < 3 || length > 25 {
		return false
	}
	for _, ch := range username {
		if '0' <= ch && ch <= '9' {
			continue
		}
		if 'a' <= ch && ch <= 'z' {
			continue
		}
		if 'A' <= ch && ch <= 'Z' {
			continue
		}
		if ch == '_' || ch == '-' || ch == '.' {
			continue
		}
		return false
	}
	return true
}

func (as *accountService) Register(ctx context.Context, request *apipb.RegisterRequest) (*apipb.RegisterResponse, error) {
	logger.Infof("AccountService.Register: %v", request)
	var errs []apipb.RegisterResponse_RegisterError
	if !validUsername(request.Username) {
		errs = append(errs, apipb.RegisterResponse_USERNAME_INVALID)
	}
	if _, err := mail.ParseAddress(request.Email); err != nil {
		logger.Infof("Invalid email: %v", err)
		errs = append(errs, apipb.RegisterResponse_EMAIL_INVALID)
	}
	request.FullName = strings.Trim(request.FullName, " ")
	if len(request.FullName) < 3 {
		errs = append(errs, apipb.RegisterResponse_FULL_NAME_INVALID)
	}
	if errs == nil {
		var users []storage.Account
		if res := storage.GormDB.Debug().Where("email = ? OR lower(username) = ?", request.Email, strings.ToLower(request.Username)).Find(&users); res.Error != nil {
			logger.Errorf("Failed looking up existing users: %v", res.Error)
			return nil, status.Error(codes.Internal, "")
		}
		for _, usr := range users {
			if usr.Email == request.Email {
				errs = append(errs, apipb.RegisterResponse_EMAIL_TAKEN)
			}
			if strings.EqualFold(usr.Username, request.Username) {
				errs = append(errs, apipb.RegisterResponse_USERNAME_TAKEN)
			}
		}
	}
	if errs != nil {
		return &apipb.RegisterResponse{
			Errors: errs,
		}, nil
	}
	acc := storage.Account{
		Username:    request.Username,
		FullName:    request.FullName,
		Email:       request.Email,
		Password:    passwordHash(request.Password),
		IsStaff:     false,
		IsSuperuser: false,
	}
	if res := storage.GormDB.Create(&acc); res.Error != nil {
		logger.Errorf("Failed creating user: %v", res.Error)
		return nil, status.Error(codes.Internal, "")
	}
	requests.GetUser(ctx).UserId = acc.AccountId
	return &apipb.RegisterResponse{}, nil
}

func passwordHash(pw string) string {
	if hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost); err != nil {
		panic(err)
	} else {
		return base64.StdEncoding.EncodeToString(hash)
	}
}

func (as *accountService) Login(ctx context.Context, request *apipb.LoginRequest) (*apipb.LoginResponse, error) {
	user := storage.Account{}
	if res := storage.GormDB.Debug().Where("lower(username) = ?", strings.ToLower(request.Username)).First(&user); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return &apipb.LoginResponse{
				Errors: []apipb.LoginResponse_LoginError{apipb.LoginResponse_INVALID_CREDENTIALS},
			}, nil
		}
		logger.Errorf("Failed looking up user for login: %v", res.Error)
		return nil, status.Error(codes.Internal, "")
	}
	if !comparePassword(request.Password, user.Password) {
		return &apipb.LoginResponse{
			Errors: []apipb.LoginResponse_LoginError{apipb.LoginResponse_INVALID_CREDENTIALS},
		}, nil
	}
	requests.GetUser(ctx).UserId = user.AccountId
	return &apipb.LoginResponse{}, nil
}

func comparePassword(pw string, hash string) bool {
	strHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		logger.Errorf("Failed decoding password hash: %v", err)
		return false
	}
	if err := bcrypt.CompareHashAndPassword(strHash, []byte(pw)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			logger.Errorf("Failed comparing password hash: %v", err)
		}
		return false
	}
	return true
}
