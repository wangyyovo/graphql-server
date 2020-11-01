package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphql-server/internal/bff"
	"graphql-server/internal/model"
)

func (r *mutationResolver) Signup(ctx context.Context) (*model.Token, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Log(ctx context.Context) (*model.Result, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns bff.MutationResolver implementation.
func (r *Resolver) Mutation() bff.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
