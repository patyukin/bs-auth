package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/bs-auth/internal/model"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"github.com/pquerna/otp/totp"
	"time"
)

type OtpHashRedis struct {
	UserId     int64  `json:"user_id"`
	OtpAuthUrl string `json:"otp_auth_url"`
	OtpSecret  string `json:"otp_secret"`
}

func (s *serv) SignIn(ctx context.Context, user *model.User, fingerprint string) (*desc.SignInResponse, error) {
	totpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "example.com",
		AccountName: "admin@admin.com",
		SecretSize:  15,
	})

	otpHashRedis := OtpHashRedis{
		UserId:     user.ID,
		OtpAuthUrl: totpKey.URL(),
		OtpSecret:  totpKey.Secret(),
	}

	// key for fingerprint
	fpKey := fmt.Sprintf("fingerprint:%s", fingerprint)
	fpValue, err := json.Marshal(otpHashRedis)
	if err != nil {
		return nil, fmt.Errorf("can't marshal fingerprint: %w", err)
	}

	err = s.cacher.Set(ctx, fpKey, string(fpValue), 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("can't set fingerprint in cache: %w", err)
	}

	return &desc.SignInResponse{
		Base32:     totpKey.Secret(),
		OtpAuthUrl: totpKey.URL(),
	}, nil
}
