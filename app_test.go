package e365_gateway

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const url = "https://q3bamvpkvrulg.elma365.ru"
const token = "9dccd775-f46a-4167-b2b0-4bc2e6d6356b"

type Product struct {
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
	validItemId := "018a2b9f-003d-2b48-7e2a-324e6fc16db8"

	ctxBg := context.Background()
	//goLimit := 10

	t.Run("single_success", func(t *testing.T) {

		t.Run("get_by_id", func(t *testing.T) {
			item, err := goods.GetByID(ctxBg, validItemId)
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
			item, err := goods.Create(ctxBg, p)
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
			item, err := goods.Update(ctxBg, validItemId, Product{
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
			item, err := goods.SetStatus(ctxBg, validItemId, "st2")
			require.NoError(t, err)
			require.Equal(t, 2, item.Status.Status)
		})

		t.Run("get_status_info", func(t *testing.T) {
			si, err := goods.GetStatusInfo(ctxBg)
			require.NoError(t, err)
			require.Equal(t, 2, len(si.StatusItems))
		})

	})

	t.Run("search", func(t *testing.T) {

		t.Run("search_first", func(t *testing.T) {
			item, err := goods.Search().First(ctxBg)
			require.NoError(t, err)
			require.Equal(t, "test1", item.Name)
		})

		t.Run("search_first_filter", func(t *testing.T) {
			item, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"__name": "Мясо",
					"price":  AppNumberFilter.From(50).To(500),
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
			}).First(ctxBg)

			require.NoError(t, err)
			require.Equal(t, "Мясо", item.Name)
		})

		t.Run("search_all", func(t *testing.T) {
			items, err := goods.Search().Size(95).All(ctxBg)
			require.NoError(t, err)
			require.Equal(t, 95, len(items))
		})

		t.Run("search_all_filter", func(t *testing.T) {
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"price":  AppNumberFilter.From(50).To(500),
					"__name": "Мясо",
				},
			}).All(ctxBg)
			require.NoError(t, err)
			require.Equal(t, 1, len(items))
		})

		t.Run("search_all_include_del", func(t *testing.T) {
			items, err := goods.Search().IncludeDeleted().Size(23).All(ctxBg)
			require.NoError(t, err)
			require.Len(t, items, 23)
		})

		t.Run("search_where", func(t *testing.T) {
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"price":  AppNumberFilter.From(500),
					"__name": "Мясо",
				},
			}).Size(1).All(ctxBg)
			require.NoError(t, err)
			require.Equal(t, 1, len(items))
		})

		t.Run("search_date", func(t *testing.T) {
			aug28, _ := time.Parse(time.DateOnly, "2023-08-28")
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"__createdAt": AppDateFilter.To(aug28),
				},
			}).AllAtOnce(ctxBg, 10)
			require.NoError(t, err)
			require.Len(t, items, 604)

		})

		t.Run("bool_filter", func(t *testing.T) {
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"isPoAktsii": true,
				},
			}).All(ctxBg)
			require.NoError(t, err)
			require.Len(t, items, 1)
		})

		t.Run("category_filter", func(t *testing.T) {
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{
					"categ": AppCategoryFilter("one"),
				},
			}).All(ctxBg)
			require.NoError(t, err)
			require.Len(t, items, 2)
		})

		t.Run("app_field_filter", func(t *testing.T) {
			items, err := goods.Search().Where(SearchFilter{
				Fields: Fields{

					"appField": AppApplicationFilter("72f60b1a-168e-412f-8cfd-e119c28e99b7"),
				},
			}).Size(100).All(ctxBg)
			require.NoError(t, err)
			require.Len(t, items, 2)
		})

	})

	t.Run("table_tests", func(t *testing.T) {
		t.Parallel()

		t.Run("create", func(t *testing.T) {
			t.Parallel()
			testCases := []struct {
				name        string
				expectedErr error
				ctx         context.Context
				item        Product
			}{
				{name: "failed_request_creation", expectedErr: ErrCreateRequest},
				{name: "success", ctx: ctxBg, expectedErr: nil, item: Product{AppCommon: AppCommon{Name: "test_creation"}, Price: 23}},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					item, err := goods.Create(tc.ctx, tc.item)
					require.ErrorIs(t, err, tc.expectedErr)
					if tc.expectedErr != nil {
						return
					}
					require.Len(t, item.ID, uuid4Len)
					require.Equal(t, item.Name, tc.item.Name)
					require.Equal(t, item.Price, tc.item.Price)
				})
			}
		})

		t.Run("get_by_id", func(t *testing.T) {
			t.Parallel()

			validItemIdNotExisted := "018a2b9f-003d-2b48-7e2a-324e6fc16db9"
			testCases := []struct {
				name        string
				expectedErr error
				ctx         context.Context
				id          string
			}{
				{name: "invalid_id", expectedErr: ErrInvalidID},
				{name: "failed_request_creation", id: validItemId, expectedErr: ErrCreateRequest},
				{name: "id_not_found", ctx: ctxBg, expectedErr: ErrResponseStatusNotOK, id: validItemIdNotExisted},
				{name: "success", ctx: ctxBg, expectedErr: nil, id: validItemId},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					item, err := goods.GetByID(tc.ctx, tc.id)
					require.ErrorIs(t, err, tc.expectedErr)
					if tc.expectedErr != nil {
						return
					}
					require.Len(t, item.ID, uuid4Len)
				})
			}
		})

		t.Run("update", func(t *testing.T) {
			t.Parallel()

			validItemIdNotExisted := "018a2b9f-003d-2b48-7e2a-324e6fc16db9"
			testCases := []struct {
				name        string
				expectedErr error
				ctx         context.Context
				id          string
				item        Product
			}{
				{name: "invalid_id", expectedErr: ErrInvalidID},
				{name: "failed_request_creation", id: validItemId, expectedErr: ErrCreateRequest},
				{name: "id_not_found", ctx: ctxBg, expectedErr: ErrResponseStatusNotOK, id: validItemIdNotExisted},
				{name: "success", ctx: ctxBg, expectedErr: nil, id: validItemId, item: Product{AppCommon: AppCommon{Name: "test_creation_changed"}, Price: 27}},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					item, err := goods.Update(tc.ctx, tc.id, tc.item)
					require.ErrorIs(t, err, tc.expectedErr)
					if tc.expectedErr != nil {
						return
					}
					require.Len(t, item.ID, uuid4Len)
					require.Equal(t, item.Name, tc.item.Name)
					require.Equal(t, item.Price, tc.item.Price)
				})
			}
		})

		t.Run("set_status", func(t *testing.T) {
			t.Parallel()

			validItemIdNotExisted := "018a2b9f-003d-2b48-7e2a-324e6fc16db9"
			validStatusCode := "st2"
			invalidStatusCode := "st22"
			testCases := []struct {
				name        string
				expectedErr error
				ctx         context.Context
				id          string
				statusCode  string
			}{
				{name: "invalid_id", expectedErr: ErrInvalidID},
				{name: "failed_request_creation", id: validItemId, expectedErr: ErrCreateRequest},
				{name: "id_not_found", ctx: ctxBg, expectedErr: ErrResponseStatusNotOK, id: validItemIdNotExisted},
				{name: "invalid_status_code", ctx: ctxBg, expectedErr: ErrResponseStatusNotOK, statusCode: invalidStatusCode, id: validItemIdNotExisted},
				{name: "success", ctx: ctxBg, expectedErr: nil, id: validItemId, statusCode: validStatusCode},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					item, err := goods.SetStatus(tc.ctx, tc.id, tc.statusCode)
					require.ErrorIs(t, err, tc.expectedErr)
					if tc.expectedErr != nil {
						return
					}
					require.Len(t, item.ID, uuid4Len)
					require.Equal(t, 2, item.Status.Status)
				})
			}
		})

		t.Run("get_status_info", func(t *testing.T) {
			t.Parallel()

			testCases := []struct {
				name        string
				expectedErr error
				ctx         context.Context
			}{
				{name: "failed_request_creation", expectedErr: ErrCreateRequest},
				{name: "success", ctx: ctxBg, expectedErr: nil},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					statusInfo, err := goods.GetStatusInfo(tc.ctx)
					require.ErrorIs(t, err, tc.expectedErr)
					if tc.expectedErr != nil {
						return
					}
					require.Len(t, statusInfo.StatusItems, 2)
				})
			}
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
		item, err := goods.Create(context.Background(), Product{
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
