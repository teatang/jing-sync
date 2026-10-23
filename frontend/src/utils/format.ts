/**
 * 将字节大小格式化为人类易读的字符串 (B, KB, MB, GB, TB, PB)
 * 与后端 utils.FormatBytes 保持一致
 * @param bytes 字节数
 * @returns 格式化后的字符串
 */
export function formatBytes(bytes: number): string {
  if (bytes < 1024) {
    return `${bytes} B`;
  }
  const units = ["B", "KB", "MB", "GB", "TB", "PB"];
  let val = bytes;
  let i = 0;
  while (val >= 1024 && i < units.length - 1) {
    val /= 1024;
    i++;
  }
  return `${val.toFixed(2)} ${units[i]}`;
}
