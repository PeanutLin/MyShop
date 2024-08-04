package base_service

import (
	"context"
	"productshop/product_shop/base_repository"
	"productshop/product_shop/middleware/mysql/gen/model"
)

type IProductService interface {
	GetProductByID(ctx context.Context, productID int64) (product *model.Product, err error)
	GetAllProduct(ctx context.Context) ([]*model.Product, error)
	// DeleteProductByID(int64) bool
	InsertProduct(ctx context.Context, product *model.Product) (int64, error)
	// UpdateProduct(product *model.Product) error
	SubNumber(ctx context.Context, productID int64, productNum int64) error
}

type ProductService struct {
	productRepository base_repository.IProductRepository
}

// 初始化函数
func NewProductService(repository base_repository.IProductRepository) IProductService {
	return &ProductService{
		productRepository: repository,
	}
}

func (p *ProductService) GetProductByID(ctx context.Context, productID int64) (product *model.Product, err error) {
	return p.productRepository.SelectByID(ctx, productID)
}

func (p *ProductService) GetAllProduct(ctx context.Context) ([]*model.Product, error) {
	return p.productRepository.SelectAll(ctx)
}

// func (p *ProductService) DeleteProductByID(productID int64) bool {
// 	return p.productRepository.Delete(productID)
// }

func (p *ProductService) InsertProduct(ctx context.Context, product *model.Product) (int64, error) {
	return p.productRepository.Insert(ctx, product)
}

// func (p *ProductService) UpdateProduct(product *model.Product) error {
// 	return p.productRepository.Update(product)
// }

func (p *ProductService) SubNumber(ctx context.Context, productID int64, productNum int64) error {
	return p.productRepository.SubProductNum(ctx, productID, productNum)
}
