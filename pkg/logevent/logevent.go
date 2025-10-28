package logevent

type Event string

func (e Event) String() string {
	return string(e)
}

/*
	日志事件命名规范：
	[Scope]_[Action]_[Outcome]
	Scope: 哪个领域或模块 (e.g., SERVER, USER, DB)
	Action: 发生了什么动作 (e.g., START, STOP, CREATE, LOGIN)
	Outcome: 动作的结果 (e.g., SUCCESS, FAIL)
*/
// --- 通用/服务器事件 ---
const (
	EventServerStart        Event = "server_start"
	EventServerShutdown     Event = "server_shutdown"
	EventPanicRecovered     Event = "panic_recovered"
	EventDBConnectionFailed Event = "db_connection_failed"
	EventGRPCRequest        Event = "grpc_request"
	EventHTTPRequest        Event = "http_request"
	EventInternerError      Event = "internal_error"
)

// --- 用户服务事件 ---
const (
	// Info Level
	EventUserCreated Event = "user_create_success"
	EventUserLogin   Event = "user_login_success"

	// Warn Level
	EventUserLoginPasswordIncorrect Event = "user_login_fail_password_incorrect"
	EventUserLoginNotFound          Event = "user_login_fail_user_not_found"

	// Error Level
	EventDBUserQueryFailed  Event = "db_user_query_failed"
	EventDBUserCreateFailed Event = "db_user_create_failed"
)
