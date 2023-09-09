package cmos.middleware.rsp.controller;

import lombok.Data;
import lombok.extern.slf4j.Slf4j;

// 添加 Slf4j 注解后, 不必再显式创建类成员字段的 getter/setter 方法.

@Data
@Slf4j
public class ResponseData<T> {
    // 0: 成功; 其他: 失败
    // 默认为成功, 如设置为 1, 则需要与 message 一起设置.
    Integer status;
    String message;
    T data;

    public ResponseData() {
        this.status = 0;
        this.message = "";
        this.data = null;
    }
    public ResponseData(Integer status, String message) {
        this.status = status;
        this.message = message;
    }
    public ResponseData(Integer status, String message, T data) {
        this.status = status;
        this.message = message;
        this.data = data;
    }

}
