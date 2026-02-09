export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: '缺少 AI Gateway 配置' })
  }

  const body = await readBody(event)
  const base = aiGateway.url.endsWith('/') ? aiGateway.url.slice(0, -1) : aiGateway.url
  const res = await fetch(`${base}/api/admin/risk/rate-limits`, {
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
        : data?.error?.message || '创建速率限制失败'
    throw createError({ statusCode: res.status, statusMessage: msg })
  }
  return data
})
