import type { Trade } from '../../types/Trade'
import { derived, writable } from 'svelte/store'
import numeral from 'numeral'
import type { Quote } from '../../types/Quote'
import type { Snapshot } from '../../types/Snapshot'
import { numberDiff } from '../../lib/numberDiff'
import type { alpaca } from '../../../wailsjs/go/models'

export const trade = writable<Trade>(undefined)

export const tradePriceFormatted = derived(trade, $trade => {
  if (!$trade) return undefined
  return numeral($trade.p).format('$0,0.00')
})

export const quote = writable<Quote>(undefined)

export const snapshot = writable<Snapshot>(undefined)

export const asset = writable<alpaca.Asset>(undefined)

const prevDailyBar = derived(snapshot, $snapshot => {
  if (!$snapshot) return undefined

  if ($snapshot.prevDailyBar?.t.substring(0, 10) === $snapshot.dailyBar.t.substring(0, 10)) {
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
