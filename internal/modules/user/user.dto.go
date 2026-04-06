package user

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      Role   `json:"role" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required,min=6,max=30"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=30"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type ChangeProfileRequest struct {
	FirstName *string `json:"first_name" binding:"required"`
	LastName  *string `json:"last_name" binding:"required"`
}

type ChangeRoleRequest struct {
	NewRole Role `json:"new_role" binding:"required"`
}

type UserResponse struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
}

type MyAccountResponse struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
	Money     int64  `json:"money"`
}

type UserPaginationResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}
