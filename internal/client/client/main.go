package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/poohvpn/gqlgo"
)
func main() {
	client := gqlgo.NewClient(`http://localhost:9020/query`, gqlgo.Option{
		Log: func(msg string) {
			fmt.Println(msg)
		},
	})
	singleRequest(client)


	client.Endpoint = `http://localhost:9020/query`
	//server not support batch query
	//batchRequest(client)

	client.Endpoint = `http://localhost:9020/query`
	client.NotCheckHTTPStatusCode200 = true
	handleGraphqlError(client)

	client = gqlgo.NewClient(`http://localhost:9020/query`, gqlgo.Option{
		Log: func(msg string) {
			fmt.Println(msg)
		},
		WebSocketOption: gqlgo.WSOption{
			Log: func(msg string) {
				fmt.Println(msg)
			},
		},
	})
	subscribe(client)

}

func singleRequest(client *gqlgo.Client) {
	fmt.Println("-----graphql single request----")
	data := struct {
		GetSchool     struct {
			Id int `json:"id"`
			Teacher struct {
				Id 		int 	`json:"id"`
				Name 	string  `json:"name"`
			} `json:"teacher"`
		}  `json:"getSchool"`
	}{}
	err := client.Do(context.Background(), &data, gqlgo.Request{
		Query: `
		query ($id:Int!) {
			getSchool(schoolId: $id){
				id
				teacher {
					id @include(if: true)
					name @skip(if: false)
				}
			}
		}
		`,
		Variables: map[string]interface{}{
			"id": 1,
		},
	})

	detailErr := &gqlgo.DetailError{}
	gqlErr := gqlgo.GraphQLErrors{}
	switch {
	case err == nil:
	case errors.As(err, &detailErr):
		fmt.Println("detail error:", detailErr.Response.StatusCode, "\n", detailErr.Content)
		return
	case errors.As(err, &gqlErr):
		fmt.Println("graphql server error:", gqlErr.Error())
		return
	default:
		fmt.Println("graphql client error:", err.Error())
		return
	}

	j, _ := json.Marshal(data)
	fmt.Println("graphql request result:\n", string(j))

}

func batchRequest(client *gqlgo.Client) {

	fmt.Println("-----graphql batch request----")
	data1 := struct {
		GetSchool     struct {
			Id int `json:"id"`
			Teacher struct {
				Id 		int 	`json:"id"`
				Name 	string  `json:"name"`
			} `json:"teacher"`
		}  `json:"getSchool"`
	}{}
	req1 := gqlgo.Request{
		Query: `
		query ($id:Int!) {
			getSchool(schoolId: $id){
				id
				teacher {
					id
					name
				}
			}
		}
		`,
		Variables: map[string]interface{}{
			"id": 1,
		},
	}

	data2 := struct {
		School     struct {
			Id int `json:"id"`
			Teacher struct {
				Id 		int 	`json:"id"`
				Name 	string  `json:"name"`
			} `json:"teacher"`
		}  `json:"School"`
	} {}
	req2 := gqlgo.Request{
		Query: `
		query{
			school{
				id
				teacher {
					id
					name
				}
			}
		}
		`,
		Variables: map[string]interface{}{
		},
	}

	data := []interface{}{&data1, &data2}
	err := client.Do(context.Background(), data, req1, req2)

	detailErr := &gqlgo.DetailError{}
	gqlErr := gqlgo.GraphQLErrors{}
	switch {
	case err == nil:
	case errors.As(err, &detailErr):
		fmt.Println("detail error:", detailErr.Response.StatusCode, "\n", detailErr.Content)
		return
	case errors.As(err, &gqlErr):
		fmt.Println("graphql server error:", gqlErr.Error())
		return
	default:
		fmt.Println("graphql client error:", err.Error())
		return
	}

	j, _ := json.Marshal(data)
	fmt.Println("graphql request result:\n", string(j))
}

func handleGraphqlError(client *gqlgo.Client) {

	fmt.Println("-----handle graphql error----")
	data := struct {
		User struct {
			GetValue struct {
				ID   string
				Name string
			}
		}
	}{}
	req := gqlgo.Request{
		Query: `
		query{
		  User(name:""){
			id
			name
		  }
		}
		`,
		Variables: map[string]interface{}{},
	}

	err := client.Do(context.Background(), &data, req)

	detailErr := &gqlgo.DetailError{}
	gqlErr := gqlgo.GraphQLErrors{}
	switch {
	case err == nil:
	case errors.As(err, &detailErr):
		fmt.Println("detail error:", detailErr.Response.StatusCode, "\n", detailErr.Content)
		return
	case errors.As(err, &gqlErr):
		gqlOneErr := gqlErr[0]
		fmt.Println("graphql server error:",
			gqlOneErr.Message,
			"on line",
			gqlOneErr.Locations[0].Line,
			"column",
			gqlOneErr.Locations[0].Column,
		)
		return
	default:
		fmt.Println("graphql client error:", err.Error())
		return
	}

	j, _ := json.Marshal(data)
	fmt.Println("graphql request result:\n", string(j))
}

func subscribe(client *gqlgo.Client) {
	fmt.Println("-----graphql subscribe----")
	req := gqlgo.Request{
		Query: `
		subscription {
		  message(typ: "test") {
			msg
		  }
		}
		`,
		Variables: map[string]interface{}{},
	}

	user := struct {
		User []struct {
			ID       string
			Username string
		}
	}{}

	received := make(chan struct{})

	id, err := client.Subscribe(req, func(rawMsg json.RawMessage, gqlErrs gqlgo.GraphQLErrors, completed bool) error {
		if completed {
			fmt.Println("server send completed")
			return nil
		}
		if gqlErrs != nil {
			fmt.Println(gqlErrs)
			return gqlErrs
		}
		if err := json.Unmarshal(rawMsg, &user); err != nil {
			fmt.Println(err)
			return err
		}
		received <- struct{}{}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("subscription id: %s\n", id)

	<-received
	if err := client.Unsubscribe(id); err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)
}