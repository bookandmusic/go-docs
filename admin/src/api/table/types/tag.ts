export interface CreateOrUpdateTagRequestData {
  id?: number
  name: string
}

export interface GetTagRequestData {
  keyword?: string
  /** 查询参数：排序字段 */
  sort?: string
}

export interface GetTagData {
  UpdatedAt: string
  ID: number
  name: string
  num: number
}

export type GetTagResponseData = ApiResponseData<GetTagData[]>
