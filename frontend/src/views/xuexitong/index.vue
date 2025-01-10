<!--
 * 学习通页面
 * @author: sugarbecomer
 * @since: 2025-01-07
 * index.vue
-->
<template>
  <div class="container">
    <!-- 搜索框 -->
    <div class="head-box">
      <div class="search-box">
        <n-form
          ref="formRef"
          inline
          :label-width="80"
          :model="queryForm"
          label-placement="left"
        >
          <n-form-item label="账号" path="username">
            <n-input
              v-model:value="queryForm.username"
              placeholder="输入姓名"
            />
          </n-form-item>
          <n-form-item>
            <n-button
              plain
              type="info"
              attr-type="button"
              @click="handleQuery"
              dashed
            >
              搜索
            </n-button>
          </n-form-item>
        </n-form>
      </div>
      <div class="suffix-box">
        <n-button type="primary" attr-type="button" @click="handleAdd">
          <template #icon>
            <n-icon>
              <AddOutline />
            </n-icon>
          </template>
          添加账号
        </n-button>
      </div>
    </div>
    <div class="data-table">
      <n-data-table
        :columns="columns"
        :data="tableData"
        bordered
      />
      <n-pagination
      :page="queryForm.page"
      :page-size="queryForm.size" 
      :page-count="Math.ceil(total/queryForm.size)" 
      :show-size-picker="true"
      :page-size-options="[10,20,30,50]"
      @update:page="(page)=>{queryForm.page = page;handleQuery()}"
      @update:page-size="(size)=>{queryForm.size = size;handleQuery()}"
>
      </n-pagination>
      <Add v-model:show="addShow" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref, render } from "vue";
import { NTag } from "naive-ui";
import { AddOutline } from "@vicons/ionicons5";
import { Xuexitong } from "@wails/service";
import Add from './add.vue';
const queryForm = ref({
  username: "",
  page: 1,
  size: 10,
});
const handleQuery = () => {
  queryForm.value.page = 1;
  getList()
};
const total = ref(0)
const getList = ()=>{
  Xuexitong.List(queryForm.value).then(res=>{
    tableData.value = res.data
    total.value = res.total
  })
}
const columns = [
  {
    title: "id",
    key: "id",
  },
  {
    title: "账号",
    key: "username",
  },
  {
    title: "学科",
    key: "courseName",
  },
  {
    title: "状态",
    key: "status",
    render(row) {
      switch (row.status) {
        case 0:
          return h(NTag,{},{default:()=>'未开始'})
        case 1:
          return h(NTag,{type:'warning'},{default:()=>'队列中'})
        case 2:
          return h(NTag,{type:'info'},{default:()=>'进行中'})
        case 3:
          return h(NTag,{type:'success'},{default:()=>'已完成'})
        case 4:
          return h(NTag,{type:'error'},{default:()=>'已失败'})
        case 5:
          return h(NTag,{type:'error'},{default:()=>'已暂停'})
        default:
          return h(NTag,{type:'warning'},{default:()=>'未知'})
      }
    },
  },
  {
    title: "操作",
    key: "oper",
  },
];
const addShow = ref(false)
const handleAdd = () => {
  addShow.value = true
}
const tableData = ref([])
onMounted(()=>{
  handleQuery()
  // 定时刷新
  setInterval(()=>{
    getList()
  },1000*30)
})
</script>

<style scoped lang="scss">
.head-box {
  display: flex;
  justify-content: space-between;
}
</style>
