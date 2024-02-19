import { request } from "@/utils/service"
import type * as Article from "./types/article"

/** 增 */
export function createArticleDataApi(data: Article.CreateOrUpdateArticleRequestData) {
  return request({
    url: "/articles",
    method: "post",
    data
  })
}

/** 删 */
export function deleteArticleDataApi(id: string) {
  return request({
    url: `/articles/${id}`,
    method: "delete"
  })
}

export function bantchDeleteArticleDataApi(ids: string[]) {
  return request({
    url: "/articles",
    method: "delete",
    data: {
      ids: ids
    }
  })
}

/** 改 */
export function updateArticleDataApi(data: Article.CreateOrUpdateArticleRequestData) {
  return request({
    url: "/articles",
    method: "post",
    data
  })
}

/** 查 */
export function getArticleDataApi(params: Article.GetArticleRequestData) {
  return request<Article.GetArticleResponseData>({
    url: "/articles",
    method: "get",
    params
  })
}

/** 查 */
export function getArticleDetailDataApi(id: string) {
  return request<Article.GetArticleDetailResponseData>({
    url: `/articles/${id}`,
    method: "get"
  })
}
