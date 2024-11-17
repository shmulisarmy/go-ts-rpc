export function add(num1: number, num2: number): number {
    return rpc_call("add", num1, num2);
}export function printNum(num: number): void {
    return rpc_call("printNum", num);
}export function load_file(filename: string): void {
    return rpc_call("load_file", filename);
}