namespace go shop.product

include "base.thrift"

struct GetProductReq {
  1: i64 userID
  2: i64 productID
  3: i64 productNum
}

struct GetProductResp {
  1: bool isSuccess

  255: base.BaseResp baseResp
}


service ProductService {
  GetProductResp GetProduct(1: GetProductReq req)
}