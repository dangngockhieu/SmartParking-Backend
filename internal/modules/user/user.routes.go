package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware, adminOnly gin.HandlerFunc) {
	group := api.Group("/users")
	group.Use(authMiddleware)
	{
		// Lấy danh sách người dùng với phân trang
		group.GET("/user-pagination", adminOnly, handler.FindWithPagination)

		// Lấy thông tin số dư trong ví hiện tại của người dùng
		group.GET("/my-wallet", handler.GetMyWallet)

		// Lấy thông tin người dùng hiện tại
		group.GET("/my-info", handler.GetMyInfo)

		// Tạo người dùng mới (chỉ admin)
		group.POST("/admin/create", adminOnly, handler.CreateByAdmin)

		// Cập nhật thông tin người dùng
		group.PATCH("/change-password", handler.ChangePassword)

		// Thay đổi vai trò người dùng
		group.PATCH("/change-role/:id", adminOnly, handler.ChangeRole)

		// Cập nhật thông tin cá nhân
		group.PATCH("/change-profile", handler.UpdateProfile)
	}
}
