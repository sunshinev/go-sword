<template>
  <div>
    <LayoutRoot ref="ch"></LayoutRoot>
  </div>
</template>
<script type="text/javascript">
  import httpVueLoader from "http-vue-loader"

  export default {
    components: {
      LayoutRoot: 'url:'+window.location.origin+window.location.pathname + "render?path=/layout/default"
    },
    updated() {
      this.route(this.$refs.ch.routes)
    },
    methods: {
      // handle route
      route(routes) {
        var list = new Array()
        for (var i in routes) {
          list.push({
            name: routes[i].name ? routes[i].name : null,
            path: routes[i].path,
            component: routes[i].url ? httpVueLoader(routes[i].url) : null,
            redirect: routes[i].redirect ? routes[i].redirect : null,
            props: (route)=>({
              ...route.query
            })
          })
        }

        // inject route
        this.$router.addRoutes(list)
      }
    }
  }
</script>