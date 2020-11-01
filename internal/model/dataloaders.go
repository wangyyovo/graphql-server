//go:generate go run github.com/99designs/gqlgen
//go:generate go run github.com/vektah/dataloaden TeacherLoader int *graphql-server/internal/model.Teacher
//go:generate go run github.com/vektah/dataloaden TeachersLoader int []*graphql-server/internal/model.Teacher

package model

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

//type ctxKeyType struct{ name string }

//var ctxKey = ctxKeyType{"userCtx"}

var ctxKey = "loader"

type loaders struct {
	TeacherLoader      *TeacherLoader
	TeachersLoader     *TeachersLoader
}

// nolint: gosec
func LoaderMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		ldrs := &loaders{}

		// set this to zero what happens without dataloading
		wait := 250 * time.Microsecond

		// simple 1:1 loader, fetch an address by its primary key
		ldrs.TeacherLoader = &TeacherLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []int) ([]*Teacher, []error) {
				fmt.Println("TeacherLoader")

				var keySql []string
				for _, key := range keys {
					keySql = append(keySql, strconv.Itoa(key))
				}

				fmt.Printf("SELECT * FROM address WHERE id IN (%s)\n", strings.Join(keySql, ","))
				time.Sleep(5 * time.Millisecond)

				teachers := make([]*Teacher, len(keys))
				errors := make([]error, len(keys))
				for i, _ := range keys {
					teachers[i] = &Teacher{
						ID:          i,
						Name:        "test",
					}
				}
				return teachers, errors
			},
		}

		// 1:M loader
		ldrs.TeachersLoader = &TeachersLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []int) ([][]*Teacher, []error) {
				var keySql []string
				for _, key := range keys {
					keySql = append(keySql, strconv.Itoa(key))
				}

				fmt.Printf("SELECT * FROM orders WHERE customer_id IN (%s)\n", strings.Join(keySql, ","))
				time.Sleep(5 * time.Millisecond)

				teacher := make([][]*Teacher, len(keys))
				errors := make([]error, len(keys))
				for i, _ := range keys {
					//id := 10 + rand.Int()%3
					teacher[i] = []*Teacher{
						{
							ID: i,
							Name: "test",
						},
					}

					// if you had another customer loader you would prime its cache here
					// by calling `ldrs.ordersByID.Prime(id, orders[i])`
				}

				return teacher, errors
			},
		}

		ctx := context.WithValue(gctx.Request.Context(), ctxKey,ldrs)
		gctx.Request = gctx.Request.WithContext(ctx)
		gctx.Next()
		//gctx.Set(ctxKey, ldrs)
		//dlCtx := context.WithValue(gctx.Context(), ctxKey, ldrs)
		//next.ServeHTTP(w, r.WithContext(dlCtx))
	}
}

func CtxLoaders(ctx context.Context) *loaders {
	return ctx.Value(ctxKey).(*loaders)
}

func LoadersFromContext(ctx context.Context) (*loaders,error) {
	ldrVal := ctx.Value(ctxKey)
	if ldrVal == nil {
		err := fmt.Errorf("could not retrieve loaders")
		return nil, err
	}

	ldr, ok := ldrVal.(*loaders)
	if !ok {
		err := fmt.Errorf("loaders has wrong type")
		return nil, err
	}
	return ldr, nil
}
