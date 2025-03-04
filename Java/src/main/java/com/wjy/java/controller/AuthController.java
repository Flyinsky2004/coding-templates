/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package com.wjy.java.controller;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpSession;
import lombok.extern.log4j.Log4j2;
import org.springframework.web.bind.annotation.*;

import com.wjy.java.entity.dto.RestBean;
import com.wjy.java.entity.pojo.User;
import com.wjy.java.service.UserService;
import com.wjy.java.utility.JwtUtil;
@Log4j2
@RestController
@RequestMapping("api/auth")
public class AuthController {
    @Resource
    UserService userService;
    @PostMapping("login")
    public RestBean<String> login(@RequestParam("username")String username,
                                  @RequestParam("password")String password,
                                  HttpSession session){
        User user = userService.login(username,password);
        if(user == null) return RestBean.failure(403,"账号或密码错误!");
        session.setAttribute("user",user);
        String token = JwtUtil.createToken(user);
        return RestBean.success("登录成功!",token);
    }

    // @PostMapping("register")
    // public RestBean<String> register(@RequestParam("username")String username,
    //                                  @RequestParam("password")String password,
    //                                  @RequestParam("code")String code,
    //                                  @RequestParam("email")String email,
    //                                  HttpSession session) {
    //     if(username.isEmpty()||password.isEmpty()||email.isEmpty()) return RestBean.failure(400,"注册信息非法");
    //     String s = userService.register(new User(username,password,email,"https://imgs.flyinsky.wiki/file/bfce749c59dadda9e24ef.png"),session,code) ;
    //     if( s == null){
    //         return RestBean.success("注册成功!欢迎您,"+username);
    //     }else{
    //         return RestBean.failure(400,s);
    //     }
    // }

    // @PostMapping("code")
    // public RestBean<String> getCode(@RequestParam("email")String email, HttpSession session) {
    //     if(userService.sendMailAuthCode(email,session.getId())==null) return RestBean.success("验证码成功发送!请前往邮箱查看!");
    //     else return RestBean.failure(500,"发生错误,请稍后重试或联系管理员");
    // }

//    @GetMapping("logout")
//    public RestBean<String> logout(HttpSession session) {
//        User user = (User)session.getAttribute("user");
//        if(user == null) return RestBean.failure(403,"登出失败!原因:没有找到用户登录信息。");
//        session.invalidate();
//        return RestBean.success("登出成功!");
//    }
}
