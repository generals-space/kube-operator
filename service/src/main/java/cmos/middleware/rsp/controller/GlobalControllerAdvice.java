package cmos.middleware.rsp.controller;

import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.validation.ConstraintViolation;
import javax.validation.ConstraintViolationException;

import org.springframework.http.HttpStatus;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.validation.BindException;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

// 如果不添加该类和注解, 参数校验失败只会在日志里打印, spring 会自动返回异常信息的响应体, 格式如下:
// {
//     "timestamp": "",
//     "status": "500",
//     "error": "Internal Server Error",
//     "trace": "异常信息栈"
// }
// 这里统一返回为 ResponseData{status, message, data}
@RestControllerAdvice
public class GlobalControllerAdvice {
    public static final String BAD_REQUEST_MSG = "客户端请求参数错误";
    // <1> 处理 form data方式调用接口校验失败抛出的异常 
    @ExceptionHandler(BindException.class)
    public ResponseData bindExceptionHandler(BindException e) {
        List<FieldError> fieldErrors = e.getBindingResult().getFieldErrors();
        List exceptionList = fieldErrors.stream().map(o -> o.getDefaultMessage()).collect(Collectors.toList());
        return new ResponseData<>(HttpStatus.BAD_REQUEST.value(), BAD_REQUEST_MSG, exceptionList);
    }
    // <2> 处理 json 请求体调用接口校验失败抛出的异常 
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseData methodArgumentNotValidExceptionHandler(MethodArgumentNotValidException e) {
        List<FieldError> fieldErrors = e.getBindingResult().getFieldErrors();
        List exceptionList = fieldErrors.stream().map(o -> o.getDefaultMessage()).collect(Collectors.toList());
        return new ResponseData<>(3, BAD_REQUEST_MSG, exceptionList);
    }
    // <3> 处理单个参数校验失败抛出的异常
    @ExceptionHandler(ConstraintViolationException.class)
    public ResponseData constraintViolationExceptionHandler(ConstraintViolationException e) {
        Set<ConstraintViolation<?>> constraintViolations = e.getConstraintViolations();
        List exceptionList = constraintViolations.stream().map(o -> o.getMessage()).collect(Collectors.toList());
        return new ResponseData<>(3, BAD_REQUEST_MSG, exceptionList);
    }

    // <4> http json 请求体解析失败(如出现非法字符等)
    @ExceptionHandler(HttpMessageNotReadableException.class)
    public ResponseData<String> invalidHttpRequestBodyExceptionHandler(HttpMessageNotReadableException e) {
        e.printStackTrace();
        return new ResponseData<>(3, "参数类型错误", e.getLocalizedMessage());
    }
    // 捕获其他所有异常
    // @ExceptionHandler(Exception.class)
    // public ResponseData<String> handleException(Exception e) {
    //     e.printStackTrace();
    //     return new ResponseData<>(3, "未知异常", e.getLocalizedMessage());
    // }
    @ExceptionHandler(Exception.class)
    public ResponseData<String> handleException(HttpServletRequest req, Exception e, HttpServletResponse resp) {
        e.printStackTrace();

        // 在 spring boot 接口返回 json 响应体时, 有时会出现序列化失败的情况, 比如 resp.setData(deploy) 返回一个 V1Deployment 对象.
        // spring boot 默认的 jackson 解析失败, 异常就被 handleException() 捕获, 原本解析出来的内容被当作异常信息再一次被封装,
        // 于是就会出现响应体中包含2个 ResponseData 结构的情况
        //
        // {
        //     "status": 0,
        //     "message": "",
        //     "data": {V1Deployment 对象}
        // }{
        //     "status": 3,
        //     "message": "未知异常",
        //     "data": "Could not write JSON: Not an integer; nested exception is com.fasterxml.jackson.databind.JsonMappingException: Not an integer (strategy.rollingUpdate.maxSurge)"
        // }
        //
        // 这里在捕获异常后, 把遗留的解析内容先清空再返回就可以了.
        // 参考文章 [spring boot 2.4.5版本「全局统一处理异常封装」导致返回双重值问题解决](https://blog.csdn.net/u010739100/article/details/118071368)
        try {
            resp.resetBuffer();
        }catch (Exception ex){
            e.printStackTrace();
        }
        return new ResponseData<>(3, "未知异常", e.getLocalizedMessage());
    }
}
