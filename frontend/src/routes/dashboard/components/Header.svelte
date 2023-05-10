<script>
  import {
    asset,
    assetNameShort,
    intradayDiff,
    priceChangeAbs,
    priceChangePercentAbs,
    tradePriceFormatted,
  } from '../dashboardStore'

  $: up = $intradayDiff?.sign === 1
  $: down = $intradayDiff?.sign === -1
  $: sign = $intradayDiff?.sign === 1 ? '+' : $intradayDiff?.sign === -1 ? '-' : ''
</script>

<header class='header'>
  <div class='symbol'>{$asset?.symbol}</div>
  <div class='price'>{$tradePriceFormatted}</div>
  <div class='price-change' class:up={up} class:down={down}>
    <div class='change'>{sign}{$priceChangeAbs}</div>
    <div class='change-percent'>({$priceChangePercentAbs})</div>
  </div>
  <div class='asset'>
    <div class='name'>{$assetNameShort}</div>
    <div class='exchange'>{$asset?.exchange}</div>
  </div>
</header>

<style lang='scss'>
  .header {
    @apply grid gap-1 justify-start;
    grid-template-columns: 8rem 10rem auto;
    grid-template-areas:
      'symbol price price-change'
      'asset asset asset';

    .symbol {
      grid-area: symbol;
      @apply font-bold tracking-widest text-info-content;
    }

    .price {
      grid-area: price;
    }

    .price-change {
      grid-area: price-change;
    }

    .asset {
      grid-area: asset;
      @apply flex flex-row gap-1 text-sm;

      .name {
        @apply text-secondary;
      }

      .exchange {
        @apply opacity-75;
      }
    }

    .symbol, .price {
      @apply text-4xl;
    }

    .price-change {
      @apply flex gap-1 items-end mb-0.5;
    }

    .price, .change, .change-percent {
      @apply tabular-nums;
    }

    .up {
      @apply text-success;
    }

    .down {
      @apply text-error;
    }
  }
</style>
