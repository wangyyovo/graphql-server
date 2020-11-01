package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphql-server/internal/bff"
	"graphql-server/internal/model"
)

func (r *entityResolver) FindUserByID(ctx context.Context, id int) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Entity returns bff.EntityResolver implementation.
func (r *Resolver) Entity() bff.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
