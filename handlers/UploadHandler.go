package handlers

import (
	"linn221/shop/services"
	"net/http"
)

func HandleImageUploadSingle(imageService services.ImageUploader) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		uri, errs := imageService.UploadSingle(r, "image")
		if errs != nil {
			handleError(errs.Respond(w), w)
			return
		}

		handleError(respondOk(w, map[string]any{
			"url": uri,
		}), w)
	}

}

func handleError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
