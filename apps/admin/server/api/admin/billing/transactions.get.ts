export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: '缺少 AI Gateway 配置' })
  }

  const base = aiGateway.url.endsWith('/') ? aiGateway.url.slice(0, -1) : aiGateway.url
  const query = getQuery(event)
  const url = new URL(`${base}/api/admin/billing/transactions`)

  for (const [key, value] of Object.entries(query)) {
    if (typeof value === 'string' && value) {
      url.searchParams.set(key, value)
    } else if (Array.isArray(value)) {
      value.filter((item) => typeof item === 'string' && item).forEach((item) => url.searchParams.append(key, item))
    }
  }

  const res = await fetch(url.toString(), {
    headers: {
      cookie: event.node.req.headers.cookie || ''
    }
  })
  const data = await res.json()
  if (!res.ok) {
    const msg =
      typeof data?.error === 'string'
        ? data.error
        : data?.error?.message || '获取交易流水失败'
    throw createError({ statusCode: res.status, statusMessage: msg })
  }
  return data
})
