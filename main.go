package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type data struct {
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

type server struct {
	Storage []data
}

func main() {
	s := server{
		Storage: []data{
			{
				ID:   "1",
				Data: "hello world",
			},
		},
	}

	mux := http.NewServeMux()

	mux.Handle("/getitem", s.getItem())
	mux.Handle("/createitem", s.createItem())
	mux.Handle("/updateitem", s.updateItem())
	mux.Handle("/deleteitem", s.deleteItem())

	log.Print("listening...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("unable to start listener")
	}

}

func (s *server) getItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := r.URL.Query().Get("id")

		dsMap := make(map[string]data, len(s.Storage))

		for _, d := range s.Storage {
			dsMap[d.ID] = d
		}

		d, ok := dsMap[reqID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		res, err := json.MarshalIndent(d, "", "\t")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, err1 := w.Write(res)
		if err1 != nil {
			log.Printf("error writing response %v", err)
		}
	}
}

func (s *server) createItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		var req data

		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		s.Storage = append(s.Storage, req)

		res, err := json.MarshalIndent(req, "", "\t")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err1 := w.Write(res)
		if err1 != nil {
			log.Printf("error writing response %v", err)
		}
	}
}

func (s *server) updateItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		var req data

		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		var index int
		var found bool

		for i, item := range s.Storage {
			if item.ID == req.ID {
				index = i
				found = true
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		s.Storage[index] = data{
			ID:   req.ID,
			Data: req.Data,
		}

		res, err := json.MarshalIndent(req, "", "\t")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, err1 := w.Write(res)
		if err1 != nil {
			log.Printf("error writing response %v", err)
		}
	}
}

func (s *server) deleteItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := r.URL.Query().Get("id")

		var index int
		var found bool

		for i, d := range s.Storage {
			if d.ID == reqID {
				index = i
				found = true
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		s.Storage = append(s.Storage[:index], s.Storage[index+1:]...)

		w.WriteHeader(http.StatusOK)
	}
}
