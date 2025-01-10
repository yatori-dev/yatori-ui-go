<script lang="ts" setup>
import { useRouter } from "vue-router";
import { NConfigProvider } from "naive-ui";
import { zhCN, dateZhCN } from "naive-ui";
const tabs = [
  {
    name: "学习通",
    url: "/xuexitong",
  },
  {
    name: "智慧树",
    url: "",
  },
  { name: "学习公社", url: "" },
  {
    name: "仓辉实训",
    url: "",
  },
  {
    name: "创能实训",
    url: "",
  },
  {
    name: "重庆工程学院CQIE",
    url: "",
  },
];
const router = useRouter();
const handleChange = (url: string) => {
  router.push(url);
};
</script>
<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN">
    <n-message-provider>
      <n-modal-provider>
        <div class="flex-container">
          <n-tabs type="card" placement="left" size="large">
            <template v-for="(tab, index) in tabs">
              <n-tab :name="tab.name" :tab="tab.name"></n-tab>
            </template>
          </n-tabs>
          <div class="view">
            <router-view v-slot="{ Component }">
              <component :is="Component" />
            </router-view>
          </div>
        </div>
      </n-modal-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style lang="scss" scoped>
.flex-container {
  display: flex;
  flex-direction: row;
  .n-tabs {
    width: 180px;
  }
  .view {
    flex: 5;
    box-sizing: border-box;
    padding: 7px;
  }
}
</style>
