package validator

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/novychok/authasvs/internal/entity"
)

type ValidatorASVS struct {
	v10 *validator.Validate
}

func New() *ValidatorASVS {
	v := &ValidatorASVS{
		v10: validator.New(),
	}

	return v
}

func (v *ValidatorASVS) Validate(
	ctx context.Context,
	userCreate *entity.UserCreate,
) error {

	// Included: 2.1.1 / 2.1.2 / 2.1.4
	err := v.v10.StructCtx(ctx, userCreate)
	if err != nil {
		return fmt.Errorf("password validation failed %v", slog.Any("error", err))
	}

	// pass := asvs_point_2_1_3(userCreate.Password)

	return nil
}

func asvs_point_2_1_3(pw string) string {
	re := regexp.MustCompile(`\s+`)
	pw = re.ReplaceAllString(pw, " ")
	pw = strings.TrimSpace(pw)
	return pw
}
