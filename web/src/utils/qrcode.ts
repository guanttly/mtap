// 核心目的：二维码工具
// 模块功能：前端二维码生成与展示
// 核心目的：二维码工具
// 模块功能：前端二维码生成与展示

/**
 * 将预约凭证数据编码为二维码内容字符串
 * 真正生成二维码图片需配合 qrcodejs2-fixes 等库
 */
export function encodeCredential(appointmentId: string, patientId: string, token: string): string {
  return JSON.stringify({ aid: appointmentId, pid: patientId, tok: token, ts: Date.now() })
}

export function decodeCredential(raw: string): { aid: string, pid: string, tok: string, ts: number } | null {
  try {
    return JSON.parse(raw)
  }
  catch { return null }
}
