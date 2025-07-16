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

export const JOB_SPEED_CONFIGS: JobSelectOptionConfig[] = [
  {
    id: 0,
    name: "标准",
    remark: "默认选项",
  },
  {
    id: 1,
    name: "快速",
    remark: "用Alist缓存扫描目标目录",
  },
  {
    id: 2,
    name: "低速",
    remark: "防止频繁被网盘限制",
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
    name: "间隔",
    remark: "每n分钟同步一次",
  },
  {
    id: 1,
    name: "cron",
    remark: "同linux cron",
  },
  {
    id: 2,
    name: "仅手动",
    remark: "不自动调用，手动触发",
  },
];
