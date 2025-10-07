package service

import "sync"

var (
	srv  Service
	once sync.Once
)

type Service interface {
	Run()
}
