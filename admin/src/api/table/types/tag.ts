export interface CreateOrUpdateTagRequestData {
  id?: string
  name: string
}

export interface GetTagRequestData {
  keyword?: string
  /** 查询参数：排序字段 */
  sort?: string
}

export interface GetTagData {
  UpdatedAt: string
  ID: string
  name: string
  num: number
}

export type GetTagResponseData = ApiResponseData<GetTagData[]>
