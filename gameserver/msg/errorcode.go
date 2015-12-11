package msg

const (
	RE_SUCCESS             = iota //! 成功
	RE_UNKNOW_ERR                 //! 未知错误
	RE_ACCOUNT_EXIST              //! 账户已存在
	RE_ACCOUNT_NOT_EXIST          //! 账户不存在
	RE_INVALID_ACCOUNTNAME        //! 无效账户名
	RE_INVALID_PASSWORD           //! 无效密码
	RE_INVALID_ACCOUNTID          //! 无效账户ID
	RE_INVALID_LOGINKEY           //! 无效登录Key
	RE_INVALID_PLAYERNAME         //! 无效用户名
	RE_NOT_LOGIN                  //! 用户未登录
	RE_ROLE_NOT_EXIST             //! 玩家尚未创建角色
	RE_ROLE_EXIST                 //! 玩家已创建角色
)
