import { request } from "@/utils/service"
import type * as Collection from "./types/collection"

/** 增 */
export function createCollectionDataApi(data: Collection.CreateOrUpdateCollectionRequestData) {
  return request({
    url: "/collections",
    method: "post",
    data
  })
}

/** 删 */
export function deleteCollectionDataApi(id: string) {
  return request({
    url: `/collections/${id}`,
    method: "delete"
  })
}

/** 改 */
export function updateCollectionDataApi(data: Collection.CreateOrUpdateCollectionRequestData) {
  return request({
    url: "/collections",
    method: "post",
    data
  })
}

/** 查 */
export function getCollectionDataApi(params: Collection.GetCollectionRequestData) {
  return request<Collection.GetCollectionResponseData>({
    url: "/collections",
    method: "get",
    params
  })
}
