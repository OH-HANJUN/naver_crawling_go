package database

import (
	"fmt"
	"context"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
	"google.golang.org/api/iterator"
	"naver/domain"
	
)


func GetClient(ctx context.Context) (*bigquery.Client, error){
	
	client, err := bigquery.NewClient(ctx, "damoa-fb351", option.WithCredentialsFile("/Users/j/Desktop/naver/database/key.json"))
	if err != nil {
        return nil, err
    }
	return client, nil

}

func GetData(w io.Writer, iter *bigquery.RowIterator) ([]domain.NaverData, error) {
	var db_data []domain.NaverData
	for {
			var row domain.NaverData
			err := iter.Next(&row)
			if err == iterator.Done {
					return db_data, nil
			}
			if err != nil {
					return db_data, fmt.Errorf("error iterating through results: %w", err)
			}
			db_data = append(db_data,row)
	}
	return db_data, nil
	
}

