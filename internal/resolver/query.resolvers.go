package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphql-server/internal/bff"
	"graphql-server/internal/model"
)

func (r *queryResolver) Signin(ctx context.Context) (*model.Token, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Shapes(ctx context.Context) (model.Shapes, error) {
	i := 1
	f := 1.1
	m := model.Square{
		Edge: &i,
		Area: &f,
	}

	return m, nil
}

func (r *queryResolver) School(ctx context.Context) (*model.School, error) {
	return &model.School{
		ID: 9999,
		Teacher: &model.Teacher{
			ID:   22,
			Name: "tony",
		},
	}, nil
}

// Query returns bff.QueryResolver implementation.
func (r *Resolver) Query() bff.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
