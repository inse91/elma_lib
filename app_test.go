package e365_gateway

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const url = "https://q3bamvpkvrulg.elma365.ru"
const token = "9dccd775-f46a-4167-b2b0-4bc2e6d6356b"

type Product struct {
	//Name  string `json:"__name"`
	Common
	Price int `json:"price,omitempty"`
}

func TestElmaApp(t *testing.T) {

	s := NewStand("https://q3bamvpkvrulg.elma365.ru", "", "33ef3e66-c1cd-4d99-9a77-ddc4af2893cf")
	goods := NewApp[Product](AppSettings{
		Stand:     s,
		Namespace: "goods",
		Code:      "goods",
	})

	t.Run("get by id", func(t *testing.T) {
		item, err := goods.GetByID("26cc4e77-0f02-44ae-a92f-0a34b8a6f4fc")
		require.NoError(t, err)
		require.Equal(t, "Мясо", item.Name)
	})

	t.Run("create_item", func(t *testing.T) {
		item, err := goods.Create(&Product{
			Common: Common{
				Name: "test1",
			},
			Price: 15,
		})
		require.NoError(t, err)
		require.NotNil(t, item)
		fmt.Println(item.ID)
	})

	t.Run("create_many", func(t *testing.T) {
		for i := 0; i < 1000; i++ {

			item, err := goods.Create(&Product{
				Common: Common{
					Name: "test1",
				},
				Price: 10,
			})
			require.NoError(t, err)
			require.NotNil(t, item)
			fmt.Println(i, item.ID)
			time.Sleep(time.Millisecond * 500)

		}

	})

	t.Run("update", func(t *testing.T) {
		item, err := goods.Update("0189f301-b06f-1ac3-f257-99628aa722de", &Product{
			Common: Common{
				Name: "test1",
			},
			Price: 25,
		})

		require.NoError(t, err)
		require.Equal(t, 25, item.Price)
	})

	t.Run("set_status", func(t *testing.T) {
		item, err := goods.SetStatus("b937afb7-df6e-4c95-9076-5018f36a6ee7", "st1") // пельмени
		require.NoError(t, err)
		require.Equal(t, 1, item.Status.Status)
	})

	t.Run("get_status", func(t *testing.T) {
		si, err := goods.GetStatusInfo() // пельмени
		require.NoError(t, err)
		require.Equal(t, 2, len(si.StatusItems))
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
					InStatuses: []string{
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
				InStatuses: []string{
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
			}).All()
			require.NoError(t, err)
			require.Equal(t, 4, len(items))
		})

		t.Run("search_all_filter", func(t *testing.T) {
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
