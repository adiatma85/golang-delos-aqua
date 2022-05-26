package repository

import (
	"testing"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/db"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/test"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PondRepositorySuite struct {
	suite.Suite
	pondRepo repository.PondRepositoryInterface
}

func TestPondRepository(t *testing.T) {
	suite.Run(t, new(PondRepositorySuite))
	defer test.TearDownHelper()
}

// Function to initialize the test suite
func (suite *PondRepositorySuite) SetupSuite() {
	// Initialize Configuration
	test.SetupInitialize("../../.env")
	db.SetupTestingDb(test.Host, test.Username, test.Password, test.Port, test.Database)
	suite.pondRepo = repository.GetPondRepository()

	// Need Farm repo because pond can not be an orphan
	farmRepo := repository.GetFarmRepository()

	// inserting dummy data for farms
	for _, farm := range farms {
		farmRepo.Create(farm)
	}

	// inserting dummy data for ponds
	for _, pond := range ponds {
		suite.pondRepo.Create(pond)
	}
}

// Create Farm instance Test
func (suite *PondRepositorySuite) TestCreatePond_Positive() {
	// Creating Pond
	createdPond, err := suite.pondRepo.Create(willBePond)

	a := suite.Assert()
	a.Equal(willBePond.Name, createdPond.Name, "both of the name from dummy data and existed pond should have the same value")
	a.NoError(err, "should have no error when creating new pond with this parameter")
}

// Get All Pond instances Test
func (suite *PondRepositorySuite) TestGetAllPond_Positive() {
	ponds, err := suite.pondRepo.GetAll()

	a := suite.Assert()
	a.NotEmpty(ponds, "ponds variable is not empty")
	a.NoError(err, "should have no error when fetching ponds (bulk fetch)")
}

// Test Get Pond from id
func (suite *PondRepositorySuite) TestGetById_Positive() {
	pond, err := suite.pondRepo.GetById("1")
	a := suite.Assert()

	a.Equal(uint(1), pond.ID, "both of the id from client data and existed farm should have the same value")
	a.Equal(ponds[0].Name, pond.Name, "both of the name from dummy data and existed farm should have the same value")
	a.Equal(ponds[0].FarmId, pond.FarmId, "both of the farm_id from dummy data and existed farm should have the same value")
	a.NoError(err, "should have no error when fetching farm (singular fetch by id)")
}

// Test Get Pond (Negative)
func (suite *PondRepositorySuite) TestGetById_Negative() {
	nonExistentPond, err := suite.pondRepo.GetById("1000")
	a := suite.Assert()

	a.Error(err, "should have an error when fetching pond (singular fetch by id)")
	a.ErrorIs(err, gorm.ErrRecordNotFound, "the type of error must be error not found")
	a.Equal(err.Error(), "record not found")
	a.Nil(nonExistentPond, "the resource shoul have not exist or nil")
}

// Test Get Pond by defined model struct
func (suite *PondRepositorySuite) TestGetByModel_Positive() {
	where := models.Pond{
		Name: "Pond 1 in Farm 1",
	}

	pond, err := suite.pondRepo.GetByModel(where)
	a := suite.Assert()

	// Assert each field that exist to make sure both of them is match
	a.Equal(where.Name, pond.Name, "both of the name from dummy data and existed pond should have the same value")
	a.NoError(err, "should have no error when fetching farm (singular fetch by defined struct)")
}

// Test Get Pond by defined model struct (negative)
func (suite *PondRepositorySuite) TestGetByModel_Negative() {
	where := models.Pond{
		Name: "lorem ipsum",
	}

	nonExistentPond, err := suite.pondRepo.GetByModel(where)
	a := suite.Assert()

	a.Error(err, "should have an error when fetching farm (singular fetch by model)")
	a.ErrorIs(err, gorm.ErrRecordNotFound, "the type of error must be error not found")
	a.Equal(err.Error(), "record not found")
	a.Nil(nonExistentPond, "the resource shoul have not exist or nil")
}

// Test Update Existing Resource
func (suite *PondRepositorySuite) TestUpdateExistedPond_Positive() {
	a := suite.Assert()
	updatePond := models.Pond{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "changing name",
	}

	err := suite.pondRepo.Update(&updatePond)

	// Equal assertion to make sure that updated attribute is updated
	updatedPond, _ := suite.pondRepo.GetById("1")
	a.Equal(updatePond.ID, updatedPond.ID, "both of the 'id' user from client and existed pond should have the same value")
	a.Equal(updatePond.Name, updatedPond.Name, "both of the 'name' user from client and existed farm should have the same value")
	a.NoError(err, "should have no error when updating pond")
}

// Test Delete Existing Resource
func (suite *PondRepositorySuite) TestDeletePond_Positive() {
	a := suite.Assert()
	pond := models.Pond{
		Model: gorm.Model{
			ID: 2,
		},
	}

	err := suite.pondRepo.Delete(&pond)
	a.NoError(err, "should have no error when deleting farm")
}
