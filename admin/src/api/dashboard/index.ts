import { request } from "@/utils/service"
import type * as Dashboard from "./types/index"

/** 获取dashboard数据 */
export function getArticleInfoApi() {
  return request<Dashboard.ArticleInfoResponseData>({
    url: "/dashboard",
    method: "get"
  })
}
