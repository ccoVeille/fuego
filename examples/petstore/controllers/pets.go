package controller

import (
	"log/slog"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/petstore/models"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

// default pagination options
var optionPagination = option.Group(
	option.QueryInt("per_page", "Number of items per page", param.Required()),
	option.QueryInt("page", "Page number", param.Default(1), param.Example("1st page", 1), param.Example("42nd page", 42), param.Example("100th page", 100)),
)

type PetsRessources struct {
	PetsService PetsService
}

func (rs PetsRessources) Routes(s *fuego.Server) {
	petsGroup := fuego.Group(s, "/pets").Header("X-Header", "header description")

	fuego.Get(petsGroup, "/", rs.filterPets,
		optionPagination,
		option.Query("name", "Filter by name", param.Example("cat name", "felix"), param.Nullable()),
		option.QueryInt("younger_than", "Only get pets younger than given age in years", param.Default(3)),
		option.Description("Filter pets"),
	)

	fuego.Get(petsGroup, "/all", rs.getAllPets,
		optionPagination,
		option.Tags("my-tag"),
		option.Description("Get all pets"),
	)

	fuego.Get(petsGroup, "/by-age", rs.getAllPetsByAge, option.Description("Returns an array of pets grouped by age"))
	fuego.Post(petsGroup, "/", rs.postPets)

	fuego.Get(petsGroup, "/{id}", rs.getPets)
	fuego.Get(petsGroup, "/by-name/{name...}", rs.getPetByName)
	fuego.Put(petsGroup, "/{id}", rs.putPets)
	fuego.Put(petsGroup, "/{id}/json", rs.putPets).
		RequestContentType("application/json")
	fuego.Delete(petsGroup, "/{id}", rs.deletePets)
}

func (rs PetsRessources) getAllPets(c fuego.ContextNoBody) ([]models.Pets, error) {
	page := c.QueryParamInt("page", 1)
	pageWithTypo := c.QueryParamInt("page-with-typo", 1) // this shows a warning in the logs because "page-with-typo" is not a declared query param
	slog.Info("query params", "page", page, "page-with-typo", pageWithTypo)
	return rs.PetsService.GetAllPets()
}

func (rs PetsRessources) filterPets(c fuego.ContextNoBody) ([]models.Pets, error) {
	return rs.PetsService.FilterPets(PetsFilter{
		Name:        c.QueryParam("name"),
		YoungerThan: c.QueryParamInt("younger_than", 3), // if the default value is not set as the same as the declared value, it will be overwritten and a warning will be logged
	})
}

func (rs PetsRessources) getAllPetsByAge(c fuego.ContextNoBody) ([][]models.Pets, error) {
	return rs.PetsService.GetAllPetsByAge()
}

func (rs PetsRessources) postPets(c *fuego.ContextWithBody[models.PetsCreate]) (models.Pets, error) {
	body, err := c.Body()
	if err != nil {
		return models.Pets{}, err
	}

	return rs.PetsService.CreatePets(body)
}

func (rs PetsRessources) getPets(c fuego.ContextNoBody) (models.Pets, error) {
	id := c.PathParam("id")

	return rs.PetsService.GetPets(id)
}

func (rs PetsRessources) getPetByName(c fuego.ContextNoBody) (models.Pets, error) {
	name := c.PathParam("name")

	return rs.PetsService.GetPetByName(name)
}

func (rs PetsRessources) putPets(c *fuego.ContextWithBody[models.PetsUpdate]) (models.Pets, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.Pets{}, err
	}

	return rs.PetsService.UpdatePets(id, body)
}

func (rs PetsRessources) deletePets(c *fuego.ContextNoBody) (any, error) {
	return rs.PetsService.DeletePets(c.PathParam("id"))
}

type PetsFilter struct {
	Name        string
	YoungerThan int
}

type PetsService interface {
	GetPets(id string) (models.Pets, error)
	GetPetByName(name string) (models.Pets, error)
	CreatePets(models.PetsCreate) (models.Pets, error)
	GetAllPets() ([]models.Pets, error)
	FilterPets(PetsFilter) ([]models.Pets, error)
	GetAllPetsByAge() ([][]models.Pets, error)
	UpdatePets(id string, input models.PetsUpdate) (models.Pets, error)
	DeletePets(id string) (any, error)
}
