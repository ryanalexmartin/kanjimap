import { JSDOM } from 'jsdom'
import fs from 'fs'
import path from 'path'

const html = fs.readFileSync(path.resolve(__dirname, '../popup.html'), 'utf8')

let dom
let document
let window

describe('Popup', () => {
  beforeEach(() => {
    dom = new JSDOM(html, { runScripts: 'dangerously' })
    document = dom.window.document
    window = dom.window

    global.browser = {
      storage: {
        local: {
          get: jest.fn(),
          set: jest.fn(),
        },
      },
      runtime: {
        sendMessage: jest.fn(),
      },
    }

    global.fetch = jest.fn()
  })

  test('login success', async () => {
    document.getElementById('username').value = 'testuser'
    document.getElementById('password').value = 'testpassword'

    global.fetch.mockResolvedValueOnce({
      json: () => Promise.resolve({ token: 'fake-token' }),
    })

    document.getElementById('loginBtn').click()

    await new Promise(resolve => setTimeout(resolve, 0))

    expect(global.browser.storage.local.set).toHaveBeenCalledWith({
      authToken: 'fake-token',
      username: 'testuser',
      password: 'testpassword',
    })
    expect(document.getElementById('status').textContent).toBe('Logged in successfully!')
  })

  test('login failure', async () => {
    document.getElementById('username').value = 'testuser'
    document.getElementById('password').value = 'wrongpassword'

    global.fetch.mockResolvedValueOnce({
      json: () => Promise.resolve({}),
    })

    document.getElementById('loginBtn').click()

    await new Promise(resolve => setTimeout(resolve, 0))

    expect(document.getElementById('status').textContent).toBe('Login failed. Please try again.')
  })
})