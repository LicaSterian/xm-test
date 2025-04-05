package repo

import (
	"companies/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName        string = "mx-auth"
	CompaniesCollection string = "companies"
)

var (
	ErrFindOne                = errors.New("findOne returned an error")
	ErrFindOneDecode          = errors.New("findOne returned an error while decode ")
	ErrFindOneAndUpdate       = errors.New("findOneAndUpdate returned an error")
	ErrFindOneAndUpdateDecode = errors.New("findOneAndUpdate returned an error while decoding")
)

type mongoCompanyRepo struct {
	client *mongo.Client
}

func NewMongoCompanyRepo(mongoClient *mongo.Client) CompanyRepo {
	return &mongoCompanyRepo{
		client: mongoClient,
	}
}

func (r *mongoCompanyRepo) CreateCompany(ctx context.Context, company models.Company) (uuid.UUID, error) {
	result, err := r.client.Database(DatabaseName).Collection(CompaniesCollection).InsertOne(ctx, company)
	if err != nil {
		return uuid.Nil, err
	}
	insertedId, err := uuid.FromBytes(result.InsertedID.(primitive.Binary).Data)
	if err != nil {
		return uuid.Nil, err
	}
	return insertedId, nil
}

func (r *mongoCompanyRepo) PatchCompany(ctx context.Context, companyId uuid.UUID, updateCompanyInput models.UpdateCompanyInput) (models.Company, error) {
	filter := bson.M{"_id": companyId}
	update := bson.M{"$set": updateCompanyInput.ToBsonM()}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var updatedCompany models.Company
	result := r.client.
		Database(DatabaseName).
		Collection(CompaniesCollection).
		FindOneAndUpdate(ctx, filter, update, opts)
	err := result.Err()
	if err != nil {
		return models.Company{}, errors.Join(ErrFindOneAndUpdate, err)
	}
	err = result.Decode(&updatedCompany)
	if err != nil {
		return models.Company{}, errors.Join(ErrFindOneAndUpdateDecode, err)
	}

	return updatedCompany, nil
}

func (r *mongoCompanyRepo) GetCompany(ctx context.Context, companyId uuid.UUID) (models.Company, error) {
	filter := bson.M{
		"_id": companyId,
	}
	result := r.client.Database(DatabaseName).Collection(CompaniesCollection).FindOne(ctx, filter)
	err := result.Err()
	if err != nil {
		return models.Company{}, errors.Join(ErrFindOne, err)
	}
	var company models.Company
	err = result.Decode(&company)
	if err != nil {
		return models.Company{}, errors.Join(ErrFindOneDecode, err)
	}
	return company, nil
}

func (r *mongoCompanyRepo) DeleteCompany(ctx context.Context, companyId uuid.UUID) error {
	return nil
}
