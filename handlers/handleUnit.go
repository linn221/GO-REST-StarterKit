package handlers

import (
	"linn221/shop/models"
	"net/http"
)

type NewUnit struct {
	Name        inputString     `json:"name" validate:"required,min=2,max=100"`
	Symbol      inputString     `json:"symbol" validate:"required,min=1,max=10"`
	Description *optionalString `json:"description" validate:"omitempty,max=500"`
}

func HandleUnitCreate(unitService *models.UnitService) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewUnit) error {

		unit := models.Unit{
			Name:        input.Name.String(),
			Symbol:      input.Symbol.String(),
			Description: input.Description.StringPtr(),
		}
		unit.ShopId = session.ShopId

		_, err := unitService.Store(r.Context(), &unit, session.ShopId)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		return nil
	})
}

func HandleUnitUpdate(unitService *models.UnitService) http.HandlerFunc {

	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewUnit) error {

		ctx := r.Context()
		_, err := unitService.Update(ctx, &models.Unit{
			Name:        input.Name.String(),
			Symbol:      input.Symbol.String(),
			Description: input.Description.StringPtr(),
		}, session.ResId, session.ShopId)
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}

func HandleUnitDelete(unitService *models.UnitService) http.HandlerFunc {
	return DeleteHandler(func(w http.ResponseWriter, r *http.Request, session Session) error {
		_, err := unitService.Delete(r.Context(), session.ResId, session.ShopId)
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}
