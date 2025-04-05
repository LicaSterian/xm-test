package repo

import (
	"companies/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DatabaseName        string = "mx-auth"
	CompaniesCollection string = "companies"
)

var (
	ErrFindOne       = errors.New("find one returned an error")
	ErrFindOneDecode = errors.New("find one error while decode ")
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

func (r *mongoCompanyRepo) PatchCompany(ctx context.Context, companyId uuid.UUID, company models.Company) (models.Company, error) {
	return models.Company{}, nil
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
