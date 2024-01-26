export interface CreateOrUpdateCategoryRequestData {
  id?: string
  name: string
}

export interface GetCategoryRequestData {
  keyword?: string
  /** 查询参数：排序字段 */
  sort?: string
}

export interface GetCategoryData {
  UpdatedAt: string
  ID: string
  name: string
  num: number
}

export type GetCategoryResponseData = ApiResponseData<GetCategoryData[]>
