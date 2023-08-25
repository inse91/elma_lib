package e365_gateway

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const url = "https://q3bamvpkvrulg.elma365.ru"
const token = "9dccd775-f46a-4167-b2b0-4bc2e6d6356b"

type Product struct {
	//Name  string `json:"__name"`
	AppCommon
	Price int `json:"price,omitempty"`
}

func TestElmaApp(t *testing.T) {

	s := NewStand(testDefaultStandSettings)
	goods := NewApp[Product](Settings{
		Stand:     s,
		Namespace: "goods",
		Code:      "goods",
	})

	t.Run("single_success", func(t *testing.T) {

		validItemId := "018a2b9f-003d-2b48-7e2a-324e6fc16db8"

		t.Run("get by id", func(t *testing.T) {
			item, err := goods.GetByID(validItemId)
			require.NoError(t, err)
			require.Equal(t, "test2", item.Name)
		})

		t.Run("create_item", func(t *testing.T) {
			now := time.Now()
			p := Product{
				AppCommon: AppCommon{
					Name: "test1",
				},
				Price: 15,
			}
			item, err := goods.Create(p)
			fmt.Println(time.Since(now).String())
			require.NoError(t, err)
			require.Len(t, item.ID, uuid4Len)
			require.Equal(t, item.Name, p.Name)
			require.Equal(t, item.Price, p.Price)
			fmt.Println(item.ID)
		})

		t.Run("update", func(t *testing.T) {

			newPrice := 25
			newName := "test2"
			item, err := goods.Update(validItemId, Product{
				AppCommon: AppCommon{
					Name: newName,
				},
				Price: newPrice,
			})

			require.NoError(t, err)
			require.Equal(t, newPrice, item.Price)
			require.Equal(t, newName, item.Name)
		})

		t.Run("set_status", func(t *testing.T) {
			item, err := goods.SetStatus(validItemId, "st2")
			require.NoError(t, err)
			require.Equal(t, 2, item.Status.Status)
		})

		t.Run("get_status", func(t *testing.T) {
			si, err := goods.GetStatusInfo()
			require.NoError(t, err)
			require.Equal(t, 2, len(si.StatusItems))
		})

	})

	t.Run("find", func(t *testing.T) {
		t.Run("find_filter_MAP", func(t *testing.T) {
			items, err := goods.find(filter{
				Active: true,
				From:   0,
				Size:   100,
				SearchFilter: SearchFilter{
					Fields: Fields{
						"__name": "Мясо",
						"price":  AppNumberFilter(50, 200),
						//"__deletedAt": AppDateFilter(),
					},
					IDs: []string{
						"26cc4e77-0f02-44ae-a92f-0a34b8a6f4fc",
						"b937afb7-df6e-4c95-9076-5018f36a6ee7",
					},
					SortExpressions: []SortExpression{
						{Ascending: true, Field: "price"},
					},
					AtStatus: []string{
						"st1",
					},
					StatusGroupId: "",
				},
			})

			//for _, i := range items {
			//	fmt.Println(*i)
			//}
			require.NoError(t, err)
			require.Equal(t, 5, len(items))
			//require.Equal(t, "Мясо", items[0].Name)
			//require.Equal(t, 50, items[0].Price)
		})
	})

	t.Run("search", func(t *testing.T) {

		t.Run("search_first", func(t *testing.T) {
			item, err := goods.Search().First()
			require.NoError(t, err)
			require.Equal(t, "Мясо", item.Name)
		})

		t.Run("search_first_filter", func(t *testing.T) {
			item, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"__name": "Мясо",
					"price":  AppNumberFilter(50, 500),
					//"__deletedAt": AppDateFilter(),
				},

				IDs: []string{
					"26cc4e77-0f02-44ae-a92f-0a34b8a6f4fc",
					"b937afb7-df6e-4c95-9076-5018f36a6ee7",
				},
				SortExpressions: []SortExpression{
					{Ascending: true, Field: "price"},
				},
				AtStatus: []string{
					"st1",
				},
				StatusGroupId: "",
			}).First()

			require.NoError(t, err)
			require.Equal(t, "Мясо", item.Name)
		})

		t.Run("search_all", func(t *testing.T) {
			items, err := goods.Search().Size(95).All()
			require.NoError(t, err)
			require.Equal(t, 95, len(items))
		})

		t.Run("search_all_filter", func(t *testing.T) {
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"price": AppNumberFilter(50, 200),
				},
			}).AllAtOnce()
			require.NoError(t, err)
			require.Equal(t, 4, len(items))
		})

		t.Run("search_all_filter_1", func(t *testing.T) {
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"price":  AppNumberFilter(50, 500),
					"__name": "Мясо",
				},
			}).All()
			require.NoError(t, err)
			require.Equal(t, 1, len(items))
		})

		t.Run("search_all_include_del", func(t *testing.T) {
			items, err := goods.Search().IncludeDeleted().Size(23).All()
			require.NoError(t, err)
			require.Equal(t, 23, len(items))
		})

		t.Run("search_where", func(t *testing.T) {
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"price":  AppNumberFilter(50, 500),
					"__name": "Мясо",
				},
			}).Size(1).All()
			require.NoError(t, err)
			require.Equal(t, 1, len(items))
		})

		t.Run("search_all_at_once", func(t *testing.T) {
			items, err := goods.Search().AllAtOnce()
			require.NoError(t, err)
			require.Equal(t, 416, len(items))

		})

	})

}

func TestApp_CreateMany(t *testing.T) {

	s := NewStand(testDefaultStandSettings)
	goods := NewApp[Product](Settings{
		Stand:     s,
		Namespace: "goods",
		Code:      "goods",
	})

	var isAnySuccess bool

	for i := 0; i < 10; i++ {
		//t.Run("create_item", func(t *testing.T) {
		now := time.Now()
		item, err := goods.Create(Product{
			AppCommon: AppCommon{
				Name: "test1",
			},
			Price: 15,
		})
		if err == nil {
			isAnySuccess = true
		}
		fmt.Println(time.Since(now).String())
		assert.ErrorIs(t, err, ErrSendRequest)
		assert.Contains(t, err.Error(), "Timeout")
		//assert.NotNil(t, item)
		fmt.Println(item.ID)
		//})
	}

	require.True(t, isAnySuccess)
}
