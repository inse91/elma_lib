package e365_gateway

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestElmaProc(t *testing.T) {

	s := NewStand("https://q3bamvpkvrulg.elma365.ru", "", "33ef3e66-c1cd-4d99-9a77-ddc4af2893cf")
	bp := NewProcess[EmptyCtx](AppSettings{
		Stand:     s,
		Namespace: "goods.goods",
		Code:      "bp1",
	})

	t.Run("run_proc", func(t *testing.T) {
		id, err := bp.Run(EmptyCtx{})
		require.NoError(t, err)
		require.NotEqual(t, "", id)
		fmt.Println(id)
	})

}
