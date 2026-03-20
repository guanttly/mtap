// 核心目的：数据脱敏工具
// 模块功能：患者姓名脱敏（张*明）、手机号脱敏
// 核心目的：数据脱敏工具
// 模块功能：患者姓名脱敏（张*明）、手机号脱敏
/** 姓名脱敏：张三 → 张*, 张三丰 → 张*丰 */
export function maskName(name) {
    if (!name)
        return '';
    if (name.length === 1)
        return name;
    if (name.length === 2)
        return `${name[0]}*`;
    return `${name[0]}${'*'.repeat(name.length - 2)}${name[name.length - 1]}`;
}
/** 手机号脱敏：13812345678 → 138****5678 */
export function maskPhone(phone) {
    if (!phone || phone.length < 8)
        return phone;
    return `${phone.slice(0, 3)}****${phone.slice(-4)}`;
}
/** 证件号脱敏：仅保留前3位和后4位 */
export function maskId(id) {
    if (!id || id.length < 8)
        return id;
    return `${id.slice(0, 3)}${'*'.repeat(id.length - 7)}${id.slice(-4)}`;
}
