import { Chat } from '@ai-sdk/vue'

export function useEditorAi() {
    const aiBusy = ref(false)
    const showAiPreview = ref(false)
    const aiPreviewOriginal = ref('')
    const aiPreviewResult = ref('')
    const aiPreviewHadSelection = ref(false)
    const aiPreviewType = ref<'polish' | 'expand' | 'summary'>('polish')
    const aiSelectionRange = ref<{ from: number; to: number } | null>(null)
    const aiStreamActive = ref(false)
    const toast = useToast()

    // 编辑器引用
    const editorRef = ref<any>(null)

    // 创建编辑器专用Chat实例
    const editorChat = new Chat({
        onError(error: any) {
            const status = error?.status || error?.cause?.status || error?.response?.status
            const rawMessage = String(error?.message || '')
            const lower = rawMessage.toLowerCase()
            let message = rawMessage || 'AI 请求失败'
            if (status === 402 || lower.includes('insufficient balance')) {
                message = '余额不足，请先充值后再试'
            } else if (status === 409 || lower.includes('ref_id conflict')) {
                message = '计费信息冲突，请稍后重试'
            }
            toast.add({ title: message, color: 'red' })
            console.error('Editor AI error:', error)
        },
        onFinish: ({ message, isError, isDisconnect, isAbort }) => {
            if (isError || isDisconnect || isAbort) {
                const parts = Array.isArray(message?.parts) ? message.parts : []
                const errorText = parts
                    .map((part: any) => (typeof part?.text === 'string' ? part.text : ''))
                    .join(' ')
                    .trim()
                const lower = errorText.toLowerCase()
                let display = errorText || 'AI 响应失败'
                if (lower.includes('payment required') || lower.includes('insufficient balance')) {
                    display = '余额不足，请先充值后再试'
                }
                toast.add({ title: display, color: 'red' })
                return
            }
            if (!showAiPreview.value) return
            
            const parts = Array.isArray(message?.parts) ? message.parts : []
            const text = parts
                .filter((part: any) => part?.type === 'text' && typeof part.text === 'string')
                .map((part: any) => part.text)
                .join('')
            
            aiPreviewResult.value = text
            editorRef.value?.updateAiPreview?.({ text, loading: false })
        },
    })

    // 监听消息流式更新
    watch(
        () => editorChat.messages,
        (messages) => {
            if (!showAiPreview.value || !aiStreamActive.value) return
            const lastAssistant = [...messages].reverse().find((item) => item?.role === 'assistant')
            if (!lastAssistant) return
            
            const parts = Array.isArray(lastAssistant?.parts) ? lastAssistant.parts : []
            const text = parts
                .filter((part: any) => part?.type === 'text' && typeof part.text === 'string')
                .map((part: any) => part.text)
                .join('')
            
            if (!text) return
            
            // 更新结果
            aiPreviewResult.value = text
            
            // 更新编辑器预览
            if (editorRef.value?.updateAiPreview) {
                console.log('[EditorAI] Updating preview with text length:', text.length, 'streaming:', editorChat.status === 'streaming')
                editorRef.value.updateAiPreview({ 
                    text, 
                    loading: editorChat.status === 'streaming' 
                })
            } else {
                console.warn('[EditorAI] editorRef.updateAiPreview not available')
            }
        },
        { deep: true, immediate: true }
    )

    // 监听Chat状态变化并更新预览
    watch(
        () => editorChat.status,
        (status) => {
            if (!aiStreamActive.value || !showAiPreview.value) return
            if (status === 'streaming' && aiPreviewResult.value) {
                if (editorRef.value?.updateAiPreview) {
                    editorRef.value.updateAiPreview({ text: aiPreviewResult.value, loading: true })
                }
            } else if (status === 'ready' || status === 'error') {
                aiBusy.value = false
                if (editorRef.value?.updateAiPreview) {
                    editorRef.value.updateAiPreview({ loading: false })
                }
            }
        }
    )

    // 处理AI操作
    async function handleAiAction(payload: { 
        type: 'polish' | 'expand' | 'summary'
        selection: string
        fullText: string 
    }) {
        if (aiBusy.value) return
        aiBusy.value = true

        try {
            editorChat.stop()
            aiSelectionRange.value = editorRef.value?.getSelectionRange?.() || null
            
            const inputText = payload.selection || payload.fullText
            const trimmed = inputText.trim()
            if (!trimmed) {
                aiBusy.value = false
                return
            }

            const clipped = trimmed.length > 6000 ? trimmed.slice(0, 6000) : trimmed
            
            const promptMap: Record<string, string> = {
                polish: '请润色以下文本，保持含义不变，仅输出润色后的文本。',
                expand: '请扩写以下文本，使其更完整流畅，仅输出扩写后的文本。',
                summary: '请为以下文本生成摘要，仅输出摘要内容。',
            }
            const prompt = `${promptMap[payload.type]}\n\n${clipped}`

            // 重置状态
            aiPreviewType.value = payload.type
            aiPreviewOriginal.value = payload.selection || payload.fullText
            aiPreviewResult.value = ''
            aiPreviewHadSelection.value = Boolean(payload.selection)
            aiStreamActive.value = true
            
            // 显示浮动框
            if (aiSelectionRange.value) {
                showAiPreview.value = true
                await nextTick()
                if (editorRef.value?.setAiPreview) {
                    editorRef.value.setAiPreview({
                        from: aiSelectionRange.value.from,
                        to: aiSelectionRange.value.to,
                        text: 'AI 生成中…',
                        loading: true,
                        type: payload.type,
                    })
                }
            } else {
                showAiPreview.value = true
            }
            
            // 清空并发送消息
            // @ts-ignore - ai-sdk Chat accepts direct assign in runtime
            editorChat.messages = []
            await nextTick()
            await editorChat.sendMessage({ text: prompt })
        } catch (error) {
            console.error('AI action error:', error)
            aiBusy.value = false
            aiStreamActive.value = false
        }
    }

    // 应用AI结果
    function applyAiResult() {
        if (!aiPreviewResult.value.trim()) return
        
        const range = aiSelectionRange.value
        if (range) {
            editorRef.value?.replaceRange?.(range.from, range.to, aiPreviewResult.value)
        } else if (aiPreviewHadSelection.value) {
            editorRef.value?.replaceSelection(aiPreviewResult.value)
        } else {
            editorRef.value?.insertAtCursor(aiPreviewResult.value)
        }
        
        showAiPreview.value = false
        editorRef.value?.clearAiPreview?.()
        aiStreamActive.value = false
    }

    // 关闭AI预览
    function closeAiPreview() {
        showAiPreview.value = false
        editorChat.stop()
        editorRef.value?.clearAiPreview?.()
        aiStreamActive.value = false
        aiBusy.value = false
    }

    // 停止AI生成（由"终止"按钮触发）
    function stopAiGeneration() {
        console.log('[EditorAI] Stopping AI generation')
        editorChat.stop()
        closeAiPreview()
    }

    return {
        // 状态
        editorRef,
        editorChat,
        aiBusy,
        showAiPreview,
        aiPreviewOriginal,
        aiPreviewResult,
        aiPreviewHadSelection,
        aiPreviewType,
        aiSelectionRange,

        // 方法
        handleAiAction,
        applyAiResult,
        closeAiPreview,
        stopAiGeneration,
    }
}
