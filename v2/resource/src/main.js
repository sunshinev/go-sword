
import Vue from 'vue'
import ViewUI from 'view-design';
import 'view-design/dist/styles/iview.css';
Vue.use(ViewUI);

import VueRouter from 'vue-router'
Vue.use(VueRouter)

import httpVueLoader from "http-vue-loader"
Vue.use(httpVueLoader)

import VueAxios from 'vue-axios';
import axios from 'axios';
Vue.use(VueAxios, axios)

import CodeDiff from 'vue-code-diff'
Vue.use(CodeDiff)

// import VueHighlightJS from 'vue-highlightjs'

// Tell Vue.js to use vue-highlightjs
// Vue.use(VueHighlightJS)


import App from './App.vue'

Vue.config.productionTip = false

const router = new VueRouter({
	routes:[]
})

new Vue({
	router:router,
	render: h => h(App)
}).$mount('#app')

axios.interceptors.request.use(function(config) {
     Vue.prototype.$Spin.show();
    return config
})

axios.interceptors.response.use(response=>{
    Vue.prototype.$Spin.hide();
    return response
})
