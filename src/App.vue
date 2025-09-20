<script>
export default {
  watch: {
    $route(to) {
      document.title = to.meta.title || 'App'
    }
  }
}
</script>
<template>
  <component :is="$route.meta.layoutComponent" />
</template>