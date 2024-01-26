import request from '@/utils/request'

export function getCategoryList(params) {
  return request({
    url: '/categories',
    method: 'get',
    params
  })
}

export function editCategory(data) {
  return request({
    url: '/categories',
    method: 'post',
    data
  })
}

export function deleteCategory(data) {
  return request({
    url: '/categories',
    method: 'delete',
    data
  })
}
