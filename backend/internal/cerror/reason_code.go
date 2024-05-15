package cerror

type ReasonCode string

const (
	RC00000 ReasonCode = "RC00000" // Message: reason code が設定されていないエラーです 対応策については実装を確認し、対応が必要であった場合は適切にコードを設定してください
	RC00001 ReasonCode = "RC00001" // Message: Panicエラーが発生しました 直ちにエンジニアは対応してください
	RC20001 ReasonCode = "RC20001" // Message: Failed to initialize firebase
	RC20002 ReasonCode = "RC20002" // Message: Authorization header is not found
	RC20003 ReasonCode = "RC20003" // Message: Authorization must be Bearer
	RC20004 ReasonCode = "RC20004" // Message: Invalid firebase id token
	RC20005 ReasonCode = "RC20005" // Message: Failed to get user information
	RC20006 ReasonCode = "RC20006" // Message: OrderNumber is not found
	RC20007 ReasonCode = "RC20007" // Message: Custom Claims is not found
	RC20008 ReasonCode = "RC20008" // Message: Unsupported security scheme
)
