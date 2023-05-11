<script lang='ts'>
  import { GetAsset, GetSnapshot, IsReady, Subscribe } from '../../../wailsjs/go/main/App'
  import { onMount } from 'svelte'
  import { EventsOn } from '../../../wailsjs/runtime'
  import type { Trade } from '@/types/Trade'
  import { asset, quote, snapshot, trade } from './dashboardStore'
  import type { Quote } from '@/types/Quote'
  import Header from './components/Header.svelte'

  let isSubscribed = false
  const symbol = 'AAPL'

  EventsOn('ready', async () => {
    isSubscribed = await Subscribe(symbol)
  })

  EventsOn('trade', (data) => {
    $trade = data satisfies Trade
  })

  EventsOn('quote', (data) => {
    $quote = data satisfies Quote
  })

  EventsOn('snapshot', (data) => {
    $snapshot = data
  })

  EventsOn('asset', (data) => {
    $asset = data
  })

  onMount(async () => {
    const isReady = await IsReady()

    if (isReady && !isSubscribed) {
      isSubscribed = await Subscribe(symbol)
    }

    if ($asset === undefined) {
      $asset = await GetAsset(symbol)
    }

    if ($snapshot === undefined) {
      $snapshot = await GetSnapshot(symbol)
    }
  })
</script>

<div class='mx-2 mt-2'>
  <Header />
</div>
