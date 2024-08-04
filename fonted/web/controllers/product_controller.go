package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"productshop/common"
	"productshop/datamodels"
	"productshop/services"
	"strconv"
	"text/template"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

var (
	htmlOutPath  = "./fonted/web/htmlProductShow/" // 生成 Html保存目录
	templatePath = "./fonted/web/views/template/"  // 静态文件模板目录
)

// 生成静态化页面
func (p *ProductController) GetGenerateHtml() {
	// 1.获取模板
	template, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 2. 获取 html 生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")

	// 3.获取模板渲染数据
	productId, err := p.Ctx.URLParamInt64("productID")
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(productId)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//4.生成静态文件
	generateStaticHtml(p.Ctx, template, fileName, product)
}

// 生成静态化页面
func generateStaticHtml(ctx iris.Context, template *template.Template,
	fileName string, product *datamodels.Product) {
	// 1.判断文件是否存在
	if common.IsFileExist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Debug(err)
		}
	}
	// 2.生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	defer file.Close()
	template.Execute(file, &product)
}

// 门户网站
// http://localhost:8080/product/
func (p *ProductController) Get() mvc.View {
	ProductArray, err := p.ProductService.GetAllProduct()
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Layout: "",
		Name:   "product/menhu.html",
		Data: iris.Map{
			"product0": ProductArray[0],
			"product1": ProductArray[1],
			"product2": ProductArray[2],
			"product3": ProductArray[3],
			"product4": ProductArray[4],
			"product5": ProductArray[5],
			"product6": ProductArray[6],
			"product7": ProductArray[7],
		},
	}
}

// 商品的详细信息
// http://localhost:8080/product/detail?productID=1
func (p *ProductController) GetDetail() mvc.View {
	productId, err := p.Ctx.URLParamInt64("productID")
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProductByID(productId)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Layout: "",
		Name:   "product/product.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// 生成订单
// http://localhost:8080/product/order?productID=1
func (p *ProductController) GetOrder() string {
	// 获取商品 ID
	productId, err := p.Ctx.URLParamInt64("productID")
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	fmt.Println("productID", productId)

	hostURL := "http://" + common.ValidateHost1 + ":" + common.ValidatePort + "/onsale?productID=" + strconv.FormatInt(productId, 10)
	fmt.Println("validate : ", hostURL)
	response, body, err := common.GetCurl(hostURL, p.Ctx.Request())
	if err != nil {
		return "server error"
	}

	// 判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return "true"
		} else {
			return "false"
		}
	} else {
		return "false"
	}

}
