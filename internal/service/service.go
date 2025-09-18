package service

type Service struct {
	storage Storage
	cache   Cache
	queue   Queue
	sender  Sender
}

func New(storage Storage, cache Cache, queue Queue, sender Sender) *Service {
	return &Service{
		storage: storage,
		cache:   cache,
		queue:   queue,
		sender:  sender,
	}
}
