"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.add = add;
exports.printNum = printNum;
exports.load_file = load_file;
function add(num1, num2) {
    return rpc_call("add", num1, num2);
}
function printNum(num) {
    return rpc_call("printNum", num);
}
function load_file(filename) {
    return rpc_call("load_file", filename);
}
