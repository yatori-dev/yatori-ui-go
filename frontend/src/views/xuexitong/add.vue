<template>
  <n-modal
    v-model:show="modalShow"
    style="width: 400px"
    title="添加账号"
    :mask-closable="false"
    :bordered="false"
    preset="dialog"
    positive-text="保存"
    negative-text="取消"
    :on-positive-click="handleOk"
    draggable
  >
    <n-form ref="formRef" :model="form" :rules="rules">
      <n-form-item path="username" label="账号">
        <n-input
          v-model:value="form.username"
          @keydown.enter.prevent
          placeholder="请输入账号"
        />
      </n-form-item>
      <n-form-item path="password" label="密码">
        <n-input-group>
          <n-input
            v-model:value="form.password"
            type="password"
            @keydown.enter.prevent
            placeholder="请输入密码"
          >
          </n-input>
          <n-button type="primary" ghost @click="handleQueryCourse">
            查询课程
          </n-button>
        </n-input-group>
      </n-form-item>
      <n-form-item path="courseId" label="课程">
        <n-select
          value-field="courseId"
          label-field="courseName"
          v-model:value="form.courseId"
          :options="courseList"
          :disabled="!courseList.length"
          :render-label="renderCourse"
        />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, h, ref, VNodeChild } from "vue";
import { FormInst, SelectOption, useMessage } from "naive-ui";
import { Xuexitong } from "@wails/service";
defineOptions({
  name: "XueXiTongAdd",
})
const props = defineProps({
  show: {
    type: Boolean,
    required: true,
  },
});
const message = useMessage();
const modalShow = computed({
  get() {
    return props.show;
  },
  set(val: boolean) {
    emit("update:show", val);
  },
});
const form = ref({
  username: "",
  password: "",
  courseId: "",
});
const formRef = ref<FormInst | null>(null);
const courseList = ref<any[]>([]);
const renderCourse = (option: SelectOption): VNodeChild => {
  return [
    h(
      "span",
      { class: "option-text" },
      {
        default: () => {
          if (option.courseName) {
            return `${option.courseName} - ${option.courseTeacher} - ${option.courseImage}`;
          } else {
            return "";
          }
        },
      }
    ),
  ];
};
const rules = {
  username: [
    { required: true, message: "请输入账号" },
  ],
  password: [
    { required: true, message: "请输入密码" },
  ],
  courseId: [{ required: true, message: "请选择课程" }],
};
const handleOk = async () => {
  formRef.value?.validate((err) => {
    if (err) {
      message.error("参数错误");
      return;
    }
    const data = {
      courseName: courseList.value.find(
        (item: any) => item.courseId === form.value.courseId
      )!.courseName,
      ...form.value,
    }
    Xuexitong.Add(data).then(() => {
      message.success("添加成功");
      modalShow.value = false;
    });
  });
  return Promise.reject();
};
const emit = defineEmits(["update:show"]);
const handleQueryCourse = () => {
  if (!form.value.password || !form.value.username) {
    message.warning("请输入账号和密码");
    return;
  }
  Xuexitong.QueryCourse(form.value).then((res) => {
    courseList.value = res.data;
  });
};
</script>
