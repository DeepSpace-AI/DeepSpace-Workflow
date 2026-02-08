export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: '缺少 AI Gateway 配置' })
  }

  const base = aiGateway.url.endsWith('/') ? aiGateway.url.slice(0, -1) : aiGateway.url
  const res = await fetch(`${base}/api/admin/models/sync`, {
    method: 'POST',
    headers: {
      cookie: event.node.req.headers.cookie || ''
    }
  })
  const data = await res.json()
  if (!res.ok) {
    const msg =
      typeof data?.error === 'string'
        ? data.error
        : data?.error?.message || '同步上游模型失败'
    throw createError({ statusCode: res.status, statusMessage: msg })
  }
  return data
})
