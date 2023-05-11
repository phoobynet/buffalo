<script lang="ts">
  import numeral from 'numeral'

  export let title: string
  export let percentChange: number

  $: percent = numeral(percentChange).format('0.00%')
  $: isUp = percentChange > 0
  $: isDown = percentChange < 0
</script>

<div class="change-since">
  <div class="up" class:is-up={isUp} />
  <div class="content" class:is-up={isUp} class:is-down={isDown}>
    <h5 class="title">{change.since}</h5>
    <div class="change-container">
      <div class="change-pct">{percent}</div>
    </div>
  </div>
  <div class="down" class:is-down={isDown} />
</div>

<style lang="scss">
  .change-since {
    @apply items-center text-center min-w-[5rem] tabular-nums text-xs text-white;
  }

  .content {
    @apply w-[99%] flex flex-col py-1;

    &.is-up {
      @apply from-green-500 to-green-950 bg-gradient-to-b;
    }

    &.is-down {
      @apply from-red-500 to-red-950 bg-gradient-to-t;
    }

    .change-container {
      @apply min-h-[2rem] flex items-center justify-center;
    }

    .title {
      @apply w-full font-extrabold;
    }

    .change-pct {
      @apply text-sm font-light leading-tight px-2;
    }
  }

  .up,
  .down {
    @apply h-6 w-full;
  }

  .up {
    &.is-up {
      @apply bg-green-500;
    }

    clip-path: polygon(0 100%, 100% 100%, 50% 0);
  }

  .down {
    &.is-down {
      @apply bg-red-500;
    }

    clip-path: polygon(0 0, 50% 100%, 100% 0);
  }
</style>
