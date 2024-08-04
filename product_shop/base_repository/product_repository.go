package base_repository

import (
	"context"
	"productshop/product_shop/middleware/mysql"
	"productshop/product_shop/middleware/mysql/gen/model"
	"productshop/product_shop/middleware/mysql/gen/query"

	"github.com/pkg/errors"
)

type IProductRepository interface {
	// 商品保存
	Insert(ctx context.Context, product *model.Product) (productID int64, err error)
	// 根据商品主键删除商品
	Delete(ctx context.Context, productID int64) (isOK bool, err error)
	// 根据主键查找商品
	SelectByID(ctx context.Context, productID int64) (product *model.Product, err error)
	// 查找所有商品
	SelectAll(ctx context.Context) (products []*model.Product, err error)
	// 减少指定商品的数量
	SubProductNum(ctx context.Context, productID int64, productNum int64) (err error)
}

type ProductRepository struct {
	Q *query.Query
}

// 新建订单管理库
func NewProductRepository() IProductRepository {
	return &ProductRepository{
		Q: mysql.QueryDB,
	}
}

// 商品插入
func (p *ProductRepository) Insert(ctx context.Context, product *model.Product) (productID int64, err error) {
	err = p.Q.Product.WithContext(ctx).Save(product)
	if err != nil {
		return -1, err
	}

	return productID, nil
}

// 根据商品主键删除商品
func (p *ProductRepository) Delete(ctx context.Context, productID int64) (isOK bool, err error) {
	panic("not implent")
}

// 根据商品 ID 查询商品
func (p *ProductRepository) SelectByID(ctx context.Context, productID int64) (product *model.Product, err error) {
	productPO := p.Q.Product
	product, err = productPO.Where(productPO.ID.Eq(int32(productID))).First()
	if err != nil {
		return nil, nil
	}
	return product, nil
}

// 获取所有商品
func (p *ProductRepository) SelectAll(ctx context.Context) (products []*model.Product, err error) {
	products, err = p.Q.Product.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	return products, nil
}

// 减少指定商品数量
func (p *ProductRepository) SubProductNum(ctx context.Context, productID int64, productNum int64) (err error) {
	product, err := p.SelectByID(ctx, productID)
	if err != nil {
		return errors.Wrap(err, "SubProductNum error")
	}
	product.ProductNum -= int32(productNum)
	_, err = p.Insert(ctx, product)
	if err != nil {
		return errors.Wrap(err, "SubProductNum error")
	}
	return nil
}
