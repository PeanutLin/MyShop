package repositories

import (
	"productshop/datamodels"
	"productshop/db"

	"gorm.io/gorm"
)

type IProduct interface {
	// 连接数据库
	Conn() error
	// 商品插入
	Insert(*datamodels.Product) (int64, error)
	// 根据商品主键删除商品
	Delete(int64) bool
	// 商品更新
	Update(*datamodels.Product) error
	// 根据主键查找商品
	SelectByKey(int64) (*datamodels.Product, error)
	// 查找所有商品
	SelectAll() ([]*datamodels.Product, error)
	// 减少指定商品数量
	SubProductNum(int64, int64) error
}

type ProductManager struct {
	sqlConn *gorm.DB
}

// 新建商品管理接口
func NewProductManager(db *gorm.DB) IProduct {
	return &ProductManager{
		sqlConn: db,
	}
}

// 数据库连接
func (p *ProductManager) Conn() error {
	if p.sqlConn == nil {
		mysql, err := db.NewMysqlConn()
		if err != nil {
			return err
		}
		p.sqlConn = mysql
	}
	return nil
}

// 商品插入
func (p *ProductManager) Insert(product *datamodels.Product) (int64, error) {
	if err := p.Conn(); err != nil {
		return -1, err
	}
	err := p.sqlConn.Create(product).Error
	if err != nil {
		return -1, err
	}
	return product.ID, nil
}

// 根据商品主键删除商品
func (p *ProductManager) Delete(productID int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}

	err := p.sqlConn.Where("ID=?", productID).Delete(&datamodels.Product{}).Error
	return err != nil
}

// 商品更新
func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}
	err := p.sqlConn.Save(product).Error
	if err != nil {
		return err
	}
	return nil
}

// 根据商品 ID 查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	productResult = &datamodels.Product{}
	err = p.sqlConn.Where("ID=?", productID).First(productResult).Error
	if err != nil {
		return nil, err
	}
	return productResult, nil
}

// 获取所有商品
func (p *ProductManager) SelectAll() (products []*datamodels.Product, err error) {
	if err = p.Conn(); err != nil {
		return nil, err
	}
	var results []datamodels.Product
	err = p.sqlConn.Find(&results).Error
	if err != nil {
		return nil, err
	}
	for _, product := range results {
		products = append(products, &product)
	}

	return products, nil
}

// 减少指定商品数量
func (p *ProductManager) SubProductNum(productID int64, num int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	product := &datamodels.Product{}
	err := p.sqlConn.Where("ID=?", productID).Take(product).Error
	if err != nil {
		return err
	}

	err = p.sqlConn.Model(product).Update("productNum", product.ProductNum-num).Error
	if err != nil {
		return err
	}
	return nil
}
