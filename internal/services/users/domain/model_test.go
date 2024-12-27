package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"poll.ant/internal/services/users/domain"
)

func TestUser(t *testing.T) {
	t.Run("New 테스트", func(t *testing.T) {
		t.Run("email, name, provider를 받아 User 객체를 만든다", func(t *testing.T) {
			user, _ := domain.New("email", "name", "kakao")
			assert.Equal(t, user.Email, "email")
			assert.Equal(t, user.Name, "name")
			assert.Equal(t, user.Provider, "kakao")
		})
	})

	t.Run("Update 테스트", func(t *testing.T) {
		t.Run("Name을 변경할 수 있다.", func(t *testing.T) {
			user, _ := domain.New("email", "name", "kakao")
			name := "new Name"
			user.Update(domain.UpdateType{Name: &name})
			assert.Equal(t, user.Name, "new Name")
		})
	})
}
