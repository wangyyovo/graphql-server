package main

//github.com/shurcooL/graphql
//github.com/machinebox/graphql
//github.com/sony/appsync-client-go  sub
//https://github.com/poohvpn/gqlgo  sub
//https://github.com/Yamashou/gqlgenc

//https://github.com/novacloudcz/graphql-orm
import (
	"bytes"
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"

	"time"
)


func main() {
	// create a client (safe to share across requests)
	//query()
	//queryWithRole()
	//mutationWithRole()
	//queryUnion()

	upload()

	// https://github.com/jaydenseric/graphql-multipart-request-spec
	//singleUploadCustom()
	//singleUploadWithPayloadCustom()
	//multipleUploadCustom()
	//multipleUpload2Custom()
	//multipleUploadWithPayloadCustom()
	//multipleUploadWithPayload2Custom()

	//batchUpload()
}


func query() {
	// create a client (safe to share across requests)
	client := graphql.NewClient("http://localhost:9020/query")

	// make a request
	req := graphql.NewRequest(`
	query {
		userInfo {
		info {
		  name
		  gender @include(if: true)
		  genre @skip(if: true)
          birthday {
			day
          }
		}
	  }
	}
`)

	// set any variables
	//req.Var("key", "value")

	// set header fields
	//req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData map[string]interface{}
	//respData = make(map[string]interface{})
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Println(respData)

}


func queryWithRole() {
	// create a client (safe to share across requests)
	client := graphql.NewClient("http://localhost:9020/query")

	// make a request
	req := graphql.NewRequest(`
	query {
		log(logid: 1) {
			log {
			  title
			}
		}
	}

`)

	// set any variables
	//req.Var("key", "value")

	// set header fields
	//req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData map[string]interface{}
	//respData = make(map[string]interface{})
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Println(respData)

}


func mutationWithRole() {
	// create a client (safe to share across requests)
	client := graphql.NewClient("http://localhost:9020/query")

	// make a request
	req := graphql.NewRequest(`
	query {
		shapes {
			__typename
			... on Square {
				area
				edge
			}
			... on Circle {
				area
				radius
			}
		}
	}
`)

	// set any variables
	//req.Var("key", "value")

	// set header fields
	//req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData map[string]interface{}
	//respData = make(map[string]interface{})
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Println(respData)

}

func queryUnion() {
	// create a client (safe to share across requests)
	client := graphql.NewClient("http://localhost:9020/query")

	// make a request
	req := graphql.NewRequest(`
	mutation {
		addLike(articleid: "1") {
			status
			errors {
				code
				message
				description
			}
		}
	}
`)

	// set any variables
	//req.Var("key", "value")

	// set header fields
	//req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData map[string]interface{}
	//respData = make(map[string]interface{})
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Println(respData)

}



func upload() {


	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		file, header, _ := r.FormFile("file")
		defer file.Close()
		fmt.Println(header.Filename)

		b, _ := ioutil.ReadAll(file)
		fmt.Println(string(b))

		_, _ = io.WriteString(w, `{"data":{"value":"some data"}}`)
	}))
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	client := graphql.NewClient(srv.URL, graphql.UseMultipartForm())
	f := strings.NewReader(`This is a file`)
	req := graphql.NewRequest("query {}")
	req.File("file", "filename.txt", f)

	// run it and capture the response
	var respData map[string]interface{}
	//respData = make(map[string]interface{})
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Println(respData)
}


func singleUploadCustom() {
	operations := `{"query": "mutation ($file: Upload!) {singleUpload(file: $file) {id,name,content,contentType}}","variables": { "file": null }}`
	mapData := `{ "0": ["variables.file"] }`
	files := []file{
		{
			mapKey:      "0",
			name:        "a.txt",
			content:     "test",
			contentType: "text/plain",
		},
	}
	client := http.Client{}
	req := createUploadRequest("http://localhost:9020/query", operations, mapData, files)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))
	return
}

// github.com/99designs/gqlgen/graphql不支持批量
func batchUploadCustom() {
	operations := `[{ "query": "mutation ($file: Upload!) { singleUpload(file: $file) { id,name,content,contentType } }", "variables": { "file": null } }, { "query": "mutation($files: [Upload!]!) { multipleUpload(files: $files) { id,name,content,contentType } }", "variables": { "files": [null, null] } }]`
	mapData := `{ "0": ["0.variables.file"], "1": ["1.variables.files.0"], "2": ["1.variables.files.1"] }`
	files := []file{
		{
			mapKey:      "0",
			name:        "a.txt",
			content:     "batch test a",
			contentType: "text/plain",
		},
		{
			mapKey:      "1",
			name:        "b.txt",
			content:     "batch test b",
			contentType: "text/plain",
		},
		{
			mapKey:      "2",
			name:        "c.txt",
			content:     "batch test c",
			contentType: "text/plain",
		},
	}
	client := http.Client{}
	req := createUploadRequest("http://localhost:9020/query", operations, mapData, files)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))
	return
}



func singleUploadWithPayloadCustom() {
	operations := `{ "query": "mutation ($req: UploadFile!) { singleUploadWithPayload(req: $req) { id, name, content, contentType } }", "variables": { "req": {"file": null, "id": 1 } } }`
	mapData := `{ "0": ["variables.req.file"] }`
	files := []file{
		{
			mapKey:      "0",
			name:        "a.txt",
			content:     "test",
			contentType: "text/plain",
		},
	}
	client := http.Client{}
	req := createUploadRequest("http://localhost:9020/query", operations, mapData, files)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))
	return

}

func multipleUploadCustom() {
	operations := `{ "query": "mutation($files: [Upload!]!) { multipleUpload(files: $files) { id, name, content, contentType } }", "variables": { "files": [null, null] } }`
	mapData := `{ "0": ["variables.files.0"], "1": ["variables.files.1"] }`
	files := []file{
		{
			mapKey:      "0",
			name:        "a.txt",
			content:     "test1",
			contentType: "text/plain",
		},
		{
			mapKey:      "1",
			name:        "b.txt",
			content:     "test2",
			contentType: "text/plain",
		},
	}
	client := http.Client{}
	req := createUploadRequest("http://localhost:9020/query", operations, mapData, files)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))
	return
}


func multipleUpload2Custom() {
	operations := `{ "query": "mutation($files: [Upload!]!) { multipleUpload(files: $files) { id, name, content, contentType } }", "variables": { "files": [null, null] } }`
	mapData := `{ "0": ["variables.files.0","variables.files.1"]}`
	files := []file{
		{
			mapKey:      "0",
			name:        "a.txt",
			content:     "test1",
			contentType: "text/plain",
		},
	}
	client := http.Client{}
	req := createUploadRequest("http://localhost:9020/query", operations, mapData, files)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))
	return
}

func multipleUploadWithPayloadCustom()  {
	operations := `{ "query": "mutation($req: [UploadFile!]!) { multipleUploadWithPayload(req: $req) { id, name, content, contentType } }", "variables": { "req": [ { "id": 1, "file": null }, { "id": 2, "file": null } ] } }`
	mapData := `{ "0": ["variables.req.0.file", "variables.req.1.file"] }`
	files := []file{
		{
			mapKey:      "0",
			name:        "a.txt",
			content:     "test1",
			contentType: "text/plain",
		},

	}
	client := http.Client{}
	req := createUploadRequest("http://localhost:9020/query", operations, mapData, files)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))
	return
}



func multipleUploadWithPayload2Custom()  {
	operations := `{ "query": "mutation($req: [UploadFile!]!) { multipleUploadWithPayload(req: $req) { id, name, content, contentType } }", "variables": { "req": [ { "id": 1, "file": null }, { "id": 2, "file": null } ] } }`
	mapData := `{ "0": ["variables.req.0.file"], "1": ["variables.req.1.file"] }`
	files := []file{
		{
			mapKey:      "0",
			name:        "a.txt",
			content:     "test1",
			contentType: "text/plain",
		},
		{
			mapKey:      "1",
			name:        "b.txt",
			content:     "test2",
			contentType: "text/plain",
		},
	}
	client := http.Client{}
	req := createUploadRequest("http://localhost:9020/query", operations, mapData, files)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))
	return
}

type file struct {
	mapKey      string
	name        string
	content     string
	contentType string
}

func createUploadRequest(url, operations, mapData string, files []file) *http.Request {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err := bodyWriter.WriteField("operations", operations)
	if err != nil {
		return nil
	}
	err = bodyWriter.WriteField("map", mapData)
	if err != nil {
		return nil
	}
	for i := range files {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, files[i].mapKey, files[i].name))
		h.Set("Content-Type", files[i].contentType)
		ff, err := bodyWriter.CreatePart(h)
		if err != nil {
			return nil
		}
		_, err = ff.Write([]byte(files[i].content))
		if err != nil {
			return nil
		}
	}
	err = bodyWriter.Close()

	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	return req
}