const hasLetter = /[A-Za-z]/
const hasDigit = /\d/

/**
 * 密码规则：至少 8 个字符、最多 72 字节（bcrypt 上限），
 * 且至少包含一个字母和一个数字；与 user-rpc 的校验保持一致。
 */
export const passwordRule = {
  validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
    const bytes = new TextEncoder().encode(value).length
    if ([...value].length < 8 || bytes > 72 || !hasLetter.test(value) || !hasDigit.test(value)) {
      callback(new Error('密码至少 8 个字符、最多 72 字节，且包含字母和数字'))
      return
    }
    callback()
  },
  trigger: 'blur' as const,
}
