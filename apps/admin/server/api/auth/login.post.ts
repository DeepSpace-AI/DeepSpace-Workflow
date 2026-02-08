export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: '缺少 AI Gateway 配置' })
  }

  const body = await readBody(event)
  const base = aiGateway.url.endsWith('/') ? aiGateway.url.slice(0, -1) : aiGateway.url

  const res = await fetch(`${base}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body ?? {})
  })

  const setCookie = res.headers.get('set-cookie')
  if (setCookie) {
    setHeader(event, 'set-cookie', setCookie)
  }

  const data = await res.json()
  if (!res.ok) {
    const msg =
      typeof data?.error === 'string'
        ? data.error
        : data?.error?.message || '登录失败'
    throw createError({ statusCode: res.status, statusMessage: msg })
  }

  return data
})
