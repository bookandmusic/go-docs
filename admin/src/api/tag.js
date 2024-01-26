import request from '@/utils/request'

export function getTagList(params) {
  return request({
    url: '/tags',
    method: 'get',
    params
  })
}

export function editTag(data) {
  return request({
    url: '/tags',
    method: 'post',
    data
  })
}

export function deleteTag(data) {
  return request({
    url: '/tags',
    method: 'delete',
    data
  })
}
