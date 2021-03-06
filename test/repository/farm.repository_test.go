package repository

import (
	"testing"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/db"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/test"
	"github.com/adiatma85/golang-rest-template-api/test/fixtures"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type FarmRepositorySuite struct {
	suite.Suite
	farmRepo repository.FarmRepositoryInterface
}

func TestFarmRepository(t *testing.T) {
	suite.Run(t, new(FarmRepositorySuite))
	defer test.TearDownHelper()
}

// Function to initialize the test suite
func (suite *FarmRepositorySuite) SetupSuite() {
	// Initialize Configuration
	test.SetupInitialize("../../.env")
	db.SetupTestingDb(test.Host, test.Username, test.Password, test.Port, test.Database)
	suite.farmRepo = repository.GetFarmRepository()

	// inserting dummy data
	for _, farm := range fixtures.Farms {
		suite.farmRepo.Create(farm)
	}
}

// Create Farm instance Test
func (suite *FarmRepositorySuite) TestCreateFarm_Positive() {
	// Creating Farm
	createdFarm, err := suite.farmRepo.Create(fixtures.WillBeFarm)

	a := suite.Assert()
	a.Equal(fixtures.WillBeFarm.Name, createdFarm.Name, "both of the name from dummy data and existed farm should have the same value")
	a.NoError(err, "should have no error when creating new farm with this parameter")
}

// Get All Farm instances Test
func (suite *FarmRepositorySuite) TestGetAllFarm_Positive() {
	farms, err := suite.farmRepo.GetAll()

	a := suite.Assert()
	a.NotEmpty(farms, "farms variable is not empty")
	a.NoError(err, "should have no error when fetching farms (bulk fetch)")
}

// Test Get Farm from id
func (suite *FarmRepositorySuite) TestGetById_Positive() {
	farm, err := suite.farmRepo.GetById("1")
	a := suite.Assert()

	a.Equal(uint(1), farm.ID, "both of the id from client data and existed farm should have the same value")
	a.Equal(fixtures.Farms[0].Name, farm.Name, "both of the name from dummy data and existed farm should have the same value")
	a.NoError(err, "should have no error when fetching farm (singular fetch by id)")
}

// Test Get Farm (Negative)
func (suite *FarmRepositorySuite) TestGetById_Negative() {
	nonExistentFarm, err := suite.farmRepo.GetById("1000")
	a := suite.Assert()

	a.Error(err, "should have an error when fetching farm (singular fetch by id)")
	a.ErrorIs(err, gorm.ErrRecordNotFound, "the type of error must be error not found")
	a.Equal(err.Error(), "record not found")
	a.Nil(nonExistentFarm, "the resource shoul have not exist or nil")
}

// Test Get Farm by defined model struct
func (suite *FarmRepositorySuite) TestGetByModel_Positive() {
	where := models.Farm{
		Name: "Farm 2",
	}

	farm, err := suite.farmRepo.GetByModel(where)
	a := suite.Assert()

	// Assert each field that exist to make sure both of them is match
	a.Equal(where.Name, farm.Name, "both of the name from dummy data and existed farm should have the same value")
	a.NoError(err, "should have no error when fetching farm (singular fetch by defined struct)")
}

// Test Get Farm by defined model struct (negative)
func (suite *FarmRepositorySuite) TestGetByModel_Negative() {
	where := models.Farm{
		Name: "lorem ipsum",
	}

	nonExistentFarm, err := suite.farmRepo.GetByModel(where)
	a := suite.Assert()

	a.Error(err, "should have an error when fetching farm (singular fetch by model)")
	a.ErrorIs(err, gorm.ErrRecordNotFound, "the type of error must be error not found")
	a.Equal(err.Error(), "record not found")
	a.Nil(nonExistentFarm, "the resource shoul have not exist or nil")
}

// Test Update Existing Resource
func (suite *FarmRepositorySuite) TestUpdateExistedFarm_Positive() {
	a := suite.Assert()
	updateFarm := models.Farm{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "changing name",
	}

	err := suite.farmRepo.Update(&updateFarm)

	// Equal assertion to make sure that updated attribute is updated
	updatedFarm, _ := suite.farmRepo.GetById("1")
	a.Equal(updateFarm.ID, updatedFarm.ID, "both of the 'id' user from client and existed farm should have the same value")
	a.Equal(updateFarm.Name, updatedFarm.Name, "both of the 'name' user from client and existed farm should have the same value")
	a.NoError(err, "should have no error when updating farm")
}

// Test Delete Existing Resource
func (suite *FarmRepositorySuite) TestDeleteFarm_Positive() {
	a := suite.Assert()
	farm := models.Farm{
		Model: gorm.Model{
			ID: 2,
		},
	}

	err := suite.farmRepo.Delete(&farm)
	a.NoError(err, "should have no error when deleting farm")
}
