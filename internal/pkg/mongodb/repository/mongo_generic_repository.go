package repository

import (
	"context"
	"fmt"
	"reflect"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/data"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/data/specification"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/mapper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/mongodb"
	reflectionHelper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/reflectionhelper"
	typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/iancoleman/strcase"

	"emperror.dev/errors"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://github.com/Kamva/mgm
// https://github.com/mongodb/mongo-go-driver
// https://blog.logrocket.com/how-to-use-mongodb-with-go/
// https://www.mongodb.com/docs/drivers/go/current/quick-reference/
// https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
// https://www.mongodb.com/docs
type mongoGenericRepository[TDataModel interface{}, TEntity interface{}] struct {
	db             *mongo.Client
	databaseName   string
	collectionName string
}

// NewGenericMongoRepositoryWithDataModel create new gorm generic repository
func NewGenericMongoRepositoryWithDataModel[TDataModel interface{}, TEntity interface{}](
	db *mongo.Client,
	databaseName string,
	collectionName string,
) data.GenericRepositoryWithDataModel[TDataModel, TEntity] {
	return &mongoGenericRepository[TDataModel, TEntity]{
		db:             db,
		collectionName: collectionName,
		databaseName:   databaseName,
	}
}

// Add add a new entity to the database
func (m *mongoGenericRepository[TDataModel, TEntity]) Add(
	ctx context.Context,
	entity TEntity,
) error {
	// get the data model type
	// ej: TDataModel = *model.User
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	// get the entity type
	// ej: TEntity = *entity.User
	modelType := typeMapper.GetGenericTypeByT[TEntity]()

	// get the collection from the database
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	// if the model type is the same as the entity type, insert the entity
	if modelType == dataModelType {
		// insert the entity into the database
		_, err := collection.InsertOne(ctx, entity, &options.InsertOneOptions{})
		if err != nil {
			return err
		}
		return nil
	} else {
		// if the model type is not the same as the entity type, map the entity to the data model
		dataModel, err := mapper.Map[TDataModel](entity)
		if err != nil {
			return err
		}
		// insert the data model into the database
		_, err = collection.InsertOne(ctx, dataModel, &options.InsertOneOptions{})
		if err != nil {
			return err
		}
		e, err := mapper.Map[TEntity](dataModel)
		if err != nil {
			return err
		}
		reflectionHelper.SetValue[TEntity](entity, e)
	}
	return nil
}

// AddAll add a list of entities to the database
func (m *mongoGenericRepository[TDataModel, TEntity]) AddAll(
	ctx context.Context,
	entities []TEntity,
) error {
	// iterate over the entities and add them to the database
	for _, entity := range entities {
		err := m.Add(ctx, entity)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetById get an entity by its id
func (m *mongoGenericRepository[TDataModel, TEntity]) GetById(
	ctx context.Context,
	id uuid.UUID,
) (TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if modelType == dataModelType {
		var model TEntity
		// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/query-document/
		// https://www.mongodb.com/docs/drivers/go/current/quick-reference/
		// https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
		// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.10.3/bson
		if err := collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&model); err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection
			if err == mongo.ErrNoDocuments {
				return *new(TEntity), customErrors.NewNotFoundErrorWrap(
					err,
					fmt.Sprintf(
						"can't find the entity with id %s into the database.",
						id.String(),
					),
				)
			}
			return *new(TEntity), errors.WrapIf(
				err,
				fmt.Sprintf(
					"can't find the entity with id %s into the database.",
					id.String(),
				),
			)
		}
		return model, nil
	} else {
		var dataModel TDataModel
		// find the data model by its id
		if err := collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&dataModel); err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection
			if err == mongo.ErrNoDocuments {
				return *new(TEntity), customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("can't find the entity with id %s into the database.", id.String()))
			}
			return *new(TEntity), errors.WrapIf(err, fmt.Sprintf("can't find the entity with id %s into the database.", id.String()))
		}
		entity, err := mapper.Map[TEntity](dataModel)
		if err != nil {
			return *new(TEntity), err
		}
		return entity, nil
	}
}

func (m *mongoGenericRepository[TDataModel, TEntity]) GetAll(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[TEntity], error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if modelType == dataModelType {
		result, err := mongodb.Paginate[TEntity](
			ctx,
			listQuery,
			collection,
			nil,
		)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		result, err := mongodb.Paginate[TDataModel](ctx, listQuery, collection, nil)
		if err != nil {
			return nil, err
		}
		models, err := utils.ListResultToListResultDto[TEntity](result)
		if err != nil {
			return nil, err
		}
		return models, nil
	}
}

// Search search entities by a search term
func (m *mongoGenericRepository[TDataModel, TEntity]) Search(
	ctx context.Context,
	searchTerm string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[TEntity], error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if modelType == dataModelType {
		fields := reflectionHelper.GetAllFields(
			typeMapper.GetGenericTypeByT[TEntity](),
		)
		var a bson.A
		for _, field := range fields {
			if field.Type.Kind() != reflect.String {
				continue
			}
			name := strcase.ToLowerCamel(field.Name)
			a = append(
				a,
				bson.D{
					{Key: name, Value: primitive.Regex{Pattern: searchTerm}},
				},
			)
		}
		filter := bson.D{
			{Key: "$or", Value: a},
		}
		result, err := mongodb.Paginate[TEntity](
			ctx,
			listQuery,
			collection,
			filter,
		)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		fields := reflectionHelper.GetAllFields(typeMapper.GetGenericTypeByT[TDataModel]())
		var a bson.A
		for _, field := range fields {
			if field.Type.Kind() != reflect.String {
				continue
			}
			name := strcase.ToLowerCamel(field.Name)
			a = append(a, bson.D{{Key: name, Value: primitive.Regex{Pattern: searchTerm}}})
		}
		filter := bson.D{
			{Key: "$or", Value: a},
		}
		result, err := mongodb.Paginate[TDataModel](ctx, listQuery, collection, filter)
		if err != nil {
			return nil, err
		}
		models, err := utils.ListResultToListResultDto[TEntity](result)
		if err != nil {
			return nil, err
		}
		return models, nil
	}
}

// GetByFilter get entities by a filter
func (m *mongoGenericRepository[TDataModel, TEntity]) GetByFilter(
	ctx context.Context,
	filters map[string]interface{},
) ([]TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	// we could use also bson.D{} for filtering, it is also a map
	cursorResult, err := collection.Find(ctx, filters)
	if err != nil {
		return nil, err
	}

	defer cursorResult.Close(ctx) // nolint: errcheck

	if modelType == dataModelType {
		var models []TEntity

		for cursorResult.Next(ctx) {
			var e TEntity
			if err := cursorResult.Decode(&e); err != nil {
				return nil, errors.WrapIf(err, "Find")
			}
			models = append(models, e)
		}

		return models, nil
	} else {
		var dataModels []TDataModel

		for cursorResult.Next(ctx) {
			var d TDataModel
			if err := cursorResult.Decode(&d); err != nil {
				return nil, errors.WrapIf(err, "Find")
			}
			dataModels = append(dataModels, d)
		}

		models, err := mapper.Map[[]TEntity](dataModels)
		if err != nil {
			return nil, err
		}
		return models, nil
	}
}

// GetByFuncFilter get entities by a function filter
func (m *mongoGenericRepository[TDataModel, TEntity]) GetByFuncFilter(
	ctx context.Context,
	filterFunc func(TEntity) bool,
) ([]TEntity, error) {
	return nil, nil
}

// FirstOrDefault get the first entity that matches the filter
func (m *mongoGenericRepository[TDataModel, TEntity]) FirstOrDefault(
	ctx context.Context,
	filters map[string]interface{},
) (TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if modelType == dataModelType {
		var model TEntity
		// we could use also bson.D{} for filtering, it is also a map
		if err := collection.FindOne(ctx, filters).Decode(&model); err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection
			if err == mongo.ErrNoDocuments {
				return *new(TEntity), nil
			}
			return *new(TEntity), err
		}

		return model, nil
	} else {
		var dataModel TDataModel
		if err := collection.FindOne(ctx, filters).Decode(&dataModel); err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection
			if err == mongo.ErrNoDocuments {
				return *new(TEntity), nil
			}
			return *new(TEntity), err
		}

		model, err := mapper.Map[TEntity](dataModel)
		if err != nil {
			return *new(TEntity), err
		}
		return model, nil
	}
}

// Update update an entity in the database
func (m *mongoGenericRepository[TDataModel, TEntity]) Update(
	ctx context.Context,
	entity TEntity,
) error {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	if modelType == dataModelType {
		var id interface{}
		id = reflectionHelper.GetFieldValueByName(entity, "Id")
		if id == nil {
			id = reflectionHelper.GetFieldValueByName(entity, "ID")
			if id == nil {
				return errors.New("id field not found")
			}
		}

		var updated TEntity
		// https://www.mongodb.com/docs/manual/reference/method/db.collection.findOneAndUpdate/
		if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": entity}, ops).Decode(&updated); err != nil {
			return err
		}
	} else {
		dataModel, err := mapper.Map[TDataModel](entity)
		if err != nil {
			return err
		}

		var id interface{}
		id = reflectionHelper.GetFieldValueByName(dataModel, "Id")
		if id == nil {
			id = reflectionHelper.GetFieldValueByName(dataModel, "ID")
			if id == nil {
				return errors.New("id field not found")
			}
		}
		// https://www.mongodb.com/docs/manual/reference/method/db.collection.findOneAndUpdate/
		if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": dataModel}, ops).Decode(&dataModel); err != nil {
			return err
		}

		e, err := mapper.Map[TEntity](dataModel)
		if err != nil {
			return err
		}
		reflectionHelper.SetValue[TEntity](entity, e)
	}

	return nil
}

// UpdateAll update a list of entities in the database
func (m *mongoGenericRepository[TDataModel, TEntity]) UpdateAll(
	ctx context.Context,
	entities []TEntity,
) error {
	for _, e := range entities {
		err := m.Update(ctx, e)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete delete an entity by its id
func (m *mongoGenericRepository[TDataModel, TEntity]) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if err := collection.FindOneAndDelete(ctx, bson.M{"_id": id.String()}).Err(); err != nil {
		return err
	}

	return nil
}

// SkipTake skip and take entities from the database
func (m *mongoGenericRepository[TDataModel, TEntity]) SkipTake(
	ctx context.Context,
	skip int,
	take int,
) ([]TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	l := int64(take)
	s := int64(skip)

	cursorResult, err := collection.Find(ctx, bson.D{}, &options.FindOptions{
		Limit: &l,
		Skip:  &s,
	})
	if err != nil {
		return nil, err
	}
	defer cursorResult.Close(ctx) // nolint: errcheck

	if modelType == dataModelType {
		var models []TEntity
		for cursorResult.Next(ctx) {
			var e TEntity
			if err := cursorResult.Decode(&e); err != nil {
				return nil, errors.WrapIf(err, "Find")
			}
			models = append(models, e)
		}

		return models, nil
	} else {
		var dataModels []TDataModel
		for cursorResult.Next(ctx) {
			var d TDataModel
			if err := cursorResult.Decode(&d); err != nil {
				return nil, errors.WrapIf(err, "Find")
			}
			dataModels = append(dataModels, d)
		}
		models, err := mapper.Map[[]TEntity](dataModels)
		if err != nil {
			return nil, err
		}
		return models, nil
	}
}

// Count count the number of entities in the database
func (m *mongoGenericRepository[TDataModel, TEntity]) Count(
	ctx context.Context,
) int64 {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0
	}
	return count
}

// Find find entities by a specification
func (m *mongoGenericRepository[TDataModel, TEntity]) Find(
	ctx context.Context,
	specification specification.Specification,
) ([]TEntity, error) {
	// TODO implement me
	panic("implement me")
}
