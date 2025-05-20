package response

import (
	"encoding/json"
	"firstproject/internal/types"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)
 const (
	StatusOk="OK"
	StatusError="Error"
 )

func WriteJson(w http.ResponseWriter,status int,data interface{})error{
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GenrealError(err error) types.Response {
 return types.Response{
	Status: StatusError,
	Error: err.Error(),
 }
}
func ValidationError( errs validator.ValidationErrors) types.Response{
	var errMesgs[]string
	for _,err:=range errs{
		switch err.ActualTag(){
	    case "required":
			errMesgs=append(errMesgs, fmt.Sprintf("filed %s is required",err.Field()))
		default:
		     errMesgs=append(errMesgs, fmt.Sprintf("filed %s is invalid",err.Field()))
		
		}}
      return types.Response{
		Status: StatusError,
		Error: strings.Join(errMesgs,",") ,
	  }
	}
