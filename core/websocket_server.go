package main

type Server struct {
	pattern string // Адрес сервера.
	clients map[int]Client
}
