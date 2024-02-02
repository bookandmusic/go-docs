export interface ArticleInfoData {
  article_count: number
  tag_count: number
  category_count: number
  collection_count: number
}

export type ArticleInfoResponseData = ApiResponseData<ArticleInfoData>
