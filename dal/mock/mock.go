//go:generate mockgen -package=mock -source=../client_provider.go -destination=client_provider.go
//go:generate mockgen -package=mock -source=../user_provider.go -destination=user_provider.go
//go:generate mockgen -package=mock -source=../user_service.go -destination=user_service.go

package mock
