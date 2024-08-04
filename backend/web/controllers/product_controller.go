package controllers

import (
	"fmt"
	"productshop/common"
	"productshop/datamodels"
	"productshop/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ProductController struct {
	// 请求上下文
	Ctx iris.Context
	// Product 服务
	ProductService services.IProductService
}

const TagName = "productshop"

func (p *ProductController) GetAll() mvc.View {
	products, err := p.ProductService.GetAllProduct()
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	fmt.Println(errMsg)
	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			//"errMsg":     errMsg,
			"productArray": products,
		},
	}
}

func (p *ProductController) Get() mvc.View {
	products, err := p.ProductService.GetAllProduct()
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	fmt.Println(errMsg)
	return mvc.View{
		Name: "product/product.html",
		Data: iris.Map{
			//"errMsg":     errMsg,
			"productArray": products,
		},
	}
}

// 修改商品
func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: TagName})

	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

// 添加商品
// http://localhost:8081/product/add
func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name: "product/add.html",
	}
}

// 删除商品
// http://localhost:8081/product/add
func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.ProductService.DeleteProductByID(id)

	p.Ctx.Redirect("/product/all")
}

// 添加商品 POST 方法
func (p *ProductController) PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "productshop"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(id)

	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}
