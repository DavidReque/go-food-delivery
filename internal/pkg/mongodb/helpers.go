package mongodb

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"

	"emperror.dev/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://stackoverflow.com/a/23650312/581476

// Paginate is a helper function to paginate a collection of items
// It returns a list of items and the total number of items
// It uses the listQuery to get the limit and offset
// It uses the filter to filter the items
// It uses the collection to get the items
// It uses the ctx to get the context
// It uses the options to get the options
func Paginate[T any](
	ctx context.Context,
	listQuery *utils.ListQuery,
	collection *mongo.Collection,
	filter interface{},
) (*utils.ListResult[T], error) {
	// if filter is nil, set it to an empty bson.D
	if filter == nil {
		filter = bson.D{}
	}

	// count the number of documents in the collection
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.WrapIf(err, "CountDocuments")
	}

	// get the limit and skip from the listQuery
	limit := int64(listQuery.GetLimit())
	skip := int64(listQuery.GetOffset())

	// find the documents in the collection
	cursor, err := collection.Find(
		ctx,
		filter,
		&options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		})
	if err != nil {
		return nil, err
	}

	// close the cursor
	defer cursor.Close(ctx)

	// create a new list of items
	var items []T

	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/cursor/#retrieve-all-documents
	// retrieve all documents from the cursor
	err = cursor.All(ctx, &items)
	if err != nil {
		return nil, err
	}

	return utils.NewListResult[T](
		items,
		listQuery.GetSize(), // size is the number of items per page
		listQuery.GetPage(), // page is the current page
		count,               // count is the total number of items
	), nil
}
