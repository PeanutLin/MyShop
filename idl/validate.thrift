namespace go shop.validate

include "base.thrift"

struct GetValidateReq {
  1: i64 userID
  2: string userCookie
}

struct GetValidateResp {
  1: bool isSuccess
  
  255: base.BaseResp baseResp
}

service ValidateService {
  GetValidateResp GetValidate(1: GetValidateReq req)
}