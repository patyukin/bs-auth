package auth

import (
	"context"
	"encoding/json"
	"fmt"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"github.com/pquerna/otp/totp"
)

func (s *serv) CheckCode(ctx context.Context, req *desc.CheckCodeRequest) (int64, error) {
	key := req.GetFingerprint()
	fpKey := fmt.Sprintf("fingerprint:%s", key)
	value, err := s.cacher.Get(ctx, fpKey)
	if err != nil {
		return 0, fmt.Errorf("can't get fingerprint from cache: %w", err)
	}

	var fpValue OtpHashRedis
	err = json.Unmarshal([]byte(value), &fpValue)
	if err != nil {
		return 0, fmt.Errorf("can't unmarshal fingerprint: %w", err)
	}

	valid := totp.Validate(req.Code, fpValue.OtpSecret)
	if !valid {
		return 0, fmt.Errorf("invalid code")
	}

	err = s.cacher.Delete(ctx, fpKey)
	if err != nil {
		return 0, fmt.Errorf("can't clear cache: %w", err)
	}

	return fpValue.UserId, nil
}
