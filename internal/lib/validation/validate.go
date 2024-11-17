package validation

import (
	ssov1 "github.com/dm1tl/protos/gen/go/sso"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateLoginData(req *ssov1.LoginRequest) error {
	return validation.ValidateStruct(
		req,
		validateEmail(&req.Email),
		validatePassword(&req.Password),
	)
}

func ValidateRegisterData(req *ssov1.RegisterRequest) error {
	return validation.ValidateStruct(
		req,
		validateEmail(&req.Email),
		validatePassword(&req.Password),
	)
}

func ValidateTokenData(req *ssov1.ValidateTokenRequest) error {
	return validation.ValidateStruct(
		req,
		validateToken(&req.Token),
	)
}

func validateToken(token *string) *validation.FieldRules {
	return validation.Field(token, validation.Required)
}

func validateEmail(email *string) *validation.FieldRules {
	return validation.Field(email, validation.Required, is.Email)
}

func validatePassword(password *string) *validation.FieldRules {
	return validation.Field(password, validation.Required, validation.Length(8, 0))
}

func validateUserId(userId *int64) *validation.FieldRules {
	return validation.Field(userId, validation.Required)
}
