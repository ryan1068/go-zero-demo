/** 密码规则：8-72 位，且至少包含一个字母和一个数字。 */
export const passwordPattern = /^(?=.*[A-Za-z])(?=.*\d).{8,72}$/

export const passwordRule = {
  pattern: passwordPattern,
  message: '密码为 8-72 位，且至少包含字母和数字',
  trigger: 'blur' as const,
}
