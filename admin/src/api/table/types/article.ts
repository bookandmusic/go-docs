export interface CreateOrUpdateArticleRequestData {
  id?: string
  title: string
  content: string
  type: string
  tags: string[]
  categoryId?: number
  collectionId?: number
}

export interface GetArticleRequestData {
  keyword?: string
  type?: string
  /** 查询参数：排序字段 */
  sort?: string
}

export interface GetArticleData {
  UpdatedAt: string
  ID: string
  categoryId: number
  categoryName: string
  collectionId: number
  collectionName: string
  tags: string[]
  tagsNum: number
  title: string
  type: string
}

export type GetArticleResponseData = ApiResponseData<{
  data: GetArticleData[]
  total: number
}>

export type GetArticleDetailResponseData = ApiResponseData<CreateOrUpdateArticleRequestData>
