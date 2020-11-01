package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"graphql-server/internal/bff"
	"graphql-server/internal/model"
)

func (r *subscriptionResolver) Message(ctx context.Context, typ string) (<-chan *model.Message, error) {
	msg := make(chan *model.Message, 1)
	r.message = msg
	return msg, nil
}

// Subscription returns bff.SubscriptionResolver implementation.
func (r *Resolver) Subscription() bff.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
