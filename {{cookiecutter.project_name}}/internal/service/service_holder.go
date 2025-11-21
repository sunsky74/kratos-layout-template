package service

type Holder struct {
	GreeterService *GreeterService
}

func NewServiceHolder(gs *GreeterService) *Holder {
	return &Holder{
		GreeterService: gs,
	}
}
