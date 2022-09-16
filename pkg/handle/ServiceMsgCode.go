package handle

type MsgCode struct {
	Code int
	Msg  string
}

var (
	/* 成功 */
	SUCCESS = MsgCode{Code: 200, Msg: "成功"}

	/* 默认失败 */
	COMMON_FAIL = MsgCode{-4396, "失败"}

	/* 请求错误 <0 */
	TOKEN_IS_EXPIRED = MsgCode{-2, "token已过期"}

	/* 内部错误 600 ~ 999*/
	INTERNAL_ERROR = MsgCode{601, "内部错误, check log"}

	/* 参数错误：1000 ~ 1999 */
	PARAM_NOT_VALID    = MsgCode{1001, "参数无效"}
	PARAM_IS_BLANK     = MsgCode{1002, "参数为空"}
	PARAM_TYPE_ERROR   = MsgCode{1003, "参数类型错误"}
	PARAM_NOT_COMPLETE = MsgCode{1004, "参数缺失"}

	/* 用户错误 2000 ~ 2999*/
	USER_NOT_LOGIN          = MsgCode{2001, "用户未登录"}
	USER_PASSWORD_DIFFERENT = MsgCode{2002, "用户两次密码输入不一致"}
	USER_CREDENTIALS_ERROR  = MsgCode{2003, "密码错误"}
	/*
	 USER_NOT_LOGIN(2001, "用户未登录"),
	 USER_ACCOUNT_EXPIRED(2002, "账号已过期"),

	 USER_CREDENTIALS_EXPIRED(2004, "密码过期"),
	 USER_ACCOUNT_DISABLE(2005, "账号不可用"),
	 USER_ACCOUNT_LOCKED(2006, "账号被锁定"),
	 USER_ACCOUNT_NOT_EXIST(2007, "账号不存在"),
	 USER_ACCOUNT_ALREADY_EXIST(2008, "账号已存在"),
	 USER_ACCOUNT_USE_BY_OTHERS(2009, "账号下线"),
	 USER_NO_PERMISSION(403,"用户无权限"),
	 USER_NO_PHONE_CODE(500,"验证码错误"),
	*/
)
