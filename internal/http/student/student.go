package student

import (
	"encoding/json"
	"errors"
	"firstproject/internal/storage"
	"firstproject/internal/types"
	"firstproject/internal/utils/response"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		err:=json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GenrealError(err))
			return
		}
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GenrealError(err))
			return
		}
		if err:=validator.New().Struct(student);err!=nil{
			validateErr:=err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,validateErr)
		}
		lastId,err:=storage.CreateStudent(student.Name,student.Email,student.Age)
        if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		} 
		slog.Info("user careated ",slog.String("userId",fmt.Sprint(lastId)))

		response.WriteJson(w,http.StatusAccepted,map[string]int64{"id":lastId})
	}
}

func GetById(storage storage.Storage)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
        id:=r.PathValue("id")
		slog.Info("gettting a student ",slog.String("id",id))
         intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenrealError(err))
			return
		}

		student, err := storage.GetstudentById(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GenrealError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}