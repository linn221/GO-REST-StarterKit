package services

// type ServiceError struct {
// 	Err  error
// 	Code int
// }

// func (se *ServiceError) Error() string {
// 	return se.Error()
// }

// func (se *ServiceError) Respond(w http.ResponseWriter) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(se.Code)
// 	return json.NewEncoder(w).Encode(map[string]any{
// 		"error": se.Err.Error(),
// 	})
// }
// func SystemErr(err error) *ServiceError {
// 	return &ServiceError{
// 		Err:  err,
// 		Code: http.StatusInternalServerError,
// 	}
// }
// func SystemErrString(s string, args ...error) *ServiceError {
// 	if len(args) > 0 {
// 		err := args[0]
// 		return SystemErr(errors.New(s + ": " + err.Error()))
// 	}
// 	return SystemErr(errors.New(s))
// }
// func ClientErr(s string) *ServiceError {
// 	return &ServiceError{
// 		Err:  errors.New(s),
// 		Code: http.StatusBadRequest,
// 	}
// }
