export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: '缺少 AI Gateway 配置' })
  }

  const base = aiGateway.url.endsWith('/') ? aiGateway.url.slice(0, -1) : aiGateway.url
  const res = await fetch(`${base}/api/admin/models/providers`, {
    headers: {
      cookie: event.node.req.headers.cookie || ''
    }
  })
  const data = await res.json()
  if (!res.ok) {
    const msg =
      typeof data?.error === 'string'
        ? data.error
        : data?.error?.message || '获取模型提供商列表失败'
    throw createError({ statusCode: res.status, statusMessage: msg })
  }
  return data
})
