import { request } from "@/utils/service"
import type * as Tag from "./types/tag"

/** 增 */
export function createTagDataApi(data: Tag.CreateOrUpdateTagRequestData) {
  return request({
    url: "/tags",
    method: "post",
    data
  })
}

/** 删 */
export function deleteTagDataApi(id: string) {
  return request({
    url: `/tags/${id}`,
    method: "delete"
  })
}

/** 改 */
export function updateTagDataApi(data: Tag.CreateOrUpdateTagRequestData) {
  return request({
    url: "/tags",
    method: "post",
    data
  })
}

/** 查 */
export function getTagDataApi(params: Tag.GetTagRequestData) {
  return request<Tag.GetTagResponseData>({
    url: "/tags",
    method: "get",
    params
  })
}
