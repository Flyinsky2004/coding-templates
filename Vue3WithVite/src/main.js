import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';
import axios from "axios";
//开发环境
axios.defaults.baseURL="http://localhost:8080"
//生产环境
axios.defaults.baseURL= window.location.protocol + '//' + window.location.host
const app = createApp(App)

app.use(createPinia()).use(router).use(ElementPlus).use(Antd)

app.mount('#app')
