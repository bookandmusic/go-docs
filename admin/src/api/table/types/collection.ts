export interface CreateOrUpdateCollectionRequestData {
  id?: string
  name: string
  author: string
}

export interface GetCollectionRequestData {
  keyword?: string
  /** 查询参数：排序字段 */
  sort?: string
}

export interface GetCollectionData {
  UpdatedAt: string
  ID: string
  name: string
  author: string
  num: number
}

export type GetCollectionResponseData = ApiResponseData<GetCollectionData[]>
