package validate

import (
	"github.com/go-playground/validator/v10"
)

// error handling

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(v validator.FieldError) string {
	switch v.Tag() {
	case "required":
		return "Data tidak boleh kosong"
	case "lte":
		return "data harus kurang dari " + v.Param()
	case "gte":
		return "data harus besar dari " + v.Param()
	case "email":
		return "email anda salah, silahkan cek kembali"
	case "min":
		return "data harus lebih besar dari" + v.Param()
	case "max":
		return "data tidak boleh lebih dari" + v.Param()
	}
	return "Input yang dimasukkan salah, silahkan ulangi"
}
