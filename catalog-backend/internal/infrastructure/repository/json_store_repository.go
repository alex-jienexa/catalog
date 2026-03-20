package repository

import (
	"catalog-backend/internal/domain/entity"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type JsonStoreRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewJsonStoreRepository(filePath string) *JsonStoreRepository {
	return &JsonStoreRepository{
		filePath: filePath,
	}
}

func (r *JsonStoreRepository) Get() (*entity.StoreInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Возвращаем дефолтные значения
			return &entity.StoreInfo{
				Title:       "Добро пожаловать в наш каталог!",
				Description: "Мы - современный онлайн-магазин, который предлагает широкий ассортимент товаров различных категорий. Наша цель - сделать покупки удобными и доступными для каждого.",
				ImageURL:    "",
				Features: []entity.Feature{
					{ID: "1", Icon: "⭐", Title: "Качество товаров", Description: "Все товары проходят тщательную проверку перед публикацией"},
					{ID: "2", Icon: "🚚", Title: "Быстрая доставка", Description: "Отправляем товары в день заказа по всей стране"},
					{ID: "3", Icon: "💬", Title: "Поддержка 24/7", Description: "Наша служба поддержки всегда готова помочь вам"},
					{ID: "4", Icon: "🔄", Title: "Легкий возврат", Description: "Простая процедура возврата товара в течение 14 дней"},
				},
			}, nil
		}
		return nil, fmt.Errorf("failed to read store file: %w", err)
	}

	var info entity.StoreInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse store file: %w", err)
	}
	return &info, nil
}

func (r *JsonStoreRepository) Update(info *entity.StoreInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Валидация: Title и Description не должны быть пустыми
	if info.Title == "" || info.Description == "" {
		return fmt.Errorf("title and description are required")
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal store info: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write store file: %w", err)
	}
	return nil
}

func (r *JsonStoreRepository) EnsureDefault(defaultInfo *entity.StoreInfo) error {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return r.Update(defaultInfo)
	}
	return nil
}
