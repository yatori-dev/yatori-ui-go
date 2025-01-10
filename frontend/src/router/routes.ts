export const routes: any[] = [
  {
    path: "/",
    name: "Index",
    redirect: "/xuexitong",
    meta: {
      title: "学习通",
    },
    children: [
      {
        path: "/xuexitong",
        name: "xuexitong",
        component: () => import("@/views/xuexitong/index.vue"),
        meta: {
          title: "学习通",
        },
      },
    ],
  },
];
