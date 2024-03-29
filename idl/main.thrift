namespace go abcp

struct GetUserInfoReq {
    1: required i64 user_id (api.body = "uid")
}

struct UserData {
    1: i64    id;
    2: string name;
}

struct GetUserInfoResp {
    1: required i32      code;
    2: optional string   message;
    3: optional UserData data;
}

service User {
    GetUserInfoResp GetUserInfoMethod(1: GetUserInfoReq req) (api.get = "/user/info", api.serializer = 'json')
}

service User1 {
    GetUserInfoResp GetUserInfo1Method(1: GetUserInfoReq req) (api.get = "/user/info1", api.serializer = 'json')
}
