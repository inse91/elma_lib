package e365_gateway

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type Product struct {
	AppCommon
	Price int `json:"price"`
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

	t.Run("documentation_tests", func(t *testing.T) {

		t.Run("full_success_scenario", func(t *testing.T) {

			now := time.Now()
			newName := "newProduct" + now.String()
			newPrice := float64(now.Second())

			// creation
			newItem := Product{AppCommon: AppCommon{Name: newName}, Price: int(newPrice)}
			item, err := goods.Create(ctxBg, newItem)
			require.NoError(t, err)
			require.Equal(t, item.Price, newItem.Price)
			require.Equal(t, item.Name, newItem.Name)

			// searching
			foundItem, err := goods.Search().
				Where(SearchFilter{
					Fields: Fields{
						"price": Field.Number().From(newPrice).To(newPrice),
					},
				}).First(ctxBg)
			require.NoError(t, err)
			require.Equal(t, item.ID, foundItem.ID)

			// updating
			updPrice := foundItem.Price + 17
			foundItem.Price = updPrice
			updItem, err := goods.Update(ctxBg, foundItem.ID, foundItem)
			require.NoError(t, err)
			require.Equal(t, updItem.Price, updPrice)

			// getting by id
			itemById, err := goods.GetByID(ctxBg, foundItem.ID)
			require.NoError(t, err)
			require.Equal(t, itemById.Name, foundItem.Name)

			// getting status info
			si, err := goods.GetStatusInfo(ctxBg)
			require.NoError(t, err)
			siLen := len(si.StatusItems)
			require.Greater(t, len(si.StatusItems), 0)

			newStatus := si.StatusItems[siLen-1]
			itemWithChangedStatus, err := goods.SetStatus(ctxBg, itemById.ID, newStatus.Code)
			require.NoError(t, err)
			require.Equal(t, newStatus.Id, itemWithChangedStatus.Status.Status)

		})

		t.Run("search", func(t *testing.T) {

			t.Run("from", func(t *testing.T) {

				firstElem, err := goods.Search().First(ctxBg)
				require.NoError(t, err)
				secondElem, err := goods.Search().From(1).First(ctxBg)
				require.NoError(t, err)

				require.NotEqual(t, firstElem.ID, secondElem.ID)

			})

			t.Run("size", func(t *testing.T) {

				items, err := goods.Search().All(ctxBg)
				require.NoError(t, err)
				require.Equal(t, len(items), 10)

				const size = 5
				itemsNew, err := goods.Search().Size(size).All(ctxBg)
				require.NoError(t, err)
				require.Equal(t, len(itemsNew), size)

			})

			t.Run("active", func(t *testing.T) {

				count, err := goods.Search().Count(ctxBg)
				require.NoError(t, err)

				countIncludedDeleted, err := goods.Search().IncludeDeleted().Count(ctxBg)
				require.NoError(t, err)

				require.NotEqual(t, count, countIncludedDeleted)

			})

			t.Run("ids", func(t *testing.T) {

				const id1 = "018a44a7-89b2-81c0-b848-d23fe5ae2023"
				const id2 = "018a44a7-dffc-3ea7-5eb4-afc28df6f507"

				items, err := goods.Search().Where(SearchFilter{
					IDs: []string{id1, id2},
				}).All(ctxBg)

				require.NoError(t, err)
				require.Equal(t, 2, len(items))
				require.Subset(t, []string{items[0].ID, items[1].ID}, []string{id1, id2})

			})

			t.Run("sort_expressions", func(t *testing.T) {

				t.Run("asc", func(t *testing.T) {

					item, err := goods.Search().Where(SearchFilter{
						SortExpressions: []SortExpression{
							{
								Ascending: true,
								Field:     "price",
							},
						},
					}).First(ctxBg)

					require.NoError(t, err)
					require.Equal(t, "theLeastExpensiveProduct", item.Name)

				})

				t.Run("desc", func(t *testing.T) {

					item, err := goods.Search().Where(SearchFilter{
						SortExpressions: []SortExpression{
							{
								Ascending: false,
								Field:     "price",
							},
						},
					}).All(ctxBg)
					_ = item
					require.NoError(t, err)
					//require.Equal(t, "theMostExpensiveProduct", item.Name)

				})

			})

			t.Run("at_status", func(t *testing.T) {

				count, err := goods.Search().Where(SearchFilter{
					AtStatus: []string{"st2"},
				}).Count(ctxBg)

				require.NoError(t, err)
				require.Equal(t, 6, count)

			})

			t.Run("fields", func(t *testing.T) {

				t.Run("bool", func(t *testing.T) {

					items, err := goods.Search().Where(SearchFilter{
						Fields: Fields{
							"isPoAktsii": true,
						},
					}).All(ctxBg)

					require.NoError(t, err)
					require.Equal(t, 7, len(items))

				})

				t.Run("string", func(t *testing.T) {

					item, err := goods.Search().Where(SearchFilter{
						Fields: Fields{
							"__name": "theMostExpensiveProduct",
						},
					}).First(ctxBg)

					require.NoError(t, err)
					require.Equal(t, 2000, item.Price)

				})

				t.Run("number", func(t *testing.T) {

					t.Run("equal", func(t *testing.T) {

						item, err := goods.Search().Where(SearchFilter{
							Fields: Fields{
								"price": Field.Number().Equal(2000),
							},
						}).First(ctxBg)

						require.NoError(t, err)
						require.Equal(t, "theMostExpensiveProduct", item.Name)

					})

					t.Run("in_range", func(t *testing.T) {

						t.Run("from", func(t *testing.T) {
							items, err := goods.Search().Where(SearchFilter{
								Fields: Fields{
									"price": Field.Number().From(200),
								},
							}).All(ctxBg)

							require.NoError(t, err)
							require.Equal(t, 4, len(items))
						})

						t.Run("from_to", func(t *testing.T) {
							items, err := goods.Search().Where(SearchFilter{
								Fields: Fields{
									"price": Field.Number().From(100).To(200),
								},
							}).All(ctxBg)

							require.NoError(t, err)
							require.Equal(t, 5, len(items))
						})
					})

				})

				t.Run("date", func(t *testing.T) {

					augustTheFirst, err := time.Parse(time.RFC3339, "2023-08-01T00:00:00Z")
					require.NoError(t, err)
					items, err := goods.Search().Where(SearchFilter{
						Fields: Fields{
							"__createdAt": Field.Date().To(augustTheFirst),
						},
					}).All(ctxBg)

					require.NoError(t, err)
					require.Equal(t, 5, len(items))

				})

				t.Run("category", func(t *testing.T) {

					items, err := goods.Search().Where(SearchFilter{
						Fields: Fields{
							"categ": Field.Category("one"),
						},
					}).All(ctxBg)

					require.NoError(t, err)
					require.Equal(t, 6, len(items))

				})

				t.Run("app", func(t *testing.T) {

					items, err := goods.Search().Where(SearchFilter{
						Fields: Fields{
							"appField": Field.App("72f60b1a-168e-412f-8cfd-e119c28e99b7"),
						},
					}).All(ctxBg)

					require.NoError(t, err)
					require.Len(t, items, 2)

				})

			})

			t.Run("complex", func(t *testing.T) {

				t.Run("number_date_status", func(t *testing.T) {

					aug15, err := time.Parse(time.DateOnly, "2023-08-15")
					require.NoError(t, err)

					aug31, err := time.Parse(time.DateOnly, "2023-08-31")
					require.NoError(t, err)

					items, err := goods.Search().Where(SearchFilter{
						Fields: Fields{
							"__createdAt": Field.Date().From(aug15).To(aug31),
							"price":       Field.Number().From(20),
						},
						AtStatus: []string{"st1"},
					}).All(ctxBg)

					require.NoError(t, err)
					require.Equal(t, 8, len(items))

				})

				t.Run("ids_sorted", func(t *testing.T) {

					item, err := goods.Search().Where(SearchFilter{
						//Fields: Fields{
						//	"categ": Field.Category("one"),
						//},
						IDs: []string{
							"26cc4e77-0f02-44ae-a92f-0a34b8a6f4fc",
							"b937afb7-df6e-4c95-9076-5018f36a6ee7",
							"6c196840-2bed-478c-9d67-018a034158d6",
							"5c57bc55-9623-4d9a-983b-de3bb0ce6343",
							"3940aa63-9c80-4d27-ba02-012de8045225",
						},
						SortExpressions: []SortExpression{
							{Field: "price", Ascending: false},
						},
					}).First(ctxBg)

					require.NoError(t, err)
					require.Equal(t, 500, item.Price)

				})

				t.Run("category_bool_includeDel", func(t *testing.T) {

					items, err := goods.Search().
						Where(SearchFilter{
							Fields: Fields{
								"categ":      Field.Category("one"),
								"isPoAktsii": true,
							},
						}).
						IncludeDeleted().
						All(ctxBg)

					require.NoError(t, err)
					require.Equal(t, 5, len(items))

				})

			})

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
					require.Len(t, statusInfo.StatusItems, 3)
				})
			}
		})

		t.Run("find", func(t *testing.T) {
			validItemIdNotExisted := "018a2b9f-003d-2b48-7e2a-324e6fc16db9"
			testCases := []struct {
				name        string
				expectedErr error
				ctx         context.Context
				filter      filter
				expectedLen int
			}{
				{name: "failed_request_creation", expectedErr: ErrCreateRequest},
				{name: "id_not_found", ctx: ctxBg, filter: filter{
					From:   0,
					Size:   1,
					Active: true,
					SearchFilter: SearchFilter{
						IDs: []string{validItemIdNotExisted},
					},
				}},
				{name: "bad_request", ctx: ctxBg, expectedErr: ErrResponseStatusNotOK, filter: filter{
					From:   0,
					Size:   1,
					Active: true,
					SearchFilter: SearchFilter{
						Fields: Fields{
							"appField": validItemId,
						},
					},
				}},
				{name: "success", ctx: ctxBg, expectedLen: 1, filter: filter{
					From:   0,
					Size:   2,
					Active: true,
					SearchFilter: SearchFilter{
						IDs: []string{validItemId},
					},
				}},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					items, _, err := goods.find(tc.ctx, tc.filter)
					require.ErrorIs(t, err, tc.expectedErr)
					if tc.expectedErr != nil {
						return
					}
					require.Len(t, items, tc.expectedLen)
				})
			}
		})

	})

}
