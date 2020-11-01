package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphql-server/internal/bff"
	"graphql-server/internal/model"
)

func (r *schoolResolver) Teacher(ctx context.Context, obj *model.School) (*model.Teacher, error) {
	ldr, err := model.LoadersFromContext(ctx)
	if err != nil {
		return nil, err
	}

	m, err := ldr.TeacherLoader.Load(obj.ID)
	if err != nil {
		return nil, err
	}
	if obj.ID == 0 {
		r.message <- &model.Message{Msg: "test"}
	}

	return m, nil
}

func (r *studentResolver) Teacher(ctx context.Context, obj *model.Student, id int) (*model.Teacher, error) {
	panic(fmt.Errorf("not implemented"))
}

// School returns bff.SchoolResolver implementation.
func (r *Resolver) School() bff.SchoolResolver { return &schoolResolver{r} }

// Student returns bff.StudentResolver implementation.
func (r *Resolver) Student() bff.StudentResolver { return &studentResolver{r} }

type schoolResolver struct{ *Resolver }
type studentResolver struct{ *Resolver }
