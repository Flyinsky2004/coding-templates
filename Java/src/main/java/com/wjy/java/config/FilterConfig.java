/*
 * @Author: Flyinsky 2084151024@qq.com
 * @Date: 2025-03-04 19:49:46
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2025-03-04 19:58:44
 * @FilePath: /coding-templates/Java/src/main/java/com/wjy/java/config/FilterConfig.java
 */
package com.wjy.java.config;

import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.wjy.java.filter.JwtFilter;

@Configuration
public class FilterConfig {

    @Bean
    public FilterRegistrationBean<JwtFilter> jwtFilterRegistration() {
        FilterRegistrationBean<JwtFilter> registration = new FilterRegistrationBean<>();
        registration.setFilter(new JwtFilter());
        registration.setOrder(1); // 设置过滤器的优先级，越小优先级越高
        registration.addUrlPatterns("/*");
        registration.setAsyncSupported(true); // 启用异步支持
        return registration;
    }
}
