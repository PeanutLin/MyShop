namespace go shop.sale

include "base.thrift"

struct GetSaleReq {
  1: i64 userID
  2: i64 productID
  3: i64 productNum
  4: string userCookie
}

struct GetSaleResp {
  1: bool isSuccess
  
  255: base.BaseResp baseResp
}

service SaleService {
  GetSaleResp GetSale(1: GetSaleReq req)
}