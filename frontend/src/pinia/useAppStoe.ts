import { defineStore } from 'pinia'
export const useAppStore = defineStore('app', {
    state: (): { drawer: boolean } => ({
        drawer: true,
    }),
    actions: {
        getDrawer():    boolean {
            return this.drawer
        }
    }
})