<script lang='ts'>
  import { Route, Router } from 'svelte-routing'
  import Dashboard from '@/routes/dashboard/Dashboard.svelte'

  export let url = '/'
</script>

<Router {url}>
  <Route path='/'>
    <Dashboard />
  </Route>
</Router>
