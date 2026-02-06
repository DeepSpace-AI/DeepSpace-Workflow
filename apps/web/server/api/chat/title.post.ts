import { generateText } from 'ai'
import { createOpenAICompatible } from '@ai-sdk/openai-compatible'

function sanitizeTitle(raw: string) {
  const cleaned = raw
    .replace(/[\r\n]+/g, ' ')
    .replace(/["'“”‘’]/g, '')
    .replace(/[。！？.!?]+$/g, '')
    .trim()

  if (!cleaned) return ''
  return cleaned.length > 30 ? cleaned.slice(0, 30) : cleaned
}

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig()
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: 'Missing AI Gateway config' })
  }

  const body = await readBody<{ text?: string; model?: string }>(event)
  const text = String(body?.text || '').trim()
  const model = String(body?.model || 'deepseek-chat').trim() || 'deepseek-chat'

  if (!text) {
    throw createError({ statusCode: 400, statusMessage: 'text is required' })
  }

  const openai = createOpenAICompatible({
    name: 'newapi',
    baseURL: aiGateway.url.endsWith('/v1') ? aiGateway.url : `${aiGateway.url}/v1`,
    headers: {
      cookie: event.node.req.headers.cookie || '',
    },
  })

  const prompt = text.length > 500 ? `${text.slice(0, 500)}...` : text

  const result = await generateText({
    model: openai(model),
    system:
      '你是对话标题生成器。请基于用户首条消息生成一个简洁中文标题。要求：8-18字，不带引号，不要句号，不要序号，只返回标题文本。',
    prompt,
    maxOutputTokens: 32,
    temperature: 0.2,
  })

  const title = sanitizeTitle(result.text || '')
  if (!title) {
    throw createError({ statusCode: 422, statusMessage: 'title generation failed' })
  }

  return { title }
})
