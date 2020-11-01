package resolver

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"graphql-server/internal/bff"
	"graphql-server/internal/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	message chan *model.Message
}

func New() bff.Config {

	conf := bff.Config{}
	conf.Resolvers = &Resolver{}
	conf.Directives = bff.DirectiveRoot{
		HasRole: func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
			//if !getCurrentUser(ctx).HasRole(role) {
			if true {
				// block calling the next resolver
				//return nil, fmt.Errorf("Access denied")
			}

			// or let it pass through
			return next(ctx)

		},

	}
	conf.Complexity.Query.Signin = func(childComplexity int) int {
		fmt.Println(childComplexity)
		return childComplexity
	}

	return conf
}