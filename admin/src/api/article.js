import request from '@/utils/request'

export function getArticleList(params) {
  return request({
    url: '/articles',
    method: 'get',
    params
  })
}

export function getArticleDetail(articleId) {
  return request({
    url: `/articles/${articleId}`,
    method: 'get'
  })
}

export function editArticle(data) {
  return request({
    url: '/articles',
    method: 'post',
    data
  })
}

export function deleteArticle(articleId) {
  return request({
    url: `/articles/${articleId}`,
    method: 'delete'
  })
}

export function importeArticleZip(data) {
  return request({
    url: '/articles/import',
    method: 'post',
    data
  })
}
