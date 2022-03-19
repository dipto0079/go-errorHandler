package handler

import (
	tpb "errorHandler/gunk/v1/blog"

	"net/http"
)

func (h *Handler) Home(rw http.ResponseWriter, r *http.Request) {
	blogData, err := h.tb.ListBlog(r.Context(), &tpb.ListBlogRequest{})
	//fmt.Printf("=============data===================",blogData)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(rw, "index.html", blogData); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
