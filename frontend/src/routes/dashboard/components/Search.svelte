<script lang="ts">
  export let searchValue = ''
</script>

<div class='search'>
  <input type='text' bind:value={searchValue} class='input'>
</div>

<style lang='scss'>
  .search {
    input {}
  }
</style>
