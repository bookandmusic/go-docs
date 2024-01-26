import request from '@/utils/request'

export function getCollectionList(params) {
  return request({
    url: '/collections',
    method: 'get',
    params
  })
}

export function editCollection(data) {
  return request({
    url: '/collections',
    method: 'post',
    data
  })
}

export function deleteCollection(data) {
  return request({
    url: '/collections',
    method: 'delete',
    data
  })
}
