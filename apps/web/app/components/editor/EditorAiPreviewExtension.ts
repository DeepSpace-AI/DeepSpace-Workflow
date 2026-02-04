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
    title.textContent = resolveLabel(state.type)

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

    const applyButton = document.createElement('button')
    applyButton.type = 'button'
    applyButton.className = 'ai-preview-btn ai-preview-btn-primary'
    applyButton.textContent = '使用'
    applyButton.addEventListener('click', (event) => {
        event.preventDefault()
        event.stopPropagation()
        options.onApply?.()
    })

    actions.appendChild(discardButton)
    actions.appendChild(applyButton)

    header.appendChild(title)
    header.appendChild(actions)

    const body = document.createElement('div')
    body.className = 'ai-preview-body'
    body.textContent = state.text || (state.loading ? 'AI 生成中…' : '')

    container.appendChild(header)
    container.appendChild(body)

    return container
}

export const AiPreviewExtension = Extension.create<AiPreviewOptions>({
    name: 'aiPreview',

    addOptions() {
        return {
            onApply: undefined,
            onDiscard: undefined,
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
                            return {
                                ...value,
                                ...meta.payload,
                                loading: meta.payload.loading ?? value.loading,
                            }
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
                        const widget = Decoration.widget(
                            pos,
                            () => createPreviewWidget(data, this.options),
                            { side: 1, key: 'ai-preview-widget' }
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
