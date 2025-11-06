package resp

// 业务错误码
const ()

type AppError struct {
	Code    int    // 错误码
	Message string // 错误信息
}

func (err AppError) Error() string {
	return err.Message
}

var (
	ErrBadRequest                  = AppError{400, "invalid or missing parameter"}              // 请求参数错误
	ErrUnauthorized                = AppError{401, "unauthorized access"}                       // 未授权
	ErrForbidden                   = AppError{403, "forbidden"}                                 // 没有权限
	ErrNotFound                    = AppError{404, "the requested resource could not be found"} // 资源未找到
	ErrMethodNotAllowed            = AppError{405, ""}                                          // 不支持的请求方法
	ErrConflict                    = AppError{409, "conflict"}                                  // 冲突
	ErrStatusUnprocessableEntity   = AppError{422, "validation Error"}                          // 参数校验失败
	ErrTooManyRequests             = AppError{429, "too Many Requests"}                         // 请求过于频繁
	ErrServerError                 = AppError{500, "internal Server Error"}                     // 服务器内部错误
	ErrGatewayTimeout              = AppError{504, "request Timeout"}                           // 请求超时
	ErrUsernameAlreadyExists       = AppError{1001, "username already exists"}                  // 用户名已存在
	ErrIncorrectUsernameOrPassword = AppError{1002, "incorrect username or password"}           // 用户名或密码错误
	ErrUserAlreadyHasFamily        = AppError{1003, "user already has a family"}                // 用户已存在家庭
	ErrUserHasBeenModified         = AppError{1004, "user has been modified"}                   // 用户已被修改
	ErrInvalidOriginalPassword     = AppError{1005, "invalid original password"}                // 原密码不正确
	ErrUserNotFound                = AppError{1006, "user not found"}                           // 用户未找到
	ErrPhoneNumberAlreadyExists    = AppError{1007, "phone number already exists"}              // 手机号已存在
	ErrSessionBlocked              = AppError{2001, "blocked session"}                          // 会话被阻止
	ErrSessionUserIDMismatch       = AppError{2002, "incorrect session user"}                   // 会话用户ID不匹配
	ErrSessionExpired              = AppError{2003, "expired session"}                          // 会话已过期
	ErrSessionNotFound             = AppError{2004, "session not found"}                        // 会话未找到
	ErrTokenExpired                = AppError{2005, "token expired"}                            // 令牌已过期
	ErrTokenInvalid                = AppError{2006, "invalid token"}                            // 令牌无效
	ErrFamilyNotFound              = AppError{3001, "family not found"}                         // 家庭未找到
	ErrFamilyHasBeenModified       = AppError{3002, "family has been modified"}                 // 家庭已被修改
	ErrFamilyHasMembers            = AppError{3003, "family still has members"}                 // 家庭仍有成员
	ErrFamilyNotOwner              = AppError{3004, "user is not the owner of the family"}      // 用户不是家庭的所有者
	ErrFamilyAddMemberFailed       = AppError{3005, "failed to add family member"}              // 添加家庭成员失败
	ErrFamilyAlreadyJoined         = AppError{3006, "user has already joined the family"}       // 用户已加入家庭
	ErrInvitationNotFound          = AppError{4001, "family invitation not found"}              // 家庭邀请未找到
	ErrInvitationExpired           = AppError{4002, "family invitation expired"}                // 家庭邀请已过期
	ErrInviteYourself              = AppError{4003, "you cannot invite yourself to the family"} // 不能邀请自己加入家庭
	ErrInvitationHasBeenHandled    = AppError{4004, "family invitation has been handled"}       // 家庭邀请已被处理
	ErrInvitationCanceled          = AppError{4005, "family invitation has been canceled"}      // 家庭邀请已被取消
	ErrMenuAlreadyExists           = AppError{5001, "menu already exists"}                      // 菜单已存在
	ErrMenuNotFound                = AppError{5002, "menu not found"}                           // 菜单未找到
	ErrMenuCategoryNotFound        = AppError{5003, "menu category not found"}                  // 菜单分类未找到
	ErrMenuHasBeenModified         = AppError{5004, "menu has been modified"}                   // 菜单已被修改
)
