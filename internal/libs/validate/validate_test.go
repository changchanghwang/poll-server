package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"poll.ant/internal/libs/types"
	"poll.ant/internal/libs/validate"
)

func TestValidate(t *testing.T) {
	validate.Init()
	t.Run("ValidateDto 테스트", func(t *testing.T) {
		t.Run("required 테스트", func(t *testing.T) {
			t.Run("없다면 에러를 반환한다.", func(t *testing.T) {
				type RequiredTest struct {
					Required string `validate:"required"`
				}

				requiredTest := RequiredTest{}

				err := validate.ValidateDto(requiredTest)
				assert.NotNil(t, err)
				assert.Equal(t, "Required는(은) 필수 값입니다.", err.Error())
			})
		})
		t.Run("email 테스트", func(t *testing.T) {
			t.Run("이메일 형식이 아니라면 에러를 반환한다.", func(t *testing.T) {
				type EmailTest struct {
					Email string `validate:"email"`
				}

				emailTest := EmailTest{
					Email: "email",
				}

				err := validate.ValidateDto(emailTest)
				assert.NotNil(t, err)
				assert.Equal(t, "Email는(은) 반드시 이메일 형식이어야 합니다.", err.Error())
			})
		})
		t.Run("gte 테스트", func(t *testing.T) {
			t.Run("값이 작다면 에러를 반환한다.", func(t *testing.T) {
				type GteTest struct {
					Gte int `validate:"gte=10"`
				}

				gteTest := GteTest{
					Gte: 5,
				}

				err := validate.ValidateDto(gteTest)
				assert.NotNil(t, err)
				assert.Equal(t, "Gte는(은) 반드시 10보다 크거나 같아야 합니다.", err.Error())
			})
		})
		t.Run("lte 테스트", func(t *testing.T) {
			t.Run("값이 크다면 에러를 반환한다.", func(t *testing.T) {
				type LteTest struct {
					Lte int `validate:"lte=10"`
				}

				lteTest := LteTest{
					Lte: 15,
				}

				err := validate.ValidateDto(lteTest)
				assert.NotNil(t, err)
				assert.Equal(t, "Lte는(은) 반드시 10보다 작거나 같아야 합니다.", err.Error())
			})
		})
		t.Run("유효하지 않은 값이라면 에러를 반환한다.", func(t *testing.T) {
			type InvalidTest struct {
				Max int `validate:"max=1"`
			}

			invalidTest := InvalidTest{
				Max: 2,
			}

			err := validate.ValidateDto(invalidTest)
			assert.NotNil(t, err)
			assert.Equal(t, "Max는(은) 유효하지 않은 값입니다.", err.Error())
		})
		t.Run("message가 여러개일 경우 \n으로 구분하여 반환한다.", func(t *testing.T) {
			type MultiMessageTest struct {
				Required string `validate:"required"`
				Email    string `validate:"email"`
			}

			multiMessageTest := MultiMessageTest{}

			err := validate.ValidateDto(multiMessageTest)
			assert.NotNil(t, err)
			assert.Equal(t, "Required는(은) 필수 값입니다.\nEmail는(은) 반드시 이메일 형식이어야 합니다.", err.Error())
		})
		t.Run("oneof 테스트", func(t *testing.T) {
			t.Run("허용되지 않는 값이라면 에러를 반환한다.", func(t *testing.T) {
				type OneOfTest struct {
					OneOf string `validate:"oneof=google kakao naver"`
				}

				oneOfTest := OneOfTest{
					OneOf: "facebook",
				}

				err := validate.ValidateDto(oneOfTest)
				assert.NotNil(t, err)
				assert.Equal(t, "OneOf는(은) 반드시 'google', 'kakao', 'naver' 중 하나여야 합니다.", err.Error())
			})
		})

		t.Run("calendardate 테스트", func(t *testing.T) {
			t.Run("YYYY-MM-DD 포맷이 아니라면 에러를 반환한다.", func(t *testing.T) {
				type calendarDateType struct {
					Date types.CalendarDate `validate:"calendardate"`
				}

				dto := calendarDateType{
					Date: types.CalendarDate("2024-09-01T00:00:00.000Z"),
				}

				err := validate.ValidateDto(dto)
				assert.NotNil(t, err)
				assert.Equal(t, "Date는(은) 반드시 YYYY-MM-DD 형식이어야 합니다.", err.Error())
			})
		})
	})
}
