package demodto

type UserDetailReq struct {
	Id int64 `json:"id" validate:"required"`
}

type UserListReq struct {
	Q        string `json:"q"`
	PageSize int    `json:"page_size"`
	PageNum  int    `json:"page_num"`
}

type UserListResp struct {
	List  []*User `json:"list"`
	Total int64   `json:"total"`
}

type User struct {
	CreatedAt string `json:"created_at"`
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Status    int    `json:"status"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
}

type UserAddReq struct {
	Phone    string `json:"phone" form:"phone"`
	Username string `json:"username" form:"username" validate:"required"`
	Avatar   string `json:"avatar" form:"avatar"`
	Name     string `json:"name" validate:"required" form:"name"`
}

type UserAddResp struct {
	Id int64 `json:"id"`
}

type UserDeleteReq struct {
	Id int64 `json:"id"`
}

type UserUpdateReq struct {
	Status   int    `json:"status"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}
