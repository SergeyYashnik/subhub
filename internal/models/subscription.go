package models

import "time"

type Subscription struct {
	ID          string     `json:"id" db:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string     `json:"service_name" db:"service_name"  example:"Yandex Plus"`
	Price       int        `json:"price" db:"price" example:"400"`
	UserID      string     `json:"user_id" db:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
}

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"400"`
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"08-2025"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"400"`
	StartDate   string `json:"start_date" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"08-2025"`
}
