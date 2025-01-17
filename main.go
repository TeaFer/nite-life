package main

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		panic(err)
	}

	if err = store.Init(); err != nil {
		panic(err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
