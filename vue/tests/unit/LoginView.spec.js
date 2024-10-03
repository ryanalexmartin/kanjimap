import { shallowMount } from '@vue/test-utils'
import LoginView from '@/components/LoginView.vue'

describe('LoginView.vue', () => {
  it('emits login event when login is successful', async () => {
    const wrapper = shallowMount(LoginView)
    const username = 'testuser'
    const password = 'testpassword'

    wrapper.setData({ username, password })

    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ token: 'fake-token' }),
      })
    )

    await wrapper.vm.login()

    expect(wrapper.emitted().login).toBeTruthy()
    expect(wrapper.emitted().login[0]).toEqual([username])
    expect(localStorage.getItem('token')).toBe('fake-token')
    expect(localStorage.getItem('username')).toBe(username)
  })

  it('handles login failure', async () => {
    const wrapper = shallowMount(LoginView)
    const consoleSpy = jest.spyOn(console, 'error')

    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: false,
        status: 401,
      })
    )

    await wrapper.vm.login()

    expect(wrapper.emitted().login).toBeFalsy()
    expect(consoleSpy).toHaveBeenCalledWith(expect.stringContaining('There was a problem with the login request'))
  })
})