import { readable } from 'svelte/store'
import { EventsOn } from '../../wailsjs/runtime'

export const trade = readable(undefined, set => {
  EventsOn('trade', data => {
    console.log(data)
    set(data)
  })
})
