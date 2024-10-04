import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import path from 'path'

export default defineConfig({
    plugins: [vue()],
    test: {
        environment: 'jsdom',
        globals: true,
        test: {
            deps: {
                optimizer: {
                    web: {
                        include: ['@vue', '@vueuse'],
                    },
                },
                resolve: {
                    alias: {
                        '@': path.resolve(__dirname, './src'),
                    },
                },
            },
        },
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
})
