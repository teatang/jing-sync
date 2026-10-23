export interface IConfig {
  loading: boolean;
}

export const langs = ['zh-CN', 'en'] as const
export type Language = typeof langs[number]

export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

export interface RequestOptions {
  method?: HttpMethod;
  headers?: Record<string, string>;
  body?: unknown;
  params?: Record<string, string | number | boolean>;
}

export interface Engine {
  id?: number;
  url: string;
  token: string;
  remark: string;
  create_time?: string;
  update_time?: string;
}

export interface User {
  id?: number;
  username: string;
  password: string;
  create_time?: string;
  update_time?: string;
}

export interface Job {
  id?: number;
  remark?: string;
  src_path: string;
  dst_path: string;
  engine_id: number;
  user_id?: number;
  speed: number;
  sync_mode: number; // 0: 全量, 1: 仅新增, 2: 移动
  method?: string;
  interval?: number;
  is_cron?: number;
  cron_expr?: string;
  min_size?: number | null;
  max_size?: number | null;
  excludes?: string;
  enabled: number; // 1: 启用, 0: 禁用
  status?: number | boolean;
  create_time?: string;
  update_time?: string;
}

export interface SyncLog {
  id?: number;
  job_id: number;
  status: number; // 1: 成功, 0: 失败
  error_msg?: string;
  files_added: number;
  start_time: string;
  end_time: string;
  details?: string;
  create_time?: string;
}

export interface LoginInfo {
  token: string;
  expire_at: string;
}

export interface InfoList<T> {
  list: T[];
  pagination: Pagination;
}

export interface JsonSuccess<T> {
  code: number;
  success: boolean;
  data: T;
}

export interface JsonError {
  code: number;
  success: boolean;
  msg: string;
}

export interface JsonReturn<T> {
  code: number;
  success: boolean;
  msg?: string;
  data?: T | undefined;
}

export interface Pagination {
  page: number;
  size: number;
  total: number;
}

export interface JobSelectOptionConfig {
  id: number;
  name: string;
  remark?: string;
}

export interface CascaderNode {
  value: string;
  label: string;
  children?: CascaderNode[];
}
