package config

import (
	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
)

func FirebaseConfirmEmailVerificationURL(firebaseAuthEmulatorHost string) string {
	if firebaseAuthEmulatorHost == "" {
		return "https://" + constant.FirebaseConfirmEmailVerificationEndpoint
	}
	return "http://" + firebaseAuthEmulatorHost + "/" + constant.FirebaseConfirmEmailVerificationEndpoint
}
