package handler

import (
	tpc "errorHandler/gunk/v1/category"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Category struct {
	ID     int64
	Title  string
	Errors map[string]string
}

func (c *Category) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required.Error("This Filed cannot be blank"), validation.Length(3, 0)),
	)
}

func (h *Handler) CategoryList(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, _ := h.tc.GetAllData(ctx, &tpc.GetAllDataCategoryRequest{})
	if err := h.templates.ExecuteTemplate(rw, "list_category.html", res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Add
func (h *Handler) CategoryCreate(rw http.ResponseWriter, r *http.Request) {

	vErrs := map[string]string{"title": ""}
	title := ""
	h.createCategoryFormData(rw, title, vErrs)
	return

}

func (h *Handler) createCategoryFormData(rw http.ResponseWriter, title string, errs map[string]string) {

	form := Category{
		Title:  title,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "create_Category.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Store
func (h *Handler) CategoryStore(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var category Category
	if err := h.decoder.Decode(&category, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if aErr := category.Validate(); aErr != nil {
		//fmt.Printf("%T", aErr)
		vErrors, ok := aErr.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[key] = value.Error()
			}
			h.createCategoryFormData(rw, category.Title, vErr)
			return
		}

		http.Error(rw, aErr.Error(), http.StatusInternalServerError)
		return
	}

	_, err := h.tc.Create(r.Context(), &tpc.CreateCategoryRequest{
		Category: &tpc.Category{
			Title: category.Title,
		},
	})
	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/category/list", http.StatusTemporaryRedirect)
}

//Delete
func (h *Handler) Delete(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	id, erre := strconv.ParseInt(Id, 10, 64)

	if erre != nil {
		http.Error(rw, erre.Error(), http.StatusInternalServerError)
		return
	}

	_, err := h.tc.Delete(r.Context(), &tpc.DeleteCategoryRequest{
		ID: id,
	})
	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/category/list", http.StatusTemporaryRedirect)
}

//Edit
func (h *Handler) Edit(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	id, erre := strconv.ParseInt(Id, 10, 64)

	if erre != nil {
		http.Error(rw, erre.Error(), http.StatusInternalServerError)
		return
	}

	res, err := h.tc.Get(r.Context(), &tpc.GetCategoryRequest{
		ID: id,
	})
	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	errs := map[string]string{}
	h.editData(rw, res.Category.ID, res.Category.Title, errs)
	return
}

func (h *Handler) editData(rw http.ResponseWriter, id int64, title string, errs map[string]string) {

	form := Category{
		ID:     id,
		Title:  title,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "edit_Category.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Update
func (h *Handler) Update(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	Id := vars["id"]

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	id, erre := strconv.ParseInt(Id, 10, 64)

	if erre != nil {
		http.Error(rw, erre.Error(), http.StatusInternalServerError)
		return
	}

	var category Category
	if err := h.decoder.Decode(&category, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := category.Validate(); err != nil {
		valError, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range valError {
				vErrs[key] = value.Error()
			}
			h.editData(rw, id, category.Title, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := h.tc.Update(r.Context(), &tpc.UpdateCategoryRequest{
		Category: &tpc.Category{
			ID:    id,
			Title: category.Title,
		},
	})
	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/category/list", http.StatusTemporaryRedirect)
}

// func (h *Handler) bookActive(rw http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	Id := vars["id"]

// 	const updateStatusTodo = `UPDATE books SET status = true WHERE id=$1`
// 	res := h.db.MustExec(updateStatusTodo, Id)

// 	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
// 		http.Error(rw, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(rw, r, "/Book/List", http.StatusTemporaryRedirect)
// }

// func (h *Handler) bookDeactivate(rw http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	Id := vars["id"]

// 	const updateStatusTodo = `UPDATE books SET status = false WHERE id=$1`
// 	res := h.db.MustExec(updateStatusTodo, Id)

// 	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
// 		http.Error(rw, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(rw, r, "/Book/List", http.StatusTemporaryRedirect)
// }
