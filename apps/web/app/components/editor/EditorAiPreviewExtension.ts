import { Extension } from '@tiptap/core'
import { Plugin, PluginKey } from 'prosemirror-state'
import { Decoration, DecorationSet } from 'prosemirror-view'

export type AiPreviewPayload = {
    from: number
    to: number
    text: string
    loading?: boolean
    type?: 'polish' | 'expand' | 'summary'
}

type AiPreviewState = {
    from: number | null
    to: number | null
    text: string
    loading: boolean
    type?: 'polish' | 'expand' | 'summary'
    visible: boolean
}

type AiPreviewMeta =
    | { type: 'set'; payload: AiPreviewPayload }
    | { type: 'update'; payload: Partial<AiPreviewPayload> }
    | { type: 'clear' }

export type AiPreviewOptions = {
    onApply?: () => void
    onDiscard?: () => void
    onStop?: () => void
}

const AiPreviewKey = new PluginKey<AiPreviewState>('ai-preview')

const defaultState: AiPreviewState = {
    from: null,
    to: null,
    text: '',
    loading: false,
    type: undefined,
    visible: false,
}

const typeLabelMap: Record<string, string> = {
    polish: '润色',
    expand: '扩写',
    summary: '摘要',
}

function resolveLabel(type?: string) {
    if (!type) return 'AI 生成'
    return typeLabelMap[type] ? `AI ${typeLabelMap[type]}` : 'AI 生成'
}

function createPreviewWidget(state: AiPreviewState, options: AiPreviewOptions) {
    const container = document.createElement('div')
    container.className = 'ai-preview-widget'
    container.contentEditable = 'false'

    const header = document.createElement('div')
    header.className = 'ai-preview-header'

    const title = document.createElement('div')
    title.className = 'ai-preview-title'
    
    const titleText = document.createElement('span')
    titleText.textContent = resolveLabel(state.type)
    title.appendChild(titleText)

    // 添加 loading 指示器（CSS 动画）
    if (state.loading) {
        const spinner = document.createElement('span')
        spinner.className = 'ai-preview-spinner'
        title.appendChild(spinner)
    }

    const actions = document.createElement('div')
    actions.className = 'ai-preview-actions'

    const discardButton = document.createElement('button')
    discardButton.type = 'button'
    discardButton.className = 'ai-preview-btn ai-preview-btn-ghost'
    discardButton.textContent = '丢弃'
    discardButton.addEventListener('click', (event) => {
        event.preventDefault()
        event.stopPropagation()
        options.onDiscard?.()
    })

    const actionButton = document.createElement('button')
    actionButton.type = 'button'
    actionButton.className = 'ai-preview-btn ai-preview-btn-primary'
    
    // 根据加载状态改变按钮行为
    if (state.loading) {
        actionButton.textContent = '终止'
        actionButton.classList.add('ai-preview-btn-stop')
        actionButton.addEventListener('click', (event) => {
            event.preventDefault()
            event.stopPropagation()
            options.onStop?.()
        })
    } else {
        actionButton.textContent = '使用'
        actionButton.addEventListener('click', (event) => {
            event.preventDefault()
            event.stopPropagation()
            options.onApply?.()
        })
    }

    actions.appendChild(discardButton)
    actions.appendChild(actionButton)

    header.appendChild(title)
    header.appendChild(actions)

    const body = document.createElement('div')
    body.className = 'ai-preview-body'
    
    // 显示内容或加载提示
    if (state.text) {
        body.textContent = state.text
        // 如果正在加载中，添加流式更新的视觉效果
        if (state.loading) {
            body.classList.add('streaming')
        }
    } else if (state.loading) {
        body.textContent = 'AI 生成中…'
        body.classList.add('ai-preview-loading')
    } else {
        body.textContent = ''
    }

    container.appendChild(header)
    container.appendChild(body)

    // 保存引用以便更新
    container.setAttribute('data-ai-preview', 'true')

    return container
}

export const AiPreviewExtension = Extension.create<AiPreviewOptions>({
    name: 'aiPreview',

    addOptions() {
        return {
            onApply: undefined,
            onDiscard: undefined,
            onStop: undefined,
        }
    },

    addCommands() {
        return {
            setAiPreview:
                (payload: AiPreviewPayload) =>
                ({ tr, dispatch }) => {
                    dispatch?.(tr.setMeta(AiPreviewKey, { type: 'set', payload } satisfies AiPreviewMeta))
                    return true
                },
            updateAiPreview:
                (payload: Partial<AiPreviewPayload>) =>
                ({ tr, dispatch }) => {
                    dispatch?.(tr.setMeta(AiPreviewKey, { type: 'update', payload } satisfies AiPreviewMeta))
                    return true
                },
            clearAiPreview:
                () =>
                ({ tr, dispatch }) => {
                    dispatch?.(tr.setMeta(AiPreviewKey, { type: 'clear' } satisfies AiPreviewMeta))
                    return true
                },
        }
    },

    addProseMirrorPlugins() {
        return [
            new Plugin<AiPreviewState>({
                key: AiPreviewKey,
                state: {
                    init: () => ({ ...defaultState }),
                    apply: (tr, value) => {
                        const meta = tr.getMeta(AiPreviewKey) as AiPreviewMeta | undefined
                        if (meta?.type === 'clear') {
                            return { ...defaultState }
                        }
                        if (meta?.type === 'set') {
                            return {
                                from: meta.payload.from,
                                to: meta.payload.to,
                                text: meta.payload.text,
                                loading: Boolean(meta.payload.loading),
                                type: meta.payload.type,
                                visible: true,
                            }
                        }
                        if (meta?.type === 'update') {
                            const updated = {
                                ...value,
                                loading: meta.payload.loading ?? value.loading,
                            }
                            // 更新文本内容
                            if (meta.payload.text !== undefined) {
                                updated.text = meta.payload.text
                            }
                            // 更新位置
                            if (meta.payload.from !== undefined) {
                                updated.from = meta.payload.from
                            }
                            if (meta.payload.to !== undefined) {
                                updated.to = meta.payload.to
                            }
                            // 更新类型
                            if (meta.payload.type !== undefined) {
                                updated.type = meta.payload.type
                            }
                            return updated
                        }
                        if (value.visible && value.from !== null && value.to !== null && tr.docChanged) {
                            const mappedFrom = tr.mapping.map(value.from)
                            const mappedTo = tr.mapping.map(value.to)
                            return {
                                ...value,
                                from: mappedFrom,
                                to: mappedTo,
                            }
                        }
                        return value
                    },
                },
                props: {
                    decorations: (state) => {
                        const data = AiPreviewKey.getState(state)
                        if (!data || !data.visible || data.to === null) {
                            return null
                        }
                        const pos = Math.min(Math.max(data.to, 0), state.doc.content.size)
                        
                        // 每次状态变化都重新创建 widget 以确保内容更新
                        // 使用状态的文本和loading作为key的一部分来强制更新
                        const widgetKey = `ai-preview-${pos}-${data.text.length}-${data.loading}`
                        
                        const widget = Decoration.widget(
                            pos,
                            () => createPreviewWidget(data, this.options),
                            { side: 1, key: widgetKey }
                        )
                        return DecorationSet.create(state.doc, [widget])
                    },
                },
            }),
        ]
    },
})

declare module '@tiptap/core' {
    interface Commands<ReturnType> {
        aiPreview: {
            setAiPreview: (payload: AiPreviewPayload) => ReturnType
            updateAiPreview: (payload: Partial<AiPreviewPayload>) => ReturnType
            clearAiPreview: () => ReturnType
        }
    }
}
