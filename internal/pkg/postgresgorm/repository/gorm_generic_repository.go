package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/data"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/data/specification"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/mapper"
	gormPostgres "github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/helpers/gormextensions"
	reflectionHelper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/reflectionhelper"
	typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"

	"emperror.dev/errors"
	"github.com/iancoleman/strcase"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// gormGenericRepository es una implementación genérica del repositorio usando GORM
type gormGenericRepository[TDataModel interface{}, TEntity interface{}] struct {
	db *gorm.DB
}

// NewGenericGormRepositoryWithDataModel create new gorm generic repository
func NewGenericGormRepositoryWithDataModel[TDataModel interface{}, TEntity interface{}](
	db *gorm.DB,
) data.GenericRepositoryWithDataModel[TDataModel, TEntity] {
	return &gormGenericRepository[TDataModel, TEntity]{
		db: db,
	}
}

// NewGenericGormRepository create new gorm generic repository
func NewGenericGormRepository[TEntity interface{}](
	db *gorm.DB,
) data.GenericRepository[TEntity] {
	return &gormGenericRepository[TEntity, TEntity]{
		db: db,
	}
}

func (r *gormGenericRepository[TDataModel, TEntity]) Add(
	ctx context.Context,
	entity TEntity,
) error {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()

	if modelType == dataModelType {
		err := r.db.WithContext(ctx).Create(entity).Error
		if err != nil {
			return err
		}

		return nil
	} else {
		dataModel, err := mapper.Map[TDataModel](entity)
		if err != nil {
			return err
		}
		err = r.db.WithContext(ctx).Create(dataModel).Error
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

func (r *gormGenericRepository[TDataModel, TEntity]) AddAll(
	ctx context.Context,
	entities []TEntity,
) error {
	for _, entity := range entities {
		err := r.Add(ctx, entity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *gormGenericRepository[TDataModel, TEntity]) GetById(
	ctx context.Context,
	id uuid.UUID,
) (TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()

	if modelType == dataModelType {
		var model TEntity
		if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
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
		if err := r.db.WithContext(ctx).First(&dataModel, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *gormGenericRepository[TDataModel, TEntity]) GetAll(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[TEntity], error) {
	result, err := gormPostgres.Paginate[TDataModel, TEntity](
		ctx,
		listQuery,
		r.db,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *gormGenericRepository[TDataModel, TEntity]) Search(
	ctx context.Context,
	searchTerm string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[TEntity], error) {
	fields := reflectionHelper.GetAllFields(
		typeMapper.GetGenericTypeByT[TDataModel](),
	)
	query := r.db

	for _, field := range fields {
		if field.Type.Kind() != reflect.String {
			continue
		}

		query = query.Or(
			fmt.Sprintf("%s LIKE ?", strcase.ToSnake(field.Name)),
			"%"+strings.ToLower(searchTerm)+"%",
		)
	}

	result, err := gormPostgres.Paginate[TDataModel, TEntity](
		ctx,
		listQuery,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *gormGenericRepository[TDataModel, TEntity]) GetByFilter(
	ctx context.Context,
	filters map[string]interface{},
) ([]TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	if modelType == dataModelType {
		var models []TEntity
		err := r.db.WithContext(ctx).Where(filters).Find(&models).Error
		if err != nil {
			return nil, err
		}
		return models, nil
	} else {
		var dataModels []TDataModel
		err := r.db.WithContext(ctx).Where(filters).Find(&dataModels).Error
		if err != nil {
			return nil, err
		}
		models, err := mapper.Map[[]TEntity](dataModels)
		if err != nil {
			return nil, err
		}
		return models, nil
	}
}

func (r *gormGenericRepository[TDataModel, TEntity]) GetByFuncFilter(
	ctx context.Context,
	filterFunc func(TEntity) bool,
) ([]TEntity, error) {
	return *new([]TEntity), nil
}

func (r *gormGenericRepository[TDataModel, TEntity]) FirstOrDefault(
	ctx context.Context,
	filters map[string]interface{},
) (TEntity, error) {
	return *new(TEntity), nil
}

func (r *gormGenericRepository[TDataModel, TEntity]) Update(
	ctx context.Context,
	entity TEntity,
) error {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	if modelType == dataModelType {
		err := r.db.WithContext(ctx).Save(entity).Error
		if err != nil {
			return err
		}
	} else {
		dataModel, err := mapper.Map[TDataModel](entity)
		if err != nil {
			return err
		}
		err = r.db.WithContext(ctx).Save(dataModel).Error
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

func (r *gormGenericRepository[TDataModel, TEntity]) UpdateAll(
	ctx context.Context,
	entities []TEntity,
) error {
	for _, e := range entities {
		err := r.Update(ctx, e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *gormGenericRepository[TDataModel, TEntity]) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	entity, err := r.GetById(ctx, id)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Delete(entity, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *gormGenericRepository[TDataModel, TEntity]) SkipTake(
	ctx context.Context,
	skip int,
	take int,
) ([]TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	if modelType == dataModelType {
		var models []TEntity
		err := r.db.WithContext(ctx).
			Offset(skip).
			Limit(take).
			Find(&models).
			Error
		if err != nil {
			return nil, err
		}
		return models, nil
	} else {
		var dataModels []TDataModel
		err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&dataModels).Error
		if err != nil {
			return nil, err
		}
		models, err := mapper.Map[[]TEntity](dataModels)
		if err != nil {
			return nil, err
		}
		return models, nil
	}
}

func (r *gormGenericRepository[TDataModel, TEntity]) Count(
	ctx context.Context,
) int64 {
	var dataModel TDataModel
	var count int64
	r.db.WithContext(ctx).Model(&dataModel).Count(&count)
	return count
}

func (r *gormGenericRepository[TDataModel, TEntity]) Find(
	ctx context.Context,
	specification specification.Specification,
) ([]TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	if modelType == dataModelType {
		var models []TEntity
		err := r.db.WithContext(ctx).
			Where(specification.GetQuery(), specification.GetValues()...).
			Find(&models).
			Error
		if err != nil {
			return nil, err
		}
		return models, nil
	} else {
		var dataModels []TDataModel
		err := r.db.WithContext(ctx).Where(specification.GetQuery(), specification.GetValues()...).Find(&dataModels).Error
		if err != nil {
			return nil, err
		}
		models, err := mapper.Map[[]TEntity](dataModels)
		if err != nil {
			return nil, err
		}
		return models, nil
	}
}
