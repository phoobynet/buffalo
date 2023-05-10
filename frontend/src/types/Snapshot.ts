export interface Snapshot {
  latestTrade: LatestTrade
  latestQuote: LatestQuote
  minuteBar: Bar
  dailyBar: Bar
  prevDailyBar: Bar
}

export interface LatestTrade {
  t: string
  p: number
  s: number
  x: string
  i: number
  c: string[]
  z: string
  u: string
}

export interface LatestQuote {
  t: string
  bp: number
  bs: number
  bx: string
  ap: number
  as: number
  ax: string
  c: string[]
  z: string
}

export interface Bar {
  t: string
  o: number
  h: number
  l: number
  c: number
  v: number
  n: number
  vw: number
}
