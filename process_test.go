package e365_gateway

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestElmaProc(t *testing.T) {

	s := NewStand(testDefaultStandSettings)

	t.Run("without_ctx", func(t *testing.T) {

		ctx := context.Background()
		settings := AppSettings{
			Stand:     s,
			Namespace: "goods.goods",
			Code:      "bp1",
		}

		bp := NewProc[EmptyCtx](settings)

		t.Run("run_proc", func(t *testing.T) {

			proc, err := bp.Run(ctx, EmptyCtx{})
			require.NoError(t, err)
			require.Len(t, proc.ID, uuid4Len)

			time.Sleep(time.Millisecond * 500)
			procInfo, err := bp.GetInstanceById(ctx, proc.ID)
			require.NoError(t, err)
			require.Equal(t, settings.Code, procInfo.Template.Code)
			require.Equal(t, proc.ID, procInfo.ID)
			require.Contains(t, []string{StateDone, StateExec}, procInfo.State)
		})

	})

	t.Run("with_ctx", func(t *testing.T) {

		ctx := context.Background()
		settings := AppSettings{
			Stand:     s,
			Namespace: "goods",
			Code:      "bp2",
		}

		type procCtx struct {
			ProcCommon
			Number int `json:"number"`
		}

		bp := NewProc[procCtx](settings)

		t.Run("run_proc", func(t *testing.T) {

			num := 2
			proc, err := bp.Run(ctx, procCtx{
				Number: num,
			})
			require.NoError(t, err)
			require.Len(t, proc.ID, uuid4Len)

			// wait for bp to get done
			time.Sleep(time.Second)

			proc, err = bp.GetInstanceById(ctx, proc.ID)
			require.NoError(t, err)
			require.Equal(t, StateDone, proc.State)
			require.Equal(t, num*2, proc.Number)

		})

	})

}
