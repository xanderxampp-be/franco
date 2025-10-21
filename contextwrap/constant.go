package contextwrap

type elapsedCtxKey string
type trxTypeCtxKey string
type ipAddressSourceCtxKey string
type agentCtxKey string
type trxObjectCtxKey string
type deviceCtxKey string
type isFinancialCtxKey string
type ctxKey string
type processIDCtxKey string
type traceCtxKey string
type accountDebetCtxKey string
type amountCtxKey string
type amountFloatCtxKey string
type feeCtxKey string
type respCtxKey string
type logRespCtxKey string
type thirdPartyCtxKey string

const ElapsedKey elapsedCtxKey = "elapsed"
const IpAddressSourceKey ipAddressSourceCtxKey = "ipAddressSource"
const AgentKey agentCtxKey = "agent"
const TrxTypeKey trxTypeCtxKey = "trxType"
const TrxObjectKey trxObjectCtxKey = "trxObject"
const DeviceKey deviceCtxKey = "device"
const IsFinancialKey deviceCtxKey = "isFinancial"
const BodyKey ctxKey = "body"
const ProcessIDKey processIDCtxKey = "procId"
const TraceKey traceCtxKey = "trace"
const AccountDebetKey accountDebetCtxKey = "accountDebet"
const AmountKey amountCtxKey = "amount"
const AmountFloatKey amountFloatCtxKey = "amountFloat"
const FeeKey feeCtxKey = "fee"
const RespKey respCtxKey = "resp"
const LogRespKey logRespCtxKey = "logResp"
const ThirdPartyKey thirdPartyCtxKey = "thirdParty"
