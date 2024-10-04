import { describe, it, expect, vi, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import LoginView from '@/views/LoginView.vue'
import { useCharacterStore } from '@/store'
import { VueWrapper } from '@vue/test-utils'
import * as vueRouter from 'vue-router'
import axios from 'axios'

vi.mock('axios')
vi.mock('vue-router')

describe('LoginView.vue', () => {
  let wrapper: VueWrapper<any>
  let store: ReturnType<typeof useCharacterStore>
  let mockRouter: { push: ReturnType<typeof vi.fn> }

  beforeEach(() => {
    const pinia = createPinia()
    setActivePinia(pinia)
    store = useCharacterStore()
    
    vi.spyOn(store, 'setIsLoggedIn')
    vi.spyOn(store, 'setUsername')
    vi.spyOn(store, 'loadCharacters').mockResolvedValue()

    mockRouter = { push: vi.fn() }
    vi.spyOn(vueRouter, 'useRouter').mockReturnValue(mockRouter as any)

    wrapper = mount(LoginView, {
      global: {
        plugins: [pinia],
        stubs: ['router-link']
      }
    })
  })

  it('renders login form correctly', () => {
    expect(wrapper.find('form').exists()).toBe(true)
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('handles successful login correctly', async () => {
    vi.useFakeTimers()
    const mockToken = 'fake-token'
    vi.mocked(axios.post).mockResolvedValueOnce({ data: { token: mockToken } })

    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('testpass')
    await wrapper.find('form').trigger('submit.prevent')

    expect(axios.post).toHaveBeenCalledWith(expect.any(String), {
      username: 'testuser',
      password: 'testpass',
    })

    // Wait for the next tick to allow for async operations to complete
    await vi.advanceTimersByTime(0)
    await flushPromises()

    expect(mockRouter.push).toHaveBeenCalledWith('/')
  })

  it('handles login failure correctly', async () => {
    const mockError = new Error('Login failed')
    vi.mocked(axios.post).mockRejectedValueOnce(mockError)

    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('wrongpass')
    await wrapper.find('form').trigger('submit.prevent')

    expect(axios.post).toHaveBeenCalledWith(expect.any(String), {
      username: 'testuser',
      password: 'wrongpass',
    })

    expect(consoleSpy).toHaveBeenCalledWith('Login failed:', mockError)

    consoleSpy.mockRestore()
  })
})