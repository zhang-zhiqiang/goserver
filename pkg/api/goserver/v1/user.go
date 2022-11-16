package v1

// LoginRequest defines the request fields for `/login`.
type LoginRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
}

// LoginResponse defines the response fields for `/login`.
type LoginResponse struct {
	Token string `json:"token"`
}

// CreateUserRequest is used to store request parameters.
type CreateUserRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
	Nickname string `json:"nickname" valid:"required,stringlength(1|255)"`
	Email    string `json:"email" valid:"required,email"`
	Phone    string `json:"phone" valid:"required,stringlength(11|11)"`
}

type GetUserResponse UserInfo

// ChangePasswordRequest defines the ChangePasswordRequest data format.
type ChangePasswordRequest struct {
	// Old password.
	// Required: true
	OldPassword string `json:"oldPassword" valid:"required,stringlength(6|18)"`

	// New password.
	// Required: true
	NewPassword string `json:"newPassword" valid:"required,stringlength(6|18)"`
}

// UpdateUserRequest specify fields can be updated for user resource.
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" valid:"required,stringlength(1|255)"`
	Email    *string `json:"email" valid:"required,email"`
	Phone    *string `json:"phone" valid:"required,stringlength(11|11)"`
}

type UserInfo struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListUserRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type ListUserResponse struct {
	TotalCount int64 `json:""`
	Users      []*UserInfo
}
