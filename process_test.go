package e365_gateway

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestElmaProc(t *testing.T) {

	s := NewStand(testDefaultStandSettings)
	ctxBg := context.Background()

	type procCtx struct {
		ProcCommon
		Number int `json:"number"`
	}

	t.Run("complex_empty_ctx", func(t *testing.T) {

		settings := Settings{
			Stand:     s,
			Namespace: "goods.goods",
			Code:      "bp1",
		}

		bp := NewProc[EmptyProcCtx](settings)

		t.Run("run_proc", func(t *testing.T) {

			proc, err := bp.Run(ctxBg, EmptyProcCtx{})
			require.NoError(t, err)
			require.Len(t, proc.ID, uuid4Len)

			time.Sleep(time.Second)
			procInfo, err := bp.GetInstanceById(ctxBg, proc.ID)
			require.NoError(t, err)
			require.Equal(t, settings.Code, procInfo.Template.Code)
			require.Equal(t, proc.ID, procInfo.ID)
			require.Contains(t, []string{StateDone, StateExec}, procInfo.State)
		})

	})

	t.Run("complex", func(t *testing.T) {

		ctx := context.Background()
		settings := Settings{
			Stand:     s,
			Namespace: "goods",
			Code:      "bp2",
		}

		bp := NewProc[procCtx](settings)

		t.Run("run_proc", func(t *testing.T) {

			num := 2
			proc, err := bp.Run(ctx, procCtx{
				Number: num,
			})
			require.NoError(t, err)
			require.Len(t, proc.ID, uuid4Len)
			//fmt.Println(proc.ID)
			// wait for bp to get done
			time.Sleep(time.Second * 2)

			proc, err = bp.GetInstanceById(ctx, proc.ID)
			require.NoError(t, err)
			require.Equal(t, StateDone, proc.State)
			require.Equal(t, num*2, proc.Number)

		})

	})

	t.Run("table_tests", func(t *testing.T) {

		t.Run("proc_run", func(t *testing.T) {

			t.Run("failed", func(t *testing.T) {

				type TC[T interface{}] struct {
					name        string
					stand       Stand
					namespace   string
					code        string
					ctx         context.Context
					procCtx     T
					bp          Proc[T]
					errExpected error
					timeout     time.Duration
				}
				testCases := []TC[interface{}]{
					{name: "failed_marshal", stand: s, ctx: ctxBg, procCtx: func() {}, errExpected: ErrEncodeRequestBody, namespace: "goods.goods", code: "bp1"},
					{name: "failed_request_creation", stand: s, ctx: nil, procCtx: nil, errExpected: ErrCreateRequest, namespace: "goods.goods", code: "bp1"},
					{name: "failed_request_sending", stand: s, ctx: ctxBg, procCtx: nil, errExpected: ErrSendRequest, timeout: time.Nanosecond, namespace: "goods.goods", code: "bp1"},
					//{name: "response_status_!ok", stand: s, ctx: ctxBg, procCtx: map[string]string{"1": "2"}, errExpected: ErrResponseStatusNotOK, namespace: "goods.goods", code: "bp1"},
					//{name: "nil_stand", stand: nil, ctx: ctxBg, procCtx: EmptyProcCtx{}, errExpected: ErrCreateRequest, namespace: "goods.goods", code: "bp1"},
					{name: "empty_namespace", stand: s, ctx: ctxBg, procCtx: EmptyProcCtx{}, errExpected: ErrResponseStatusNotOK, namespace: "", code: "bp1"},
					{name: "empty_code", stand: s, ctx: ctxBg, procCtx: EmptyProcCtx{}, errExpected: ErrResponseStatusNotOK, namespace: "goods.goods", code: ""},
				}

				for _, tc := range testCases {
					t.Run(tc.name, func(t *testing.T) {
						p := NewProc[interface{}](Settings{
							Stand:     tc.stand,
							Namespace: tc.namespace,
							Code:      tc.code,
						})
						if tc.timeout > 0 {
							p.SetClientTimeout(tc.timeout)
						}
						_, err := p.Run(tc.ctx, tc.procCtx)

						require.ErrorIs(t, err, tc.errExpected)
						fmt.Println(err)
					})
				}
			})

			t.Run("success", func(t *testing.T) {

				type TC[T interface{}] struct {
					name      string
					stand     Stand
					namespace string
					code      string
					ctx       context.Context
					procCtx   T
				}
				testCases1 := []TC[EmptyProcCtx]{
					{name: "success_empty_ctx", stand: s, ctx: ctxBg, procCtx: EmptyProcCtx{}, namespace: "goods.goods", code: "bp1"},
				}

				for _, tc := range testCases1 {
					t.Run(tc.name, func(t *testing.T) {
						p := NewProc[EmptyProcCtx](Settings{
							Stand:     tc.stand,
							Namespace: tc.namespace,
							Code:      tc.code,
						})
						procItem, err := p.Run(tc.ctx, tc.procCtx)

						require.NoError(t, err)
						require.Len(t, procItem.ID, uuid4Len)

					})
				}

				testCases2 := []TC[procCtx]{
					{name: "success_ctx", stand: s, ctx: ctxBg, procCtx: procCtx{Number: 2}, namespace: "goods", code: "bp2"},
				}
				for _, tc := range testCases2 {
					t.Run(tc.name, func(t *testing.T) {
						p := NewProc[procCtx](Settings{
							Stand:     tc.stand,
							Namespace: tc.namespace,
							Code:      tc.code,
						})
						procItem, err := p.Run(tc.ctx, tc.procCtx)

						require.NoError(t, err)
						require.Len(t, procItem.ID, uuid4Len)

					})
				}

			})

		})

		t.Run("proc_get_by_id", func(t *testing.T) {

			validProcId := "bde45d3c-50f0-43bd-b273-9dbeb6d714f6"
			validProcIdNotExisted := "bde45d3c-50f0-43bd-b273-9dbeb6d714f7"
			type TC struct {
				name        string
				ctx         context.Context
				procInstId  string
				errExpected error
				timeout     time.Duration
			}
			testCases := []TC{
				{name: "invalid_proc_id", ctx: ctxBg, procInstId: "invalid_id", errExpected: ErrInvalidID},
				{name: "failed_request_creation", ctx: nil, procInstId: validProcId, errExpected: ErrCreateRequest},
				{name: "failed_request_sending", ctx: ctxBg, procInstId: validProcId, errExpected: ErrSendRequest, timeout: time.Nanosecond},
				{name: "failed_resp_status_!ok", ctx: ctxBg, procInstId: validProcIdNotExisted, errExpected: ErrResponseStatusNotOK},
				{name: "success", ctx: ctxBg, procInstId: validProcId, errExpected: nil},
			}

			bp := NewProc[procCtx](Settings{
				Stand:     s,
				Namespace: "goods",
				Code:      "bp2",
			})

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					bp.SetClientTimeout(time.Second * 5)
					if tc.timeout > 0 {
						bp.SetClientTimeout(tc.timeout)
					}
					proc, err := bp.GetInstanceById(tc.ctx, tc.procInstId)

					require.ErrorIs(t, err, tc.errExpected)
					if err != nil {
						return
					}
					require.Equal(t, tc.procInstId, proc.ID)
					//fmt.Println(err.Error())
				})
			}
		})

	})

}
