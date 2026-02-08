export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: '缺少 AI Gateway 配置' })
  }

  const id = getRouterParam(event, 'id')
  if (!id) {
    throw createError({ statusCode: 400, statusMessage: '缺少用户ID' })
  }

  const base = aiGateway.url.endsWith('/') ? aiGateway.url.slice(0, -1) : aiGateway.url
  const method = event.node.req.method || 'GET'
  const url = `${base}/api/admin/users/${id}`

  if (method === 'GET') {
    const res = await fetch(url, {
      headers: {
        cookie: event.node.req.headers.cookie || ''
      }
    })
    const data = await res.json()
    if (!res.ok) {
      const msg =
        typeof data?.error === 'string'
          ? data.error
          : data?.error?.message || '获取用户详情失败'
      throw createError({ statusCode: res.status, statusMessage: msg })
    }
    return data
  }

  if (method === 'PATCH') {
    const body = await readBody(event)
    const res = await fetch(url, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        cookie: event.node.req.headers.cookie || ''
      },
      body: JSON.stringify(body ?? {})
    })
    const data = await res.json()
    if (!res.ok) {
      const msg =
        typeof data?.error === 'string'
          ? data.error
          : data?.error?.message || '更新用户失败'
      throw createError({ statusCode: res.status, statusMessage: msg })
    }
    return data
  }

  if (method === 'DELETE') {
    const res = await fetch(url, {
      method: 'DELETE',
      headers: {
        cookie: event.node.req.headers.cookie || ''
      }
    })
    const data = await res.json()
    if (!res.ok) {
      const msg =
        typeof data?.error === 'string'
          ? data.error
          : data?.error?.message || '删除用户失败'
      throw createError({ statusCode: res.status, statusMessage: msg })
    }
    return data
  }

  throw createError({ statusCode: 405, statusMessage: '不支持的请求方法' })
})
