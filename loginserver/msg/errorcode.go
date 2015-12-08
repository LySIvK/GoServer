package msg

const (
	RE_SUCCESS             = iota //! 成功返回码
	RE_UNKNOW_ERR                 //! 未知错误
	RE_ACCOUNT_EXIST              //! 账户已存在
	RE_ACCOUNT_NOT_EXIST          //! 账户不存在
	RE_INVALID_ACCOUNTNAME        //! 无效账户名
	RE_INVALID_PASSWORD           //! 无效密码
	RE_NOT_LOGIN                  //! 用户未登录
)
