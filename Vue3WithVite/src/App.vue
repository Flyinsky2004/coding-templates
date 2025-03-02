<script setup>
import {RouterView} from "vue-router";
import {useThemeStore} from "@/stores/theme.js";
import {theme} from "ant-design-vue";
const themeStore = useThemeStore()

themeStore.initTheme()
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
dayjs.locale('zh-cn');
</script>

<template>
<a-config-provider
      :theme="{
      algorithm: themeStore.isDark ?theme.darkAlgorithm:theme.defaultAlgorithm,
    }"
      :locale="zhCN"
  >
  </a-config-provider>
  <div class="container1 select-none h-screen bkg-theme-switch">
<!--       :class="themeStore.currentTheme === 'light' ? 'lightBKG' : 'darkBKG'">-->
    <div class="h-full backdrop-blur-sm overflow-y-auto font-sans">
      <RouterView/>
    </div>
  </div>
</template>

<style scoped>
.container1{
  font-family: 'srchan', serif;
  background-size: 100% 100%; /* 确保背景图片覆盖整个容器 */
  background-attachment: fixed; /* 背景图片固定不动 */
  background-position: center; /* 背景图片居中 */
  background-blend-mode: multiply;
}
</style>
