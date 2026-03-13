namespace go microdemo

struct ProfileRequest {
  1: required i64 userId
  2: optional string scenario
}

struct ProfileResponse {
  1: required i64 userId
  2: required string name
  3: required string city
  4: required string bio
}

struct RecommendRequest {
  1: required i64 userId
  2: optional i32 limit = 3
  3: optional string scenario
}

struct RecommendItem {
  1: required i64 postId
  2: required string title
  3: required string reason
}

struct RecommendResponse {
  1: required i64 userId
  2: required list<RecommendItem> items
}

service ProfileService {
  ProfileResponse GetProfile(1: ProfileRequest req)
}

service RecommendService {
  RecommendResponse GetRecommendations(1: RecommendRequest req)
}
