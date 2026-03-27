package auth

import "github.com/gin-gonic/gin"

type CurrentUser struct {
	UserID      int      `json:"user_id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	RoleID      int      `json:"role_id"`
	RoleCode    string   `json:"role_code"`
	Permissions []string `json:"permissions"`
}

const currentUserKey = "current_user"

func SetCurrentUser(c *gin.Context, user CurrentUser) {
	c.Set(currentUserKey, user)
}

func GetCurrentUser(c *gin.Context) (CurrentUser, bool) {
	value, ok := c.Get(currentUserKey)
	if !ok {
		return CurrentUser{}, false
	}
	user, ok := value.(CurrentUser)
	return user, ok
}

func HasPermission(user CurrentUser, permission string) bool {
	if permission == "" {
		return true
	}
	for _, item := range user.Permissions {
		if item == "*" || item == permission {
			return true
		}
	}
	return false
}
