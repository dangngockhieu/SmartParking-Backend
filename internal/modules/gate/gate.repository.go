package gate

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create thêm mới một cổng vào database
func (r *Repository) Create(gate *Gate) error {
	return r.db.Create(gate).Error
}

// FindByID tìm cổng theo ID
func (r *Repository) FindByID(id uint64) (*Gate, error) {
	var gate Gate
	err := r.db.First(&gate, id).Error
	if err != nil {
		return nil, err
	}
	return &gate, nil
}

// Update cập nhật thông tin cổng
func (r *Repository) Update(gate *Gate) error {
	return r.db.Save(gate).Error
}
