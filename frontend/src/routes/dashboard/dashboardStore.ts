import { derived, writable } from 'svelte/store'
import numeral from 'numeral'
import { numberDiff } from '@/lib/numberDiff'
import type { alpaca } from '../../../wailsjs/go/models'
import type { marketdata } from '../../../wailsjs/go/models'
import type { StreamQuote, StreamTrade } from '@/lib/types'
import { formatISO } from 'date-fns'

export const symbol = writable<string>('')

export const trade = writable<StreamTrade>()

export const tradePriceFormatted = derived(trade, $trade => {
  if (!$trade) return undefined

  return numeral($trade.p).format('$0,0.00')
})

export const quote = writable<StreamQuote>()

export const snapshot = writable<marketdata.Snapshot>()

export const latestTrade = derived([trade, snapshot, symbol], ([$trade, $snapshot, $symbol]) => {
  if (!$trade || !$snapshot || !$symbol) return undefined

  if ($trade) {
    return $trade
  }

  return {
    S: $symbol,
    ...$snapshot.latestTrade,
  }
})

export const asset = writable<alpaca.Asset>()

const prevDailyBar = derived(snapshot, $snapshot => {
  if (!$snapshot) return undefined

  const date = formatISO(Date.now(), {
    representation: 'date',
  })

  if ($snapshot.dailyBar.t.substring(0, 10) !== date) {
    return $snapshot.dailyBar
  }

  return $snapshot.prevDailyBar
})

export const assetNameShort = derived(asset, $asset => {
  if (!$asset) return undefined

  return $asset?.name?.replace('Common Stock', '').trim()
})

export const intradayDiff = derived([prevDailyBar, trade], ([$prevDailyBar, $trade]) => {
  if (!$prevDailyBar || !$trade) return undefined

  return numberDiff($prevDailyBar.c, $trade.p)
})

export const priceChangeAbs = derived(intradayDiff, $intradayDiff => {
  if (!$intradayDiff) return undefined

  return numeral(Math.abs($intradayDiff.change)).format('0.00')
})

export const priceChangePercentAbs = derived(intradayDiff, $intradayDiff => {
  if (!$intradayDiff) return undefined

  return numeral(Math.abs($intradayDiff.changePercent)).format('0.00%')
})

export const previousCloseFormatted = derived(prevDailyBar, $prevDailyBar => {
  if (!$prevDailyBar) return undefined

  return numeral($prevDailyBar.c).format('$0,0.00')
})
