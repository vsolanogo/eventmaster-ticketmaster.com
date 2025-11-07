package services

import (
    "errors"

    "eventmaster-go/internal/models"
    "eventmaster-go/internal/repositories"

    "gorm.io/gorm"
)

const ticketmasterSystemEmail = "ticketmaster@eventmaster.local"
const ticketmasterSystemPassword = "ticketmaster123"

// EnsureTicketmasterSystemUser makes sure a dedicated system user exists for imported events.
func EnsureTicketmasterSystemUser(userRepo repositories.UserRepository) (string, error) {
    user, err := userRepo.FindByEmail(ticketmasterSystemEmail)
    if err == nil {
        return user.ID, nil
    }

    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return "", err
    }

    systemUser := &models.User{
        Email:    ticketmasterSystemEmail,
        Password: ticketmasterSystemPassword,
    }

    if err := userRepo.Create(systemUser); err != nil {
        return "", err
    }

    return systemUser.ID, nil
}
