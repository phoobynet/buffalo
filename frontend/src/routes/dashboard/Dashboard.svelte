<script lang='ts'>
  import { IsReady, Subscribe } from '../../../wailsjs/go/main/App'
  import { onMount } from 'svelte'
  import { EventsOn } from '../../../wailsjs/runtime'
  import { asset, quote, snapshot, symbol, trade } from './dashboardStore'
  import Header from './components/Header.svelte'
  import { alpaca, marketdata } from '../../../wailsjs/go/models'
  import Search from '@/routes/dashboard/components/Search.svelte'
  import type { StreamQuote, StreamTrade } from '@/lib/types'

  let isReady = false

  function sleep(t: number = 1_000): Promise<void> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve()
      }, t)
    })
  }

  let isSubscribed = false
  $symbol = 'AAPL'

  EventsOn('trade', (data) => {
    $trade = data satisfies StreamTrade
  })

  EventsOn('quote', (data) => {
    $quote = data satisfies StreamQuote
  })

  EventsOn('snapshot', (data) => {
    $snapshot = data satisfies marketdata.Snapshot
  })

  EventsOn('asset', (data) => {
    $asset = data satisfies alpaca.Asset
  })

  onMount(async () => {
    isReady = await IsReady()

    if (!isReady) {
      for (let i = 0; i < 10; i++) {
        await sleep()
        isReady = await IsReady()
        if (isReady) {
          await Subscribe($symbol)
          break
        }
      }
    }
  })
</script>

<div>
  <Search />
  <div class='mx-2 mt-2'>
    <Header />

    <pre>{JSON.stringify($snapshot, null, 2)}</pre>
  </div>
</div>
