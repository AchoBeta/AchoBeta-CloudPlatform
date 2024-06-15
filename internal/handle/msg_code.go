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

	/* 内部错误 600 ~ 999 */
	INTERNAL_ERROR             = MsgCode{601, "内部错误, check log"}
	INTERNAL_FILE_UPLOAD_ERROR = MsgCode{602, "文件上传失败"}
	/* 参数错误：1000 ~ 1999 */
	PARAM_NOT_VALID    = MsgCode{1001, "参数无效"}
	PARAM_IS_BLANK     = MsgCode{1002, "参数为空"}
	PARAM_TYPE_ERROR   = MsgCode{1003, "参数类型错误"}
	PARAM_NOT_COMPLETE = MsgCode{1004, "参数缺失"}

	PARAM_FILE_SIZE_TOO_BIG = MsgCode{1010, "文件过大"}

	/* 用户错误 2000 ~ 2999 */
	USER_NOT_LOGIN             = MsgCode{2001, "用户未登录"}
	USER_PASSWORD_DIFFERENT    = MsgCode{2002, "用户两次密码输入不一致"}
	USER_ACCOUNT_NOT_EXIST     = MsgCode{2003, "账号不存在"}
	USER_CREDENTIALS_ERROR     = MsgCode{2004, "密码错误"}
	USER_ACCOUNT_ALREADY_EXIST = MsgCode{2008, "账号已存在"}
	CAPTCHA_ERROR              = MsgCode{2100, "验证码错误"}
	INSUFFICENT_PERMISSIONS    = MsgCode{2200, "权限不足"}

	/* 镜像错误 3000 ~ 3999 */
	IMAGE_NOT_FIND    = MsgCode{3001, "镜像未找到"}
	IMAGE_CREATE_FAIL = MsgCode{3002, "镜像构建失败"}
	IMAGE_PULL_FAIL   = MsgCode{3003, "镜像拉取失败"}
	IMAGE_PUSH_FAIL   = MsgCode{3004, "镜像上传失败"}
	IMAGE_REMOVE_FAIL = MsgCode{3005, "镜像删除失败"}
	IMAGE_EXIST       = MsgCode{3006, "镜像已存在"}

	/* 容器错误 4000 ~ 4999 */
	CONTAINER_NOT_FOUND    = MsgCode{4001, "容器未找到"}
	CONTAINER_CREATE_FAIL  = MsgCode{4002, "容器创建失败"}
	CONTAINER_START_FAIL   = MsgCode{4003, "容器启动失败"}
	CONTAINER_STOP_FAIL    = MsgCode{4004, "容器停止失败"}
	CONTAINER_RESTART_FAIL = MsgCode{4005, "容器重启失败"}
	CONTAINER_REMOVE_FAIL  = MsgCode{4006, "容器删除失败"}
	CONTAINER_IS_DESTORY   = MsgCode{4007, "容器已销毁"}

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
