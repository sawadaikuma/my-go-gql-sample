package dataloader

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"my-go-gql-sample/graph/model"
)

const loadersKey = "dataloaders"

type Loaders struct {
	TagTodoById TagLoader
}

func LogAndQuery(db *sql.DB, query string, args ...interface{}) *sql.Rows {
	fmt.Printf("[SQL][query]%v", query)
	fmt.Printf("[SQL][args]%v", args...)
	res, err := db.Query(query, args...)
	if err != nil {
		panic(err)
	}
	return res
}

func Middleware(conn *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
				TagTodoById: TagLoader{
					maxBatch: 100,                   // 最大待ち数
					wait:     10 * time.Millisecond, // 待ち時間
					fetch: func(ids []string) ([]*model.Tag, []error) {
						return fetchTagTodoById(conn, ids)
					},
				},
			})
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func fetchTagTodoById(conn *sql.DB, ids []string) ([]*model.Tag, []error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		placeholders[i] = "?"
		args[i] = ids[i]
	}

	// fmt.Printf("[DATALOADER][ids] ==> %+v\n", ids)

	// TODO: fix to gorm
	res := LogAndQuery(conn,
		"SELECT id, name, todo_id from localdb.tag WHERE todo_id IN ("+strings.Join(placeholders, ",")+")",
		args...,
	)
	defer res.Close()

	tagByTodoId := map[string]*model.Tag{}
	for res.Next() {
		tag := model.Tag{}
		err := res.Scan(&tag.ID, &tag.Name, &tag.TodoID)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("[DATALOADER][tag Scanned] ==> %+v\n", tag)
		tagByTodoId[tag.TodoID] = &tag
	}

	tags := make([]*model.Tag, len(ids))
	for i, id := range ids {
		tags[i] = tagByTodoId[id]
	}
	// fmt.Printf("[DATALOADER][Tag][Response] ==> %+v\n", tags)
	return tags, nil
}
