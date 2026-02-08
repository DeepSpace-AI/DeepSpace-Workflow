export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: '缺少 AI Gateway 配置' })
  }

  const id = getRouterParam(event, 'id')
  if (!id) {
    throw createError({ statusCode: 400, statusMessage: '缺少模型ID' })
  }

  const base = aiGateway.url.endsWith('/') ? aiGateway.url.slice(0, -1) : aiGateway.url
  const body = await readBody(event)
  const res = await fetch(`${base}/api/admin/models/${id}`, {
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
        : data?.error?.message || '更新模型失败'
    throw createError({ statusCode: res.status, statusMessage: msg })
  }
  return data
})
