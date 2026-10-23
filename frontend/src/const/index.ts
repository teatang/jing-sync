import type { JobSelectOptionConfig } from "@/types";

export const JOB_IS_ENABLE_CONFIG: JobSelectOptionConfig[] = [
  {
    id: 0,
    name: "禁用",
  },
  {
    id: 1,
    name: "启用",
  },
];

export const SYNC_MODE_CONFIGS: JobSelectOptionConfig[] = [
  {
    id: 1,
    name: "仅新增 (增量)",
    remark: "仅新增目标端缺少或修改的文件，安全无覆盖",
  },
  {
    id: 0,
    name: "完全同步 (镜像)",
    remark: "目标端与源端完全对齐",
  },
  {
    id: 2,
    name: "移动文件 (剪切)",
    remark: "同步到目标端后从源端删除",
  },
];

export const JOB_SPEED_CONFIGS: JobSelectOptionConfig[] = [
  {
    id: 0,
    name: "标准速率",
    remark: "默认平衡选项",
  },
  {
    id: 1,
    name: "快速扫描",
    remark: "使用OpenList/AList缓存扫描",
  },
  {
    id: 2,
    name: "低速限频",
    remark: "防止请求过快被网盘限制",
  },
];

export const JOB_METHOD_CONFIGS: JobSelectOptionConfig[] = [
  {
    id: 0,
    name: "仅新增",
    remark: "仅新增目标目录没有的文件",
  },
  {
    id: 1,
    name: "全同步",
    remark: "目标目录比源目录多的文件将被删除",
  },
];

export const JOB_IS_CRON_CONFIGS: JobSelectOptionConfig[] = [
  {
    id: 0,
    name: "定时间隔",
    remark: "每隔指定分钟自动同步一次",
  },
  {
    id: 1,
    name: "Cron表达式",
    remark: "标准Linux 5段式Cron格式 (如: 0 2 * * *)",
  },
  {
    id: 2,
    name: "仅手动触发",
    remark: "不自动调用，仅手动按需触发",
  },
];
