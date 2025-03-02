import {defineStore} from 'pinia'

const useThemeStore = defineStore('theme', {
    state: () => ({
        isDark: false,
    }),

    actions: {
        initTheme() {
            // 检查本地存储中是否有保存的主题
            const savedTheme = localStorage.getItem('theme')
            if (savedTheme) {
                this.isDark = savedTheme === 'dark'
            } else {
                // 如果没有保存的主题,检查系统偏好
                this.isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
            }
            this.applyTheme()
        },

        toggleTheme() {
            this.isDark = !this.isDark
            this.applyTheme()
        },

        applyTheme() {
            const root = document.documentElement
            if (this.isDark) {
                root.classList.add('dark')
            } else {
                root.classList.remove('dark')
            }
            // 添加过渡类
            root.classList.add('theme-transition')
            // 保存主题到本地存储
            localStorage.setItem('theme', this.isDark ? 'dark' : 'light')

            // 可选：在过渡完成后移除过渡类
            setTimeout(() => {
                root.classList.remove('theme-transition')
            }, 300) // 与 CSS 中的 duration 保持一致
        }
    },

    getters: {
        currentTheme: (state) => state.isDark ? 'dark' : 'light'
    }
})

export {useThemeStore}