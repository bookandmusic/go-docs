import { request } from "@/utils/service"
import type * as Category from "./types/category"

/** 增 */
export function createCategoryDataApi(data: Category.CreateOrUpdateCategoryRequestData) {
  return request({
    url: "/categories",
    method: "post",
    data
  })
}

/** 删 */
export function deleteCategoryDataApi(id: string) {
  return request({
    url: `/categories/${id}`,
    method: "delete"
  })
}

export function bantchDeleteCategoryDataApi(ids: string[]) {
  return request({
    url: "/categories",
    method: "delete",
    data: {
      ids: ids
    }
  })
}

/** 改 */
export function updateCategoryDataApi(data: Category.CreateOrUpdateCategoryRequestData) {
  return request({
    url: "/categories",
    method: "post",
    data
  })
}

/** 查 */
export function getCategoryDataApi(params: Category.GetCategoryRequestData) {
  return request<Category.GetCategoryResponseData>({
    url: "/categories",
    method: "get",
    params
  })
}
