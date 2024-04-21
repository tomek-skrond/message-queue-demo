package main

type APIServer struct {
	DB         *Storage
	listenPort string
}

func NewAPIServer(lp string, db *Storage) (*APIServer, error) {
	return &APIServer{
		DB:         db,
		listenPort: lp,
	}, nil
}
