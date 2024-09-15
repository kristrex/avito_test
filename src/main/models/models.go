package models

import (
	"time"

	"gorm.io/gorm"
)

var Db *gorm.DB

// Основные структуры:
// Employee - cтруктура, содержащая данные о пользователе
type Employee struct {
	ID        string `json:"id"`
	Username  string `json:"username" gorm:"column:username"`
	FirstName string `json:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name" gorm:"column:last_name"`
	CreatedAt string `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt string `json:"updatedAt" gorm:"column:updated_at"`
}

// OrganizationResponsible - структура, содержащая данные из связующей таблицы
type OrganizationResponsible struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId" gorm:"column:organization_id"`
	UserID         string `json:"user_Id" gorm:"column:user_id"`
}

// Tender - структура, содержащая данные о тендерах
type Tender struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceType    string    `json:"serviceType" gorm:"column:servicetype"`
	Status         string    `json:"status"`
	OrganizationID string    `json:"organizationId" gorm:"column:organizationid"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"createdAt" gorm:"column:createdat"`
}

// Bids - структура, содержащая данные о предложениях
type Bids struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TenderID    string    `json:"tenderId" gorm:"column:tenderid"`
	AuthorType  string    `json:"authorType" gorm:"column:authortype"`
	AuthorID    string    `json:"authorId" gorm:"column:authorid"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:createdat"`
}

// Вспомогательные структуры
// CreationTender - Структура для создания тендера
type CreationTender struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"serviceType" gorm:"column:servicetype"`
	Status          string    `json:"status"`
	OrganizationID  string    `json:"organizationId" gorm:"column:organizationid"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:createdat"`
	CreatorUsername string    `json:"creatorUsername"`
}

// TenderResponse - структура, содержащая данные для ответа
type TenderResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ServiceType string `json:"serviceType"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
}

// BidsResponse - структура, содержащая данные для ответа
type BidsResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	AuthorType string `json:"authorType" gorm:"column:authortype"`
	AuthorID   string `json:"authorId" gorm:"column:authorid"`
	Version    int    `json:"version"`
	CreatedAt  string `json:"createdAt"`
}
