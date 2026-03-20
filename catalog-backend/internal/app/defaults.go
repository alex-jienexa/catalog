package app

import "catalog-backend/internal/domain/entity"

func defaultStoreInfo() *entity.StoreInfo {
	return &entity.StoreInfo{
		Title:       "🛒 Добро пожаловать в наш каталог!",
		Description: "Мы - современный онлайн-магазин, который предлагает широкий ассортимент товаров различных категорий. Наша цель - сделать покупки удобными и доступными для каждого.",
		ImageURL:    "https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?ixlib=rb-1.2.1&auto=format&fit=crop&w=700&q=80",
		Features: []entity.Feature{
			{Icon: "⭐", Title: "Качество товаров", Description: "Все товары проходят тщательную проверку перед публикацией"},
			{Icon: "🚚", Title: "Быстрая доставка", Description: "Отправляем товары в день заказа по всей стране"},
			{Icon: "💬", Title: "Поддержка 24/7", Description: "Наша служба поддержки всегда готова помочь вам"},
			{Icon: "🔄", Title: "Легкий возврат", Description: "Простая процедура возврата товара в течение 14 дней"},
		},
	}
}
