export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: '缺少 AI Gateway 配置' })
  }

  const base = aiGateway.url.endsWith('/') ? aiGateway.url.slice(0, -1) : aiGateway.url
  const method = event.node.req.method || 'GET'

  if (method === 'GET') {
    const res = await fetch(`${base}/api/admin/plans`, {
      headers: {
        cookie: event.node.req.headers.cookie || ''
      }
    })
    const data = await res.json()
    if (!res.ok) {
      const msg =
        typeof data?.error === 'string'
          ? data.error
          : data?.error?.message || '获取套餐列表失败'
      throw createError({ statusCode: res.status, statusMessage: msg })
    }
    return data
  }

  if (method === 'POST') {
    const body = await readBody(event)
    const res = await fetch(`${base}/api/admin/plans`, {
      method: 'POST',
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
          : data?.error?.message || '创建套餐失败'
      throw createError({ statusCode: res.status, statusMessage: msg })
    }
    return data
  }

  throw createError({ statusCode: 405, statusMessage: '不支持的请求方法' })
})
