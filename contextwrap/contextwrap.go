package contextwrap

import (
	"context"
	"net/http"
	"time"

	"github.com/xanderxampp-be/franco/dto/response"
	"github.com/xanderxampp-be/franco/log/entity"
)

func GetFinancialFlagFromContext(ctx context.Context) bool {
	lr := ctx.Value(IsFinancialKey)
	if l, ok := lr.(bool); ok {
		return l
	} else {
		return false
	}
}

func GetTrxTypeFromContext(ctx context.Context) string {
	lr := ctx.Value(TrxTypeKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func SetTrxTypeFromContext(ctx context.Context, trxType string) context.Context {
	ctx = context.WithValue(ctx, TrxTypeKey, trxType)
	return ctx
}

func GetTrxObjectFromContext(ctx context.Context) string {
	lr := ctx.Value(TrxObjectKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func SetTrxObjectFromContext(ctx context.Context, trxObject string) context.Context {
	ctx = context.WithValue(ctx, TrxObjectKey, trxObject)
	return ctx
}

func GetAccountDebetFromContext(ctx context.Context) string {
	lr := ctx.Value(AccountDebetKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func SetAccountDebetFromContext(ctx context.Context, accountDebet string) context.Context {
	ctx = context.WithValue(ctx, AccountDebetKey, accountDebet)
	return ctx
}

func GetAmountFromContext(ctx context.Context) int {
	lr := ctx.Value(AmountKey)
	if l, ok := lr.(int); ok {
		return l
	} else {
		return 0
	}
}

func SetAmountFromContext(ctx context.Context, amount int) context.Context {
	ctx = context.WithValue(ctx, AmountKey, amount)
	return ctx
}

func GetAmountFloatFromContext(ctx context.Context) float64 {
	lr := ctx.Value(AmountFloatKey)
	if l, ok := lr.(float64); ok {
		return l
	} else {
		return 0
	}
}

func SetAmountFloatFromContext(ctx context.Context, amountFloat float64) context.Context {
	ctx = context.WithValue(ctx, AmountFloatKey, amountFloat)
	return ctx
}

func GetFeeFromContext(ctx context.Context) int {
	lr := ctx.Value(FeeKey)
	if l, ok := lr.(int); ok {
		return l
	} else {
		return 0
	}
}

func SetFeeFromContext(ctx context.Context, fee int) context.Context {
	ctx = context.WithValue(ctx, FeeKey, fee)
	return ctx
}

func GetIpAddressSourceFromContext(ctx context.Context) string {
	lr := ctx.Value(IpAddressSourceKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func GetAgentFromContext(ctx context.Context) string {
	lr := ctx.Value(AgentKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func GetLogResponse(r *http.Request) *entity.Responselog {
	lr := r.Context().Value(LogRespKey)
	if l, ok := lr.(*entity.Responselog); ok {
		return l
	} else {
		return &entity.Responselog{}
	}
}

func GetLogResponseFromContext(ctx context.Context) *entity.Responselog {
	lr := ctx.Value(LogRespKey)
	if lr == nil {
		return &entity.Responselog{}
	}
	if l, ok := lr.(*entity.Responselog); ok {
		return l
	} else {
		return &entity.Responselog{}
	}
}

func GetTraceFromContext(ctx context.Context) []interface{} {
	lr := ctx.Value(TraceKey)
	if l, ok := lr.([]interface{}); ok {
		return l
	} else {
		return []interface{}{}
	}
}

func SetTraceFromContext(ctx context.Context, trace []interface{}) context.Context {
	ctx = context.WithValue(ctx, TraceKey, trace)
	return ctx
}

func GetProcessID(r *http.Request) string {
	lr := r.Context().Value(ProcessIDKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func GetProcessIDFromContext(ctx context.Context) string {
	lr := ctx.Value(ProcessIDKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func GetBody(r *http.Request) []byte {
	lr := r.Context().Value(BodyKey)
	if l, ok := lr.([]byte); ok {
		return l
	} else {
		return []byte("")
	}
}

func GetBodyFromContext(ctx context.Context) []byte {
	lr := ctx.Value(BodyKey)
	if l, ok := lr.([]byte); ok {
		return l
	} else {
		return []byte("")
	}
}

func GetStartFromContext(ctx context.Context) time.Time {
	lr := ctx.Value(ElapsedKey)
	if l, ok := lr.(time.Time); ok {
		return l
	} else {
		return time.Time{}
	}
}

func GetResponseFromContext(ctx context.Context) *response.Response {
	lr := ctx.Value(RespKey)
	if lr == nil {
		return &response.Response{}
	}
	if l, ok := lr.(*response.Response); ok {
		return l
	} else {
		return &response.Response{}
	}
}

func GetDeviceFromContext(ctx context.Context) *entity.Device {
	lr := ctx.Value(DeviceKey)
	if lr == nil {
		return &entity.Device{}
	}
	if l, ok := lr.(*entity.Device); ok {
		return l
	} else {
		return &entity.Device{}
	}
}

func GetThirdPartyFromContext(ctx context.Context) string {
	lr := ctx.Value(ThirdPartyKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func SetThirdPartyFromContext(ctx context.Context, thirdParty string) context.Context {
	ctx = context.WithValue(ctx, ThirdPartyKey, thirdParty)
	return ctx
}
