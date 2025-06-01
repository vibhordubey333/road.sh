package usecases

import (
    "vibhordubey333/road.sh/broadcast-server/internal/domain/entities"
    "vibhordubey333/road.sh/broadcast-server/internal/domain/repositories"
)

// BroadcastUseCase handles broadcasting messages to clients
type BroadcastUseCase struct {
    clientRepo repositories.ClientRepository
}

// NewBroadcastUseCase creates a new broadcast use case
func NewBroadcastUseCase(clientRepo repositories.ClientRepository) *BroadcastUseCase {
    return &BroadcastUseCase{
        clientRepo: clientRepo,
    }
}

// BroadcastMessage broadcasts a message to all connected clients
func (uc *BroadcastUseCase) BroadcastMessage(message entities.Message) error {
    return uc.clientRepo.Broadcast(message)
}

// RegisterClient registers a new client
func (uc *BroadcastUseCase) RegisterClient(client repositories.Client) error {
    return uc.clientRepo.Add(client)
}

// UnregisterClient removes a client
func (uc *BroadcastUseCase) UnregisterClient(client repositories.Client) error {
    return uc.clientRepo.Remove(client)
}

// ClientCount returns the number of connected clients
func (uc *BroadcastUseCase) ClientCount() int {
    return uc.clientRepo.Count()
}