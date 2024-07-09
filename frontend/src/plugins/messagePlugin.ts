// src/plugins/messagePlugin.ts
import { App, reactive } from 'vue'

interface Message {
    content: string
    type: 'info' | 'success' | 'warning' | 'error'
}

export interface MessageState {
    queue: Message[]
    currentMessage: Message | null
    show: boolean
    timeout: number
}

const messageState = reactive<MessageState>({
    queue: [],
    currentMessage: null,
    show: false,
    timeout: 3000
})

const showMessage = (content: string, type: Message['type'] = 'info') => {
    messageState.queue.push({ content, type })
    console.log("show:",messageState.show)
    if (!messageState.show) {
        displayNextMessage()
    }
}

const displayNextMessage = () => {
    if (messageState.queue.length > 0) {
        messageState.currentMessage = messageState.queue.shift() || null
        messageState.show = true
        console.log("messageState:",messageState)
    }
}

const messagePlugin = {
    install(app: App) {
        app.config.globalProperties.$message = {
            info(content: string) {
                showMessage(content, 'info')
            },
            success(content: string) {
                showMessage(content, 'success')
            },
            warning(content: string) {
                showMessage(content, 'warning')
            },
            error(content: string) {
                console.log("content:",content)
                showMessage(content, 'error')
            }
        }

        console.log("provide.messageState:",messageState)

        app.provide('messageState', messageState)
        app.provide('displayNextMessage', displayNextMessage)
    }
}

export default messagePlugin

declare module '@vue/runtime-core' {
    interface ComponentCustomProperties {
        $message: {
            info(content: string): void
            success(content: string): void
            warning(content: string): void
            error(content: string): void
        }
    }
}
