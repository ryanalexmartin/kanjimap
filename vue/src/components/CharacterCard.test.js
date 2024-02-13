// CharacterCard.test.js
import { mount } from '@vue/test-utils'
import { expect, test } from 'vitest'
import CharacterCard from './CharacterCard.vue'

test('emits update-learned event when markLearned is called', async () => {
    const wrapper = mount(CharacterCard, {
        props: {
            character: {
                id: 1,
                character: '我',
                pinyin: 'wǒ',
                definition: 'I, me', // Add a comma here
            }
        }
    });
    await wrapper.vm.markLearned();
    expect(wrapper.emitted('update-learned')).toBeTruthy();
})
