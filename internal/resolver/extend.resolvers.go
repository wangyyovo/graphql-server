package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphql-server/internal/model"
	"io/ioutil"

	"github.com/99designs/gqlgen/graphql"
)

func (r *mutationResolver) SingleUpload(ctx context.Context, file graphql.Upload) (*model.File, error) {
	fmt.Println("SingleUpload")
	content, err := ioutil.ReadAll(file.File)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(content))

	return &model.File{
		ID:          1,
		Name:        file.Filename,
		Content:     string(content),
		ContentType: file.ContentType,
	}, nil
}

func (r *mutationResolver) SingleUploadWithPayload(ctx context.Context, req model.UploadFile) (*model.File, error) {
	content, err := ioutil.ReadAll(req.File.File)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(content))

	return &model.File{
		ID:          1,
		Name:        req.File.Filename,
		Content:     string(content),
		ContentType: req.File.ContentType,
	}, nil
}

func (r *mutationResolver) MultipleUpload(ctx context.Context, files []*graphql.Upload) ([]*model.File, error) {
	fmt.Println("MultipleUpload")
	var contents []string
	var resp []*model.File
	for i := range files {
		content, err := ioutil.ReadAll(files[i].File)
		if err != nil {
			return nil, err
		}
		contents = append(contents, string(content))
		resp = append(resp, &model.File{
			ID:          i + 1,
			Name:        files[i].Filename,
			Content:     string(content),
			ContentType: files[i].ContentType,
		})
	}
	return resp, nil
}

func (r *mutationResolver) MultipleUploadWithPayload(ctx context.Context, req []*model.UploadFile) ([]*model.File, error) {
	var ids []int
	var contents []string
	var resp []*model.File
	for i := range req {
		content, err := ioutil.ReadAll(req[i].File.File)
		if err != nil {
			return nil, err
		}
		ids = append(ids, req[i].ID)
		contents = append(contents, string(content))
		resp = append(resp, &model.File{
			ID:          i + 1,
			Name:        req[i].File.Filename,
			Content:     string(content),
			ContentType: req[i].File.ContentType,
		})
	}
	return resp, nil
}

func (r *queryResolver) GetSchool(ctx context.Context, schoolID int) (*model.School, error) {
	return &model.School{
		ID: schoolID,
		Teacher: &model.Teacher{
			ID:   0,
			Name: "Teacher",
		},
	}, nil
}

func (r *queryResolver) GetTeachers(ctx context.Context, id []int) ([]*model.Teacher, error) {
	var ts []*model.Teacher
	for i := 0; i < len(id); i++ {
		v := &model.Teacher{
			ID:   i,
			Name: fmt.Sprintf("%d", i),
		}
		ts = append(ts, v)
	}
	return ts, nil
}

func (r *queryResolver) Schools(ctx context.Context) ([]*model.School, error) {
	var schools []*model.School
	for i := 0; i < 10; i++ {
		m := &model.School{
			ID: i,
		}
		schools = append(schools, m)
	}
	return schools, nil
}
